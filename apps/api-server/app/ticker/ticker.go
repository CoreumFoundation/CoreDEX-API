package ticker

import (
	"context"
	"fmt"
	"sync"
	"time"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	ohlcgrpclient "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc/client"
	"github.com/CoreumFoundation/CoreDEX-API/domain/rates"
	tradesclient "github.com/CoreumFoundation/CoreDEX-API/domain/trade/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Application struct {
	client      ohlcgrpc.OHLCServiceClient
	rates       *rates.Fetchers
	rateCache   *cache
	tickerCache *cache
}

type cache struct {
	mutex *sync.RWMutex
	data  map[string]*dmn.LockableCache
}

func NewApplication() *Application {
	ohclClient := ohlcgrpclient.Client()
	rf := rates.NewFetcher(tradesclient.Client(), ohclClient)
	app := &Application{
		client: ohclClient,
		rates:  rf,
		rateCache: &cache{
			mutex: &sync.RWMutex{},
			data:  make(map[string]*dmn.LockableCache),
		},
		tickerCache: &cache{
			mutex: &sync.RWMutex{},
			data:  make(map[string]*dmn.LockableCache),
		},
	}
	go dmn.CleanCache(app.rateCache.data, app.rateCache.mutex, 60*time.Minute)
	go dmn.CleanCache(app.tickerCache.data, app.tickerCache.mutex, 30*time.Minute)
	return app
}

func (s *Application) GetTickers(ctx context.Context, opt *dmn.TickerReadOptions) *dmn.USDTicker {
	retvals := s.getTickers(ctx, opt)
	tickers := tickersToHTTP(retvals, opt)
	usdRetvals := s.GetUSDRates(ctx, opt)
	usdTickers := tickersToUSD(tickers, usdRetvals)
	return &dmn.USDTicker{
		Tickers:    tickers,
		USDTickers: usdTickers,
	}
}

/*
To convert the tickers to USD, we need to calculate the USD value. For that purpose we have a package
fs-utils-lib/go/rates.

Prices calculated uses the close price, we have to standardize the prices to this close price.
*/
func tickersToUSD(tickers *dmn.Tickers, usdRates map[string]float64) *dmn.Tickers {
	// Create new dmn.Tickers object (the input is a set of pointers, which we do not want to modify)
	m := make(map[string]*dmn.TickerPoint)
	// Loop over the tickers and replace the price values with the USD values.
	for symbol, t := range *tickers {
		ticker := *t
		if usdRate, ok := usdRates[symbol]; ok {
			ticker.OpenPrice = (ticker.OpenPrice / ticker.LastPrice) * usdRate
			ticker.FirstPrice = (ticker.FirstPrice / ticker.LastPrice) * usdRate
			ticker.HighPrice = (ticker.HighPrice / ticker.LastPrice) * usdRate
			ticker.LowPrice = (ticker.LowPrice / ticker.LastPrice) * usdRate
			ticker.LastPrice = usdRate
			if usdRate == 0.0 {
				ticker.OpenPrice = 0.0
				ticker.LastPrice = 0.0
				ticker.FirstPrice = 0.0
				ticker.HighPrice = 0.0
				ticker.LowPrice = 0.0
			}
		}
		m[symbol] = &ticker
	}
	return (*dmn.Tickers)(&m)
}

// Tickers to http evaluates the symbols to be either non inverted or inverted and switches the volume and invertedVolume accordingly
func tickersToHTTP(tickers *dmn.Tickers, opt *dmn.TickerReadOptions) *dmn.Tickers {
	retvals := make(dmn.Tickers)
	for _, symbol := range opt.Symbols {
		if _, ok := (*tickers)[symbol]; ok {
			retvals[symbol] = (*tickers)[symbol]
			continue
		}
	}

	return &retvals
}

// Tickers are once calculated valid for up to refreshInterval seconds, however requests might come in in parallel and the cache might be empty, leading to multiple retrieval requests and subsequent caching of the same data.
// We do however do not want to serialize the requests, since that could be slow or blocking.
// The alternative used here, is that the actual cache is 15 seconds, and synchronized with the clock at refreshInterval to allow refreshes.
// As long as there is demand for the data, somewhere in the 5 second interval beyond the wanted caching period, the data will be refreshed by placing a request a refresh channel.
// This refresh channel is connect to blocking go routines per request type (symbol or group of symbols).
//
// Note: Since multiple instances of the service are running, the refresh is not 100% deduplicated: Another instances might attempt to refresh the cache at about the same time.
func (s *Application) getTickers(ctx context.Context, opt *dmn.TickerReadOptions) *dmn.Tickers {
	// Retrieve the tickers from the OHLC service:
	tickerPoints := make(map[string]*dmn.TickerPoint)

	ohlcs := make([]*ohlcgrpc.OHLCs, 0, len(opt.Symbols))
	var mutex sync.Mutex
	var wg sync.WaitGroup
	for _, symbol := range opt.Symbols {
		// Cache check for the symbol:
		s.tickerCache.mutex.RLock()
		if cache, ok := s.tickerCache.data[symbol]; ok {
			v := cache.Value.(*dmn.TickerPoint)
			tickerPoints[symbol] = v
			s.tickerCache.mutex.RUnlock()
			continue
		}
		s.tickerCache.mutex.RUnlock()
		wg.Add(1)
		go func(symbol string) {
			baseOHLCS, err := s.getOHLC(ctx, symbol, opt)
			if err != nil {
				logger.Errorf("(no cache) Error getting ohlc data for %s: %s", symbol, err.Error())
				wg.Done()
				return
			}
			mutex.Lock()
			ohlcs = append(ohlcs, baseOHLCS)
			mutex.Unlock()
			wg.Done()
		}(symbol)
	}
	wg.Wait()
	tickersResp := s.ohlcsToTickers(ohlcs, opt, tickerPoints)
	return (*dmn.Tickers)(&tickersResp)
}

// Retrieve the OHLC data from the source
func (s *Application) getOHLC(ctx context.Context, symbol string, opt *dmn.TickerReadOptions) (*ohlcgrpc.OHLCs, error) {
	// Get the OHLC data from the source:
	// Temporary until we know what load the cache can handle:
	// 10% of the traffic goes to allowCache=true based on the first character of the symbol (so hypothetically always the same symbols and with that tickers will use the cache).
	loadSymbol := &ohlcgrpc.OHLCFilter{
		Symbol:     symbol,
		Network:    opt.Network,
		Period:     &ohlcgrpc.Period{PeriodType: ohlcgrpc.PeriodType_PERIOD_TYPE_HOUR, Duration: 1},
		To:         timestamppb.New(time.Unix(0, opt.To.UnixNano())),
		From:       timestamppb.New(time.Unix(0, opt.To.Add(-opt.Period).UnixNano())),
		Backfill:   true,
		AllowCache: true,
	}
	baseOHLCS, err := s.client.Get(ohlcgrpclient.AuthCtx(ctx), loadSymbol)
	if err != nil {
		logger.Errorf("(source) Error getting ohlc data for %s: %s", symbol, err.Error())
		return nil, err
	}
	// Prevent downstream failures on empty arrays.
	if len(baseOHLCS.OHLCs) == 0 {
		return nil, fmt.Errorf("no ohlc data found for %s", symbol)
	}
	return baseOHLCS, nil
}

// ohlcsToTickers converts the ohlc data to ticker data
// The values calculated are cached for refreshInterval max (clock rounded to refreshInterval intervals) with an allowed stale period of 5s, giving an always retrieval of values
// from the cache for performance and data cost reasons (we can refresh in a go blocking routine while the data is still being served quickly)
func (s *Application) ohlcsToTickers(ohlcs []*ohlcgrpc.OHLCs, domainOptions *dmn.TickerReadOptions,
	tickerPoints map[string]*dmn.TickerPoint) map[string]*dmn.TickerPoint {
	for _, ohlc := range ohlcs {
		tickerPoint := calculateTickerOHLC(ohlc, domainOptions)
		tickerPoints[ohlc.OHLCs[0].Symbol] = tickerPoint
		s.tickerCache.mutex.Lock()
		s.tickerCache.data[ohlc.OHLCs[0].Symbol] = &dmn.LockableCache{
			Value:       tickerPoint,
			LastUpdated: time.Now(),
		}
		s.tickerCache.mutex.Unlock()
	}
	return tickerPoints
}

// calculate the open, high, low, close, volume and invertedVolume
// Returns a single OHLC with the calculated values.
func calculateTickerOHLC(ohlcs *ohlcgrpc.OHLCs, domainOptions *dmn.TickerReadOptions) *dmn.TickerPoint {
	fromTime := domainOptions.To.Add(-domainOptions.Period)
	// get the base ohlcs for the requested period calculation.
	// The input data contains the base data to calculate the requested period over.
	// Calculate the volume:
	var open, close, low, high, volume, invertedVolume, firstPrice float64
	// Assumption is that the data might not be ordered by time.
	var tStart, tEnd time.Time
	for _, baseOHLCS := range ohlcs.OHLCs {
		if low == 0.0 || low > baseOHLCS.Low {
			low = baseOHLCS.Low
		}
		if high < baseOHLCS.High {
			high = baseOHLCS.High
		}
		volume += baseOHLCS.Volume

		// calculate the open:
		if tStart.IsZero() || tStart.After(baseOHLCS.Timestamp.AsTime()) {
			open = baseOHLCS.Open
			tStart = baseOHLCS.Timestamp.AsTime()
		}
		// Calculate the first price: This is the first price in the time period (so timestamp > From)
		if tStart.After(fromTime) && (firstPrice == 0.0 || firstPrice > baseOHLCS.Open) {
			firstPrice = baseOHLCS.Open
		}
		// calculate the close:
		if tEnd.IsZero() || tEnd.Before(baseOHLCS.Timestamp.AsTime()) {
			close = baseOHLCS.Close
			tEnd = baseOHLCS.Timestamp.AsTime()
		}
	}
	// The calculated values might cover the requested time period or might be from before the requested time period:
	// If they are from before the requested time period, volume is 0, high, low and close are the same as the open.
	if tStart.Before(fromTime) {
		volume = 0
		high = open
		low = open
		close = open
		firstPrice = 0.0
	}

	t := &dmn.TickerPoint{
		OpenTime:       fromTime.Unix(),
		CloseTime:      domainOptions.To.Unix(),
		OpenPrice:      open,
		HighPrice:      high,
		LowPrice:       low,
		LastPrice:      close,
		FirstPrice:     firstPrice,
		Volume:         volume,
		InvertedVolume: invertedVolume,
	}
	return t
}

func key(symbol string, network metadata.Network) string {
	return fmt.Sprintf("%s:%d", symbol, network)
}

// Returns the standardized USD result for a single token.
func (s *Application) getRate(ctx context.Context, symbol string, network metadata.Network) (float64, error) {
	s.rateCache.mutex.RLock()
	if cache, ok := s.rateCache.data[key(symbol, network)]; ok {
		v := cache.Value.(float64)
		s.rateCache.mutex.RUnlock()
		return v, nil
	}
	s.rateCache.mutex.RUnlock()
	f := *s.rates
	fetcher := f[network]
	// Get the USD rate:
	denom1, err := denom.NewDenom(symbol)
	if err != nil {
		return 0.0, err
	}
	usd, err := fetcher.ParseExchangeRate(ctx, denom1, dmn.QUOTE_ASSET)
	if err != nil || usd == 0.0 {
		return 0.0, err
	}
	// Cache the rate:
	s.rateCache.mutex.Lock()
	s.rateCache.data[key(symbol, network)] = &dmn.LockableCache{
		Value:       usd,
		LastUpdated: time.Now(),
	}
	s.rateCache.mutex.Unlock()
	return usd, nil
}

func (s *Application) GetUSDRates(ctx context.Context, opt *dmn.TickerReadOptions) map[string]float64 {
	// Key is the symbol.String in the input order (e.g. no inversion of the symbol required)
	rates := make(map[string]float64)
	var mutex sync.Mutex
	var wg sync.WaitGroup
	for _, symbol := range opt.Symbols {
		wg.Add(1)
		go func(symbol string) {
			usd, err := s.getRate(ctx, symbol, opt.Network)
			if err != nil {
				logger.Errorf("Error getting rate for %s: %s", symbol, err.Error())
				wg.Done()
				return
			}
			mutex.Lock()
			rates[symbol] = usd
			mutex.Unlock()
			wg.Done()
		}(symbol)
	}
	wg.Wait()
	return rates
}
