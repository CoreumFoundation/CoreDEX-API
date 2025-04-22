package http

import (
	"encoding/json"
	"net/http"

	networklib "github.com/CoreumFoundation/CoreDEX-API/domain/network"
	dmnsymbol "github.com/CoreumFoundation/CoreDEX-API/domain/symbol"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
)

func (s *httpServer) getMarket() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		network, err := networklib.Network(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		// Convert symbol into the two denoms:
		sym, err := dmnsymbol.NewSymbol(symbol)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		// Query the market data:
		marketData, err := s.app.Trade.GetMarket(r.Context(), sym, network)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}
		return json.NewEncoder(w).Encode(marketData)
	}
}
