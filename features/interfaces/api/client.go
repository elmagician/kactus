package api

import (
	"net/http"

	"github.com/elmagician/kactus/internal/api"
	internalPicker "github.com/elmagician/kactus/internal/picker"
)

// Client provides methods to test HTTP Rest API endpoints.
type Client struct {
	store *internalPicker.Store
	cli   *api.Client

	request api.RequestPreparation

	autoResetRequest bool
	resetAutoRequest bool
}

// New initializes an HTTP API tester.
func New(store *internalPicker.Store, autoReset bool) (*Client, error) {
	cli, err := api.NewClient(http.DefaultClient)
	if err != nil {
		return nil, err
	}

	initStatusCode(store)

	return &Client{
		store:            store,
		cli:              cli,
		autoResetRequest: autoReset,
		resetAutoRequest: autoReset,
	}, nil
}

// ResetRequest resets client request.
// If cli.autoResetRequest is enabled,
// request will automatically be cleared
// after response parsing.
func (cli *Client) ResetRequest() {
	cli.request = api.RequestPreparation{}
}

// Reset resets client instance.
// Client will be reset to default.
// Response/request will be forgotten.
func (cli *Client) Reset() {
	cli.cli.Reset()
	cli.ResetRequest()
	cli.autoResetRequest = cli.resetAutoRequest
}

func (cli *Client) DisableAutoResetRequest() {
	cli.autoResetRequest = false
}

func (cli *Client) EnableAutoResetRequest() {
	cli.autoResetRequest = true
}

// Trace activates client request tracing.
func (cli *Client) Trace() {
	cli.cli.SetTrace(true)
}

// DisableTrace disables client request tracing.
func (cli *Client) DisableTrace() {
	cli.cli.SetTrace(false)
}

// FollowRedirect enables redirection.
func (cli *Client) FollowRedirect() {
	cli.cli.SetFollowRedirection(true)
}

// DisableRedirect disables redirection.
func (cli *Client) DisableRedirect() {
	cli.cli.SetFollowRedirection(false)
}

// ExecuteRequest builds and executes request through http client.
func (cli *Client) ExecuteRequest() error {
	if err := cli.cli.EmitRequest(cli.request); err != nil {
		return err
	}

	if cli.autoResetRequest {
		cli.ResetRequest()
	}

	return nil
}

// nolint: lll
func initStatusCode(store *internalPicker.Store) {
	saveStatus := func(key string, code int) {
		store.Pick("status."+key, code, internalPicker.PersistentValue)
	}

	// 100
	saveStatus("Continue", http.StatusContinue)
	saveStatus("SwitchingProtocols", http.StatusSwitchingProtocols)
	saveStatus("Processing", http.StatusProcessing)
	saveStatus("EarlyHints", http.StatusEarlyHints)

	// 200
	saveStatus("Ok", http.StatusOK)
	saveStatus("Created", http.StatusCreated)
	saveStatus("Accepted", http.StatusAccepted)
	saveStatus("NonAuthoritativeInfo", http.StatusNonAuthoritativeInfo)
	saveStatus("NoContent", http.StatusNoContent)
	saveStatus("ResetContent", http.StatusResetContent)
	saveStatus("PartialContent", http.StatusPartialContent)
	saveStatus("MultiStatus", http.StatusMultiStatus)
	saveStatus("AlreadyReported", http.StatusAlreadyReported)
	saveStatus("IMUsed", http.StatusIMUsed)

	// 300
	saveStatus("MultipleChoices", http.StatusMultipleChoices)
	saveStatus("MovedPermanently", http.StatusMovedPermanently)
	saveStatus("Found", http.StatusFound)
	saveStatus("SeeOther", http.StatusSeeOther)
	saveStatus("NotModified", http.StatusNotModified)
	saveStatus("UseProxy", http.StatusUseProxy)
	saveStatus("TemporaryRedirect", http.StatusTemporaryRedirect)
	saveStatus("PermanentRedirect", http.StatusPermanentRedirect)

	// 400
	saveStatus("BadRequest", http.StatusBadRequest)
	saveStatus("Unauthorized", http.StatusUnauthorized)
	saveStatus("PaymentRequired", http.StatusPaymentRequired)
	saveStatus("Forbidden", http.StatusForbidden)
	saveStatus("NotFound", http.StatusNotFound)
	saveStatus("MethodNotAllowed", http.StatusMethodNotAllowed)
	saveStatus("NotAcceptable", http.StatusNotAcceptable)
	saveStatus("ProxyAuthRequired", http.StatusProxyAuthRequired)
	saveStatus("RequestTimeout", http.StatusRequestTimeout)
	saveStatus("Conflict", http.StatusConflict)
	saveStatus("Gone", http.StatusGone)
	saveStatus("LengthRequired", http.StatusLengthRequired)
	saveStatus("PreconditionFailed", http.StatusPreconditionFailed)
	saveStatus("RequestEntityTooLarge", http.StatusRequestEntityTooLarge)
	saveStatus("RequestURITooLong", http.StatusRequestURITooLong)
	saveStatus("UnsupportedMediaType", http.StatusUnsupportedMediaType)
	saveStatus("RequestedRangeNotSatisfiable", http.StatusRequestedRangeNotSatisfiable)
	saveStatus("ExpectationFailed", http.StatusExpectationFailed)
	saveStatus("Teapot", http.StatusTeapot)
	saveStatus("MisdirectedRequest", http.StatusMisdirectedRequest)
	saveStatus("UnprocessableEntity", http.StatusUnprocessableEntity)
	saveStatus("Locked", http.StatusLocked)
	saveStatus("FailedDependency", http.StatusFailedDependency)
	saveStatus("TooEarly", http.StatusTooEarly)
	saveStatus("UpgradeRequired", http.StatusUpgradeRequired)
	saveStatus("PreconditionRequired", http.StatusPreconditionRequired)
	saveStatus("TooManyRequests", http.StatusTooManyRequests)
	saveStatus("RequestHeaderFieldsTooLarge", http.StatusRequestHeaderFieldsTooLarge)
	saveStatus("UnavailableForLegalReasons", http.StatusUnavailableForLegalReasons)

	// 500
	saveStatus("InternalServerError", http.StatusInternalServerError)
	saveStatus("NotImplemented", http.StatusNotImplemented)
	saveStatus("BadGateway", http.StatusBadGateway)
	saveStatus("ServiceUnavailable", http.StatusServiceUnavailable)
	saveStatus("GatewayTimeout", http.StatusGatewayTimeout)
	saveStatus("HTTPVersionNotSupported", http.StatusHTTPVersionNotSupported)
	saveStatus("VariantAlsoNegotiates", http.StatusVariantAlsoNegotiates)
	saveStatus("InsufficientStorage", http.StatusInsufficientStorage)
	saveStatus("LoopDetected", http.StatusLoopDetected)
	saveStatus("NotExtended", http.StatusNotExtended)
	saveStatus("NetworkAuthenticationRequired", http.StatusNetworkAuthenticationRequired)
}
