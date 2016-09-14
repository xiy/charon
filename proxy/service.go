package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// Service represents a unique backend service behind the gateway.
type Service struct {
	Name          string       `json:"name"`
	RoutingPrefix string       `json:"routing_prefix"`
	URL           *url.URL     `json:"url"`
	Handler       http.Handler `json:"handler"`
}

// NewService creates and returns a new Service instance with the given name,
// routing prefix and URL.
func NewService(name string, prefix string, url *url.URL, connectionTimeout time.Duration) (s *Service) {
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = newServiceTransport(connectionTimeout)

	defaultDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		defaultDirector(req)
		req.Host = url.Host
		req.URL.RawQuery = url.RawQuery
	}

	service := &Service{
		Name:          name,
		RoutingPrefix: prefix,
		URL:           url,
		Handler:       proxy,
	}

	return service
}
