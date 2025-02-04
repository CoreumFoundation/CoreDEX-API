package dex

import (
	"strconv"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	cmtypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"

	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

type EventOrderPlacedHandler struct{}

func (e *EventOrderPlacedHandler) EventType() proto.Message {
	return &dextypes.EventOrderPlaced{}
}

func (e *EventOrderPlacedHandler) Parse(event cmtypes.Event) proto.Message {
	res := &dextypes.EventOrderPlaced{}
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
		}
	}
	return res
}
