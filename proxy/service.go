package proxy

import (
	"charon/proxy/transport"
	"charon/proxy/transport/httprp"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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
func NewService(
	name string,
	prefix string,
	url *url.URL,
	connectionTimeout time.Duration,
) (s *Service) {
	ctx := context.Background()
	proxyHandler := newCustomReverseProxyHandler(url, connectionTimeout)
	proxy := httprp.NewServer(ctx, url, proxyHandler)

	service := &Service{
		Name:          name,
		RoutingPrefix: prefix,
		URL:           url,
		Handler:       proxy,
	}

	return service
}

func newCustomReverseProxyHandler(
	baseURL *url.URL,
	connectionTimeout time.Duration,
) (proxy *httputil.ReverseProxy) {
	proxy = httputil.NewSingleHostReverseProxy(baseURL)

	proxy.Transport = transport.NewServiceTransport(connectionTimeout)

	targetQuery := baseURL.RawQuery

	proxy.Director = func(req *http.Request) {
		req.Host = baseURL.Host
		req.URL.Scheme = baseURL.Scheme
		req.URL.Host = baseURL.Host
		req.URL.Path = singleJoiningSlash(baseURL.Path, req.URL.Path)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}

	return
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
