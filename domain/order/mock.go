package order

import (
	"context"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockOrderServiceClient struct {
	seq int
	db  map[string]*orderWrapper
}

type orderWrapper struct {
	seq   int
	order *Order
}

func NewMockOrderServiceClient() OrderServiceClient {
	return &MockOrderServiceClient{
		db: make(map[string]*orderWrapper),
	}
}

func (c *MockOrderServiceClient) Upsert(ctx context.Context, in *Order, opts ...grpc.CallOption) (*emptypb.Empty, error) {
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

func (c *MockOrderServiceClient) Get(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Order, error) {
	key := fmt.Sprintf("%d-%s", in.Sequence, in.Network.String())
	order, exists := c.db[key]
	if !exists {
		return nil, errors.New("not found")
	}
	return order.order, nil
}

func (c *MockOrderServiceClient) GetAll(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*Orders, error) {
	wrappers := make([]*orderWrapper, 0, len(c.db))
	for _, wrapper := range c.db {
		wrappers = append(wrappers, wrapper)
	}
	sort.SliceStable(wrappers, func(i, j int) bool {
		return wrappers[i].seq <= wrappers[j].seq
	})
	res := make([]*Order, len(wrappers))
	for i := range wrappers {
		res[i] = wrappers[i].order
	}
	return &Orders{
		Orders: res,
	}, nil
}

func (c *MockOrderServiceClient) BatchUpsert(ctx context.Context, in *Orders, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	panic("not implemented")
}
