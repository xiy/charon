package proxy

import (
	"charon/proxy/transport/httprp"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"
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
	ctx := context.Background()
	proxy := httprp.NewServer(ctx, url, connectionTimeout)

	service := &Service{
		Name:          name,
		RoutingPrefix: prefix,
		URL:           url,
		Handler:       proxy,
	}

	return service
}
