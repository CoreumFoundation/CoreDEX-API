package http

import (
	"context"
	"io"
	nethttp "net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/cors"

	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type routerOption func(r *mux.Router) *mux.Router

type allowedMethod struct {
	method map[string]bool
}

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

/*
Router contains all the usable routers for the application
The routers are pre-configured by providing the configuration yaml

router: Basic router, not limited of auth applied in any way
AuthorizedRouter: Router that requires authentication
UnauthorizedRouter: Router that does not require authentication
AdminRouter: Router that requires admin authentication
*/
type Server struct {
	router    *mux.Router
	AppRouter *mux.Router
	Router    Router
	server    nethttp.Server
}

// To be provided implementations by the application
type Router interface {
	Health() nethttp.Handler
}

type AuthInterface interface {
	Auth(h nethttp.Handler) handler.Handler
}

// Route description. Used to add route to the http server
type Route struct {
	Path    string
	Handler nethttp.Handler
	Method  HttpMethod
}

var (
	conf           *httpConfig
	HTTPServer     *Server
	allowedMethods = &allowedMethod{
		method: make(map[string]bool),
	}
)

// init starts the related infrastructure components (e.g. the always present base setup)
func init() {
	parseConfig()
	HTTPServer = &Server{
		router: mux.NewRouter(),
	}
}

// Specific health check setup for the application
// Typical health checks are configure as an http endpoint "/api/v1/healthz"
// and are of the method GET.
func InitHealth(r Route) {
	allowedMethods.method[string(r.Method)] = true
	HTTPServer.router.Handle(r.Path, r.Handler).Methods(string(r.Method))
}

func InitRoutes(r []Route) {
	HTTPServer.configureAPIRouterBase()

	for _, route := range r {
		// Make certain CORS has the correct settings for the methods:
		allowedMethods.method[string(route.Method)] = true
		// Add the routes:
		HTTPServer.router.Path(route.Path).Methods(string(route.Method)).Handler(route.Handler)
		logger.Infof("Added route: %s %s", route.Method, route.Path)
	}
}

func (s *Server) configureAPIRouterBase() error {
	s.AppRouter = s.setupAPIRouter(func(r *mux.Router) *mux.Router {
		return r
	})
	return nil
}

func (s *Server) setupAPIRouter(options ...routerOption) *mux.Router {
	apiRouter := s.router.PathPrefix("/").Subrouter()
	for _, opt := range options {
		opt(s.router)
	}
	return apiRouter
}

func logFormatter(_ io.Writer, pp handlers.LogFormatterParams) {
	logger.WithFields(logger.Fields{
		"method":        pp.Request.Method,
		"path":          pp.URL.Path,
		"query":         pp.URL.Query(),
		"status":        pp.StatusCode,
		"size":          pp.Size,
		"addr":          pp.Request.RemoteAddr,
		"sg-session-id": pp.Request.Header.Get("sg-session-id"),
	})
}

func (s *Server) setupHTTPServer() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: conf.CORS.AllowedOrigins,
		AllowedMethods: allowedMethods.toArray(),
		AllowedHeaders: []string{"*"},
		Debug:          false,
	}).Handler(s.router)
	s.server = nethttp.Server{
		Addr:           conf.Port,
		Handler:        handlers.CustomLoggingHandler(os.Stdout, handlers.ProxyHeaders(corsHandler), logFormatter),
		ReadTimeout:    conf.TimeOut.Read.Duration,
		WriteTimeout:   conf.TimeOut.Write.Duration,
		IdleTimeout:    conf.TimeOut.Idle.Duration,
		MaxHeaderBytes: 1 << 20,
	}
}

func (s *Server) Start(ctx context.Context) error {
	logger.Infof("Starting HTTP server on port %s", conf.Port)
	// Finalize setup of the http server
	HTTPServer.setupHTTPServer()

	errs := make(chan error)
	go func() {
		errs <- s.server.ListenAndServe()
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		ctxShutDown, cancel := context.WithTimeout(context.Background(), conf.TimeOut.Shutdown.Duration)
		defer cancel()

		err := s.server.Shutdown(ctxShutDown)
		if err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			return err
		}

		return nil
	}
}

func (s *allowedMethod) toArray() []string {
	t := make([]string, len(s.method))
	for v := range s.method {
		t = append(t, v)
	}
	return t
}
