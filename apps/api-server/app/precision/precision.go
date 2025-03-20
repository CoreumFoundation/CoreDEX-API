package precision

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/samber/lo"
	dec "github.com/shopspring/decimal"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
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

// Normalize order to have the precision of the currencies applied
// Add SymbolAmount, RemainingSymbolAmount, HumanReadablePrice
func (app *Application) NormalizeOrder(ctx context.Context, order *ordergrpc.Order) (*coreum.OrderBookOrder, error) {
	baseDenomPrecision, quoteDenomPrecision, err := app.precisions(ctx, order.MetaData.Network, order.BaseDenom, order.QuoteDenom)
	if err != nil {
		return nil, err
	}

	price := dec.NewFromFloat(order.Price)
	quoteAmountSubunit := dec.New(order.Quantity.Value, order.Quantity.Exp)

	remainingQuantity := dec.New(order.RemainingQuantity.Value, order.RemainingQuantity.Exp)
	symbolAmount := toSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &quoteAmountSubunit, order.Side)

	remainingSymbolAmount:=toSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &remainingQuantity, order.Side)

	return &coreum.OrderBookOrder{
		Price:                 fmt.Sprintf("%f", price.InexactFloat64()),
		HumanReadablePrice:    toSymbolPrice(baseDenomPrecision, quoteDenomPrecision, price.InexactFloat64(), &quoteAmountSubunit, order.Side).String(),
		Amount:                quoteAmountSubunit.String(),
		SymbolAmount:          symbolAmount.String(),
		Sequence:              uint64(order.Sequence),
		Account:               order.Account,
		OrderID:               order.OrderID,
		RemainingAmount:       remainingQuantity.String(),
		RemainingSymbolAmount: remainingSymbolAmount.String(),
	}, nil
}

func toSymbolPrice(baseDenomPrecision, quoteDenomPrecision int32, subunitPrice float64, quantity *dec.Decimal, side orderproperties.Side) dec.Decimal {
	price := dec.NewFromFloat(subunitPrice)
	quoteAmountSubunit := quantity
	baseAmountSubunit := quoteAmountSubunit.Mul(price)
	var humanReadablePrice dec.Decimal
	switch side {
	case orderproperties.Side_SIDE_SELL:
		humanReadablePrice = quoteAmountSubunit.Div(dec.New(1, quoteDenomPrecision)).
			Div(baseAmountSubunit.Div(dec.New(1, baseDenomPrecision)))
	case orderproperties.Side_SIDE_BUY:
		humanReadablePrice = baseAmountSubunit.Div(dec.New(1, baseDenomPrecision)).
			Div(quoteAmountSubunit.Div(dec.New(1, quoteDenomPrecision)))
	}
	return humanReadablePrice
}

func toSymbolAmount(baseDenomPrecision, quoteDenomPrecision int32, quantity *dec.Decimal, side orderproperties.Side) dec.Decimal {
	symbolAmount := *quantity
	switch side {
	case orderproperties.Side_SIDE_SELL:
		symbolAmount = symbolAmount.Div(dec.New(1, int32(baseDenomPrecision)))
		// remainingSymbolAmount = remainingQuantity.Div(dec.New(1, int32(baseDenomPrecision)))
	case orderproperties.Side_SIDE_BUY:
		symbolAmount = symbolAmount.Div(dec.New(1, int32(quoteDenomPrecision)))
		// remainingSymbolAmount = remainingQuantity.Div(dec.New(1, int32(quoteDenomPrecision)))
	}
	return symbolAmount
}

// Returns human readable price and amount
func (app *Application) NormalizeTrade(ctx context.Context, trade *tradegrpc.Trade) (*dmn.Trade, error) {
	tr := &dmn.Trade{}
	baseDenomPrecision, quoteDenomPrecision, err := app.precisions(ctx, trade.MetaData.Network, trade.Denom1, trade.Denom2)
	if err != nil {
		return nil, err
	}
	amount := dec.New(trade.Amount.Value, trade.Amount.Exp)
	humanReadablePrice := toSymbolPrice(baseDenomPrecision, quoteDenomPrecision, trade.Price, &amount, trade.Side)
	// TODO Price has to be corrected by the divider of baseDenomPrecision/quoteDenomPrecision (or the other way around?)
	af := amount.InexactFloat64() / math.Pow(10, float64(baseDenomPrecision))

	tr.Price = humanReadablePrice.InexactFloat64()
	tr.Amount = decimal.FromFloat64(af)
	tr.HumanReadablePrice = fmt.Sprintf("%f", trade.Price)
	tr.SymbolAmount = fmt.Sprintf("%f", trade.Amount.Float64())

	return tr, nil
}

// Returns human readable price and amount
func (app *Application) OrderToTrade(ctx context.Context, order *ordergrpc.Order) (*dmn.Trade, error) {
	tr := &dmn.Trade{}
	denom1Precision, denom2Precision, err := app.precisions(ctx, order.MetaData.Network, order.BaseDenom, order.QuoteDenom)
	if err != nil {
		return nil, err
	}
	amount := *decimal.ToSDec(order.Quantity)
	// TODO Price has to be corrected by the divider of baseDenomPrecision/quoteDenomPrecision (or the other way around?)
	price, _ := dec.NewFromFloat(order.Price).Div(amount).Float64()
	price = price * (math.Pow(10, float64(denom1Precision)) / math.Pow(10, float64(denom2Precision)))
	af := amount.InexactFloat64() / math.Pow(10, float64(denom1Precision))

	tr.Trade = &tradegrpc.Trade{
		Price:  order.Price,
		Amount: order.Quantity,
	}
	tr.Price = price
	tr.Amount = decimal.FromFloat64(af)
	tr.HumanReadablePrice = fmt.Sprintf("%f", price)
	tr.SymbolAmount = fmt.Sprintf("%f", af)

	return tr, nil
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
	cur, ok := app.currencyCache.data[cacheKey(den.String(), network)]
	app.currencyCache.mutex.RUnlock()
	if ok {
		precision := *cur.Value.(currencygrpc.Currency).Denom.Precision
		return precision, nil
	}
	denomCurrency, err := app.currencyClient.Get(currencyclient.AuthCtx(ctx), &currencygrpc.ID{
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
	denomCurrency, err := app.currencyClient.Get(currencyclient.AuthCtx(ctx), &currencygrpc.ID{
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
