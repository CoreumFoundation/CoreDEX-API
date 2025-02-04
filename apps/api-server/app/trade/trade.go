package trade

import (
	"context"
	"fmt"

	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	tradegrpclient "github.com/CoreumFoundation/CoreDEX-API/domain/trade/client"
)

type Application struct {
	client tradegrpc.TradeServiceClient
}

type Trade struct {
	*tradegrpc.Trade
	HumanReadablePrice string
	SymbolAmount       string
}

type Trades []*Trade

func NewApplication() *Application {
	app := &Application{
		client: tradegrpclient.Client(),
	}
	return app
}

func (app *Application) GetTrades(ctx context.Context, filter *tradegrpc.Filter) (*Trades, error) {
	trades, err := app.client.GetAll(tradegrpclient.AuthCtx(ctx), filter)
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
	}
	return &trs, nil
}
