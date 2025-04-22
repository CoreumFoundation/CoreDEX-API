package trade

import (
	"context"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockTradeServiceClient struct {
	seq int
	db  map[string]*orderWrapper
}

func (c *MockTradeServiceClient) GetTradePairs(ctx context.Context, in *TradePairFilter, opts ...grpc.CallOption) (*TradePairs, error) {
	//TODO implement me
	panic("implement me")
}

func (c *MockTradeServiceClient) UpsertTradePair(ctx context.Context, in *TradePair, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

type orderWrapper struct {
	seq   int
	order *Trade
}

func NewMockTradeServiceClient() TradeServiceClient {
	return &MockTradeServiceClient{
		db: make(map[string]*orderWrapper),
	}
}

func (c *MockTradeServiceClient) Upsert(ctx context.Context, in *Trade, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	key := fmt.Sprintf("%d-%s", in.Sequence, in.MetaData.Network.String())
	var seq int
	if wrapper, exists := c.db[key]; exists {
		seq = wrapper.seq
	} else {
		c.seq += 1
		seq = c.seq
	}
	c.db[key] = &orderWrapper{
		seq:   seq,
		order: in,
	}
	return out, nil
}

func (c *MockTradeServiceClient) Get(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Trade, error) {
	key := fmt.Sprintf("%d-%s", in.Sequence, in.Network.String())
	order, exists := c.db[key]
	if !exists {
		return nil, errors.New("not found")
	}
	return order.order, nil
}

func (c *MockTradeServiceClient) GetAll(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*Trades, error) {
	wrappers := make([]*orderWrapper, 0, len(c.db))
	for _, wrapper := range c.db {
		wrappers = append(wrappers, wrapper)
	}
	sort.SliceStable(wrappers, func(i, j int) bool {
		return wrappers[i].seq <= wrappers[j].seq
	})
	res := make([]*Trade, len(wrappers))
	for i := range wrappers {
		res[i] = wrappers[i].order
	}
	return &Trades{
		Trades: res,
	}, nil
}

func (c *MockTradeServiceClient) BatchUpsert(ctx context.Context, in *Trades, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	panic("not implemented")
}
