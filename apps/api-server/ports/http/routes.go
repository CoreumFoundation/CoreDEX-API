package http

import (
	"net/http"

	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app"
	behttp "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/http"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type httpServer struct {
	app *app.Application
}

// NewHttpServer sets up the routes and returns a startable http server.
func NewHttpServer(app *app.Application) *behttp.Server {
	s := httpServer{app: app}
	behttp.InitHealth(behttp.Route{
		Path: "/healthz", Method: behttp.GET, Handler: s.Health(),
	})
	behttp.InitRoutes([]behttp.Route{
		{Path: "/ohlc", Method: behttp.GET, Handler: s.getOHLC()},
		{Path: "/tickers", Method: behttp.GET, Handler: s.getTickers()},
		{Path: "/trades", Method: behttp.GET, Handler: s.getTrades()},
		{Path: "/currencies", Method: behttp.GET, Handler: s.getCurrencies()},
		{Path: "/order/create", Method: behttp.POST, Handler: s.createOrder()},
		{Path: "/order/cancel", Method: behttp.POST, Handler: s.cancelOrder()},
		{Path: "/order/submit", Method: behttp.POST, Handler: s.submitOrder()},
		{Path: "/order/orderbook", Method: behttp.GET, Handler: s.getOrders()},
		{Path: "/wallet/assets", Method: behttp.GET, Handler: s.getAssets()},
		{Path: "/ws", Method: behttp.GET, Handler: s.wsEndpoint()},
	})
	return behttp.HTTPServer
}

// Health is application the application specific health check. Any interface health can be checked here.
func (s *httpServer) Health() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := s.app.Health(); err != nil {
			logger.Errorf("healthcheck failure. Error %v", err)
			return err
		}
		return nil
	}
}
