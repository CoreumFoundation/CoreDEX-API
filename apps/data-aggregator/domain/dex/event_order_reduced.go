package dex

import (
	"encoding/json"
	"strconv"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	cmtypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

type EventOrderReducedHandler struct{}

func (e *EventOrderReducedHandler) EventType() proto.Message {
	return &dextypes.EventOrderReduced{}
}

func (e *EventOrderReducedHandler) Parse(event cmtypes.Event) proto.Message {
	res := &dextypes.EventOrderReduced{}
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
		case "received_coin":
			var coin types.Coin
			if err := json.Unmarshal([]byte(attribute.Value), &coin); err == nil {
				res.ReceivedCoin = coin
			}
		case "sent_coin":
			var coin types.Coin
			if err := json.Unmarshal([]byte(attribute.Value), &coin); err == nil {
				res.SentCoin = coin
			}
		}
	}
	return res
}
