package http

import (
	"encoding/json"
	"net/http"

	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	networklib "github.com/CoreumFoundation/CoreDEX-API/domain/network"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
)

func (s *httpServer) getCurrencies() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		network, err := networklib.Network(r)
		if err != nil {
			return handler.NewAPIError(401, "network.invalid")
		}
		currencies, err := s.app.Currency.GetCurrencies(r.Context(), network)
		if err != nil {
			return json.NewEncoder(w).Encode(&currencygrpc.Currencies{})
		}
		return json.NewEncoder(w).Encode(currencies)
	}
}
