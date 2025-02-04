package currency

import (
	"context"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockCurrencyServiceClient struct {
	seq int
	db  map[string]*currencyWrapper
}

type currencyWrapper struct {
	seq      int
	currency *Currency
}

func NewMockCurrencyServiceClient() CurrencyServiceClient {
	return &MockCurrencyServiceClient{
		db: make(map[string]*currencyWrapper),
	}
}

func (c *MockCurrencyServiceClient) Upsert(ctx context.Context, in *Currency, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	key := fmt.Sprintf("%s-%s", in.Denom.Denom, in.MetaData.Network.String())
	var seq int
	if wrapper, exists := c.db[key]; exists {
		seq = wrapper.seq
	} else {
		c.seq += 1
		seq = c.seq
	}
	c.db[key] = &currencyWrapper{
		seq:      seq,
		currency: in,
	}
	return out, nil
}

func (c *MockCurrencyServiceClient) Get(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Currency, error) {
	key := fmt.Sprintf("%s-%s", in.Denom, in.Network.String())
	currency, exists := c.db[key]
	if !exists {
		return nil, errors.New("not found")
	}
	return currency.currency, nil
}

func (c *MockCurrencyServiceClient) GetAll(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*Currencies, error) {
	wrappers := make([]*currencyWrapper, 0, len(c.db))
	for _, wrapper := range c.db {
		wrappers = append(wrappers, wrapper)
	}
	sort.SliceStable(wrappers, func(i, j int) bool {
		return wrappers[i].seq <= wrappers[j].seq
	})
	res := make([]*Currency, len(wrappers))
	for i := range wrappers {
		res[i] = wrappers[i].currency
	}
	return &Currencies{
		Currencies: res,
	}, nil
}

func (c *MockCurrencyServiceClient) BatchUpsert(ctx context.Context, in *Currencies, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	panic("not implemented")
}
