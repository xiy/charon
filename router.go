package main

import (
	"charon/logging"
	"charon/proxy"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/husobee/vestigo"
)

// Router is a wrapper around the Vestigo HTTP Muxer that handles route
// creation and management for the gateway.
type Router struct {
	Muxer          *vestigo.Router
	ServiceTimeout time.Duration
	Services       map[string]*proxy.Service
	Logger         log.Logger
}

// NewRouter creates and returns a new instance of Router with a defined
// timeout for non-responding services.
func NewRouter(serviceTimeout time.Duration, allowTrace bool) (r *Router) {
	vestigo.AllowTrace = allowTrace

	return &Router{
		Muxer:          vestigo.NewRouter(),
		ServiceTimeout: serviceTimeout,
		Services:       make(map[string]*proxy.Service),
		Logger:         logging.NewStdoutLogger(),
	}
}

// ServeHTTP delegates calls to the underlying muxer of the Router instance.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Logger.Log("method", req.Method, "path", req.URL.Path, "user_agent", req.UserAgent())
	r.Muxer.ServeHTTP(w, req)
}

// AddService creates a new instance of a Servicewith a given URL and
// routing prefix and attaches it to the router.
func (r *Router) AddService(name, uri, prefix string) *proxy.Service {
	var (
		serviceURL, _ = url.Parse(uri)
		service       = proxy.NewService(name, prefix, serviceURL, r.ServiceTimeout)
	)

	// Create a path-prefix route in the muxer, handled by a
	// reverse proxy handeler. This means any request to the
	// path-prefix of '/service/*' will result in it being proxied
	// to the backend service as `/service/*`
	r.Muxer.Handle(prefix, service.Handler)
	r.Services[name] = service

	return service
}

// GetServices returns the currently registered Services and their details
func (r *Router) GetServices() {
	for _, v := range r.Services {
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(v)
		// info("Service Definition: Name=%s URL=%s", v.Name, v.URL)
	}
}
