package store

import (
	"github.com/CoreumFoundation/CoreDEX-API/apps/store/store/currency"
	ohlc "github.com/CoreumFoundation/CoreDEX-API/apps/store/store/ohlc"
	"github.com/CoreumFoundation/CoreDEX-API/apps/store/store/order"
	"github.com/CoreumFoundation/CoreDEX-API/apps/store/store/state"
	"github.com/CoreumFoundation/CoreDEX-API/apps/store/store/trade"
	storebase "github.com/CoreumFoundation/CoreDEX-API/utils/mysqlstore"
)

type StoreBase struct {
	Trade    *trade.Application
	State    *state.Application
	Order    *order.Application
	OHLC     *ohlc.Application
	Currency *currency.Application
}

func NewStore() *StoreBase {
	client := storebase.Client()

	s := &StoreBase{
		trade.NewApplication(client),
		state.NewApplication(client),
		order.NewApplication(client),
		ohlc.NewApplication(client),
		currency.NewApplication(client),
	}
	s.index()
	return s
}

func (s *StoreBase) index() {
}
