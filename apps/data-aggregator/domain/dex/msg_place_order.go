package dex

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/gogoproto/proto"
	dec "github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/domain"
	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	"github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	orderclient "github.com/CoreumFoundation/CoreDEX-API/domain/order/client"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	tradeclient "github.com/CoreumFoundation/CoreDEX-API/domain/trade/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

type MsgPlaceOrderHandler struct {
	registry *domain.Registry
}

func NewMsgPlaceOrderHandler(interfaceRegistry types.InterfaceRegistry, parserRegistry *domain.Registry) {
	registry := domain.NewRegistry(interfaceRegistry)
	registry.RegisterEventHandler(&EventOrderPlacedHandler{})
	registry.RegisterEventHandler(&EventOrderReducedHandler{})

	parserRegistry.RegisterMsgHandler(&MsgPlaceOrderHandler{
		registry: registry,
	})
}

func (e *MsgPlaceOrderHandler) MsgType() proto.Message {
	return &dextypes.MsgPlaceOrder{}
}

func (e *MsgPlaceOrderHandler) Parse(txBytes []byte, meta domain.Metadata) proto.Message {
	msg := &dextypes.MsgPlaceOrder{}

	err := proto.Unmarshal(txBytes, msg)
	if err != nil {
		logger.Errorf("Error unmarshalling MsgPlaceOrder %s: %v", string(txBytes), err)
		return nil
	}

	tr := &ordergrpc.Order{}
	tr.Account = msg.Sender
	tr.Type = ordergrpc.OrderType(msg.Type)
	tr.OrderID = msg.ID
	tr.BaseDenom, err = denom.NewDenom(msg.BaseDenom)
	if err != nil {
		logger.Errorf("Error parsing BaseDenom %s: %v", tr.BaseDenom, err)
		return nil
	}
	tr.QuoteDenom, err = denom.NewDenom(msg.QuoteDenom)
	if err != nil {
		logger.Errorf("Error parsing QuoteDenom %s: %v", tr.QuoteDenom, err)
		return nil
	}
	if msg.Type != dextypes.ORDER_TYPE_MARKET {
		tr.Price, _ = msg.Price.Rat().Float64()
	}
	tr.Quantity = decimal.FromDec(dec.NewFromBigInt(msg.Quantity.BigInt(), 0))
	tr.RemainingQuantity = decimal.FromDec(dec.NewFromBigInt(msg.Quantity.BigInt(), 0))
	tr.Side = orderproperties.Side(msg.Side)
	if msg.Type == dextypes.ORDER_TYPE_LIMIT {
		tr.TimeInForce = ordergrpc.TimeInForce(msg.TimeInForce)
	} else {
		tr.TimeInForce = ordergrpc.TimeInForce_TIME_IN_FORCE_UNSPECIFIED
	}
	tr.BlockTime = timestamppb.New(meta.BlockTime)
	tr.MetaData = &metadata.MetaData{
		Network:   meta.Network,
		UpdatedAt: timestamppb.Now(),
		CreatedAt: timestamppb.Now(),
	}
	tr.TXID = &meta.TxHash
	if msg.GoodTil != nil {
		tr.GoodTil = &ordergrpc.GoodTil{
			BlockHeight: int64(msg.GoodTil.GoodTilBlockHeight),
		}
		if msg.GoodTil.GoodTilBlockTime != nil {
			tr.GoodTil.BlockTime = timestamppb.New(*msg.GoodTil.GoodTilBlockTime)
		}
	}
	tr.BlockHeight = meta.BlockHeight

	return tr
}

func enrichDenoms(ctx context.Context, currencyClient currencygrpc.CurrencyServiceClient,
	meta domain.Metadata,
	order *ordergrpc.Order,
	enriched bool,
) bool {
	if !enriched {
		denom1Currency, err1 := currencyClient.Get(ctx, &currencygrpc.ID{
			Network: meta.Network,
			Denom:   order.BaseDenom.Denom,
		})
		denom2Currency, err2 := currencyClient.Get(ctx, &currencygrpc.ID{
			Network: meta.Network,
			Denom:   order.QuoteDenom.Denom,
		})
		if err1 == nil && err2 == nil {
			enriched = true
			if denom1Currency.Denom != nil && denom1Currency.Denom.Precision != nil {
				order.BaseDenom.Precision = denom1Currency.Denom.Precision
			}
			if denom2Currency.Denom != nil && denom2Currency.Denom.Precision != nil {
				order.QuoteDenom.Precision = denom2Currency.Denom.Precision
			}
		}
	}
	return enriched
}

func (e *MsgPlaceOrderHandler) Handle(
	ctx context.Context,
	orderClient ordergrpc.OrderServiceClient,
	tradeClient tradegrpc.TradeServiceClient,
	currencyClient currencygrpc.CurrencyServiceClient,
	action domain.Action,
	message proto.Message,
	meta domain.Metadata,
	tradeChan chan *tradegrpc.Trade,
) error {

	for _, ev := range action.Events {
		tr := e.registry.ParseEvent(ev.Type, ev)
		if tr == nil {
			continue
		}

		switch event := tr.(type) {
		case *dextypes.EventOrderPlaced:
			enriched := false
			var err error
			order := message.(*ordergrpc.Order)
			if order.Account != event.Creator || order.OrderID != event.ID {
				continue
			}
			order.Sequence = int64(event.Sequence)
			order.OrderStatus = ordergrpc.OrderStatus_ORDER_STATUS_OPEN
			order.Enriched = enrichDenoms(ctx, currencyClient, meta, order, enriched)

			if order.Enriched {
				price := dec.NewFromFloat(order.Price)
				quantity := dec.New(order.Quantity.Value, order.Quantity.Exp)
				remainingExp := order.RemainingQuantity.Exp
				if order.RemainingQuantity.Value == 0 {
					remainingExp = 0
				}
				remainingQuantity := dec.New(order.RemainingQuantity.Value, remainingExp)

				order.Quantity = decimal.FromDec(quantity)
				order.RemainingQuantity = decimal.FromDec(remainingQuantity)
				order.Price, _ = price.Float64()
				order.Enriched = true
			}

			_, err = orderClient.Upsert(orderclient.AuthCtx(ctx), order)
			if err != nil {
				return err
			}
		case *dextypes.EventOrderReduced:
			enriched := false
			order, err := orderClient.Get(orderclient.AuthCtx(ctx), &ordergrpc.ID{
				Network:  meta.Network,
				Sequence: int64(event.Sequence),
			})
			if err != nil {
				logger.Errorf("the original order with sequence %d was not found in the local database while processing height %d: %v", event.Sequence, meta.BlockHeight, err.Error())
				continue
			}

			enriched = enrichDenoms(ctx, currencyClient, meta, order, enriched)

			var amount dec.Decimal
			var price float64
			switch order.Side {
			case orderproperties.Side_SIDE_SELL:
				val, err := dec.NewFromString(event.SentCoin.Amount.String())
				if err != nil {
					return err
				}
				order.RemainingQuantity = decimal.FromDec(
					dec.New(order.RemainingQuantity.Value, order.RemainingQuantity.Exp).Sub(val),
				)
				amount, err = dec.NewFromString(event.SentCoin.Amount.String())
				if err != nil {
					return err
				}
				d, err := dec.NewFromString(event.ReceivedCoin.Amount.String())
				if err != nil {
					return err
				}
				price, _ = d.Div(amount).Float64()
			case orderproperties.Side_SIDE_BUY:
				val, err := dec.NewFromString(event.ReceivedCoin.Amount.String())
				if err != nil {
					return err
				}
				order.RemainingQuantity = decimal.FromDec(
					dec.New(order.RemainingQuantity.Value, order.RemainingQuantity.Exp).Sub(val),
				)
				amount, err = dec.NewFromString(event.ReceivedCoin.Amount.String())
				if err != nil {
					return err
				}
				d, err := dec.NewFromString(event.SentCoin.Amount.String())
				if err != nil {
					return err
				}
				price, _ = d.Div(amount).Float64()
			default:
				logger.Errorf("unexpected side %s", order.Side.String())
				continue
			}
			if order.RemainingQuantity.Value == 0 {
				order.RemainingQuantity.Exp = 0
			}

			if dec.New(order.RemainingQuantity.Value, order.RemainingQuantity.Exp).IsZero() {
				order.OrderStatus = ordergrpc.OrderStatus_ORDER_STATUS_FILLED
			}
			_, err = orderClient.Upsert(orderclient.AuthCtx(ctx), order)
			if err != nil {
				return err
			}

			// store trade
			trade := &tradegrpc.Trade{
				Account:   event.Creator,
				OrderID:   event.ID,
				Sequence:  int64(event.Sequence),
				Amount:    decimal.FromDec(amount),
				Price:     price,
				Denom1:    order.BaseDenom,
				Denom2:    order.QuoteDenom,
				Side:      order.Side,
				BlockTime: timestamppb.New(meta.BlockTime),
				MetaData: &metadata.MetaData{
					Network:   meta.Network,
					UpdatedAt: timestamppb.Now(),
					CreatedAt: timestamppb.Now(),
				},
				TXID:        &meta.TxHash,
				BlockHeight: meta.BlockHeight,
				USD:         nil, // TODO
				Enriched:    enriched,
				Processed:   false,
			}
			_, err = tradeClient.Upsert(tradeclient.AuthCtx(ctx), trade)
			if err != nil {
				return err
			}
			tradeChan <- trade
		case *dextypes.EventOrderClosed:
			order, err := orderClient.Get(orderclient.AuthCtx(ctx), &ordergrpc.ID{
				Network:  meta.Network,
				Sequence: int64(event.Sequence),
			})
			if err != nil {
				return err
			}
			order.OrderStatus = ordergrpc.OrderStatus_ORDER_STATUS_FILLED
			val, err := dec.NewFromString(event.RemainingBaseQuantity.String())
			if err != nil {
				return err
			}
			order.RemainingQuantity = decimal.FromDec(val)
			if order.RemainingQuantity.Value == 0 {
				order.RemainingQuantity.Exp = 0
			}
			_, err = orderClient.Upsert(orderclient.AuthCtx(ctx), order)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
