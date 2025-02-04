package coreum

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"net/url"
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
	Network metadata.Network
	// Read transactions get dumped into this socket for processing thus decoupling the block reader from any business logic
	ProcessBlockChannel chan *ScannedBlock
	BlockHeight         int64
	LastBlockTime       time.Time
	BlockProductionTime time.Duration
}

const MinimumBlockProductionTime = 100 * time.Millisecond
const MaximumBlockProductionTime = 10 * time.Second

func NewReader(network metadata.Network) *Reader {
	return &Reader{
		Network:             network,
		ProcessBlockChannel: make(chan *ScannedBlock, 1000),
		LastBlockTime:       time.Now(),
		BlockProductionTime: MinimumBlockProductionTime,
	}
}

type ScannedBlock struct {
	BlockEvents  []types.Event
	Transactions []*txtypes.GetTxResponse
	BlockHeight  int64
	BlockTime    time.Time
}

var nodeConnections map[metadata.Network]*client.Context

// Provide Readers without a blockheight to start from
func InitReaders() Readers {
	readers := make(Readers)
	// Setup grpc connections:
	nodeConnections = NewNodeConnections()
	for network := range nodeConnections {
		reader := NewReader(network)
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
			currentHeight := reader.BlockHeight
			if currentHeight < 1 {
				panic("block height should be at least 1")
			}
			var err error
			for {
				currentHeight, err = reader.processBlock(txClient, rpcClient, currentHeight) // Process the block and increment the height
				if err != nil {
					if isTemporaryError(err) {
						logger.Errorf("error processing block %d. will retry: %v", currentHeight, err)
						time.Sleep(1 * time.Second)
						continue
					}
					panic(errors.Wrapf(err, "error processing block %d", currentHeight))
				}
				currentHeight++
				reader.BlockHeight = currentHeight
			}
		}(reader)
	}
}

func isTemporaryError(err error) bool {
	var netOpError *net.OpError
	if errors.As(err, &netOpError) {
		return true
	}

	var urlError *url.Error
	if errors.As(err, &urlError) {
		if urlError.Err.Error() == "EOF" {
			return true
		}
	}

	if statusError, ok := status.FromError(err); ok {
		if statusError.Code() == codes.Unavailable {
			return true
		}
	}
	return false
}

var (
	measureTotalThroughputTime = time.Now()
	measureBlockLoadTime       = time.Now()
	measureTXLoadTime          = time.Now()
	previousHeight             int64
	measureTotalTransactions   int // Measure of total transacions loaded in the last 100 blocks (for performance monitoring)
)

func (r *Reader) processBlock(txClient txtypes.ServiceClient, rpcClient sdkclient.CometRPC, currentHeight int64) (int64, error) {
	ctx := context.Background()
	sb := &ScannedBlock{
		BlockHeight:  currentHeight,
		Transactions: make([]*txtypes.GetTxResponse, 0),
		BlockEvents:  make([]types.Event, 0),
	}
	var goroutineError error
	wg := sync.WaitGroup{}
	mutex := &sync.Mutex{}
	wg.Add(2)
	// Querying block from Coreum to get transactions
	go func(m *sync.Mutex) {
		tStart := time.Now()
		bhr, err := txClient.GetBlockWithTxs(ctx, &txtypes.GetBlockWithTxsRequest{Height: currentHeight})
		if err != nil && strings.Contains(err.Error(), "height must not be less than 1 or greater than the current height") {
			for strings.Contains(err.Error(), "height must not be less than 1 or greater than the current height") {
				<-time.After(r.BlockProductionTime)
				bhr, err = txClient.GetBlockWithTxs(ctx, &txtypes.GetBlockWithTxsRequest{Height: currentHeight})
			}
		}
		if err != nil {
			logger.Errorf("grpc error: %v", err)
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
		measureBlockLoadTime = measureBlockLoadTime.Add(time.Since(tStart))
		wg2 := sync.WaitGroup{}
		wg2.Add(len(bhr.Block.Data.Txs))
		tStart = time.Now() // Start of TX loading
		for _, tx := range bhr.Block.Data.Txs {
			go func(tx []byte) {
				hr := hash(tx)
				v, err := txClient.GetTx(context.Background(), &txtypes.GetTxRequest{
					Hash: hr,
				})
				if err != nil {
					logger.Errorf("error getting tx %s: %v", hr, err)
					m.Lock()
					goroutineError = err
					m.Unlock()
					wg2.Done()
					return
				}
				m.Lock()
				sb.Transactions = append(sb.Transactions, v)
				m.Unlock()
				wg2.Done()
			}(tx)
		}
		wg2.Wait()
		wg.Done()
		// Add the time processed to the measureBlockLoadTime
		measureTXLoadTime = measureTXLoadTime.Add(time.Since(tStart))
	}(mutex)
	// Querying block results from Tendermint to get block events
	go func(m *sync.Mutex) {
		br, err := rpcClient.BlockResults(ctx, &currentHeight)
		if err != nil && strings.Contains(err.Error(), "must be less than or equal to the current blockchain height") {
			for strings.Contains(err.Error(), "must be less than or equal to the current blockchain height") {
				<-time.After(r.BlockProductionTime)
				br, err = rpcClient.BlockResults(ctx, &currentHeight)
			}
		}
		if err != nil {
			logger.Errorf("rpc error for block %d: %v", currentHeight, err)
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
	}(mutex)
	wg.Wait()
	if goroutineError != nil {
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
	measureTotalTransactions += len(sb.Transactions)

	if currentHeight%100 == 0 {
		if previousHeight != 0 {
			channelCapacity := 100 - 100*float64(len(r.ProcessBlockChannel))/1000 // Percentage of channel capacity used: capacity is 1000
			logger.Infof("TotalTime: %2.f seconds. Loading %d blocks using %2.f seconds, loading %d TX using %2.f seconds, channel capacity left %2.f (percentage) (indicates blocking on processing of TX)",
				time.Since(measureTotalThroughputTime).Seconds(),
				currentHeight-previousHeight,
				measureBlockLoadTime.Sub(measureTotalThroughputTime).Seconds(),
				measureTotalTransactions,
				measureTXLoadTime.Sub(measureTotalThroughputTime).Seconds(),
				channelCapacity)
		}
		previousHeight = currentHeight
		measureBlockLoadTime = time.Now()
		measureTXLoadTime = time.Now()
		measureTotalTransactions = 0
		measureTotalThroughputTime = time.Now()
	}
	return currentHeight, nil
}

func hash(txRaw []byte) string {
	h := sha256.New()
	h.Write(txRaw)
	return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
}
