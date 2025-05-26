package order

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/samber/lo"
	dec "github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	currency "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/currency"
	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	dmncache "github.com/CoreumFoundation/CoreDEX-API/domain/cache"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	ordergrpcclient "github.com/CoreumFoundation/CoreDEX-API/domain/order/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	"github.com/CoreumFoundation/coreum/v5/pkg/client"
)

type cache struct {
	mutex *sync.RWMutex
	data  map[string]*dmncache.LockableCache
}

type Application struct {
	TxEncoder      map[metadata.Network]txClient
	orderClient    ordergrpc.OrderServiceClient
	orderbookCache *cache
	currencyClient currency.Application
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

type OrderBookOrder struct {
	OrderBookOrder *coreum.OrderBookOrder
	BaseDenom      *denom.Denom
	QuoteDenom     *denom.Denom
	Network        metadata.Network
	Side           orderproperties.Side
}

// Different way of cache expiration since we are updating keys all the time in the orderbook cache (unusual usage pattern)
// Which in the case a a skipped record can lead to "hanging" orders (never clearing until server restarts)
var (
	cacheExpire      = make(map[string]time.Time)
	cacheExpireMutex = &sync.Mutex{}
	orderBookMutex   = &sync.Mutex{}
)

func NewApplication(currencyClient *currency.Application) *Application {
	orderbookClient := ordergrpcclient.Client()
	return NewApplicationWithClients(orderbookClient, currencyClient)
}

func NewApplicationWithClients(orderClient ordergrpc.OrderServiceClient,
	currencyClient *currency.Application) *Application {

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
		data:  make(map[string]*dmncache.LockableCache),
	}
	go dmncache.CleanCache(orderbookCache.data, orderbookCache.mutex, 15*time.Minute)
	return &Application{txEncoders, orderClient, orderbookCache, *currencyClient}
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

func orderbookCacheKey(network metadata.Network, denom1, denom2 string) string {
	return fmt.Sprintf("%s-%s-%d", denom1, denom2, network)
}

// Cache the orderbooks so that the subsequent data can be gotten from the database which holds the latest orders
// The database is a lost faster than the blockchain when it comes to reading data.
// The read from the database reads with an overlap in time such that eventually skipped orders in a previous read,
// will be read in the next read.
func (a *Application) OrderBookRelevantOrders(network metadata.Network, denom1, denom2 string, limit int, aggregate bool) (*coreum.OrderBookOrders, error) {
	processStart := time.Now() // Determine what time we need to retrieve data for from the	database
	key := orderbookCacheKey(network, denom1, denom2)
	orderbook := &coreum.OrderBookOrders{}
	a.orderbookCache.mutex.RLock()
	if cache, ok := a.orderbookCache.data[key]; ok {
		orderbook = cache.Value.(*coreum.OrderBookOrders)
		processStart = cache.LastUpdated
	}
	a.orderbookCache.mutex.RUnlock()
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
	var err error
	orderbook, err = a.fetchOrderBookFromChain(orderbook, network, denom1, denom2, limit)
	if err != nil {
		return nil, err
	}
	// If the process has taken more than a second, update the orderbook with the latest orders from the database
	// Or if the data was retrieved from the cache (and more than 1 second has passed since the last update), update the orderbook
	if time.Since(processStart) > time.Second {
		orderbook, err = a.fetchOrderbookFromDatabase(orderbook, network, denom1, denom2, processStart)
		if err != nil {
			return nil, err
		}
	}
	// Order the buys and sales descending
	sort.Slice(orderbook.Buy, func(i, j int) bool {
		p1, _ := dec.NewFromString(orderbook.Buy[i].Price)
		p2, _ := dec.NewFromString(orderbook.Buy[j].Price)
		return p1.GreaterThan(p2)
	})
	sort.Slice(orderbook.Sell, func(i, j int) bool {
		p1, _ := dec.NewFromString(orderbook.Sell[i].Price)
		p2, _ := dec.NewFromString(orderbook.Sell[j].Price)
		return p1.GreaterThan(p2)
	})
	// Set the orderbook into the cache:
	a.orderbookCache.mutex.Lock()
	a.orderbookCache.data[key] = &dmncache.LockableCache{
		LastUpdated: tStartUpdate,
		Value:       orderbook,
	}
	a.orderbookCache.mutex.Unlock()
	if aggregate {
		// Clone the orderbook so that the original orderbook is not modified
		orderbookClone := &coreum.OrderBookOrders{
			Buy:  make([]*coreum.OrderBookOrder, 0),
			Sell: make([]*coreum.OrderBookOrder, 0),
		}
		// Orders are aggregated by price so that only one record exists for a given price (can reduce the number of records to be displayed)
		// This is done by summing up the quantities of orders with the same price
		orderbookClone.Sell = aggregateOrders(orderbookClone.Sell)
		orderbookClone.Buy = aggregateOrders(orderbookClone.Buy)

		orderbookClone.Buy = make([]*coreum.OrderBookOrder, len(orderbook.Buy))
		copy(orderbookClone.Buy, orderbook.Buy)
		orderbookClone.Sell = make([]*coreum.OrderBookOrder, len(orderbook.Sell))
		copy(orderbookClone.Sell, orderbook.Sell)
		// Orders are aggregated by price so that only one record exists for a given price (can reduce the number of records to be displayed)
		// This is done by summing up the quantities of orders with the same price
		orderbookClone.Sell = aggregateOrders(orderbookClone.Sell)
		orderbookClone.Buy = aggregateOrders(orderbookClone.Buy)
		return orderbookClone, nil
	}
	return orderbook, nil
}

func (a *Application) fetchOrderBookFromChain(orderbook *coreum.OrderBookOrders, network metadata.Network, denom1, denom2 string, limit int) (*coreum.OrderBookOrders, error) {
	cacheExpireMutex.Lock()
	v, ok := cacheExpire[orderbookCacheKey(network, denom1, denom2)]
	if !ok {
		cacheExpire[orderbookCacheKey(network, denom1, denom2)] = time.Now()
	}
	cacheExpireMutex.Unlock()
	if orderbook == nil || (len(orderbook.Buy) == 0 && len(orderbook.Sell) == 0) || !v.Add(5*time.Minute).After(time.Now()) {
		ctx := context.Background()

		var err error
		orderbook, err = a.TxEncoder[network].reader.QueryOrderBookRelevantOrders(ctx, denom1, denom2, uint64(limit))
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return nil, fmt.Errorf("there is no orderbook for %s - %s", denom1, denom2)
			}
			return nil, err
		}
		// The orders from the on chain orderbook need to be normalized to the orderbook format
		denom1Currency, err := a.currencyClient.GetCurrency(ctx, network, denom1)
		if err != nil {
			return nil, err
		}
		denom2Currency, err := a.currencyClient.GetCurrency(ctx, network, denom2)
		if err != nil {
			return nil, err
		}

		for _, order := range orderbook.Buy {
			o, err := a.Normalize(ctx, &OrderBookOrder{OrderBookOrder: order,
				BaseDenom: denom1Currency.Denom, QuoteDenom: denom2Currency.Denom, Network: network, Side: orderproperties.Side_SIDE_BUY})
			if err != nil {
				logger.Errorf("Error normalizing order %d: %v", order.Sequence, err)
				continue
			}
			order.Amount = o.Amount
			order.Price = o.Price
		}
		for _, order := range orderbook.Sell {
			o, err := a.Normalize(ctx, &OrderBookOrder{OrderBookOrder: order,
				BaseDenom: denom1Currency.Denom, QuoteDenom: denom2Currency.Denom, Network: network, Side: orderproperties.Side_SIDE_SELL})
			if err != nil {
				logger.Errorf("Error normalizing order %d: %v", order.Sequence, err)
				continue
			}
			order.Amount = o.Amount
			order.Price = o.Price
		}
		cacheExpireMutex.Lock()
		cacheExpire[orderbookCacheKey(network, denom1, denom2)] = time.Now()
		cacheExpireMutex.Unlock()
	}
	return orderbook, nil
}

// Process orders from the database as updates on the actual orderbook loaded from the chain
// (Orders are loaded from the chain to be able to start displaying orders while the database
// is not yet complete (e.g. the app has just started so no data in the database yet))
// Fetch additional orders from the database
// Query the database for the order status of the orders in the orderbook
func (a *Application) fetchOrderbookFromDatabase(orderbook *coreum.OrderBookOrders, network metadata.Network, denom1, denom2 string, processStart time.Time) (*coreum.OrderBookOrders, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	denom1Currency, err := a.currencyClient.GetCurrency(ctx, network, denom1)
	if err != nil {
		return nil, err
	}
	denom2Currency, err := a.currencyClient.GetCurrency(ctx, network, denom2)
	if err != nil {
		return nil, err
	}

	// Get (new) open orders from the database for the given network and denoms
	orders, err := a.orderClient.GetAll(ctx, &ordergrpc.Filter{
		Network: network,
		Denom1:  denom1Currency.Denom,
		Denom2:  denom2Currency.Denom,
		// The process writing the data is writing for a "previous" block, and we use block time to determine the time to read from
		// Also there can be a processing delay on the websocket (basic timing of the websocket updates), which can cause a further delay)
		// Lastly this might be called after a blockchain timeout has occurred (on retrieval of the orderbook): Compensate for that too
		// So we need to get some more data to compensate for that: Query up to 10 minute in the past (was 1 minute but somehow that was too short)
		From:        timestamppb.New(processStart.Add(-10 * time.Minute)),
		To:          timestamppb.Now(),
		OrderStatus: lo.ToPtr(ordergrpc.OrderStatus_ORDER_STATUS_OPEN),
	})
	if err != nil {
		return nil, err
	}
	wg := &sync.WaitGroup{}
	buyMap := a.verifyOrderBookSide(ctx, orderbook.Buy, network, wg)
	sellMap := a.verifyOrderBookSide(ctx, orderbook.Sell, network, wg)
	wg.Wait()
	// Merge the sets with the orders from the database
	for _, order := range orders.Orders {
		// Normalize the order to the orderbook format
		o, err := a.Normalize(ctx, order)
		if err != nil {
			logger.Errorf("Error normalizing order %d (%s): %v", order.Sequence, network.String(), err)
			continue
		}
		switch order.Side {
		case orderproperties.Side_SIDE_BUY:
			if _, exists := buyMap[o.Sequence]; !exists {
				buyMap[o.Sequence] = o
			}
		case orderproperties.Side_SIDE_SELL:
			if _, exists := sellMap[o.Sequence]; !exists {
				sellMap[o.Sequence] = o
			}
		}
	}
	// Convert the maps back to slices
	buySide := make([]*coreum.OrderBookOrder, 0, len(buyMap))
	sellSide := make([]*coreum.OrderBookOrder, 0, len(sellMap))
	for _, o := range buyMap {
		buySide = append(buySide, o)
	}
	for _, o := range sellMap {
		sellSide = append(sellSide, o)
	}
	orderbook.Buy = buySide
	orderbook.Sell = sellSide
	return orderbook, nil
}

// Get the orders in the orderbook from the database for verification if they are still valid
func (a *Application) verifyOrderBookSide(ctx context.Context, orderbookSide []*coreum.OrderBookOrder, network metadata.Network, wg *sync.WaitGroup) map[uint64]*coreum.OrderBookOrder {
	retVal := make([]*coreum.OrderBookOrder, 0)
	for _, o := range orderbookSide {
		wg.Add(1)
		go func(o *coreum.OrderBookOrder) {
			order, err := a.orderClient.Get(ctx, &ordergrpc.ID{Sequence: int64(o.Sequence), Network: network})
			if err != nil {
				logger.Errorf("Error getting order %d (%s) from the database: %v", o.Sequence, network.String(), err)
				wg.Done()
				return
			}
			if order.OrderStatus == ordergrpc.OrderStatus_ORDER_STATUS_CANCELED ||
				order.OrderStatus == ordergrpc.OrderStatus_ORDER_STATUS_FILLED ||
				order.OrderStatus == ordergrpc.OrderStatus_ORDER_STATUS_EXPIRED ||
				(*order.RemainingQuantity).IsZero() {
				wg.Done()
				return
			}
			orderBookMutex.Lock()
			retVal = append(retVal, o)
			orderBookMutex.Unlock()
			wg.Done()
		}(o)
	}
	retMap := make(map[uint64]*coreum.OrderBookOrder)
	for _, o := range retVal {
		retMap[o.Sequence] = o
	}
	return retMap
}

func aggregateOrders(orders []*coreum.OrderBookOrder) []*coreum.OrderBookOrder {
	aggregatedOrders := make([]*coreum.OrderBookOrder, 0)
	if len(orders) == 0 {
		return aggregatedOrders
	}
	aggregatedOrders = append(aggregatedOrders, orders[0])
	for i := 1; i < len(orders); i++ {
		if orders[i].PriceDec.Equal(orders[i-1].PriceDec) {
			s, err := dec.NewFromString(orders[i].Amount)
			if err != nil {
				continue
			}
			r, err := dec.NewFromString(aggregatedOrders[len(aggregatedOrders)-1].Amount)
			if err != nil {
				continue
			}
			aggregatedOrders[len(aggregatedOrders)-1].Amount = s.Add(r).String()
			s, err = dec.NewFromString(orders[i].SymbolAmount)
			if err != nil {
				continue
			}
			r, err = dec.NewFromString(aggregatedOrders[len(aggregatedOrders)-1].SymbolAmount)
			if err != nil {
				continue
			}
			aggregatedOrders[len(aggregatedOrders)-1].SymbolAmount = s.Add(r).String()
		} else {
			aggregatedOrders = append(aggregatedOrders, orders[i])
		}
	}
	return aggregatedOrders
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
		denomCurrency, err := a.currencyClient.GetCurrency(context.Background(), network, coin.Denom)
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
			SymbolAmount: dec.NewFromBigInt(coin.Amount.BigInt(), 0).Div(dec.New(1, precision)).String(),
		})
	}
	return walletAssets, nil
}

// Normalize order to have the precision of the currencies applied
// Add SymbolAmount, RemainingSymbolAmount, HumanReadablePrice
func (app *Application) Normalize(ctx context.Context, inputOrder interface{}) (*coreum.OrderBookOrder, error) {
	switch order := inputOrder.(type) {
	case *ordergrpc.Order:
		baseDenomPrecision, quoteDenomPrecision, err := app.currencyClient.Precisions(ctx, order.MetaData.Network, order.BaseDenom, order.QuoteDenom)
		if err != nil {
			return nil, err
		}

		price := dec.NewFromFloat(order.Price)
		quoteAmountSubunit := dec.New(order.Quantity.Value, order.Quantity.Exp)
		remainingQuantity := dec.New(order.RemainingQuantity.Value, order.RemainingQuantity.Exp)

		return &coreum.OrderBookOrder{
			Price:                 fmt.Sprintf("%f", price.InexactFloat64()),
			HumanReadablePrice:    dmn.ToSymbolPrice(baseDenomPrecision, quoteDenomPrecision, price.InexactFloat64(), &quoteAmountSubunit, orderproperties.Side_SIDE_BUY).String(),
			Amount:                quoteAmountSubunit.String(),
			SymbolAmount:          dmn.ToSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &quoteAmountSubunit, order.Side).String(),
			Sequence:              uint64(order.Sequence),
			Account:               order.Account,
			OrderID:               order.OrderID,
			RemainingAmount:       remainingQuantity.String(),
			RemainingSymbolAmount: dmn.ToSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &remainingQuantity, order.Side).String(),
		}, nil
	case *OrderBookOrder:
		baseDenomPrecision, quoteDenomPrecision, err := app.currencyClient.Precisions(ctx, order.Network, order.BaseDenom, order.QuoteDenom)
		if err != nil {
			return nil, err
		}
		price, err := dec.NewFromString(order.OrderBookOrder.Price)
		if err != nil {
			return nil, err
		}
		quoteAmountSubunit, err := dec.NewFromString(order.OrderBookOrder.Amount)
		if err != nil {
			return nil, err
		}
		remainingQuantity, err := dec.NewFromString(order.OrderBookOrder.RemainingAmount)
		if err != nil {
			return nil, err
		}
		order.OrderBookOrder.HumanReadablePrice = dmn.ToSymbolPrice(baseDenomPrecision, quoteDenomPrecision, price.InexactFloat64(), &quoteAmountSubunit, orderproperties.Side_SIDE_BUY).String()
		order.OrderBookOrder.SymbolAmount = dmn.ToSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &quoteAmountSubunit, order.Side).String()
		order.OrderBookOrder.RemainingSymbolAmount = dmn.ToSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &remainingQuantity, order.Side).String()
		return order.OrderBookOrder, nil
	}
	return nil, fmt.Errorf("unknown order type")
}
