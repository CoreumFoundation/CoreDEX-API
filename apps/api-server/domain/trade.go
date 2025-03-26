package domain

import (
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
)

type Trade struct {
	*tradegrpc.Trade
	HumanReadablePrice string
	SymbolAmount       string
	Status             ordergrpc.OrderStatus
}
