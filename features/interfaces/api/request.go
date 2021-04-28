package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elmagician/godog"

	"github.com/elmagician/kactus/internal/api"
)

var supportedMethod = []string{"GET", "POST", "OPTION", "DELETE", "PUT", "PATCH"}

// ErrInvalidMethod is thrown when an unsupported method is provided.
var ErrInvalidMethod = errors.New("invalid method provided")

// InitRequest starts a new request with default parameter.
func (cli *Client) InitRequest(withCookie bool) {
	cli.request = api.PrepareRequest(withCookie)
}

// SetEndpoint adds Endpoint to request.
func (cli *Client) SetEndpoint(endpoint string) {
	if cli.request.Empty() {
		cli.InitRequest(true)
	}

	cli.request = cli.request.SetEndpoint(endpoint)
}

// SetMethod adds HTTP method to request.
func (cli *Client) SetMethod(method string) error {
	if cli.request.Empty() {
		cli.InitRequest(true)
	}

	method = strings.ToUpper(method)

	if !validateMethod(method) {
		return fmt.Errorf("%w: method %s is not in %v", ErrInvalidMethod, method, supportedMethod)
	}

	cli.request = cli.request.SetMethod(method)

	return nil
}

// EnableCookie enables cookie storage for request.
func (cli *Client) EnableCookie() {
	if cli.request.Empty() {
		cli.InitRequest(true)
	}

	cli.request = cli.request.AllowCookies(true)
}

// DisableCookie disables cookie storage for request.
func (cli *Client) DisableCookie() {
	if cli.request.Empty() {
		cli.InitRequest(true)
	}

	cli.request = cli.request.AllowCookies(false)
}

// ClearBody resets request body.
func (cli *Client) ClearBody() {
	if cli.request.Empty() {
		cli.InitRequest(true)
	}

	cli.request.ResetBody()
}

// SetJSONBody replaces current request body with new JSON body.
func (cli *Client) SetJSONBody(body *godog.DocString) {
	cli.ClearBody()
	cli.request = cli.request.SetJSONBody(body)
}

// SetFormBody replaces current request body with new Form body.
func (cli *Client) SetFormBody(body *godog.Table) {
	cli.ClearBody()
	cli.request = cli.request.SetFORMBody(body)
}

// SetQueryParams replaces query parameters with new ones.
func (cli *Client) SetQueryParams(args *godog.Table) error {
	if cli.request.Empty() {
		cli.InitRequest(true)
	}

	cli.request = cli.request.ResetArguments()
	cli.AddQueryParams(args)

	return nil
}

// AddQueryParams adds query parameters.
func (cli *Client) AddQueryParams(args *godog.Table) {
	for _, row := range args.Rows {
		cli.AddQueryParam(row.Cells[0].Value, row.Cells[1].Value)
	}
}

// AddQueryParam adds a single query parameter.
func (cli *Client) AddQueryParam(key, val string) {
	if cli.request.Empty() {
		cli.InitRequest(true)
	}

	cli.request.AddArgument(key, val)
}

// AddCookie adds a new cookie to requests.
func (cli *Client) AddCookie(key, value string, config *godog.Table) error {
	if cli.request.Empty() {
		cli.InitRequest(true)
	}

	var err error
	cli.request, err = cli.request.AddCookie(key, value, config)

	return err
}

// ResetCookie resets requests cookies.
func (cli *Client) ResetCookie() {
	if cli.request.Empty() {
		cli.InitRequest(true)
	}

	cli.request = cli.request.ResetCookies()
}

func validateMethod(candidate string) bool {
	for _, method := range supportedMethod {
		if candidate == method {
			return true
		}
	}

	return false
}
