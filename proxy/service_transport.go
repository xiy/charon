package proxy

import (
	"charon/logging"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
)

type serviceTransport struct {
	transport *http.Transport
	logger    log.Logger
}

func newServiceTransport(connectionTimeout time.Duration) (t *serviceTransport) {
	t = &serviceTransport{
		transport: &http.Transport{},
		logger:    logging.NewStdoutLogger(),
	}

	t.transport.Dial = func(network, address string) (net.Conn, error) {
		return net.DialTimeout(network, address, connectionTimeout)
	}

	t.transport.MaxIdleConnsPerHost = 20
	t.transport.ResponseHeaderTimeout = 30 * time.Second

	return
}

func (st *serviceTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = st.transport.RoundTrip(req)
	if err == nil {
		setViaHeader(req)
	}

	headers, err := json.Marshal(resp.Header)
	if err != nil {
		st.logger.Log("error", err)
	}

	st.logger.Log("headers", headers, "status", resp.StatusCode)

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
