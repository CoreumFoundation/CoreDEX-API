package app

import (
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/currency"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/ohlc"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/order"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/ticker"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/trade"
)

type Application struct {
	Trade    *trade.Application
	Ticker   *ticker.Application
	OHLC     *ohlc.Application
	Order    *order.Application
	Currency *currency.Application
}

func NewApplication() *Application {
	return &Application{
		Trade:    trade.NewApplication(),
		Ticker:   ticker.NewApplication(),
		OHLC:     ohlc.NewApplication(),
		Order:    order.NewApplication(),
		Currency: currency.NewApplication(),
	}
}

func (app *Application) Health() error {
	return nil
}
