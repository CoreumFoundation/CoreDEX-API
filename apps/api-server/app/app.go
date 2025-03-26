package app

import (
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/currency"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/ohlc"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/order"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/ticker"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/trade"
	currencyclient "github.com/CoreumFoundation/CoreDEX-API/domain/currency/client"
)

type Application struct {
	Trade    *trade.Application
	Ticker   *ticker.Application
	OHLC     *ohlc.Application
	Order    *order.Application
	Currency *currency.Application
}

func NewApplication() *Application {
	currencyClient := currencyclient.Client()
	currencyApp := currency.NewApplication(currencyClient)
	ohlcApp := ohlc.NewApplication(currencyApp)

	return &Application{
		Trade:    trade.NewApplication(currencyApp),
		Ticker:   ticker.NewApplication(ohlcApp),
		OHLC:     ohlcApp,
		Order:    order.NewApplication(currencyApp),
		Currency: currency.NewApplication(currencyClient),
	}
}

func (app *Application) Health() error {
	return nil
}
