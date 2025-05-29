package coreum

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cometbft/cometbft/abci/types"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	"github.com/CoreumFoundation/coreum/v5/pkg/client"
)

type Readers map[metadata.Network]*Reader

type Reader struct {
	Network       metadata.Network
	ClientContext *client.Context
	// Read transactions get dumped into this socket for processing thus decoupling the block reader from any business logic
	ProcessBlockChannel        chan *ScannedBlock
	BlockHeight                int64
	LastBlockTime              time.Time
	BlockProductionTime        time.Duration
	previousHeight             int64
	currentHeight              int64
	measureTotalThroughputTime time.Time
	measureBlockLoadTime       time.Time
	measureTXLoadTime          time.Time
	measureTotalTransactions   int // Measure of total transactions loaded in the last 100 blocks (for performance monitoring)
	mutex                      *sync.Mutex
	atEndOfChain               bool // Flag to indicate if we are at the end of the chain scanning realtime (determines waits to prevent querying the chain for not yet produced blocks)
}

const MinimumBlockProductionTime = 100 * time.Millisecond
const MaximumBlockProductionTime = 10 * time.Second

func NewReader(network metadata.Network, clientContext *client.Context) *Reader {
	return &Reader{
		Network:                    network,
		ClientContext:              clientContext,
		ProcessBlockChannel:        make(chan *ScannedBlock, 1000),
		LastBlockTime:              time.Now(),
		BlockProductionTime:        MinimumBlockProductionTime,
		measureTotalThroughputTime: time.Now(),
		measureBlockLoadTime:       time.Now(),
		measureTXLoadTime:          time.Now(),
		mutex:                      &sync.Mutex{},
		atEndOfChain:               false,
	}
}

type ScannedBlock struct {
	BlockEvents  []types.Event
	Transactions []*txtypes.GetTxResponse
	BlockHeight  int64
	BlockTime    time.Time
}

var (
	nodeConnections map[metadata.Network]*client.Context
	heightRegex     = regexp.MustCompile(`height (\d+) is not available, lowest height is (\d+)`)
)

// Provide Readers without a blockheight to start from
func InitReaders() Readers {
	readers := make(Readers)
	// Setup grpc connections:
	nodeConnections = NewNodeConnections()
	for network, clientCtx := range nodeConnections {
		reader := NewReader(network, clientCtx)
		readers[network] = reader
	}
	return readers
}

// Reads the history in a controlled fashion. To be used if more history is required than the start of the application indicated
func (*Reader) ReadHistory(startBlockHeight, endBlockHeight int64) {
	// Read the blocks at a reasonable pace but such that we do not overload the memory: The trick will be in sizing of the buffered channel we dump the transactions in
	// Since that is the same channel as for real time blocks, we need to make certain our realtime blocks have preference over these historical blocks
}

// Reads blocks from a certain height and stays up to date with the last block as produced by the chain
func (r Readers) Start() {
	for _, reader := range r {
		go func(reader *Reader) {
			txClient := txtypes.NewServiceClient(nodeConnections[reader.Network])
			rpcClient := nodeConnections[reader.Network].RPCClient()
			reader.currentHeight = reader.BlockHeight
			if reader.currentHeight < 1 {
				panic("block height should be at least 1")
			}
			go reader.Logger()
			var err error
			for {
				reader.currentHeight, err = reader.processBlock(txClient, rpcClient, reader.currentHeight) // Process the block and increment the height
				if err != nil {
					if isTemporaryError(err) {
						logger.Errorf("error processing block %d. will retry: %v", reader.currentHeight, err)
						time.Sleep(1 * time.Second)
						continue
					}
					v, err := getValidBlockHeight(err)
					if err == nil {
						reader.currentHeight = v
						logger.Warnf("setting block height to %d", reader.currentHeight, err)
						continue
					}
					// We allow processing to continue (blockchain is responding, this code just has an issue with the data: This leads to dataloss)
					if isIgnorableError(err) {
						reader.currentHeight++
						reader.BlockHeight = reader.currentHeight
						continue
					}
					// panic: error processing block 6840526: rpc error: code = DeadlineExceeded desc = received context error while waiting for new LB policy update: context deadline exceeded
					panic(errors.Wrapf(err, "error processing block %d", reader.currentHeight))
				}
				reader.currentHeight++
				reader.BlockHeight = reader.currentHeight
			}
		}(reader)
	}
}

/*
isIgnorableError checks if the error is ignorable but can lead to missing data (so it is advisable to resolve it)
Types of errors managed:
* tx parse error: unable to resolve type URL /coreum.nft.v1beta1.MsgSend
*/
func isIgnorableError(err error) bool {
	if err == nil {
		return true
	}
	if strings.Contains(err.Error(), "unable to resolve type URL") {
		// This is a parse error thcd at can be ignored, but it is advisable to resolve it
		logger.Warnf("DATA LOSS MIGHT OCCUR: ignoring parse error: %s", err.Error())
		return true
	}
	return false
}

// There is the possibility of an init error for the blockheight
// => height 6599262 is not available, lowest height is 6603501
// This function parses the lowest height from this string if the error is in the correct format
func getValidBlockHeight(err error) (int64, error) {
	matches := heightRegex.FindStringSubmatch(err.Error())
	if len(matches) > 1 {
		// Extract the captured groups
		h, err := strconv.ParseInt(matches[2], 10, 64) // This is the height that is not available
		if err != nil {
			return -1, errors.Wrapf(err, "error parsing block height from error: %s", err.Error())
		}
		return h, nil
	}
	return -1, err
}

// isTemporaryError checks if the error is a temporary error that can be retried
func isTemporaryError(err error) bool {
	var netOpError *net.OpError
	if errors.As(err, &netOpError) {
		return true
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var urlError *url.Error
	if errors.As(err, &urlError) {
		if urlError.Err.Error() == "EOF" {
			return true
		}
	}

	if statusError, ok := status.FromError(err); ok {
		switch statusError.Code() {
		case codes.Unavailable, codes.NotFound:
			return true
		}
	}
	return false
}

func (r *Reader) processBlock(txClient txtypes.ServiceClient, rpcClient sdkclient.CometRPC, currentHeight int64) (int64, error) {
	sb := &ScannedBlock{
		BlockHeight:  currentHeight,
		Transactions: make([]*txtypes.GetTxResponse, 0),
		BlockEvents:  make([]types.Event, 0),
	}
	var goroutineError error
	wg := sync.WaitGroup{}
	wg.Add(2) // 2 main go routines
	if r.atEndOfChain {
		// Wait for the block to be produced by waiting the for block production time since the last block
		// By wiating like this we can manage the delay between the blocks (The time is since the last reported blocktime, not the clock of the server)
		// This is not perfect: There is the 3x block production time to wait for while 1x should be the exact perfect timing (but it is not for yet unknown reasons)
		// 2x reduces the error of not yet having the block produced from 50% of the time to 25% of the time
		if time.Now().Before(r.LastBlockTime.Add(3 * r.BlockProductionTime)) {
			<-time.After(r.LastBlockTime.Add(3 * r.BlockProductionTime).Sub(time.Now()))
		}
	}
	// context with timeout to counter slow chain response on mainly devnet:
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Querying block from Coreum to get transactions
	go func(m *sync.Mutex) {
		tStart := time.Now()
		bhr, err := txClient.GetBlockWithTxs(ctx, &txtypes.GetBlockWithTxsRequest{Height: currentHeight})
		if isBlockWithTxEndErr(err) {
			for isBlockWithTxEndErr(err) {
				r.atEndOfChain = true
				logger.Warnf("%s: Problem getting block %d: %v", r.Network.String(), currentHeight, err)
				<-time.After(r.BlockProductionTime)
				bhr, err = txClient.GetBlockWithTxs(ctx, &txtypes.GetBlockWithTxsRequest{Height: currentHeight})
			}
		}
		if err != nil {
			logger.Errorf("%s: processBlock error: %v", r.Network.String(), err)
			m.Lock()
			goroutineError = err
			m.Unlock()
			wg.Done()
			return
		}
		m.Lock()
		sb.BlockTime = bhr.Block.Header.Time
		m.Unlock()
		// Add the time processed to the measureBlockLoadTime for aggregate timing overview
		r.measureBlockLoadTime = r.measureBlockLoadTime.Add(time.Since(tStart))
		tStart = time.Now() // Start of TX loading
		for _, tx := range bhr.Block.Data.Txs {
			wg.Add(1)
			go func(tx []byte, m *sync.Mutex) {
				hr := hash(tx)
				v, err := txClient.GetTx(context.Background(), &txtypes.GetTxRequest{
					Hash: hr,
				})
				if err != nil {
					logger.Errorf("%s: error getting tx %s: %v", r.Network.String(), hr, err)
					m.Lock()
					goroutineError = err
					m.Unlock()
					wg.Done()
					return
				}
				m.Lock()
				sb.Transactions = append(sb.Transactions, v)
				m.Unlock()
				wg.Done()
			}(tx, m)
		}
		wg.Done()
		// Add the time processed to the measureBlockLoadTime
		r.measureTXLoadTime = r.measureTXLoadTime.Add(time.Since(tStart))
	}(r.mutex)
	// Querying block results from Tendermint to get block events
	go func(m *sync.Mutex) {
		br, err := rpcClient.BlockResults(ctx, &currentHeight)
		if isBlockResultEndErr(err) {
			for isBlockResultEndErr(err) {
				<-time.After(r.BlockProductionTime)
				br, err = rpcClient.BlockResults(ctx, &currentHeight)
			}
		}
		if err != nil {
			m.Lock()
			goroutineError = err
			m.Unlock()
			wg.Done()
			return
		}
		m.Lock()
		sb.BlockEvents = br.FinalizeBlockEvents
		m.Unlock()
		wg.Done()
	}(r.mutex)
	wg.Wait()
	if goroutineError != nil {
		logger.Errorf("%s: error processing block %d: %v", r.Network.String(), currentHeight, goroutineError)
		return currentHeight, goroutineError
	}
	r.BlockProductionTime = sb.BlockTime.Sub(r.LastBlockTime)
	if r.BlockProductionTime < MinimumBlockProductionTime {
		r.BlockProductionTime = MinimumBlockProductionTime
	} else if r.BlockProductionTime > MaximumBlockProductionTime {
		r.BlockProductionTime = MaximumBlockProductionTime
	}
	r.LastBlockTime = sb.BlockTime

	r.ProcessBlockChannel <- sb
	r.measureTotalTransactions += len(sb.Transactions)
	return currentHeight, nil
}

// Selective logging to keep insight in the data aggregators activity
func (r *Reader) Logger() {
	for {
		time.Sleep(5 * time.Second)
		r.mutex.Lock()
		if r.previousHeight != 0 {
			channelCapacity := 100 - 100*float64(len(r.ProcessBlockChannel))/1000 // Percentage of channel capacity used: capacity is 1000
			logger.Infof("%s: BlockHeight %d. TotalTime: %2.f seconds. Loading %d blocks using %2.f seconds, loading %d TX using %2.f seconds, channel capacity left %2.f (percentage) (indicates blocking on processing of TX)",
				r.Network.String(),
				r.currentHeight,
				time.Since(r.measureTotalThroughputTime).Seconds(),
				r.currentHeight-r.previousHeight,
				r.measureBlockLoadTime.Sub(r.measureTotalThroughputTime).Seconds(),
				r.measureTotalTransactions,
				r.measureTXLoadTime.Sub(r.measureTotalThroughputTime).Seconds(),
				channelCapacity)
		}
		r.previousHeight = r.currentHeight
		r.measureBlockLoadTime = time.Now()
		r.measureTXLoadTime = time.Now()
		r.measureTotalTransactions = 0
		r.measureTotalThroughputTime = time.Now()
		r.mutex.Unlock()
	}
}

func hash(txRaw []byte) string {
	h := sha256.New()
	h.Write(txRaw)
	return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
}

func isBlockWithTxEndErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "height must not be less than 1 or greater than the current height")
}

func isBlockResultEndErr(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "must be less than or equal to the current blockchain height") ||
		strings.Contains(err.Error(), "could not find results for height"))
}
