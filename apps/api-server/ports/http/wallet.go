package http

import (
	"encoding/json"
	"net/http"

	networklib "github.com/CoreumFoundation/CoreDEX-API/domain/network"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
)

func (s *httpServer) getAssets() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()
		q := r.URL.Query()
		address := q.Get("address")
		network, err := networklib.Network(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		res, err := s.app.Order.WalletAssets(network, address)
		if err != nil {
			return json.NewEncoder(w).Encode(res)
		}
		return json.NewEncoder(w).Encode(res)
	}
}
