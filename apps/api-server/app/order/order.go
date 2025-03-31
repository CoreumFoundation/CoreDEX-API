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
	defer a.orderbookCache.mutex.Unlock()
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

		var err error
		orderbook, err = a.TxEncoder[network].reader.QueryOrderBookRelevantOrders(ctx, denom1, denom2, uint64(limit), aggregate)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				a.orderbookCache.mutex.Unlock()
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
	}
	// If the process has taken more than a second, update the orderbook with the latest orders from the database
	// Or if the data was retrieved from the cache (and more than 1 second has passed since the last update), update the orderbook
	if time.Since(processStart) > time.Second {
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

		orders, err := a.orderClient.GetAll(ctx, &ordergrpc.Filter{
			Network: network,
			Denom1:  denom1Currency.Denom,
			Denom2:  denom2Currency.Denom,
			// This causes a slight overlap in data retrieved with the previous read, which is on purpose:
			// The process writing the data is writing for a "previous" block, and we use block time to determine the time to read from
			From: timestamppb.New(processStart.Add(-30 * time.Second)),
			To:   timestamppb.Now(), // This causes a slight overlap in data retrieved with the next read, which is on purpose
		})
		if err != nil {
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
		// Process orders from the database as updates on the actual orderbook loaded from the chain
		for _, order := range orders.Orders {
			buySideRemove, buySideAppend = a.processOrderForOrderBook(ctx, buySide, order, buySideRemove, buySideAppend, orderproperties.Side_SIDE_BUY)
			sellSideRemove, sellSideAppend = a.processOrderForOrderBook(ctx, sellSide, order, sellSideRemove, sellSideAppend, orderproperties.Side_SIDE_SELL)
		}
		for _, removeID := range buySideRemove {
			for i, buyOrder := range buySide {
				if buyOrder.Sequence == removeID {
					buySide = append(buySide[:i], buySide[i+1:]...)
				}
			}
		}
		buySide = append(buySide, buySideAppend...)
		for _, removeID := range sellSideRemove {
			for i, o := range sellSide {
				if o.Sequence == removeID {
					sellSide = append(sellSide[:i], sellSide[i+1:]...)
				}
			}
		}
		sellSide = append(sellSide, sellSideAppend...)

		orderbook.Buy = buySide
		orderbook.Sell = sellSide
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
	a.orderbookCache.data[key] = &dmncache.LockableCache{
		LastUpdated: tStartUpdate,
		Value:       orderbook,
	}
	return orderbook, nil
}

// Function generates 2 lists: A remove list and an append list
// Apply the remove list first, then apply the append list on the orderbook to get to a correct orderbook
func (a *Application) processOrderForOrderBook(ctx context.Context, orderbook []*coreum.OrderBookOrder,
	order *ordergrpc.Order,
	removeList []uint64,
	appendList []*coreum.OrderBookOrder,
	side orderproperties.Side,
) ([]uint64, []*coreum.OrderBookOrder) {
	if order.Side != side {
		return removeList, appendList
	}
	// Remove all orders which we already have in the orderbook (if it was required, we add it back later)
	for _, o := range orderbook {
		logger.Infof("Orderbook order: sequences %d %d", o.Sequence, uint64(order.Sequence))
		if o.Sequence == uint64(order.Sequence) {
			removeList = append(removeList, o.Sequence)
			// Once the order is found, we can exit the loop
			logger.Infof("Removelist: Orderbook order: sequences %d %d", o.Sequence, uint64(order.Sequence))
			break
		}
	}
	// Remove the order before adding the order again (Order might or might not be present)
	removeList = append(removeList, uint64(order.Sequence))
	// Check if order is already in append list:
	for _, o := range appendList {
		if o.Sequence == uint64(order.Sequence) {
			return removeList, appendList
		}
	}
	// Do not add orders which are not valid anymore
	if (*order.RemainingQuantity).IsZero() ||
		order.OrderStatus == ordergrpc.OrderStatus_ORDER_STATUS_CANCELED ||
		order.OrderStatus == ordergrpc.OrderStatus_ORDER_STATUS_FILLED {
		logger.Infof("Do not keep: Orderbook order: sequences %d", uint64(order.Sequence))
		return removeList, appendList
	}
	o, err := a.Normalize(ctx, order)
	if err != nil {
		logger.Errorf("Error normalizing order %d: %v", order.Sequence, err)
		return removeList, appendList
	}
	appendList = append(appendList, o)
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
