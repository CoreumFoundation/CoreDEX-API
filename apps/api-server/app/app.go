package app

import (
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/currency"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/ohlc"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/order"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/precision"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/ticker"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/trade"
	currencygrpclient "github.com/CoreumFoundation/CoreDEX-API/domain/currency/client"
)

type Application struct {
	Trade    *trade.Application
	Ticker   *ticker.Application
	OHLC     *ohlc.Application
	Order    *order.Application
	Currency *currency.Application
}

func NewApplication() *Application {
	currencyClient := currencygrpclient.Client()
	precisionClient := precision.NewApplication(currencyClient)

	return &Application{
		Trade:    trade.NewApplication(precisionClient),
		Ticker:   ticker.NewApplication(),
		OHLC:     ohlc.NewApplication(precisionClient),
		Order:    order.NewApplication(precisionClient),
		Currency: currency.NewApplication(),
	}
}

func (app *Application) Health() error {
	return nil
}
