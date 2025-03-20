/*
The tests in this package are to verify the applicaiton of the precision on the data which is represented by the denoms.
Case are to check if:
denom1.Precision == denom2.Precision => price, amount, humanReadable price and humanReadable SymbolAmount are correct
denom2.Precision =denom2.Precision+1 => price, amount, humanReadable price and humanReadable SymbolAmount are correct
denom2.Precision =denom2.Precision-1 => price, amount, humanReadable price and humanReadable SymbolAmount are correct
denom1.Precision = 6, denom2.Precision = 0 => price, amount, humanReadable price and humanReadable SymbolAmount are correct
denom1.Precision = 0, denom2.Precision = 6 => price, amount, humanReadable price and humanReadable SymbolAmount are correct

Inputs are:
* ordergrpc.Order
* tradegrpc.Trade
* ohlcgrpc.OHLC

Which have slightly different structures, however the same precision rules apply to all of them.
*/

package precision

import (
	"context"
	"os"
	"testing"

	"github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	"github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// the denoms used in this test:
var (
	baseDenom, quoteDenom *denom.Denom
)

func newApplicationMock(t *testing.T) *Application {
	// Set the environment variable CURRENCY_STORE
	os.Setenv("CURRENCY_STORE", "localhost:50051")
	currencyService := currency.NewMockCurrencyServiceClient()

	// The order of the input is alphabetical. In trade there can be an inversion of the denoms beased on alphabetical order.
	// The actual inversion is only stored as a flag in the trade data, actual order is not checked in retrieval.
	baseDenom, _ = denom.NewDenom("baseDenom-issuerstring")
	quoteDenom, _ = denom.NewDenom("quoteDenom")
	ctx := context.Background()
	_, err := currencyService.Upsert(ctx, &currency.Currency{
		Denom: baseDenom,
		MetaData: &metadata.MetaData{
			Network:   metadata.Network_DEVNET,
			UpdatedAt: timestamppb.Now(),
			CreatedAt: timestamppb.Now(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	p2 := int32(0)
	quoteDenom.Precision = &p2
	_, err = currencyService.Upsert(ctx, &currency.Currency{
		Denom: quoteDenom,
		MetaData: &metadata.MetaData{
			Network:   metadata.Network_DEVNET,
			UpdatedAt: timestamppb.Now(),
			CreatedAt: timestamppb.Now(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	return NewApplication(currencyService)
}

type normalizedOrderInput struct {
	Price               float64
	Quantity            float64
	RemainingQuantity   float64
	Side                orderproperties.Side
	BaseDenomPrecision  int32
	QuoteDenomPrecision int32
}

type normalizedOrderResult struct {
	Price                 string
	HumanReadablePrice    string
	Amount                string
	SymbolAmount          string
	RemainingAmount       string
	RemainingSymbolAmount string
}

type normalizedOrderTest struct {
	name   string
	Input  normalizedOrderInput
	Result normalizedOrderResult
}

func Test_NormalizeOrder(t *testing.T) {
	app := newApplicationMock(t)
	order := &ordergrpc.Order{
		BaseDenom:         baseDenom,
		QuoteDenom:        quoteDenom,
		Price:             1.0,
		Quantity:          decimal.FromFloat64(1.0),
		RemainingQuantity: decimal.FromFloat64(1.0),
		MetaData: &metadata.MetaData{
			Network: metadata.Network_DEVNET,
		},
	}
	normalizedOrderTests := []normalizedOrderTest{
		{
			name: "Base and quote denom precision are the same (BUY)",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1.0,
				RemainingQuantity:   1.0,
				Side:                orderproperties.Side_SIDE_BUY,
				BaseDenomPrecision:  0,
				QuoteDenomPrecision: 0,
			},
			Result: normalizedOrderResult{
				Price:                 "1.000000",
				HumanReadablePrice:    "1",
				Amount:                "1",
				SymbolAmount:          "1",
				RemainingAmount:       "1",
				RemainingSymbolAmount: "1",
			},
		},
		{
			name: "Base and quote denom precision are the same (SELL)",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1.0,
				RemainingQuantity:   1.0,
				Side:                orderproperties.Side_SIDE_SELL,
				BaseDenomPrecision:  0,
				QuoteDenomPrecision: 0,
			},
			Result: normalizedOrderResult{
				Price:                 "1.000000",
				HumanReadablePrice:    "1",
				Amount:                "1",
				SymbolAmount:          "1",
				RemainingAmount:       "1",
				RemainingSymbolAmount: "1",
			},
		},
		{
			name: "Base larger then quote denom precision (BUY)",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1.0,
				RemainingQuantity:   1.0,
				Side:                orderproperties.Side_SIDE_BUY,
				BaseDenomPrecision:  1,
				QuoteDenomPrecision: 0,
			},
			Result: normalizedOrderResult{
				Price:                 "1.000000",
				HumanReadablePrice:    "0.1",
				Amount:                "1",
				SymbolAmount:          "1",
				RemainingAmount:       "1",
				RemainingSymbolAmount: "1",
			},
		},
		{
			name: "Base larger then quote denom precision (SELL)",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1.0,
				RemainingQuantity:   1.0,
				Side:                orderproperties.Side_SIDE_SELL,
				BaseDenomPrecision:  1,
				QuoteDenomPrecision: 0,
			},
			Result: normalizedOrderResult{
				Price:                 "1.000000",
				HumanReadablePrice:    "10",
				Amount:                "1",
				SymbolAmount:          "0.1",
				RemainingAmount:       "1",
				RemainingSymbolAmount: "0.1",
			},
		},
		{
			name: "Base smaller then quote denom precision (BUY)",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1.0,
				RemainingQuantity:   1.0,
				Side:                orderproperties.Side_SIDE_BUY,
				BaseDenomPrecision:  0,
				QuoteDenomPrecision: 1,
			},
			Result: normalizedOrderResult{
				Price:                 "1.000000",
				HumanReadablePrice:    "10",
				Amount:                "1",
				SymbolAmount:          "0.1",
				RemainingAmount:       "1",
				RemainingSymbolAmount: "0.1",
			},
		},
		{
			name: "Base smaller then quote denom precision (SELL)",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1.0,
				RemainingQuantity:   1.0,
				Side:                orderproperties.Side_SIDE_SELL,
				BaseDenomPrecision:  0,
				QuoteDenomPrecision: 1,
			},
			Result: normalizedOrderResult{
				Price:                 "1.000000",
				HumanReadablePrice:    "0.1",
				Amount:                "1",
				SymbolAmount:          "1",
				RemainingAmount:       "1",
				RemainingSymbolAmount: "1",
			},
		},
		{
			name: "Base 2 smaller then quote denom precision (BUY)",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1.0,
				RemainingQuantity:   1.0,
				Side:                orderproperties.Side_SIDE_SELL,
				BaseDenomPrecision:  0,
				QuoteDenomPrecision: 2,
			},
			Result: normalizedOrderResult{
				Price:                 "1.000000",
				HumanReadablePrice:    "0.01",
				Amount:                "1",
				SymbolAmount:          "1",
				RemainingAmount:       "1",
				RemainingSymbolAmount: "1",
			},
		},
		{
			name: "Base 2 smaller then quote denom precision (BUY), 10x quantity",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            10.0,
				RemainingQuantity:   10.0,
				Side:                orderproperties.Side_SIDE_SELL,
				BaseDenomPrecision:  0,
				QuoteDenomPrecision: 2,
			},
			Result: normalizedOrderResult{
				Price:                 "1.000000",
				HumanReadablePrice:    "0.01",
				Amount:                "10",
				SymbolAmount:          "10",
				RemainingAmount:       "10",
				RemainingSymbolAmount: "10",
			},
		},
		{
			name: "Base and quote at typical subunit precision of 6",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1000000.0,
				RemainingQuantity:   1000000.0,
				Side:                orderproperties.Side_SIDE_SELL,
				BaseDenomPrecision:  6,
				QuoteDenomPrecision: 6,
			},
			Result: normalizedOrderResult{
				Price:                 "1.000000",
				HumanReadablePrice:    "1",
				Amount:                "1000000",
				SymbolAmount:          "1",
				RemainingAmount:       "1000000",
				RemainingSymbolAmount: "1",
			},
		},
	}
	for _, test := range normalizedOrderTests {
		order.Price = test.Input.Price
		order.Quantity = decimal.FromFloat64(test.Input.Quantity)
		order.RemainingQuantity = decimal.FromFloat64(test.Input.RemainingQuantity)
		order.Side = test.Input.Side
		baseDenom.Precision = &test.Input.BaseDenomPrecision
		quoteDenom.Precision = &test.Input.QuoteDenomPrecision
		normalizedOrder, err := app.NormalizeOrder(context.Background(), order)
		if err != nil {
			t.Fatalf("%s Error: %v", test.name, err)
		}
		if normalizedOrder == nil {
			t.Fatalf("%s Error: normalizedOrder is nil", test.name)
		}
		if normalizedOrder.Price != test.Result.Price {
			t.Errorf("%s Error: normalizedOrder.Price is %s, expected %s", test.name, normalizedOrder.Price, test.Result.Price)
		}
		if normalizedOrder.HumanReadablePrice != test.Result.HumanReadablePrice {
			t.Errorf("%s Error: normalizedOrder.HumanReadablePrice is %s, expected %s", test.name, normalizedOrder.HumanReadablePrice, test.Result.HumanReadablePrice)
		}
		if normalizedOrder.Amount != test.Result.Amount {
			t.Errorf("%s Error: normalizedOrder.Amount is %s, expected %s", test.name, normalizedOrder.Amount, test.Result.Amount)
		}
		if normalizedOrder.SymbolAmount != test.Result.SymbolAmount {
			t.Errorf("%s Error: normalizedOrder.SymbolAmount is %s, expected %s", test.name, normalizedOrder.SymbolAmount, test.Result.SymbolAmount)
		}
		if normalizedOrder.RemainingAmount != test.Result.RemainingAmount {
			t.Errorf("%s Error: normalizedOrder.RemainingAmount is %s, expected %s", test.name, normalizedOrder.RemainingAmount, test.Result.RemainingAmount)
		}
		if normalizedOrder.RemainingSymbolAmount != test.Result.RemainingSymbolAmount {
			t.Errorf("%s Error: normalizedOrder.RemainingSymbolAmount is %s, expected %s", test.name, normalizedOrder.RemainingSymbolAmount, test.Result.RemainingSymbolAmount)
		}
	}
}

type normalizedTradeResult struct {
	Price              float64
	HumanReadablePrice string
	Amount             *decimal.Decimal
	SymbolAmount       string
}

type normalizedTradeTest struct {
	name   string
	Input  normalizedOrderInput
	Result normalizedTradeResult
}

func Test_NormalizeTrade(t *testing.T) {
	app := newApplicationMock(t)
	trade := &tradegrpc.Trade{
		Denom1: baseDenom,
		Denom2: quoteDenom,
		Price:  1.0,
		Amount: decimal.FromFloat64(1.0),
		MetaData: &metadata.MetaData{
			Network: metadata.Network_DEVNET,
		},
	}
	normalizedTradeTests := []normalizedTradeTest{
		{
			name: "Base and quote denom precision are the same (BUY)",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1.0,
				RemainingQuantity:   1.0,
				Side:                orderproperties.Side_SIDE_BUY,
				BaseDenomPrecision:  0,
				QuoteDenomPrecision: 0,
			},
			Result: normalizedTradeResult{
				Price:              1.0,
				HumanReadablePrice: "1",
				Amount:             &decimal.Decimal{Value: 1, Exp: 0},
				SymbolAmount:       "1",
			},
		},
		{
			name: "Base and quote denom precision are the same (SELL)",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1.0,
				RemainingQuantity:   1.0,
				Side:                orderproperties.Side_SIDE_SELL,
				BaseDenomPrecision:  0,
				QuoteDenomPrecision: 0,
			},
			Result: normalizedTradeResult{
				Price:              1.0,
				HumanReadablePrice: "1",
				Amount:             &decimal.Decimal{Value: 1, Exp: 0},
				SymbolAmount:       "1",
			},
		},
		{
			name: "Base > quote denom precision (BUY)",
			Input: normalizedOrderInput{
				Price:               1.0,
				Quantity:            1.0,
				RemainingQuantity:   1.0,
				Side:                orderproperties.Side_SIDE_SELL,
				BaseDenomPrecision:  1,
				QuoteDenomPrecision: 0,
			},
			Result: normalizedTradeResult{
				Price:              1.0,
				HumanReadablePrice: "10",
				Amount:             &decimal.Decimal{Value: 1, Exp: 0},
				SymbolAmount:       "0.1",
			},
		},
	}
	for _, test := range normalizedTradeTests {
		trade.Price = test.Input.Price
		trade.Amount = decimal.FromFloat64(test.Input.Quantity)
		trade.Side = test.Input.Side
		baseDenom.Precision = &test.Input.BaseDenomPrecision
		quoteDenom.Precision = &test.Input.QuoteDenomPrecision
		normalizedTrade, err := app.NormalizeTrade(context.Background(), trade)
		if err != nil {
			t.Fatalf("%s Error: %v", test.name, err)
		}
		if normalizedTrade == nil {
			t.Fatalf("%s Error: normalizedTrade is nil", test.name)
		}
		if normalizedTrade.Price != test.Result.Price {
			t.Errorf("%s Error: normalizedTrade.Price is %2.f, expected %2.f", test.name, normalizedTrade.Price, test.Result.Price)
		}
		if normalizedTrade.HumanReadablePrice != test.Result.HumanReadablePrice {
			t.Errorf("%s Error: normalizedTrade.HumanReadablePrice is %s, expected %s", test.name, normalizedTrade.HumanReadablePrice, test.Result.HumanReadablePrice)
		}
		if !decCompare(normalizedTrade.Amount, test.Result.Amount) {
			t.Errorf("%s Error: normalizedTrade.Amount is %s, expected %s", test.name, normalizedTrade.Amount, test.Result.Amount)
		}
		if normalizedTrade.SymbolAmount != test.Result.SymbolAmount {
			t.Errorf("%s Error: normalizedTrade.SymbolAmount is %s, expected %s", test.name, normalizedTrade.SymbolAmount, test.Result.SymbolAmount)
		}
	}
}

func decCompare(a, b *decimal.Decimal) bool {
	r := decimal.ToSDec(a)
	s := decimal.ToSDec(b)
	return r.Sub(*s).IsZero()
}
