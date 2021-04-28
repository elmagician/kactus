package api

import (
	"net/http"
	"net/http/httptrace"

	"go.uber.org/zap"
)

type debugTransport struct {
	current *http.Request
}

func (dt *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	dt.current = req
	return http.DefaultTransport.RoundTrip(req)
}

// GotConn prints whether the connection has been used previously
// for the current request.
func (dt *debugTransport) GotConn(info httptrace.GotConnInfo) {
	log.Debug(
		"Connection reused",
		zap.Any("url", dt.current.URL),
		zap.Bool("reused", info.Reused),
	)
}

// GotConn prints whether the connection has been used previously
// for the current request.
func (dt *debugTransport) WroteHeaders() {
	log.Debug("Headers written")
}

// GotConn prints whether the connection has been used previously
// for the current request.
func (dt *debugTransport) WroteHeaderField(key string, value []string) {
	log.Debug(
		"Header written",
		zap.String("header", key),
		zap.Any("value", value))
}
