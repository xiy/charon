// Package httrp is a modified version of the reverse proxy transport written
// for the Go-Kit project: https://github.com/go-kit/kit/blob/master/transport/httprp/server.go
//
// This version has been modified to use a custom http.Transport object.
package httprp

import (
	"charon/proxy/transport"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"
)

// RequestFunc may take information from an HTTP request and put it into a
// request context. BeforeFuncs are executed prior to invoking the
// endpoint.
type RequestFunc func(context.Context, *http.Request) context.Context

// Server is a proxying request handler.
type Server struct {
	ctx          context.Context
	proxy        http.Handler
	before       []RequestFunc
	errorEncoder func(w http.ResponseWriter, err error)
}

// NewServer constructs a new server that implements http.Server and will proxy
// requests to the given base URL using its scheme, host, and base path.
// If the target's path is "/base" and the incoming request was for "/dir",
// the target request will be for /base/dir.
func NewServer(
	ctx context.Context,
	baseURL *url.URL,
	connectionTimeout time.Duration,
	options ...ServerOption,
) *Server {
	s := &Server{
		ctx:   ctx,
		proxy: newCustomReverseProxyHandler(baseURL, connectionTimeout),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// ServerOption sets an optional parameter for servers.
type ServerOption func(*Server)

// ServerBefore functions are executed on the HTTP request object before the
// request is decoded.
func ServerBefore(before ...RequestFunc) ServerOption {
	return func(s *Server) { s.before = before }
}

// ServeHTTP implements http.Handler.
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := s.ctx

	for _, f := range s.before {
		ctx = f(ctx, r)
	}

	s.proxy.ServeHTTP(w, r)
}

func newCustomReverseProxyHandler(baseURL *url.URL, connectionTimeout time.Duration) (proxy *httputil.ReverseProxy) {
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
