package domain

import (
	"context"

	"github.com/cosmos/gogoproto/proto"

	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	"github.com/CoreumFoundation/CoreDEX-API/domain/trade"
)

type MsgHandler interface {
	MsgType() proto.Message
	Parse(txBytes []byte, meta Metadata) proto.Message
	Handle(
		ctx context.Context,
		orderClient ordergrpc.OrderServiceClient,
		tradeClient trade.TradeServiceClient,
		currencyClient currencygrpc.CurrencyServiceClient,
		action Action,
		message proto.Message,
		meta Metadata,
		tradeChan chan *trade.Trade,
	) error
}
