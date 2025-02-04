package http

import (
	"encoding/json"
	"net/http"

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
			return err
		}
		return json.NewEncoder(w).Encode(currencies)
	}
}
