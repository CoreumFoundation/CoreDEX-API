package domain

import (
	"errors"

	"github.com/CoreumFoundation/CoreDEX-API/domain/symbol"
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
