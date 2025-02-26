package http

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	networklib "github.com/CoreumFoundation/CoreDEX-API/domain/network"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
)

func (s *httpServer) getTickers() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		network, err := networklib.Network(r)
		if err != nil {
			return err
		}
		opt, err := newTickerReadOptions(r)
		if err != nil {
			return err
		}

		opt.Network = network
		return json.NewEncoder(w).Encode(s.app.Ticker.GetTickers(r.Context(), opt))
	}
}

func newTickerReadOptions(r *http.Request) (*dmn.TickerReadOptions, error) {
	var symbols []string

	// Decode the input:
	base64Symbols := r.URL.Query().Get("symbols")
	if base64Symbols != "" {
		decodedSymbols, err := base64.StdEncoding.DecodeString(base64Symbols)
		if err != nil {
			return nil, handler.NewAPIError(422, "symbols.invalid")
		}
		// The decodedSymbols is a json array: Decode the array into symbols:
		if err := json.Unmarshal(decodedSymbols, &symbols); err != nil {
			return nil, handler.NewAPIError(422, "symbols.invalid")
		}
	}

	tickerReadOptions := dmn.NewTickerReadOptions(symbols, time.Now().Truncate(time.Second), 24*time.Hour)
	if err := tickerReadOptions.Validate(); err != nil {
		var name string
		switch err {
		case dmn.ErrSymbolInvalid:
			name = "symbols.invalid"
		case dmn.ErrTickerTooManySymbols:
			name = "symbols.too_many"
		case dmn.ErrTickerEmptySymbols:
			name = "symbols.empty"
		case dmn.ErrTickerPeriodInvalid:
			name = "period.invalid"
		}
		return nil, handler.NewAPIError(422, name)
	}

	return tickerReadOptions, nil
}
