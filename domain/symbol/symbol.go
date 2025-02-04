package symbol

import (
	"fmt"
	"strings"

	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
)

type Symbol struct {
	Denom1 *denom.Denom
	Denom2 *denom.Denom
}

func NewSymbol(symb string) (*Symbol, error) {
	// Split the string by _ in the 2 denoms:
	denoms := strings.Split(symb, "_")
	if len(denoms) != 2 {
		return nil, fmt.Errorf("invalid symbol string: %s", symb)
	}
	denom1, err := denom.NewDenom(denoms[0])
	if err != nil {
		return nil, err
	}
	denom2, err := denom.NewDenom(denoms[1])
	if err != nil {
		return nil, err
	}
	return &Symbol{
		Denom1: denom1,
		Denom2: denom2,
	}, nil
}

// Symbol uses _ as a separator between the two denominations: / and - where already used in ibc and base currency annotations
// (so _ sidesteps any potential issues with whatever is passed in)
func (s *Symbol) ToString() string {
	return fmt.Sprintf("%s_%s", s.Denom1.ToString(), s.Denom2.ToString())
}
