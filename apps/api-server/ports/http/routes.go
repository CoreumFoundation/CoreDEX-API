package http

import (
	"net/http"

	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app"
	behttp "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/http"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

const routePrepend = "/api"

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
		{Path: routePrepend + "/ohlc", Method: behttp.GET, Handler: s.getOHLC()},
		{Path: routePrepend + "/tickers", Method: behttp.GET, Handler: s.getTickers()},
		{Path: routePrepend + "/trades", Method: behttp.GET, Handler: s.getTrades()},
		{Path: routePrepend + "/currencies", Method: behttp.GET, Handler: s.getCurrencies()},
		{Path: routePrepend + "/order/create", Method: behttp.POST, Handler: s.createOrder()},
		{Path: routePrepend + "/order/cancel", Method: behttp.POST, Handler: s.cancelOrder()},
		{Path: routePrepend + "/order/submit", Method: behttp.POST, Handler: s.submitOrder()},
		{Path: routePrepend + "/order/orderbook", Method: behttp.GET, Handler: s.getOrders()},
		{Path: routePrepend + "/wallet/assets", Method: behttp.GET, Handler: s.getAssets()},
		{Path: routePrepend + "/ws", Method: behttp.GET, Handler: s.wsEndpoint()},
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
