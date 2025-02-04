package domain

import (
	cmtypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"
)

type EventHandler interface {
	EventType() proto.Message
	Parse(event cmtypes.Event) proto.Message
}
