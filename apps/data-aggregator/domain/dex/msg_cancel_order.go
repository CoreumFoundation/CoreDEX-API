package dex

import (
	"context"

	"github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/domain"
	"github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	ordermodel "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	orderclient "github.com/CoreumFoundation/CoreDEX-API/domain/order/client"
	"github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	tradeclient "github.com/CoreumFoundation/CoreDEX-API/domain/trade/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/gogoproto/proto"

	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

type MsgCancelOrderHandler struct {
	registry *domain.Registry
}

func NewMsgCancelOrderHandler(interfaceRegistry types.InterfaceRegistry, parserRegistry *domain.Registry) {
	registry := domain.NewRegistry(interfaceRegistry)
	registry.RegisterEventHandler(&EventOrderClosedHandler{})

	parserRegistry.RegisterMsgHandler(&MsgCancelOrderHandler{
		registry: registry,
	})
}

func (e *MsgCancelOrderHandler) MsgType() proto.Message {
	return &dextypes.MsgCancelOrder{}
}

func (e *MsgCancelOrderHandler) Parse(txBytes []byte, meta domain.Metadata) proto.Message {
	msg := &dextypes.MsgCancelOrder{}

	err := proto.Unmarshal(txBytes, msg)
	if err != nil {
		logger.Errorf("Error unmarshalling MsgCancelOrder %s: %v", string(txBytes), err)
		return nil
	}

	return msg
}

func (e *MsgCancelOrderHandler) Handle(
	ctx context.Context,
	orderClient ordermodel.OrderServiceClient,
	_ trade.TradeServiceClient,
	_ currency.CurrencyServiceClient,
	action domain.Action,
	_ proto.Message,
	meta domain.Metadata,
	_ chan *trade.Trade,
) error {
	for _, ev := range action.Events {
		tr := e.registry.ParseEvent(ev.Type, ev)
		if tr == nil {
			continue
		}
		event, ok := tr.(*dextypes.EventOrderClosed)
		if !ok {
			continue
		}
		order, err := orderClient.Get(tradeclient.AuthCtx(ctx), &ordermodel.ID{
			Network:  meta.Network,
			Sequence: int64(event.Sequence), // TODO: decide if we want to change both to unsigned
		})
		if err != nil {
			return err
		}

		if meta.IsEndBlockEvent {
			order.OrderStatus = ordermodel.OrderStatus_ORDER_STATUS_EXPIRED
		} else {
			order.OrderStatus = ordermodel.OrderStatus_ORDER_STATUS_CANCELED
		}
		_, err = orderClient.Upsert(orderclient.AuthCtx(ctx), order)
		if err != nil {
			return err
		}
	}
	return nil
}
