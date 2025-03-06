package ohlc

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	decimal "github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	ohlcclient "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc/client"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type Application struct {
	tradeChan          chan *tradegrpc.Trade
	ohlcClient         ohlcgrpc.OHLCServiceClient
	ohlcCache          []*ohlcgrpc.OHLC
	ohlcCacheResetTime time.Time
	mutex              *sync.RWMutex
}

func NewApplication(ctx context.Context, tradeChan chan *tradegrpc.Trade) *Application {
	app := &Application{
		tradeChan:          tradeChan,
		ohlcClient:         ohlcclient.Client(),
		ohlcCache:          make([]*ohlcgrpc.OHLC, 0),
		ohlcCacheResetTime: time.Now(),
		mutex:              &sync.RWMutex{},
	}
	app.StartOHLCProcessor()
	return app
}

func (a *Application) StartOHLCProcessor() {
	trades := map[string][]*tradegrpc.Trade{}
	for {
		// Select n trade or x seconds into a map of symbols (denom1-denom2)
		select {
		case trade := <-a.tradeChan:
			// Skip trades that are not enriched (e.g. would have borked data)
			if !trade.Enriched {
				continue
			}
			symbol := symbol(trade)
			/*
				Trades can be buy or sell.
				The amounts and price are stored for the associated buy or sell,
				To be able to build the associated OHLCs we can apply the BUY to one side and the SELL to the other side.
				Side is defined as Denom1-Denom2 and Denom2-Denom1
				The associated price and amount need to be inverted for the other side.
			*/
			switch trade.Side {
			case orderproperties.Side_SIDE_BUY:
				if _, ok := trades[symbol]; !ok {
					trades[symbol] = []*tradegrpc.Trade{}
				}
				trades[symbol] = append(trades[symbol], trade)
			case orderproperties.Side_SIDE_SELL:
				// Invert the trade
				trade.Denom1, trade.Denom2 = trade.Denom2, trade.Denom1
				r := trade.Amount.Mul(trade.Price)
				trade.Amount = decimal.FromFloat64(r)
				trade.Price = 1 / trade.Price
				// Invert symbol:
				s := strings.Split(symbol, "_")
				symbol = fmt.Sprintf("%s_%s", s[1], s[0])
				if _, ok := trades[symbol]; !ok {
					trades[symbol] = []*tradegrpc.Trade{}
				}
				trades[symbol] = append(trades[symbol], trade)
			}
			if a.len(trades) > 100 { // batch 100 trades
				a.calculateOHLCS(trades)
				trades = map[string][]*tradegrpc.Trade{}
			}
		case <-time.After(5 * time.Second):
			// Process the trades
			a.calculateOHLCS(trades)
			trades = map[string][]*tradegrpc.Trade{}
		}
	}
}

func (*Application) len(trades map[string][]*tradegrpc.Trade) int {
	l := 0
	for _, v := range trades {
		l += len(v)
	}
	return l
}

// Symbol uses _ as a separator between the two denominations: / and - where already used in ibc and base currency annotations
// (so _ sidesteps any potential issues with whatever is passed in)
func symbol(trade *tradegrpc.Trade) string {
	return fmt.Sprintf("%s_%s", trade.Denom1.ToString(), trade.Denom2.ToString())
}

// Retrieve OHLCs for the given symbol
func (a *Application) getSymbol(base *timestamppb.Timestamp, symbol string) map[string]*ohlcgrpc.OHLC {
	// The ohlc PeriodsList contains all the periods in the stored notation
	// These periods represent the buckets which need to be calculated
	pb := make([]*ohlcgrpc.PeriodBucket, 0)
	m := make(map[string]*ohlcgrpc.OHLC)
	a.mutex.RLock()
	for _, v := range ohlcgrpc.PeriodsList {
		skip := false
		for _, ohlc := range a.ohlcCache {
			// Filter out intervals we already have in the cachedOHLC, and prepopulate the required result set with cached data
			if strings.Compare(v.String(), ohlc.Period.String()) == 0 {
				if v.ToOHLCKeyTimestamppb(base).AsTime().Unix() == ohlc.Timestamp.AsTime().Unix() &&
					ohlc.Symbol == symbol {
					m[ohlc.Period.String()] = ohlc
					skip = true
					break
				}
			}
		}
		if skip {
			continue
		}
		pb = append(pb, &ohlcgrpc.PeriodBucket{
			Period:    v,
			Timestamp: v.ToOHLCKeyTimestamppb(base),
		})
	}
	a.mutex.RUnlock()
	o, err := a.ohlcClient.GetOHLCsForPeriods(context.Background(), &ohlcgrpc.PeriodsFilter{
		Symbol:  symbol,
		Periods: pb,
	})
	if err != nil {
		logger.Errorf("Error getting ohlc for symbol %s at %v: %v", symbol, base, err)
	}
	if o == nil {
		o = &ohlcgrpc.OHLCs{}
	}
	if o.OHLCs == nil {
		o.OHLCs = make([]*ohlcgrpc.OHLC, 0)
	}
	// Process the returned ohlcs into a map of ohlc period:
	for _, ohlc := range o.OHLCs {
		m[ohlc.Period.String()] = ohlc
	}
	// Check if all periods are present, if not, add the remainder
	for _, v := range ohlcgrpc.PeriodsList {
		if _, ok := m[v.String()]; !ok {
			m[v.String()] = &ohlcgrpc.OHLC{
				Symbol:    symbol,
				Period:    v,
				Timestamp: v.ToOHLCKeyTimestamppb(base),
				MetaData: &metadata.MetaData{
					CreatedAt: timestamppb.Now(),
					UpdatedAt: timestamppb.Now(),
				},
			}
		}
	}
	return m
}

func (a *Application) calculateOHLCS(inputTrades map[string][]*tradegrpc.Trade) {
	wg := &sync.WaitGroup{}
	// Dump the cache every 15 minutes (very simple way of managing the cache)
	if len(a.ohlcCache) > 0 && time.Since(a.ohlcCacheResetTime) > 15*time.Minute {
		a.ohlcCache = make([]*ohlcgrpc.OHLC, 0)
		a.ohlcCacheResetTime = time.Now()
	}
	wg.Add(len(inputTrades))
	for symbol, trades := range inputTrades {
		go a.calculateOHLC(trades, symbol, wg)
	}
	wg.Wait()
}

func (a *Application) calculateOHLC(inputTrades []*tradegrpc.Trade, symbol string, wg *sync.WaitGroup) {
	tStart := time.Now()
	// Sort trades by block time ascending so that we can handle the same minute optimized
	// (leads to all other OHLC associated to the same block time to only to be retrieved and written once)
	sort.Slice(inputTrades, func(i, j int) bool {
		return inputTrades[i].BlockTime.AsTime().Before(inputTrades[j].BlockTime.AsTime())
	})

	previousMinute := int64(0)
	var symbolData map[string]*ohlcgrpc.OHLC
	toPersistOHLCs := make([]*ohlcgrpc.OHLC, 0)
	for i, trade := range inputTrades {
		// Retrieve the symbol data but only if we do not have symbolData and the minute has changed
		p := ohlcgrpc.Period{
			PeriodType: ohlcgrpc.PeriodType_PERIOD_TYPE_MINUTE,
			Duration:   1,
		}
		currentMinute := p.ToOHLCKeyTimestamppb(trade.BlockTime).Seconds
		logger.Infof("Processing trade %d for symbol %s, minute: %d", i, symbol, currentMinute)
		// This location of the comparison if the minute has changed and the use of the pointer
		// prevent the edge case of missing the first or last set of records (and having to handle those separately)
		if previousMinute != currentMinute {
			symbolData = a.getSymbol(trade.BlockTime, symbol)
			// Add the pointers to the ohlc data to the toPersistOHLCs set
			for _, ohlc := range symbolData {
				// Check if the ohlc is already in the toPersistOHLCs set
				// If it is, skip it
				skip := false
				for _, p := range toPersistOHLCs {
					if p.Period.String() == ohlc.Period.String() && p.Timestamp.AsTime().Equal(ohlc.Timestamp.AsTime()) &&
						p.Symbol == ohlc.Symbol {
						skip = true
						break
					}
				}
				if skip {
					continue
				}
				toPersistOHLCs = append(toPersistOHLCs, ohlc)
			}
			// Add the pointers to the cache:
			a.mutex.Lock()
			for _, ohlc := range symbolData {
				// Check if the ohlc is already in the cache
				// If it is, skip it
				skip := false
				for _, p := range a.ohlcCache {
					if p.Period.String() == ohlc.Period.String() && p.Timestamp.AsTime().Equal(ohlc.Timestamp.AsTime()) &&
						p.Symbol == ohlc.Symbol {
						skip = true
						break
					}
				}
				if skip {
					continue
				}
				a.ohlcCache = append(a.ohlcCache, ohlc)
			}
			a.mutex.Unlock()
			previousMinute = currentMinute
		}
		// For all the periods, calculate the OHLC
		// Symbols contains all the periods the trades need to be applied for (until we hit a new minute)
		// Applying is determining the high,low,open,close for each period, plus summing the volume and incrementing the number
		// of trades for each period
		// A period can be empty, which has to be taken into account when comparing the low value with a potential 0 value
		// If no values are present, set the open.
		// Close is set if the time is larger than the last recorded record in that period only (trades are not guaranteed to be in order)
		for _, ohlc := range symbolData {
			price := trade.Price
			if ohlc.Open == 0 || ohlc.OpenTime == nil || trade.BlockTime.AsTime().Before(ohlc.OpenTime.AsTime()) {
				ohlc.OpenTime = trade.BlockTime
				ohlc.Open = price
			}
			if ohlc.Close == 0 || ohlc.CloseTime == nil || trade.BlockTime.AsTime().After(ohlc.CloseTime.AsTime()) {
				ohlc.CloseTime = trade.BlockTime
				ohlc.Close = price
			}
			if price > ohlc.High {
				ohlc.High = price
			}
			if ohlc.Low == 0 || price < ohlc.Low {
				ohlc.Low = price
			}
			ohlc.NumberOfTrades++
			ohlc.Volume += trade.Amount.Float64()
			ohlc.QuoteVolume += trade.Amount.Mul(trade.Price)
			ohlc.MetaData.UpdatedAt = timestamppb.Now()
			ohlc.MetaData.Network = trade.MetaData.Network
		}
	}
	// Persist the data
	_, err := a.ohlcClient.BatchUpsert(context.Background(), &ohlcgrpc.OHLCs{
		OHLCs: toPersistOHLCs,
	})
	if err != nil {
		logger.Errorf("Error upserting ohlcs for symbol %s: %v", symbol, err)
	}
	logger.Infof("Processed %d trades for symbol %s in %d microseconds", len(inputTrades), symbol, time.Since(tStart).Microseconds())
	wg.Done()
}
