// Package state Manages the state.
// State is per app and network.
// The registered final value is in int64 height
package state

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	stategrpc "github.com/CoreumFoundation/CoreDEX-API/domain/state"
	stateclient "github.com/CoreumFoundation/CoreDEX-API/domain/state/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

const defaultStateType = stategrpc.StateType_BLOCK_HEIGHT

type State struct {
	state       map[metadata.Network]int64
	stateMutex  *sync.Mutex
	stateClient stategrpc.StateServiceClient
	stateChan   map[string]chan string
}

type Content struct {
	Height int64
}

// Connect to the state store and mandatory load the data in there to be used by the app
// To be able to load the state store data, the networks config is required
func NewApplication(ctx context.Context) *State {
	s := &State{
		state:       make(map[metadata.Network]int64),
		stateMutex:  &sync.Mutex{},
		stateClient: stateclient.Client(),
		stateChan:   make(map[string]chan string),
	}
	go s.updateState()
	return s
}

func (s *State) GetState(_ context.Context, network metadata.Network) int64 {
	sq := stategrpc.StateQuery{
		Network:   network,
		StateType: defaultStateType,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	d, err := s.stateClient.Get(stateclient.AuthCtx(ctx), &sq)
	if err != nil { // Not having the record is not an error, it is just a clean new state (so no error is thrown in record not found)
		logger.Fatalf("failed to get state for record type %s and network %s: %v", defaultStateType.String(), network, err)
	}
	v := &Content{}
	err = json.Unmarshal([]byte(d.Content), v)
	if err != nil {
		logger.Errorf("no height in retrieved record type %s, content '%s' and network %s: %v. Setting default height 1", defaultStateType.String(), d.Content, network, err)
		v.Height = 1
	}
	s.state[network] = v.Height
	// v.Height = 2506618 // A block with a processable Trade
	return v.Height
}

func (s *State) SetState(network metadata.Network, height int64) {
	s.stateMutex.Lock()
	s.state[network] = height
	s.stateMutex.Unlock()
}

// Go routine for running the state update to the state store every 15 minutes
func (s *State) updateState() {
	for {
		time.Sleep(5 * time.Minute)
		s.FlushState()
	}
}

// Flush the state to the state store
func (s *State) FlushState() {
	s.stateMutex.Lock()
	for network, height := range s.state {
		c := &Content{
			Height: height,
		}
		b, err := json.Marshal(c)
		if err != nil {
			logger.Errorf("failed to marshal state for record type %s and network %s: %v", defaultStateType.String(), network.String(), err)
			continue
		}
		state := &stategrpc.State{
			MetaData: &metadata.MetaData{
				Network:   network,
				UpdatedAt: timestamppb.Now(),
			},
			StateType: defaultStateType,
			Content:   string(b),
		}
		_, err = s.stateClient.Upsert(stateclient.AuthCtx(context.Background()), state)
		if err != nil {
			logger.Errorf("failed to set state for record type %s and network %s: %v", defaultStateType.String(), network, err)
		}
	}
	s.stateMutex.Unlock()
}
