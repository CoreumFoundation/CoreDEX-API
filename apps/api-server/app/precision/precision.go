package precision

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	dec "github.com/shopspring/decimal"

	dmncache "github.com/CoreumFoundation/CoreDEX-API/domain/cache"
	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	currencyclient "github.com/CoreumFoundation/CoreDEX-API/domain/currency/client"
	"github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
)

type cache struct {
	mutex *sync.RWMutex
	data  map[string]*dmncache.LockableCache
}

type Application struct {
	currencyClient currencygrpc.CurrencyServiceClient
	currencyCache  *cache
}

func NewApplication(currencyClient currencygrpc.CurrencyServiceClient) *Application {
	app := &Application{
		currencyClient: currencyClient,
	}
	app.currencyCache = &cache{
		mutex: &sync.RWMutex{},
		data:  make(map[string]*dmncache.LockableCache),
	}
	go dmncache.CleanCache(app.currencyCache.data, app.currencyCache.mutex, 15*time.Minute)
	return app
}

func (app *Application) NormalizeOrder(ctx context.Context, order *ordergrpc.Order) (*ordergrpc.Order, error) {
	denom1Precision, denom2Precision, err := app.precisions(ctx, order.MetaData.Network, order.BaseDenom, order.QuoteDenom)
	if err != nil {
		return nil, err
	}

	exp := int32(denom1Precision - denom2Precision)
	if order.Side == orderproperties.Side_SIDE_BUY {
		exp = int32(denom2Precision - denom1Precision)
	}
	price := dec.NewFromFloat(order.Price).Mul(dec.New(1, exp))
	quantity := dec.New(order.Quantity.Value, order.Quantity.Exp-int32(denom1Precision))
	remainingExp := order.RemainingQuantity.Exp - int32(denom1Precision)
	if order.RemainingQuantity.Value == 0 {
		remainingExp = 0
	}
	remainingQuantity := dec.New(order.RemainingQuantity.Value, remainingExp)

	order.Quantity = decimal.FromDec(quantity)
	order.RemainingQuantity = decimal.FromDec(remainingQuantity)
	order.Price, _ = price.Float64()
	return order, nil
}

// Returns human readable price and amount
func (app *Application) NormalizeTrade(ctx context.Context, trade *tradegrpc.Trade) (*tradegrpc.Trade, error) {
	denom1Precision, denom2Precision, err := app.precisions(ctx, trade.MetaData.Network, trade.Denom1, trade.Denom2)
	if err != nil {
		return nil, err
	}
	amount := dec.New(trade.Amount.Value, trade.Amount.Exp)
	// TODO Price has to be corrected by the divider of baseDenomPrecision/quoteDenomPrecision (or the other way around?)
	price, _ := dec.NewFromFloat(trade.Price).Div(amount).Float64()
	price = price * (math.Pow(10, float64(denom1Precision)) / math.Pow(10, float64(denom2Precision)))
	af := amount.InexactFloat64() / math.Pow(10, float64(denom1Precision))
	trade.Price = price
	trade.Amount = decimal.FromFloat64(af)
	return trade, nil
}

func cacheKey(denom string, network metadata.Network) string {
	return fmt.Sprintf("%s-%d", denom, network)
}

// Get the currencies from the currency service to be able to present the correct precision to the user
func (app *Application) precisions(ctx context.Context, network metadata.Network, denom1, denom2 *denom.Denom) (int32, int32, error) {
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

func (app *Application) getPrecision(ctx context.Context, den *denom.Denom, network metadata.Network) (int32, error) {
	app.currencyCache.mutex.RLock()
	cur, ok := app.currencyCache.data[cacheKey(den.String())]
	app.currencyCache.mutex.RUnlock()
	if ok {
		precision := *cur.Value.(currencygrpc.Currency).Denom.Precision
		return precision, nil
	}
	denomCurrency, err := app.currencyClient.Get(currencyclient.AuthCtx(ctx), &currencygrpc.ID{
		Network: network,
		Denom:   den.Denom,
	})
	if err == nil {
		if denomCurrency.Denom != nil && denomCurrency.Denom.Precision != nil {
			app.currencyCache.mutex.Lock()
			app.currencyCache.data[cacheKey(den, network)] = &dmncache.LockableCache{
				Value:       *denomCurrency,
				LastUpdated: time.Now(),
			}
			return *denomCurrency.Denom.Precision, nil
		}
	}
	return 0, err
}

func (app *Application) GetCurrency(network metadata.Network, denom string) (*currencygrpc.Currency, error) {
	app.currencyCache.mutex.RLock()
	cur, ok := app.currencyCache.data[cacheKey(denom, network)]
	app.currencyCache.mutex.RUnlock()

	return nil
}
