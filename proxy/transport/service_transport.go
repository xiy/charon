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
		logger: logging.NewCoLogLogger("charon"),
	}

	t.transport = &http.Transport{
		MaxIdleConnsPerHost:   100000,
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
	reqStart := time.Now()
	st.logger.Printf("info: Proxying request: request=%s", req.URL)
	resp, err = st.transport.RoundTrip(req)

	if err != nil {
		st.logger.Printf("error: error during round trip request: %s (%v)", err, req)
	} else {
		st.logger.Printf("info: [<-] status=%s request_time=%v", resp.Status, time.Since(reqStart))
		setViaHeader(req)
		setUserAgentHeader(req)
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
