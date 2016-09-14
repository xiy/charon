package proxy

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

type serviceTransport struct {
	transport *http.Transport
	logger    *log.Logger
}

func newServiceTransport(connectionTimeout time.Duration) (t *serviceTransport) {
	t = &serviceTransport{
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

func (st *serviceTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
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
