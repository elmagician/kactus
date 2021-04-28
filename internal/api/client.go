package api

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"

	"go.uber.org/zap"
)

var ErrNoRequest = errors.New("trying to emit empty request")

type Client struct {
	client        *http.Client
	trace         *httptrace.ClientTrace
	initialClient *http.Client

	// Current storage
	request      *http.Request
	httpResponse *http.Response
	Response     *Response
	tracing      bool
}

func NewClient(cli *http.Client) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: nil})
	if err != nil {
		return nil, err
	}

	debug := &debugTransport{}
	trace := &httptrace.ClientTrace{
		GetConn:              nil,
		GotConn:              debug.GotConn,
		PutIdleConn:          nil,
		GotFirstResponseByte: nil,
		Got100Continue:       nil,
		Got1xxResponse:       nil,
		DNSStart:             nil,
		DNSDone:              nil,
		ConnectStart:         nil,
		ConnectDone:          nil,
		TLSHandshakeStart:    nil,
		TLSHandshakeDone:     nil,
		WroteHeaderField:     debug.WroteHeaderField,
		WroteHeaders:         debug.WroteHeaders,
		Wait100Continue:      nil,
		WroteRequest:         nil,
	}

	cli.Transport = debug
	cli.Jar = jar

	defaultCli := *cli

	return &Client{client: cli, initialClient: &defaultCli, trace: trace}, nil
}

func (cli *Client) Reset() {
	cli.request = nil
	cli.httpResponse = nil
	cli.Response = nil

	newCli, err := NewClient(cli.initialClient)
	if err != nil {
		log.Error("could not reset client", zap.Error(err))
	}

	cli.client = newCli.client
}

func (cli *Client) SetTrace(activate bool) {
	cli.tracing = activate
}

func (cli *Client) SetFollowRedirection(follow bool) {
	if follow {
		cli.client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			log.Debug("follow redirect")
			return nil
		}
	} else {
		cli.client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			log.Debug("no redirect")
			return http.ErrUseLastResponse
		}
	}
}

func (cli *Client) EmitRequest(req RequestPreparation) (err error) {
	if req.Empty() {
		return ErrNoRequest
	}

	cli.request, err = req.GenerateRequest(cli.client.Jar)
	if err != nil {
		return err
	}

	if cli.tracing {
		cli.request = cli.request.WithContext(
			httptrace.WithClientTrace(cli.request.Context(), cli.trace),
		)
	}

	// nolint: bodyclose
	cli.httpResponse, err = cli.client.Do(cli.request)
	if err != nil {
		return
	}

	defer func() {
		_, copyErr := io.Copy(ioutil.Discard, cli.httpResponse.Body) // avoid leaking.
		closeErr := cli.httpResponse.Body.Close()

		if copyErr != nil {
			log.Warn("could not discard response body", zap.Error(copyErr))
		}

		if closeErr != nil {
			log.Warn("could not close response body", zap.Error(closeErr))
		}
	}()

	body, errBody := ioutil.ReadAll(cli.httpResponse.Body)

	if errBody != nil {
		return errBody
	}

	cli.Response = NewResponse(cli.httpResponse.StatusCode, body, cli.httpResponse.Cookies(), cli.httpResponse.Header)

	return
}
