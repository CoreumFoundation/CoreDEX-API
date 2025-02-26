package ohlc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	ohlcclient "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc/client"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type Application struct {
	tradeChan  chan *tradegrpc.Trade
	ohlcClient ohlcgrpc.OHLCServiceClient
}

func NewApplication(ctx context.Context, tradeChan chan *tradegrpc.Trade) *Application {
	app := &Application{
		tradeChan:  tradeChan,
		ohlcClient: ohlcclient.Client(),
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
			// Process the trade
			symbol := symbol(trade)
			if _, ok := trades[symbol]; !ok {
				trades[symbol] = []*tradegrpc.Trade{}
			}
			/*
				Trades can be buy or sell.
				The amounts and price are stored for the associated buy or sell,
				To be able to build the associated OHLCs we can apply the BUY to one side and the SELL to the other side.
				Side is defined as Denom1-Denom2 and Denom2-Denom1
				The associated price and amount need to be inverted for the other side.

			*/
			switch trade.Side {
			case orderproperties.Side_SIDE_BUY:
				trades[symbol] = append(trades[symbol], trade)
				// case orderproperties.Side_SIDE_SELL:
				// 	// Invert the trade
				// 	trade.Denom1, trade.Denom2 = trade.Denom2, trade.Denom1
				// 	r := trade.Amount.Mul(trade.Price)
				// 	trade.Amount = decimal.FromFloat64(r)
				// 	trade.Price = 1 / trade.Price
				// 	// Invert symbol:
				// 	s := strings.Split(symbol, "_")
				// 	symbol = fmt.Sprintf("%s_%s", s[1], s[0])
				// 	trades[symbol] = append(trades[symbol], trade)
			}
			a.calculateOHLCS(trades)
			trades = map[string][]*tradegrpc.Trade{}
		case <-time.After(5 * time.Second):
			// Process the trades
			a.calculateOHLCS(trades)
			trades = map[string][]*tradegrpc.Trade{}
		}
	}
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
	for _, v := range ohlcgrpc.PeriodsList {
		pb = append(pb, &ohlcgrpc.PeriodBucket{
			Period:    v,
			Timestamp: v.ToOHLCKeyTimestamppb(base),
		})
	}
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
	// Process the return ohlcs into a map of ohlc period:
	m := make(map[string]*ohlcgrpc.OHLC)
	for _, ohlc := range o.OHLCs {
		m[ohlc.Period.String()] = ohlc
	}
	// Check if all periods are present, if not, add them
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
	for symbol, trades := range inputTrades {
		symbolData := a.getSymbol(trades[0].BlockTime, symbol)
		// For all the periods, calculate the OHLC
		// Symbols contains all the periods the trades need to be applied for
		// Applying is determining the high,low,open,close for each period, plus summing the volume and incrementing the number
		// of trades for each period
		// A period can be empty, which has to be taken into account when comparing the low value with a potential 0 value
		// If no values are present, set the open.
		// Close is set if the time is larger than the last recorded record in that period only (trades are not guaranteed to be in order)
		for _, ohlc := range symbolData {
			for _, trade := range trades {
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
		// Transform the symbolData into ohlcgrpc.OHLCs (array of ohlcgrpc.OHLC)
		upsert := make([]*ohlcgrpc.OHLC, 0)
		for _, ohlc := range symbolData {
			upsert = append(upsert, ohlc)
		}
		_, err := a.ohlcClient.BatchUpsert(context.Background(), &ohlcgrpc.OHLCs{
			OHLCs: upsert,
		})
		if err != nil {
			logger.Errorf("Error upserting ohlcs for symbol %s: %v", symbol, err)
		}
	}
}
