package main

import (
	"charon/proxy"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/husobee/vestigo"
)

// Router is a wrapper around the Vestigo HTTP Muxer that handles route
// creation and management for the gateway.
type Router struct {
	Muxer          *vestigo.Router
	ServiceTimeout time.Duration
	Services       map[string]*proxy.Service
	Logger         *log.Logger
}

// NewRouter creates and returns a new instance of Router with a defined
// timeout for non-responding services.
func NewRouter(serviceTimeout time.Duration, logger *log.Logger) (r *Router) {

	return &Router{
		Muxer:          vestigo.NewRouter(),
		ServiceTimeout: serviceTimeout,
		Services:       make(map[string]*proxy.Service),
		Logger:         logger,
	}
}

// ServeHTTP delegates calls to the underlying muxer of the Router instance.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Logger.Printf("info: method=%s path=%s user_agent=%s", req.Method, req.URL.Path+req.URL.RawQuery, req.UserAgent())
	r.Muxer.ServeHTTP(w, req)
}

// AddService creates a new instance of a Service with a given URL and
// routing prefix and attaches it to the router.
func (r *Router) AddService(name, uri, prefix string) *proxy.Service {
	serviceURL, err := url.Parse(uri)
	if err != nil {
		log.Fatalf("error: couldn't parse service addres for service `%s`", name)
	}

	service := proxy.NewService(name, prefix, serviceURL, r.ServiceTimeout)

	// Create a path-prefix route in the muxer, handled by a
	// reverse proxy handler. This means any request to the
	// path-prefix of '/service/*' will result in it being proxied
	// to the backend service as `/service/*`
	r.Muxer.Handle(prefix, service.Handler)
	r.Services[name] = service

	return service
}
