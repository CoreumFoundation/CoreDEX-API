package precision

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"
	dec "github.com/shopspring/decimal"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	dmncache "github.com/CoreumFoundation/CoreDEX-API/domain/cache"
	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	currencyclient "github.com/CoreumFoundation/CoreDEX-API/domain/currency/client"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	"github.com/CoreumFoundation/CoreDEX-API/domain/symbol"
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

type OrderBookOrder struct {
	OrderBookOrder *coreum.OrderBookOrder
	BaseDenom      *denom.Denom
	QuoteDenom     *denom.Denom
	Network        metadata.Network
	Side           orderproperties.Side
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
func (app *Application) NormalizeOrder(ctx context.Context, inputOrder interface{}) (*coreum.OrderBookOrder, error) {
	switch order := inputOrder.(type) {
	case *ordergrpc.Order:
		baseDenomPrecision, quoteDenomPrecision, err := app.precisions(ctx, order.MetaData.Network, order.BaseDenom, order.QuoteDenom)
		if err != nil {
			return nil, err
		}

		price := dec.NewFromFloat(order.Price)
		quoteAmountSubunit := dec.New(order.Quantity.Value, order.Quantity.Exp)
		remainingQuantity := dec.New(order.RemainingQuantity.Value, order.RemainingQuantity.Exp)

		return &coreum.OrderBookOrder{
			Price:                 fmt.Sprintf("%f", price.InexactFloat64()),
			HumanReadablePrice:    toSymbolPrice(baseDenomPrecision, quoteDenomPrecision, price.InexactFloat64(), &quoteAmountSubunit, orderproperties.Side_SIDE_BUY).String(),
			Amount:                quoteAmountSubunit.String(),
			SymbolAmount:          toSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &quoteAmountSubunit, order.Side).String(),
			Sequence:              uint64(order.Sequence),
			Account:               order.Account,
			OrderID:               order.OrderID,
			RemainingAmount:       remainingQuantity.String(),
			RemainingSymbolAmount: toSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &remainingQuantity, order.Side).String(),
		}, nil
	case *OrderBookOrder:
		baseDenomPrecision, quoteDenomPrecision, err := app.precisions(ctx, order.Network, order.BaseDenom, order.QuoteDenom)
		if err != nil {
			return nil, err
		}
		price, err := dec.NewFromString(order.OrderBookOrder.Price)
		if err != nil {
			return nil, err
		}
		quoteAmountSubunit, err := dec.NewFromString(order.OrderBookOrder.Amount)
		if err != nil {
			return nil, err
		}
		remainingQuantity, err := dec.NewFromString(order.OrderBookOrder.RemainingAmount)
		if err != nil {
			return nil, err
		}
		order.OrderBookOrder.HumanReadablePrice = toSymbolPrice(baseDenomPrecision, quoteDenomPrecision, price.InexactFloat64(), &quoteAmountSubunit, orderproperties.Side_SIDE_BUY).String()
		order.OrderBookOrder.SymbolAmount = toSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &quoteAmountSubunit, order.Side).String()
		order.OrderBookOrder.RemainingSymbolAmount = toSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &remainingQuantity, order.Side).String()
		return order.OrderBookOrder, nil
	}
	return nil, fmt.Errorf("unknown order type")
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
	case orderproperties.Side_SIDE_BUY:
		symbolAmount = symbolAmount.Div(dec.New(1, int32(quoteDenomPrecision)))
	}
	return symbolAmount
}

// Returns human readable price and amount
func (app *Application) NormalizeTrade(ctx context.Context, trade *tradegrpc.Trade) (*dmn.Trade, error) {
	baseDenomPrecision, quoteDenomPrecision, err := app.precisions(ctx, trade.MetaData.Network, trade.Denom1, trade.Denom2)
	if err != nil {
		return nil, err
	}
	tr := &dmn.Trade{
		Trade: trade,
	}
	quoteAmountSubunit := dec.New(trade.Amount.Value, trade.Amount.Exp)
	tr.HumanReadablePrice = toSymbolPrice(baseDenomPrecision, quoteDenomPrecision, trade.Price,
		&quoteAmountSubunit, trade.Side).String()
	tr.SymbolAmount = toSymbolAmount(baseDenomPrecision, quoteDenomPrecision, &quoteAmountSubunit, trade.Side).String()
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

/*
OHLC data is stored in the subunit price and volume notation of the orders.
This function converts the subunit price and volume to human readable price and volume.
*/
func (app *Application) NormalizeOHLC(ctx context.Context, ohlc *ohlcgrpc.OHLC) (*ohlcgrpc.OHLC, error) {
	// ohlc symbol to denoms base and quote:
	sym, err := symbol.NewSymbol(ohlc.Symbol)
	if err != nil {
		return nil, err
	}
	baseDenomPrecision, quoteDenomPrecision, err := app.precisions(ctx, ohlc.MetaData.Network, sym.Denom1, sym.Denom2)
	if err != nil {
		return nil, err
	}
	// Price is in subunit notation (subunitBase/subunitQuote)
	// We need the prices in unit notation: (base/quote) => price * 10^basePrecision/10^quotePrecision
	mult := dec.New(1, baseDenomPrecision).Div(dec.New(1, quoteDenomPrecision)).InexactFloat64()
	ohlc.Close = ohlc.Close * mult
	ohlc.Open = ohlc.Open * mult
	ohlc.High = ohlc.High * mult
	ohlc.Low = ohlc.Low * mult
	// Volume is in subunit notation
	// We need the volume in unit notation: volume * 10^-baseDenomPrecision
	ohlc.Volume = ohlc.Volume * dec.New(1, -baseDenomPrecision).InexactFloat64()
	// Inverted volume is in subunit notation
	// We need the quote volume in unit notation: volume * 10^-quoteDenomPrecision
	ohlc.QuoteVolume = ohlc.QuoteVolume * dec.New(1, -quoteDenomPrecision).InexactFloat64()
	return ohlc, nil
}
