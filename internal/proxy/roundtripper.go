package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/odpf/shield/internal/proxy/hook"

	"github.com/odpf/salt/log"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
)

type h2cTransportWrapper struct {
	// Defining two different RoundTripper
	// - httptransport: for http, h2, h2c
	// - http2transport: h2c, grpc
	// this is because &http2.Transport is not supporting
	// proxy for http & h2
	// Reference: https://sourcegraph.com/github.com/tsenart/vegeta/-/blob/lib/attack.go?L206:6#tab=references

	httpTransport *http.Transport
	grpcTransport *http2.Transport

	log  log.Logger
	hook hook.Service
}

func (t *h2cTransportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	// we need to apply errors if it failed in Director
	if err, ok := req.Context().Value(ctxRequestErrorKey).(error); ok {
		return nil, err
	} else if req.Context().Err() != nil {
		return nil, req.Context().Err()
	}

	req.Header.Del("Accept-Encoding")
	var transport http.RoundTripper = t.httpTransport
	if req.Header.Get("Content-Type") == "application/grpc" {
		transport = t.grpcTransport
	}

	req = req.WithContext(WithoutCancel(req.Context()))

	t.log.Debug("backend", req.URL.Host, "request_forwarded")

	res, err := transport.RoundTrip(req)
	if err != nil {
		return res, err
	}

	t.log.Debug("backend", req.URL.Host, "request_completed", zap.String("status", res.Status))

	return t.hook.ServeHook(res, nil)
}

func NewH2cRoundTripper(log log.Logger, hook hook.Service) http.RoundTripper {
	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 1 * time.Minute,
	}
	return &h2cTransportWrapper{
		httpTransport: &http.Transport{
			DialContext:        dialer.DialContext,
			DisableCompression: true,
		},
		grpcTransport: &http2.Transport{
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
			AllowHTTP:          true,
			DisableCompression: true,
		},
		log:  log,
		hook: hook,
	}
}
