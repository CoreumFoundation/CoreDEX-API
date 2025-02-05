package coreum

import (
	"context"
	"fmt"
	gobig "math/big"
	"sort"
	"sync"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/shopspring/decimal"

	"github.com/CoreumFoundation/coreum/v5/pkg/math/big"
	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

const Precision = 16 // TODO

// QueryOrderBooks returns list of available order books. the paginationKey should nil for the first page.
// the nextPaginationKey will be nil if there are no more pages.
func (r *Reader) QueryOrderBooks(
	ctx context.Context, paginationKey []byte,
) (data []dextypes.OrderBookData, nextPaginationKey []byte, err error) {
	dexClient := dextypes.NewQueryClient(nodeConnections[r.Network])
	res, err := dexClient.OrderBooks(ctx, &dextypes.QueryOrderBooksRequest{
		Pagination: &query.PageRequest{Key: paginationKey},
	})
	if err != nil {
		return nil, nil, err
	}
	return res.OrderBooks, res.Pagination.NextKey, nil
}

// QueryOrderBookOrders returns orders inside an order book. the paginationKey should nil for the first page.
// the nextPaginationKey will be nil if there are no more pages.
func (r *Reader) QueryOrderBookOrders(
	ctx context.Context, denom1, denom2 string, side dextypes.Side, limit uint64, reverse bool,
) (orders []dextypes.Order, nextPaginationKey []byte, err error) {
	dexClient := dextypes.NewQueryClient(nodeConnections[r.Network])
	res, err := dexClient.OrderBookOrders(ctx, &dextypes.QueryOrderBookOrdersRequest{
		BaseDenom:  denom1,
		QuoteDenom: denom2,
		Side:       side,
		Pagination: &query.PageRequest{Limit: limit, Reverse: reverse},
	})
	if err != nil {
		return nil, nil, err
	}
	return res.Orders, res.Pagination.NextKey, nil
}

type OrderBookOrder struct {
	priceDec           *gobig.Rat
	Price              string
	HumanReadablePrice string
	Amount             string
	SymbolAmount       string
	Sequence           uint64
	Account            string
	OrderID            string
}

type OrderBookOrders struct {
	Buy  []*OrderBookOrder
	Sell []*OrderBookOrder
}

func (r *Reader) QueryOrderBookBySide(ctx context.Context,
	denom1, denom2 string,
	denom1Precision, denom2Precision int64,
	limit uint64,
	side dextypes.Side,
	reverse bool,
	invert bool) ([]*OrderBookOrder, error) {
	res, _, err := r.QueryOrderBookOrders(ctx, denom1, denom2, side, limit, reverse)
	if err != nil {
		return nil, err
	}
	orders := make([]*OrderBookOrder, 0)
	switch invert {
	case false:
		for _, order := range res {
			price, ok := new(gobig.Rat).SetString(order.Price.String())
			if !ok {
				return nil, fmt.Errorf("could not parse order price %s as big.Rat", order.Price.String())
			}
			var precision *gobig.Rat
			precisionDiff := denom1Precision - denom2Precision
			if precisionDiff < 0 {
				precision = big.RatInv(decimal.New(1, int32(-precisionDiff)).Rat())
			} else if precisionDiff > 0 {
				precision = decimal.New(1, int32(-precisionDiff)).Rat()
			} else {
				precision = big.NewRatFromInt64(1)
			}
			humanReadablePrice := big.RatMul(price, precision)
			quantity := new(gobig.Rat).SetInt(order.Quantity.BigInt())
			symbolAmount := quantity
			if denom1Precision != 0 {
				symbolAmount = big.RatMul(quantity, big.RatInv(big.NewRatFromInt64(denom1Precision)))
			}
			orders = append(orders, &OrderBookOrder{
				priceDec:           price,
				Price:              ratToString(price),
				HumanReadablePrice: ratToString(humanReadablePrice),
				Amount:             ratToString(quantity),
				SymbolAmount:       ratToString(symbolAmount),
				Sequence:           order.Sequence,
				Account:            order.Creator,
				OrderID:            order.ID,
			})
		}
	case true:
		for _, order := range res {
			orderPrice, ok := new(gobig.Rat).SetString(order.Price.String())
			if !ok {
				return nil, fmt.Errorf("could not parse order price %s as big.Rat", order.Price.String())
			}
			invPrice := big.RatInv(orderPrice)
			var precision *gobig.Rat
			precisionDiff := denom2Precision - denom1Precision
			if precisionDiff < 0 {
				precision = big.RatInv(decimal.New(1, int32(-precisionDiff)).Rat())
			} else if precisionDiff > 0 {
				precision = decimal.New(1, int32(-precisionDiff)).Rat()
			} else {
				precision = big.NewRatFromInt64(1)
			}
			humanReadablePrice := big.RatMul(invPrice, precision)
			quantity := big.RatMul(new(gobig.Rat).SetInt(order.Quantity.BigInt()), invPrice)
			symbolAmount := quantity
			if denom1Precision != 0 {
				symbolAmount = big.RatMul(quantity, big.RatInv(decimal.New(1, int32(denom1Precision)).Rat()))
			}
			orders = append(orders, &OrderBookOrder{
				priceDec:           invPrice,
				Price:              ratToString(invPrice),
				HumanReadablePrice: ratToString(humanReadablePrice),
				Amount:             ratToString(quantity),
				SymbolAmount:       ratToString(symbolAmount),
				Sequence:           order.Sequence,
				Account:            order.Creator,
				OrderID:            order.ID,
			})
		}
	}
	return orders, nil
}

// QueryOrderBookRelevantOrders returns orders inside an order book around the spread.
func (r *Reader) QueryOrderBookRelevantOrders(ctx context.Context, denom1, denom2 string, denom1Precision, denom2Precision int64, limit uint64) (orders *OrderBookOrders, err error) {
	orderBookOrders := &OrderBookOrders{
		Buy:  make([]*OrderBookOrder, 0),
		Sell: make([]*OrderBookOrder, 0),
	}

	var queryError error
	lock := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		orders, err := r.QueryOrderBookBySide(ctx, denom1, denom2, denom1Precision, denom2Precision,
			limit, dextypes.SIDE_SELL, false, false)
		if err != nil {
			queryError = err
			return
		}
		lock.Lock()
		defer lock.Unlock()
		orderBookOrders.Sell = append(orderBookOrders.Sell, orders...)
	}()

	go func() {
		defer wg.Done()
		orders, err := r.QueryOrderBookBySide(ctx, denom1, denom2, denom1Precision, denom2Precision,
			limit, dextypes.SIDE_BUY, true, false)
		if err != nil {
			queryError = err
			return
		}
		lock.Lock()
		defer lock.Unlock()
		orderBookOrders.Buy = append(orderBookOrders.Buy, orders...)
	}()

	go func() {
		defer wg.Done()
		orders, err := r.QueryOrderBookBySide(ctx, denom2, denom1, denom1Precision, denom2Precision,
			limit, dextypes.SIDE_SELL, true, true)
		if err != nil {
			queryError = err
			return
		}
		lock.Lock()
		defer lock.Unlock()
		orderBookOrders.Buy = append(orderBookOrders.Buy, orders...)
	}()

	go func() {
		defer wg.Done()
		orders, err := r.QueryOrderBookBySide(ctx, denom2, denom1, denom1Precision, denom2Precision,
			limit, dextypes.SIDE_BUY, false, true)
		if err != nil {
			queryError = err
			return
		}
		lock.Lock()
		defer lock.Unlock()
		orderBookOrders.Sell = append(orderBookOrders.Sell, orders...)
	}()

	wg.Wait()
	if queryError != nil {
		return nil, queryError
	}

	sort.SliceStable(orderBookOrders.Sell, func(i, j int) bool {
		return big.RatLT(orderBookOrders.Sell[i].priceDec, orderBookOrders.Sell[j].priceDec)
	})
	if uint64(len(orderBookOrders.Sell)) > limit {
		orderBookOrders.Sell = orderBookOrders.Sell[0:limit]
	}

	sort.SliceStable(orderBookOrders.Buy, func(i, j int) bool {
		//return big.RatGT(orderBookOrders.Buy[i].priceDec, orderBookOrders.Buy[j].priceDec)
		return orderBookOrders.Buy[i].priceDec.Cmp(orderBookOrders.Buy[j].priceDec) == 1
	})
	if uint64(len(orderBookOrders.Buy)) > limit {
		orderBookOrders.Buy = orderBookOrders.Buy[0:limit]
	}

	return orderBookOrders, nil
}

func ratToString(num *gobig.Rat) string {
	return decimal.NewFromBigRat(num, Precision).String()
}
