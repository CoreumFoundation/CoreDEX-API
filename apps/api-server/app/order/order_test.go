/*
The tests in this package are to verify the applicaiton of the precision on the data which is represented by the denoms.
*/
package order

import (
	"context"
	"os"
	"testing"

	currencyapp "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/currency"
	"github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	"github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
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
	currencyApp := currencyapp.NewApplication(currencyService)
	return NewApplication(currencyApp)
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
		normalizedOrder, err := app.Normalize(context.Background(), order)
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

func decCompare(a, b *decimal.Decimal) bool {
	r := decimal.ToSDec(a)
	s := decimal.ToSDec(b)
	return r.Sub(*s).IsZero()
}
