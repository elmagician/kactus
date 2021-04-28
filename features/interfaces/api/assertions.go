package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elmagician/godog"

	"github.com/elmagician/kactus/internal/api"
	match "github.com/elmagician/kactus/internal/matchers"
)

var (
	// ErrInvalidStatus is thrown when response status code does not match
	// expected status.
	ErrInvalidStatus = errors.New("status code does not match expected")

	// ErrExpectedEmptyBody is thrown when response expected body to be empty
	// but it was not.
	ErrExpectedEmptyBody = errors.New("response body should be empty")
	// ErrExpectedBody is thrown when response expected body to not be empty
	// but it was.
	ErrExpectedBody = errors.New("response should have a body")

	// ErrUnexpectedCookie is thrown when response or client should not have
	// specified cookies but at least one exists.
	ErrUnexpectedCookie = errors.New("should not have cookie")
	// ErrExpectedCookie is thrown when response should have specified cookies
	// but at least one does not exists.
	ErrExpectedCookie = errors.New("should have cookie")
	// ErrInvalidCookieSecurity is thrown when cookies security does not match
	// expected.
	ErrInvalidCookieSecurity = errors.New("cookie security does not match expected")
	// ErrInvalidCookieHTTPSetting is thrown when cookies HTTP setting does
	// not match expected.
	ErrInvalidCookieHTTPSetting = errors.New("cookie HTTP setting does not match expected")
	// ErrInvalidCookieDomain is thrown when cookie domain does not match expected.
	ErrInvalidCookieDomain = errors.New("cookie domain does not match expected")

	// ErrInvalidArgNumber is thrown when assertions methods receive an
	// unexpected number of arguments on periodic argument.
	ErrInvalidArgNumber = errors.New("invalid quantity for periodic arguments")
)

// Exposes api errors
var (
	// ErrNoBody is thrown when asserting on response expected a body to exists
	// but none exists.
	ErrNoBody = api.ErrNoBody

	// ErrUnknownKey is thrown when provided key does not exists in JSON.
	ErrUnknownKey = api.ErrUnknownKey

	// ErrNoMatch is thrown when assertions fails to match expected value.
	ErrNoMatch = api.ErrNoMatch

	// ErrNotFullyMatch is thrown when actual asserted object contains
	// more values than expected.
	ErrNotFullyMatch = api.ErrNotFullyMatch

	// ErrNoPath is thrown when provided assertion path does not exists
	// in actual object.
	ErrNoPath = api.ErrNoPath

	// ErrNoRequest is thrown when assertion expected a request to exists
	// but none exists.
	ErrNoRequest = api.ErrNoRequest
)

// ResponseHasStatus asserts Response has expected status.
func (cli *Client) ResponseHasStatus(expectedStatus int) error {
	if !cli.cli.Response.HasStatus(expectedStatus) {
		return fmt.Errorf("%w: expected %d - got %d", ErrInvalidStatus, expectedStatus, cli.cli.Response.Status)
	}

	return nil
}

// EmptyResponseBody asserts response body should or
// should not be empty depending on not parameter.
//
// Second argument is not used. It is present for
// interface.AsNot simplification in step definitions..
func (cli *Client) EmptyResponseBody(not bool, _ ...string) error {
	if !not && !cli.cli.Response.HasEmptyBody() {
		return fmt.Errorf("%w: body %v", ErrExpectedEmptyBody, cli.cli.Response.Body)
	}

	if not && !cli.cli.Response.HasEmptyBody() {
		return fmt.Errorf("%w: body %v", ErrExpectedBody, cli.cli.Response.Body)
	}

	return nil
}

// ClientHasCookies asserts client has or does not have
// specified cookie list depending on not parameter.
func (cli *Client) ClientHasCookies(not bool, names ...string) error {
	if len(names) == 0 {
		return fmt.Errorf("%w: expected at least 1 name to be provided", ErrInvalidArgNumber)
	}

	for _, name := range names {
		has := cli.cli.HasCookie(strings.TrimSpace(name))

		if !not && !has {
			return fmt.Errorf("client %w %s", ErrExpectedCookie, name)
		}

		if not && has {
			return fmt.Errorf("client %w %s", ErrUnexpectedCookie, name)
		}
	}

	return nil
}

// ResponseHasCookies asserts response has or does not have
// specified cookie list depending on not parameter.
func (cli *Client) ResponseHasCookies(not bool, names ...string) error {
	if len(names) == 0 {
		return fmt.Errorf("%w: expected at least 1 name to be provided", ErrInvalidArgNumber)
	}

	for _, name := range names {
		has := cli.cli.Response.HasCookie(strings.TrimSpace(name))
		if !not && !has {
			return fmt.Errorf("response %w %s", ErrExpectedCookie, name)
		}

		if not && has {
			return fmt.Errorf("response %w %s", ErrUnexpectedCookie, name)
		}
	}

	return nil
}

// ResponseCookiesShouldOrShouldNotBeSecure asserts response cookies are/aren't secure.
func (cli *Client) ResponseCookiesShouldOrShouldNotBeSecure(not bool, names ...string) error {
	for _, name := range names {
		cookie := cli.cli.Response.GetCookie(strings.TrimSpace(name))

		if cookie == nil {
			return fmt.Errorf("response %w %s", ErrExpectedCookie, name)
		}

		if cookie.Secure && not {
			return fmt.Errorf("%w: %s should not be secure", ErrInvalidCookieSecurity, name)
		}

		if !not && !cookie.Secure {
			return fmt.Errorf("%w: %s should be secure", ErrInvalidCookieSecurity, name)
		}
	}

	return nil
}

// ResponseCookiesShouldOrShouldNotBeHTTPOnly asserts response cookies are/aren't HTTP Only.
func (cli *Client) ResponseCookiesShouldOrShouldNotBeHTTPOnly(not bool, params ...string) error {
	for _, name := range params {
		cookie := cli.cli.Response.GetCookie(strings.TrimSpace(name))

		if cookie == nil {
			return fmt.Errorf("response %w %s", ErrExpectedCookie, name)
		}

		if cookie.HttpOnly && not {
			return fmt.Errorf("%w: %s should not be HTTP only", ErrInvalidCookieHTTPSetting, name)
		}

		if !not && !cookie.HttpOnly {
			return fmt.Errorf("%w: %s should be HTTP only", ErrInvalidCookieHTTPSetting, name)
		}
	}

	return nil
}

// ResponseCookieDomainShouldOrShouldNotMatch  asserts response cookie domain does/doesn't match expected domains.
func (cli *Client) ResponseCookieDomainShouldOrShouldNotMatch(not bool, params ...string) error {
	if len(params) == 0 || len(params) > 2 {
		return fmt.Errorf("%w: expected 2 parameters (cookie name, expected domain)", ErrInvalidArgNumber)
	}

	name := params[0]
	expectedDomain := params[1]

	cookie := cli.cli.Response.GetCookie(name)

	if cookie == nil {
		return fmt.Errorf("response %w %s", ErrExpectedCookie, name)
	}

	if not && cookie.Domain == expectedDomain {
		return fmt.Errorf("%w: %s does match %s", ErrInvalidCookieDomain, cookie.Domain, expectedDomain)
	}

	if !not && cookie.Domain != expectedDomain {
		return fmt.Errorf("%w: %s does not match %s", ErrInvalidCookieDomain, cookie.Domain, expectedDomain)
	}

	return nil
}

// ResponseJSONShouldBeEquivalent asserts response body is a JSON resembling provided JSON.
func (cli *Client) ResponseJSONShouldBeEquivalent(expected *godog.DocString) error {
	return cli.cli.Response.JSONResemble(expected)
}

// ResponseJSONShouldContain asserts response body is a JSON having provided keys.
// If fully is true, it also ensures no key exists besides provided one.
// Keys are provided as Path using `.` separators
// To match insides JSON Array, use Index.
//  test.0.has matches a path inside: {"test": [{"matches": val}}
func (cli *Client) ResponseJSONShouldContain(fully bool, matchPaths *godog.Table) error {
	return cli.cli.Response.JSONContains(fully, matchPaths)
}

// ResponseHTMLShouldBeEquivalent asserts response body is a HTML resembling provided.
func (cli *Client) ResponseHTMLShouldBeEquivalent(body *godog.DocString) error {
	return cli.cli.Response.HTMLResemble(body)
}

// ResponseHTMLShouldContains asserts response HTML body contains provided strings.
func (cli *Client) ResponseHTMLShouldContains(elements *godog.Table) error {
	return cli.cli.Response.HTMLContain(elements)
}

// ResponseHeaderShouldOrShouldNotMatch asserts response header
// match or does not match provided value using provided matcher.
func (cli *Client) ResponseHeaderShouldOrShouldNotMatch(not bool, params ...string) error {
	if len(params) == 0 || len(params) > 3 {
		return fmt.Errorf("%w: expected 3 parameters (header name, matcher, value)", ErrInvalidArgNumber)
	}

	name := params[0]
	matcher := params[1]
	value := params[2]

	header := cli.cli.Response.RetrieveHeader(name)

	return match.Assert(matcher, header, value)
}
