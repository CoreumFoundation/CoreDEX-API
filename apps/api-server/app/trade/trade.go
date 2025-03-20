package trade

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/samber/lo"

	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/precision"
	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	decimal "github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	ordergrpcclient "github.com/CoreumFoundation/CoreDEX-API/domain/order/client"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	tradegrpclient "github.com/CoreumFoundation/CoreDEX-API/domain/trade/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type Application struct {
	tradeClient     tradegrpc.TradeServiceClient
	orderClient     ordergrpc.OrderServiceClient
	precisionClient precision.Application
}

type Trades []*dmn.Trade

func NewApplication(precisionClient *precision.Application) *Application {
	app := &Application{
		tradeClient:     tradegrpclient.Client(),
		orderClient:     ordergrpcclient.Client(),
		precisionClient: *precisionClient,
	}
	return app
}

// The trades list consists of trades that are filled and cancelled.
// Cancelled trades are in the order data, while executed trades are in the trade data. Both need to be combined.
// Infinite scroll is done by using the from/to avoiding offset and data join issues
//
// Data is calculated where required: The Trades can be provided inverted compared to the requested values
// The values which are not in the order of the requested values are recalculated.
func (app *Application) GetTrades(ctx context.Context, filter *tradegrpc.Filter) (*Trades, error) {
	trades, err := app.getTrades(ctx, filter)
	if err != nil {
		return nil, err
	}

	cancelledOrders, err := app.GetCancelledOrders(ctx, filter)
	if err != nil {
		return nil, err
	}
	// append the cancelled orders to the trades list
	for _, order := range cancelledOrders {
		*trades = append(*trades, order)
	}
	// Order the trades by blockheight descending
	sort.Slice(*trades, func(i, j int) bool {
		return (*trades)[i].BlockHeight > (*trades)[j].BlockHeight
	})
	return trades, nil
}

func (app *Application) getTrades(ctx context.Context, filter *tradegrpc.Filter) (*Trades, error) {
	trades, err := app.tradeClient.GetAll(tradegrpclient.AuthCtx(ctx), filter)
	if err != nil {
		return nil, err
	}
	trs := Trades(make([]*dmn.Trade, 0))
	// cast trs into Trades type:
	// Take into account that the data can be inverted (Denom1-Denom2 vs Denom2-Denom1)
	for _, trade := range trades.Trades {
		tr, err := app.precisionClient.NormalizeTrade(ctx, trade)
		if err != nil {
			logger.Errorf("Error normalizing trade %s: %v", *trade.TXID, err)
			continue
		}
		tr.Status = ordergrpc.OrderStatus_ORDER_STATUS_FILLED
		trs = append(trs, tr)
	}
	return &trs, nil
}

// GetCancelledOrders returns all orders that are cancelled. It transforms the trade filter into an order filter for correct results.
// The filter is only active if we have an Account in the filter
func (app *Application) GetCancelledOrders(ctx context.Context, filter *tradegrpc.Filter) ([]*dmn.Trade, error) {
	if filter.Account == nil || *filter.Account == "" {
		return nil, nil
	}
	orderFilter := &ordergrpc.Filter{
		Account:     filter.Account,
		OrderStatus: lo.ToPtr(ordergrpc.OrderStatus_ORDER_STATUS_CANCELED),
		From:        filter.From,
		To:          filter.To,
		Network:     filter.Network,
	}
	orders, err := app.orderClient.GetAll(ctx, orderFilter)
	if err != nil {
		return nil, err
	}
	// Map the orders into trades:
	trades := make([]*dmn.Trade, 0)
	for _, order := range orders.Orders {
		tr := &dmn.Trade{}
		tr.Trade = &tradegrpc.Trade{
			Price:  order.Price,
			Amount: order.Quantity,
		}
		if strings.Compare(tr.Trade.Denom1.Denom, filter.Denom1.Denom) != 0 {
			tr.Trade.Denom1, tr.Trade.Denom2 = tr.Trade.Denom2, tr.Trade.Denom1
			r := tr.Trade.Amount.Mul(tr.Trade.Price)
			tr.Trade.Amount = decimal.FromFloat64(r)
			tr.Trade.Price = 1 / tr.Trade.Price
		}

		tr.HumanReadablePrice = fmt.Sprintf("%f", order.Price)
		tr.SymbolAmount = fmt.Sprintf("%f", order.Quantity.Float64())
		tr.Status = order.OrderStatus
		trades = append(trades, tr)
	}
	return trades, nil
}
