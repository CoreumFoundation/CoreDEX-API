package order

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	sdecimal "github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	currencygrpclient "github.com/CoreumFoundation/CoreDEX-API/domain/currency/client"
	decimal "github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	ordergrpcclient "github.com/CoreumFoundation/CoreDEX-API/domain/order/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	"github.com/CoreumFoundation/coreum/v5/pkg/client"
)

type cache struct {
	mutex *sync.RWMutex
	data  map[string]*dmn.LockableCache
}

type Application struct {
	TxEncoder      map[metadata.Network]txClient
	currencyClient currencygrpc.CurrencyServiceClient
	orderClient    ordergrpc.OrderServiceClient
	orderbookCache *cache
}

type txClient struct {
	txFactory     tx.Factory
	clientContext client.Context
	reader        *coreum.Reader
}

type WalletAsset struct {
	Denom        string
	Amount       string
	SymbolAmount string
}

func NewApplication() *Application {
	currencyClient := currencygrpclient.Client()
	orderbookClient := ordergrpcclient.Client()
	return NewApplicationWithClients(currencyClient, orderbookClient)
}

func NewApplicationWithClients(currencyClient currencygrpc.CurrencyServiceClient,
	orderClient ordergrpc.OrderServiceClient) *Application {
	txEncoders := make(map[metadata.Network]txClient)
	coreum.InitReaders()
	nodeConnections := coreum.NewNodeConnections()
	for network, clientCtx := range nodeConnections {
		txFactory := client.Factory{}.
			WithKeybase(clientCtx.Keyring()).
			WithChainID(clientCtx.ChainID()).
			WithTxConfig(clientCtx.TxConfig()).
			WithSimulateAndExecute(true).
			WithGasAdjustment(2)
		txEncoders[network] = txClient{
			txFactory:     txFactory,
			clientContext: *clientCtx,
			reader:        coreum.NewReader(network, clientCtx),
		}
	}
	orderbookCache := &cache{
		mutex: &sync.RWMutex{},
		data:  make(map[string]*dmn.LockableCache),
	}
	go dmn.CleanCache(orderbookCache.data, orderbookCache.mutex, 15*time.Minute)
	return &Application{txEncoders, currencyClient, orderClient, orderbookCache}
}

func (a *Application) EncodeTx(network metadata.Network, from sdk.AccAddress, msgs ...sdk.Msg) ([]byte, error) {
	unsignedTx, err := client.GenerateUnsignedTx(
		context.Background(),
		a.TxEncoder[network].clientContext.WithFromAddress(from).WithUnsignedSimulation(true),
		a.TxEncoder[network].txFactory,
		msgs...,
	)
	if err != nil {
		return nil, err
	}

	encoder := a.TxEncoder[network].clientContext.TxConfig().TxEncoder()
	if encoder == nil {
		return nil, errors.New("cannot print unsigned tx: tx encoder is nil")
	}

	return encoder(unsignedTx.GetTx())
}

func (a *Application) SubmitTx(network metadata.Network, rawTx []byte) (*sdk.TxResponse, error) {
	return client.BroadcastRawTx(context.Background(), a.TxEncoder[network].clientContext, rawTx)
}

func (a *Application) AccountSequence(network metadata.Network, address string) (uint64, error) {
	clientCtx := a.TxEncoder[network].clientContext

	req := &authtypes.QueryAccountRequest{
		Address: address,
	}
	authQueryClient := authtypes.NewQueryClient(clientCtx)
	ctx := context.Background()
	res, err := authQueryClient.Account(ctx, req)
	if err != nil {
		logger.Errorf("Error querying account %s: %v", address, err)
		return 0, err
	}

	var acc sdk.AccountI
	if err := clientCtx.InterfaceRegistry().UnpackAny(res.Account, &acc); err != nil {
		logger.Errorf("Error unpacking account: %v", err)
		return 0, err
	}

	return acc.GetSequence(), nil
}

func orderbookCacheKey(denom1, denom2 string) string {
	return fmt.Sprintf("%s-%s", denom1, denom2)
}

// Cache the orderbooks so that the subsequent data can be gotten from the database which holds the latest orders
// The database is a lost faster than the blockchain when it comes to reading data.
// The read from the database reads with an overlap in time such that eventually skipped orders in a previous read,
// will be read in the next read.
func (a *Application) OrderBookRelevantOrders(network metadata.Network, denom1, denom2 string, limit int, aggregate bool) (*coreum.OrderBookOrders, error) {
	processStart := time.Now() // Determine what time we need to retrieve data for from the	database
	key := orderbookCacheKey(denom1, denom2)
	orderbook := &coreum.OrderBookOrders{}
	a.orderbookCache.mutex.Lock()
	if cache, ok := a.orderbookCache.data[key]; ok {
		orderbook = cache.Value.(*coreum.OrderBookOrders)
		processStart = cache.LastUpdated
	}
	/* 2 scenarios:
	* cache is empty
	* cache is not empty

	In case of empty, get from the source, and then update with the database latest state since we started the read
	In case of not empty, apply the database state from the last refresh moment.
	Applying the database state can lead to both adding and removing orders from the cache.
	It is assumed that the database does not contain all the orders (due to the inception time of the database possibly
	being after the inception of the given orderbook), so only orders with remaining quantity of 0 are removed from the orderbook.
	*/
	// Set the time to be used for the moment of the update of the orderbook. If updates are "slow" this takes care of filling the gap
	tStartUpdate := time.Now()
	if orderbook == nil || (len(orderbook.Buy) == 0 && len(orderbook.Sell) == 0) {
		processStart = time.Now()
		ctx := context.Background()
		denom1Currency, err := a.currencyClient.Get(ctx, &currencygrpc.ID{
			Network: network,
			Denom:   denom1,
		})
		if err != nil {
			a.orderbookCache.mutex.Unlock()
			return nil, err
		}
		denom2Currency, err := a.currencyClient.Get(ctx, &currencygrpc.ID{
			Network: network,
			Denom:   denom2,
		})
		if err != nil {
			a.orderbookCache.mutex.Unlock()
			return nil, err
		}
		denom1Precision := int64(0)
		if denom1Currency.Denom != nil && denom1Currency.Denom.Precision != nil {
			denom1Precision = int64(*denom1Currency.Denom.Precision)
		}
		denom2Precision := int64(0)
		if denom2Currency.Denom != nil && denom2Currency.Denom.Precision != nil {
			denom2Precision = int64(*denom2Currency.Denom.Precision)
		}

		orderbook, err = a.TxEncoder[network].reader.QueryOrderBookRelevantOrders(ctx, denom1, denom2, denom1Precision, denom2Precision, uint64(limit), aggregate)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				a.orderbookCache.mutex.Unlock()
				return nil, fmt.Errorf("there is no orderbook for %s - %s", denom1, denom2)
			}
			a.orderbookCache.mutex.Unlock()
			return nil, err
		}
	}
	// If the process has taken more than a second, update the orderbook with the latest orders from the database
	// Or if the data was retrieved from the cache (and more than 1 second has passed since the last update), update the orderbook
	if time.Since(processStart) > time.Second {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		denom1Currency, err := a.currencyClient.Get(ctx, &currencygrpc.ID{
			Network: network,
			Denom:   denom1,
		})
		if err != nil {
			a.orderbookCache.mutex.Unlock()
			return nil, err
		}
		denom2Currency, err := a.currencyClient.Get(ctx, &currencygrpc.ID{
			Network: network,
			Denom:   denom2,
		})
		if err != nil {
			a.orderbookCache.mutex.Unlock()
			return nil, err
		}

		orders, err := a.orderClient.GetAll(ctx, &ordergrpc.Filter{
			Network: network,
			Denom1:  denom1Currency.Denom,
			Denom2:  denom2Currency.Denom,
			// This causes a slight overlap in data retrieved with the previous read, which is on purpose:
			// The process writing the data is writing for a "previous" block, and we use block time to determine the time to read from
			From: timestamppb.New(processStart.Add(-5 * time.Second)),
			To:   timestamppb.Now(), // This causes a slight overlap in data retrieved with the next read, which is on purpose
		})
		if err != nil {
			a.orderbookCache.mutex.Unlock()
			return nil, err
		}
		// Orders have a status, and a remaining quantity. If the remaining quantity is 0, the order is removed from the orderbook
		// If the order is not in the orderbook, it is added to the orderbook
		// If the order is in the orderbook, it is updated with the new data
		buySide := orderbook.Buy
		sellSide := orderbook.Sell
		buySideRemove := make([]uint64, 0)
		buySideAppend := make([]*coreum.OrderBookOrder, 0)
		sellSideRemove := make([]uint64, 0)
		sellSideAppend := make([]*coreum.OrderBookOrder, 0)
		for _, order := range orders.Orders {
			a.processOrderForOrderBook(buySide, order, denom1Currency, denom2Currency, buySideRemove, buySideAppend)
			a.processOrderForOrderBook(sellSide, order, denom1Currency, denom2Currency, sellSideRemove, sellSideAppend)
		}
		for _, removeID := range buySideRemove {
			for i, buyOrder := range buySide {
				if buyOrder.Sequence == removeID {
					buySide = append(buySide[:i], buySide[i+1:]...)
					break // IDs only appear once in the orderbook
				}
			}
		}
		buySide = append(buySide, buySideAppend...)
		for _, removeID := range sellSideRemove {
			for i, o := range buySide {
				if o.Sequence == removeID {
					sellSide = append(sellSide[:i], sellSide[i+1:]...)
					break // IDs only appear once in the orderbook
				}
			}
		}
		sellSide = append(sellSide, sellSideAppend...)
		orderbook.Buy = buySide
		orderbook.Sell = sellSide
	}
	// Set the orderbook into the cache:
	a.orderbookCache.data[key] = &dmn.LockableCache{
		LastUpdated: tStartUpdate,
		Value:       orderbook,
	}
	a.orderbookCache.mutex.Unlock()
	return orderbook, nil
}

func (*Application) processOrderForOrderBook(side []*coreum.OrderBookOrder,
	order *ordergrpc.Order,
	denom1Currency *currencygrpc.Currency,
	denom2Currency *currencygrpc.Currency,
	removeList []uint64,
	appendList []*coreum.OrderBookOrder,
) ([]uint64, []*coreum.OrderBookOrder) {
	for _, o := range side {
		if o.Sequence == uint64(order.Sequence) {
			if (*order.RemainingQuantity).IsZero() || order.OrderStatus == ordergrpc.OrderStatus_ORDER_STATUS_CANCELED || order.OrderStatus == ordergrpc.OrderStatus_ORDER_STATUS_FILLED {
				removeList = append(removeList, o.Sequence)
			} else {
				denom1Precision := *denom1Currency.Denom.Precision
				denom2Precision := *denom2Currency.Denom.Precision
				price := sdecimal.NewFromFloat(order.Price)
				var precision sdecimal.Decimal
				precisionDiff := denom1Precision - denom2Precision
				if precisionDiff < 0 {
					precision = sdecimal.NewFromInt(1).Div(sdecimal.New(1, int32(-precisionDiff)))
				} else if precisionDiff > 0 {
					precision = sdecimal.New(1, int32(-precisionDiff))
				} else {
					precision = sdecimal.NewFromInt(1)
				}

				humanReadablePrice := price.Mul(precision)
				symbolAmount := *decimal.ToSDec(order.Quantity)
				symbolAmount = symbolAmount.Div(sdecimal.New(1, int32(denom1Precision)))
				remainingSymbolAmount := *decimal.ToSDec(order.RemainingQuantity)
				remainingSymbolAmount = remainingSymbolAmount.Div(sdecimal.New(1, int32(denom1Precision)))

				appendList = append(appendList, &coreum.OrderBookOrder{
					Price:                 fmt.Sprintf("%f", order.Price),
					HumanReadablePrice:    humanReadablePrice.String(),
					Amount:                order.Quantity.String(),
					SymbolAmount:          symbolAmount.String(),
					Sequence:              uint64(order.Sequence),
					Account:               order.Account,
					OrderID:               order.OrderID,
					RemainingAmount:       order.RemainingQuantity.String(),
					RemainingSymbolAmount: remainingSymbolAmount.String(),
				})
			}
		}
	}
	return removeList, appendList
}

func (a *Application) OrderBookRelevantOrdersForAccount(network metadata.Network, denom1, denom2, account string) (*coreum.OrderBookOrders, error) {
	relevantOrdersForAllAccounts, err := a.OrderBookRelevantOrders(network, denom1, denom2, 100, false) // Get more orders so that we can retrieve all potential orders for this account
	if err != nil {
		return nil, err
	}
	// Iterate over all orders and filter out the orders for the account
	relevantOrdersForAccount := &coreum.OrderBookOrders{
		Buy:  make([]*coreum.OrderBookOrder, 0),
		Sell: make([]*coreum.OrderBookOrder, 0),
	}
	for _, order := range relevantOrdersForAllAccounts.Buy {
		if order.Account == account {
			relevantOrdersForAccount.Buy = append(relevantOrdersForAccount.Buy, order)
		}
	}
	for _, order := range relevantOrdersForAllAccounts.Sell {
		if order.Account == account {
			relevantOrdersForAccount.Sell = append(relevantOrdersForAccount.Sell, order)
		}
	}
	return relevantOrdersForAccount, nil

}

func (a *Application) WalletAssets(network metadata.Network, address string) ([]WalletAsset, error) {
	coins := sdk.Coins{}
	bankClient := banktypes.NewQueryClient(a.TxEncoder[network].clientContext)
	var paginationKey []byte = nil
	for {
		res, err := bankClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{
			Address:      address,
			Pagination:   &query.PageRequest{Key: paginationKey},
			ResolveDenom: false,
		})
		if err != nil {
			return nil, err
		}
		coins = coins.Add(res.Balances...)
		paginationKey = res.Pagination.NextKey
		if paginationKey == nil {
			break
		}
	}
	// Transform the coins to WalletAsset and add the symbol amount (apply precision)
	walletAssets := make([]WalletAsset, 0)
	for _, coin := range coins {
		denomCurrency, err := a.currencyClient.Get(context.Background(), &currencygrpc.ID{
			Network: network,
			Denom:   coin.Denom,
		})
		if err != nil {
			return nil, err
		}
		precision := int32(0)
		if denomCurrency.Denom != nil && denomCurrency.Denom.Precision != nil {
			precision = *denomCurrency.Denom.Precision
		}
		walletAssets = append(walletAssets, WalletAsset{
			Denom:        coin.Denom,
			Amount:       coin.Amount.String(),
			SymbolAmount: sdecimal.NewFromBigInt(coin.Amount.BigInt(), 0).Div(sdecimal.New(1, precision)).String(),
		})
	}
	return walletAssets, nil
}
