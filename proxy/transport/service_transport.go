package transport

import (
	"charon/logging"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// ServiceTransport wraps an http.Transport, also attached a logger.
type ServiceTransport struct {
	transport *http.Transport
	logger    *log.Logger
}

// NewServiceTransport creates a new custom HTTP transport for use with an HTTP handler.
func NewServiceTransport(connectionTimeout time.Duration) (t *ServiceTransport) {
	t = &ServiceTransport{
		logger: logging.NewCoLogLogger(),
	}

	t.transport = &http.Transport{
		DisableKeepAlives:     true,
		MaxIdleConnsPerHost:   100000,
		DisableCompression:    true,
		ResponseHeaderTimeout: 30 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}

	return
}

// RoundTrip wraps the underlying transport.RoundTrip function, adding logging and custom headers.
func (st *ServiceTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = st.transport.RoundTrip(req)

	log.Printf("info: %s", req.URL.Query().Encode())

	if err == nil {
		setViaHeader(req)
	}

	return
}

func setUserAgentHeader(req *http.Request) {
	if _, present := req.Header["User-Agent"]; !present {
		req.Header.Set("User-Agent", "")
	}
}

func setViaHeader(req *http.Request) {
	via := fmt.Sprintf("%d.%d", req.ProtoMajor, req.ProtoMinor) + " charon"

	if prior, ok := req.Header["Via"]; ok {
		via = strings.Join(prior, ", ") + ", " + via
	}

	req.Header.Set("Via", via)
}

func newErrorResponse(status int) (resp *http.Response) {
	resp = &http.Response{StatusCode: status}
	resp.Body = ioutil.NopCloser(strings.NewReader(""))
	return
}
