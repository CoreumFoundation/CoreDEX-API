package domain

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	cmtypes "github.com/cometbft/cometbft/abci/types"
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	"github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type Registry struct {
	InterfaceRegistry ctypes.InterfaceRegistry
	msgs              map[string]MsgHandler
	events            map[string]EventHandler
}

func NewRegistry(interfaceRegistry ctypes.InterfaceRegistry) *Registry {
	return &Registry{
		InterfaceRegistry: interfaceRegistry,
		msgs:              make(map[string]MsgHandler),
		events:            make(map[string]EventHandler),
	}
}

func (r *Registry) RegisterMsgHandler(handler MsgHandler) {
	// events don't start with slash, but types.MsgTypeURL always adds slash as a prefix
	// so we get rid of the first character before looking up the event handler
	r.msgs[types.MsgTypeURL(handler.MsgType())] = handler
}

func (r *Registry) RegisterEventHandler(handler EventHandler) {
	r.events[types.MsgTypeURL(handler.EventType())[1:]] = handler
}

func (r *Registry) ParseMsg(typeURL string, txBytes []byte, meta Metadata) proto.Message {
	if handler, ok := r.msgs[typeURL]; ok {
		return handler.Parse(txBytes, meta)
	}
	return nil
}

func (r *Registry) ParseEvent(typeURL string, event cmtypes.Event) proto.Message {
	if handler, ok := r.events[typeURL]; ok {
		return handler.Parse(event)
	}
	return nil
}

func (r *Registry) HandleAction(
	ctx context.Context,
	orderClient ordergrpc.OrderServiceClient,
	tradeClient trade.TradeServiceClient,
	currencyClient currencygrpc.CurrencyServiceClient,
	typeURL string,
	tx proto.Message,
	action Action,
	meta Metadata,
	tradeChan chan *trade.Trade,
) error {
	if handler, ok := r.msgs[typeURL]; ok {
		return handler.Handle(ctx, orderClient, tradeClient, currencyClient, action, tx, meta, tradeChan)
	}
	return nil
}

type Action struct {
	TypeURL  string
	Messages []types.Msg
	Events   []cmtypes.Event
}

func (r *Registry) ParseActions(
	ctx context.Context,
	orderClient ordergrpc.OrderServiceClient,
	tradeClient trade.TradeServiceClient,
	currencyClient currencygrpc.CurrencyServiceClient,
	message types.Msg,
	events []cmtypes.Event,
	meta Metadata,
	tradeChan chan *trade.Trade,
) {
	currentAction := Action{Events: make([]cmtypes.Event, 0)}
	actions := make([]Action, 0)
	for _, event := range events {
		if event.Type == "message" {
			for _, attribute := range event.Attributes {
				if attribute.Key == "action" {
					if currentAction.TypeURL != "" {
						actions = append(actions, currentAction)
					}
					currentAction = Action{
						TypeURL:  attribute.Value,
						Messages: nil,
						Events:   make([]cmtypes.Event, 0),
					}
					msg, err := r.InterfaceRegistry.Resolve(attribute.Value)
					if err == nil {
						currentAction.Messages = append(currentAction.Messages, msg)
					}
				}
			}
		}
		currentAction.Events = append(currentAction.Events, normalizeEvent(event))
	}
	if currentAction.TypeURL != "" {
		actions = append(actions, currentAction)
	}

	// Process the filtered Events
	for _, action := range actions {
		if err := r.HandleAction(ctx, orderClient, tradeClient, currencyClient, action.TypeURL, message, action, meta, tradeChan); err != nil {
			logger.Errorf("couldn't handle action %s : %v", action.TypeURL, err)
		}
	}
}

type Metadata struct {
	Network         metadata.Network
	BlockHeight     int64
	BlockTime       time.Time
	TxHash          string
	IsEndBlockEvent bool
	GasUsed         int64
}

func (r *Registry) HandleBlockEvent(
	ctx context.Context,
	orderClient ordergrpc.OrderServiceClient,
	tradeClient trade.TradeServiceClient,
	currencyClient currencygrpc.CurrencyServiceClient,
	events []cmtypes.Event,
	meta Metadata,
	tradeChan chan *trade.Trade,
) {
	normalizedEvents := make([]cmtypes.Event, len(events))
	for i, event := range events {
		normalizedEvents[i] = normalizeEvent(event)
	}
	action := Action{
		TypeURL: "/coreum.dex.v1.MsgCancelOrder",
		Events:  normalizedEvents,
	}
	if err := r.HandleAction(ctx, orderClient, tradeClient, currencyClient, action.TypeURL, nil, action, meta, tradeChan); err != nil {
		logger.Errorf("couldn't handle action %s : %v", action.TypeURL, err)
	}
}

func normalizeEvent(event cmtypes.Event) cmtypes.Event {
	for i := range event.Attributes {
		attr := event.Attributes[i]
		if strings.HasPrefix(attr.Value, "\"") && strings.HasSuffix(attr.Value, "\"") {
			var str string
			if err := json.Unmarshal([]byte(attr.Value), &str); err == nil {
				attr.Value = str
			}
		}
		event.Attributes[i] = attr
	}
	return event
}
