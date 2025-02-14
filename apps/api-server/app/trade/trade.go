package trade

import (
	"context"
	"fmt"
	"sort"

	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	ordergrpcclient "github.com/CoreumFoundation/CoreDEX-API/domain/order/client"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	tradegrpclient "github.com/CoreumFoundation/CoreDEX-API/domain/trade/client"
	"github.com/samber/lo"
)

type Application struct {
	tradeClient tradegrpc.TradeServiceClient
	orderClient ordergrpc.OrderServiceClient
}

type Trade struct {
	*tradegrpc.Trade
	HumanReadablePrice string
	SymbolAmount       string
	Status             ordergrpc.OrderStatus
}

type Trades []*Trade

func NewApplication() *Application {
	app := &Application{
		tradeClient: tradegrpclient.Client(),
		orderClient: ordergrpcclient.Client(),
	}
	return app
}

// The trades list consists of trades that are filled and cancelled.
// Cancelled trades are in the order data, while executed trades are in the trade data. Both need to be combined.
// Infinite scroll is done by using the from/to avoiding offset and data join issues
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
		tr := &Trade{}
		tr.Trade = &tradegrpc.Trade{
			Price:  order.Price,
			Amount: order.Amount,
		}
		tr.HumanReadablePrice = fmt.Sprintf("%f", order.Price)
		tr.SymbolAmount = fmt.Sprintf("%f", order.Amount.Float64())
		tr.Status = order.Status
		*trades = append(*trades, tr)
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
	trs := Trades(make([]*Trade, 0))
	// cast trs into Trades type:
	for _, trade := range trades.Trades {
		tr := &Trade{}
		tr.Trade = trade
		tr.HumanReadablePrice = fmt.Sprintf("%f", trade.Price)
		tr.SymbolAmount = fmt.Sprintf("%f", trade.Amount.Float64())
		trs = append(trs, tr)
		tr.Status = ordergrpc.OrderStatus_ORDER_STATUS_FILLED
	}
	return &trs, nil
}

// GetCancelledOrders returns all orders that are cancelled. It transforms the trade filter into an order filter for correct results.
// The filter is only active if we have an Account in the filter
func (app *Application) GetCancelledOrders(ctx context.Context, filter *tradegrpc.Filter) ([]*Trade, error) {
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
	trades := make([]*Trade, 0)
	for _, order := range orders.Orders {
		tr := &Trade{}
		tr.Trade = &tradegrpc.Trade{
			Price:  order.Price,
			Amount: order.Quantity,
		}
		tr.HumanReadablePrice = fmt.Sprintf("%f", order.Price)
		tr.SymbolAmount = fmt.Sprintf("%f", order.Quantity.Float64())
		tr.Status = order.OrderStatus
		trades = append(trades, tr)
	}
	return trades, nil
}
