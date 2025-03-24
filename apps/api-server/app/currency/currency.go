package currency

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"

	dmncache "github.com/CoreumFoundation/CoreDEX-API/domain/cache"
	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	currencygrpclient "github.com/CoreumFoundation/CoreDEX-API/domain/currency/client"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
)

type cache struct {
	mutex *sync.RWMutex
	data  map[string]*dmncache.LockableCache
}

type Application struct {
	client        currencygrpc.CurrencyServiceClient
	currencyCache *cache
}

// currency client is passed in to support Mock for testing
func NewApplication(currencyClient currencygrpc.CurrencyServiceClient) *Application {
	app := &Application{
		client: currencyClient,
	}
	app.currencyCache = &cache{
		mutex: &sync.RWMutex{},
		data:  make(map[string]*dmncache.LockableCache),
	}
	go dmncache.CleanCache(app.currencyCache.data, app.currencyCache.mutex, 15*time.Minute)
	return app
}

func (app *Application) GetCurrencies(ctx context.Context, network metadata.Network) (*currencygrpc.Currencies, error) {
	return app.client.GetAll(currencygrpclient.AuthCtx(ctx), &currencygrpc.Filter{Network: network})
}

func (app *Application) getPrecision(ctx context.Context, den *denom.Denom, network metadata.Network) (int32, error) {
	app.currencyCache.mutex.RLock()
	cur, ok := app.currencyCache.data[cacheKey(den.String(), network)]
	app.currencyCache.mutex.RUnlock()
	if ok {
		precision := *cur.Value.(currencygrpc.Currency).Denom.Precision
		return precision, nil
	}
	denomCurrency, err := app.client.Get(currencygrpclient.AuthCtx(ctx), &currencygrpc.ID{
		Network: network,
		Denom:   den.Denom,
	})
	if err != nil {
		return 0, err
	}
	if denomCurrency.Denom == nil || denomCurrency.Denom.Precision == nil {
		return 0, fmt.Errorf("precision not found for %s", den.String())
	}
	app.currencyCache.mutex.Lock()
	app.currencyCache.data[cacheKey(den.String(), network)] = &dmncache.LockableCache{
		Value:       *denomCurrency,
		LastUpdated: time.Now(),
	}
	app.currencyCache.mutex.Unlock()
	return *denomCurrency.Denom.Precision, nil
}

func (app *Application) GetCurrency(ctx context.Context, network metadata.Network, denom string) (*currencygrpc.Currency, error) {
	app.currencyCache.mutex.RLock()
	cur, ok := app.currencyCache.data[cacheKey(denom, network)]
	app.currencyCache.mutex.RUnlock()
	if ok {
		return lo.ToPtr(cur.Value.(currencygrpc.Currency)), nil
	}
	denomCurrency, err := app.client.Get(currencygrpclient.AuthCtx(ctx), &currencygrpc.ID{
		Network: network,
		Denom:   denom,
	})
	if err != nil {
		return nil, err
	}
	app.currencyCache.mutex.Lock()
	app.currencyCache.data[cacheKey(denom, network)] = &dmncache.LockableCache{
		Value:       *denomCurrency,
		LastUpdated: time.Now(),
	}
	app.currencyCache.mutex.Unlock()
	return denomCurrency, nil
}

func cacheKey(denom string, network metadata.Network) string {
	return fmt.Sprintf("%s-%d", denom, network)
}

// Get the currencies from the currency service to be able to present the correct precision to the user
func (app *Application) Precisions(ctx context.Context, network metadata.Network, denom1, denom2 *denom.Denom) (int32, int32, error) {
	denom1Precision, err := app.getPrecision(ctx, denom1, network)
	if err != nil {
		return 0, 0, err
	}
	denom2Precision, err := app.getPrecision(ctx, denom2, network)
	if err != nil {
		return 0, 0, err
	}
	return denom1Precision, denom2Precision, nil
}
