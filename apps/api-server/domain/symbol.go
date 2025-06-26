package domain

import (
	"errors"

	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	"github.com/CoreumFoundation/CoreDEX-API/domain/symbol"
	dec "github.com/shopspring/decimal"
)

// Business logic errors.
type Error interface {
	error
}

type Symbol struct {
	Base, Quote string
}

var ErrSymbolInvalid Error = errors.New("symbol is invalid")

func ValidSymbol(strSymbol string) bool {
	_, err := NewSymbolFromString(strSymbol)
	return err == nil
}

func NewSymbolFromString(symbStr string) (*Symbol, error) {
	symb, err := symbol.NewSymbol(symbStr)
	if err != nil {
		return nil, ErrSymbolInvalid
	}

	return &Symbol{
		Base:  symb.Denom1.ToString(),
		Quote: symb.Denom2.ToString(),
	}, nil
}

func ToSymbolPrice(baseDenomPrecision, quoteDenomPrecision int32, subunitPrice float64, quantity *dec.Decimal, side orderproperties.Side) dec.Decimal {
	price := dec.NewFromFloat(subunitPrice)
	quoteAmountSubunit := quantity
	baseAmountSubunit := quoteAmountSubunit.Mul(price)
	var humanReadablePrice dec.Decimal
	switch side {
	case orderproperties.Side_SIDE_SELL:
		// Round the output to the baseDenomPrecision to counter floating point deviations:
		humanReadablePrice = quoteAmountSubunit.
			Div(dec.New(1, quoteDenomPrecision)).
			Div(baseAmountSubunit.Div(dec.New(1, baseDenomPrecision))).
			Round(int32(baseDenomPrecision))
	case orderproperties.Side_SIDE_BUY:
		humanReadablePrice = price.Round(int32(quoteDenomPrecision))
	}
	return humanReadablePrice
}

// Trade are received ....
func ToSymbolTradeAmount(baseDenomPrecision, quoteDenomPrecision int32, subunitPrice float64, quantity *dec.Decimal, side orderproperties.Side) dec.Decimal {
	symbolAmount := *quantity
	switch side {
	case orderproperties.Side_SIDE_SELL:
		symbolAmount = symbolAmount.
			Mul(dec.NewFromFloat(subunitPrice)).
			Div(dec.New(1, int32(baseDenomPrecision))).
			Round(int32(baseDenomPrecision))
	case orderproperties.Side_SIDE_BUY:
		symbolAmount = symbolAmount.Div(dec.New(1, int32(quoteDenomPrecision))).Round(int32(quoteDenomPrecision))
	}
	return symbolAmount
}

func ToSymbolOrderAmount(baseDenomPrecision, quoteDenomPrecision int32, quantity *dec.Decimal, side orderproperties.Side) dec.Decimal {
	symbolAmount := *quantity
	switch side {
	case orderproperties.Side_SIDE_SELL:
		symbolAmount = symbolAmount.Div(dec.New(1, int32(baseDenomPrecision))).Round(int32(baseDenomPrecision))
	case orderproperties.Side_SIDE_BUY:
		symbolAmount = symbolAmount.Div(dec.New(1, int32(quoteDenomPrecision))).Round(int32(quoteDenomPrecision))
	}
	return symbolAmount
}
