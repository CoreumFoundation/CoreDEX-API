package domain

import (
	"errors"
	"time"

	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
)

const (
	MaxTickerSymbolsNumber = 40
	DefaultTickerPeriod    = 24 * time.Hour
	QUOTE_ASSET            = "USD"
	QUOTE_PRECISION        = 0
)

var (
	ErrTickerTooManySymbols Error = errors.New("too many symbols")
	ErrTickerEmptySymbols         = errors.New("empty tickers symbols")
	ErrTickerPeriodInvalid        = errors.New("invalid tickers period")
)

type TickerPoint struct {
	OpenTime   int64
	CloseTime  int64
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	LastPrice  float64 // Actual first trade in the time window
	FirstPrice float64 // Actual last trade in the time window
	// Based on the order of the currencies in the symbol the volume and invertedVolume are calculated.
	Volume         float64
	InvertedVolume float64
	Inverted       bool // Indicates if the original symbol was inverted
}

type Tickers map[string]*TickerPoint

// Currently only 24h tickers are used.
type TickerReadOptions struct {
	Symbols []string
	To      time.Time
	Period  time.Duration
	Network metadata.Network
}

type USDTicker struct {
	Tickers    *Tickers
	USDTickers *Tickers
}

func NewTickerReadOptions(symbols []string, to time.Time, period time.Duration) *TickerReadOptions {
	return &TickerReadOptions{
		Symbols: uniqueSymbols(symbols),
		To:      to,
		Period:  period,
	}
}

func (opt *TickerReadOptions) Validate() Error {
	if len(opt.Symbols) == 0 {
		return ErrTickerEmptySymbols
	} else if len(opt.Symbols) > MaxTickerSymbolsNumber {
		return ErrTickerTooManySymbols
	}

	for _, symb := range opt.Symbols {
		if !ValidSymbol(symb) {
			return ErrSymbolInvalid
		}
	}
	return nil
}

func (opt *TickerReadOptions) From() time.Time {
	return opt.To.Add(-opt.Period)
}

func uniqueSymbols(symbols []string) []string {
	keys := make(map[string]bool, len(symbols))
	var res []string

	for _, symbol := range symbols {
		if _, value := keys[symbol]; !value {
			keys[symbol] = true
			res = append(res, symbol)
		}
	}
	return res
}
