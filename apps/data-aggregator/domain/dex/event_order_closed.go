package dex

import (
	"strconv"

	"cosmossdk.io/math"
	cmtypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

type EventOrderClosedHandler struct{}

func (e *EventOrderClosedHandler) EventType() proto.Message {
	return &dextypes.EventOrderClosed{}
}

func (e *EventOrderClosedHandler) Parse(event cmtypes.Event) proto.Message {
	res := &dextypes.EventOrderClosed{}
	for _, attribute := range event.Attributes {
		switch attribute.Key {
		case "creator":
			res.Creator = attribute.Value
		case "id":
			res.ID = attribute.Value
		case "sequence":
			id, err := strconv.ParseUint(attribute.Value, 10, 64)
			if err != nil {
				logger.Errorf("failed to parse EventOrderClosed sequence %s as uint64 : %v", attribute.Value, err)
				continue
			}
			res.Sequence = id
		case "remaining_balance":
			value, ok := math.NewIntFromString(attribute.Value)
			if !ok {
				logger.Errorf("failed to parse EventOrderClosed remaining_balance %s as integer", attribute.Value)
				return nil
			}
			res.RemainingBalance = value
		case "remaining_quantity":
			value, ok := math.NewIntFromString(attribute.Value)
			if !ok {
				logger.Errorf("failed to parse EventOrderClosed remaining_quantity %s as integer", attribute.Value)
				return nil
			}
			res.RemainingQuantity = value
		}
	}
	return res
}
