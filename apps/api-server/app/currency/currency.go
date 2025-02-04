package currency

import (
	"context"

	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	currencygrpclient "github.com/CoreumFoundation/CoreDEX-API/domain/currency/client"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
)

type Application struct {
	client currencygrpc.CurrencyServiceClient
}

func NewApplication() *Application {
	return &Application{
		client: currencygrpclient.Client(),
	}
}

func (app *Application) GetCurrencies(ctx context.Context, network metadata.Network) (*currencygrpc.Currencies, error) {
	return app.client.GetAll(currencygrpclient.AuthCtx(ctx), &currencygrpc.Filter{Network: network})
}

func (app *Application) GetCurrency(ctx context.Context, network metadata.Network, denom string) (*currencygrpc.Currency, error) {
	return app.client.Get(currencygrpclient.AuthCtx(ctx), &currencygrpc.ID{
		Network: network,
		Denom:   denom,
	})
}
