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
	return &Application{txEncoders, orderClient, *currencyClient}
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

func (a *Application) OrderBookRelevantOrders(network metadata.Network, denom1, denom2 string, limit int, aggregate bool) (*coreum.OrderBookOrders, error) {
	var err error
	orderbook, err := a.fetchOrderBookFromChain(network, denom1, denom2, limit)
	if err != nil {
		return nil, err
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

func (a *Application) fetchOrderBookFromChain(network metadata.Network, denom1, denom2 string, limit int) (*coreum.OrderBookOrders, error) {
	ctx, timeout := context.WithTimeout(context.Background(), 60*time.Second)
	defer timeout()

	var err error
	orderbook, err := a.TxEncoder[network].reader.QueryOrderBookRelevantOrders(ctx, denom1, denom2, uint64(limit))
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
	return orderbook, nil
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
			return walletAssets, err
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
			HumanReadablePrice:    dmn.ToSymbolPrice(baseDenomPrecision, quoteDenomPrecision, price.InexactFloat64(), &quoteAmountSubunit, order.Side).String(),
			Amount:                quoteAmountSubunit.String(),
			SymbolAmount:          dmn.ToSymbolOrderAmount(baseDenomPrecision, quoteDenomPrecision, &quoteAmountSubunit, order.Side).String(),
			Sequence:              uint64(order.Sequence),
			Account:               order.Account,
			OrderID:               order.OrderID,
			RemainingAmount:       remainingQuantity.String(),
			RemainingSymbolAmount: dmn.ToSymbolOrderAmount(baseDenomPrecision, quoteDenomPrecision, &remainingQuantity, order.Side).String(),
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
		order.OrderBookOrder.SymbolAmount = dmn.ToSymbolOrderAmount(baseDenomPrecision, quoteDenomPrecision, &quoteAmountSubunit, order.Side).String()
		order.OrderBookOrder.RemainingSymbolAmount = dmn.ToSymbolOrderAmount(baseDenomPrecision, quoteDenomPrecision, &remainingQuantity, order.Side).String()
		return order.OrderBookOrder, nil
	}
	return nil, fmt.Errorf("unknown order type")
}
