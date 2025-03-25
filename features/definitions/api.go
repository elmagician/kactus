package definitions

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/elmagician/kactus/features/interfaces"
	"github.com/elmagician/kactus/features/interfaces/api"
)

func InstallAPI(s *godog.ScenarioContext, client *api.Client) {
	// HEADERS ----------------
	// Set request headers from a godog table. Previous headers will be forgotten.
	s.Step(`^(?:I )?set(?:ting)? request headers:$`, client.SetRequestHeaders)
	// Update existing headers from a godog table.
	s.Step(`^(?:I )?assign(?:ing)? request headers:$`, client.AddHeaders)
	// Set a single header value.
	s.Step(`^(?:I )?set(?:ing)? ([a-zA-Z0-9-]+) request header to (.+)$`, client.SetHeader)

	// BEHAVIORS ----------------
	// Disable follow redirection
	s.Step(`^(?:I )?do not follow redirect$`, client.DisableRedirect)
	// Enable follow redirection
	s.Step(`^(?:I )?follow redirect$`, client.FollowRedirect)
	// Enable cookies
	s.Step(`(?:I )?enabl(?:e|ing) cookie$`, client.EnableCookie)
	// Disable cookies
	s.Step(`(?:I )?disabl(?:e|ing) cookie$`, client.DisableCookie)
	// Reset http client
	s.Step(`(?:I )?reset(?:ing)? http client$`, client.Reset)
	s.Step(`(?:I )?reset(?:ing)? http request$`, client.ResetRequest)
	s.Step(`(?:I )?disable auto reset for request$`, client.DisableAutoResetRequest)
	s.Step(`(?:I )?enable auto reset for request$`, client.EnableAutoResetRequest)

	// REQUEST ---------------------
	s.Step(`(?:I )?execut(?:e|ing) request$`, client.ExecuteRequest)
	// Set up request
	s.Step(
		`^(?:I )?want(?:ing)? to (GET|PUT|POST|DELETE) (.*)$`,
		func(method, endpoint string) error {
			client.SetEndpoint(endpoint)
			return client.SetMethod(method)
		},
	)
	s.Step(
		`^(?:I )?(GET|PUT|POST|DELETE) (.*)$`,
		func(method, endpoint string) error {
			client.SetEndpoint(endpoint)
			if err := client.SetMethod(method); err != nil {
				return err
			}

			// InitRequest mandatory cause we wish to ensure the call is on the correct value
			return client.ExecuteRequest()
		},
	)

	// BODY ---------------------
	// form && json are mutually exclusive. If both are defined, only JSON will be used
	s.Step(`(?:I )?set(?:ing)? request form body:$`, client.SetFormBody)
	s.Step(`(?:I )?set(?:ing)? request json body:$`, client.SetJSONBody)
	s.Step(`(?:I )?clear(?:ing)? request body$`, client.ClearBody)

	// QUERY PARAMS -------------
	s.Step(`(?:I )?set(?:ing)? request query$`, client.SetQueryParams)
	s.Step(`(?:I )?add(?:ing)? request argument ([a-zA-Z0-9-]+) to (.+)`, client.AddQueryParam)

	// COOKIES ------------------
	// Simulate browser/app side cookie
	s.Step(`(?:I )?set(?:ing)? cookie ([a-zA-Z0-9-]+) to (.+)(?: with options:)?$`, client.AddCookie)

	// Picking
	// Pick response header value
	s.Step(`^(?:I )?pick response header ([a-zA-Z1-9_-]+) as ([a-zA-Z0-9]+)$`, client.PickResponseHeader)
	// Pick key from URL Arg
	s.Step(`^(?:I )?pick key ([a-zA-Z1-9_-]+) from url ([^ ]+) as ([a-zA-Z0-9]+)$`, client.PickArgumentFromURLArg)
	// Pick json key as key from json response
	// TODO: allow to pick from Path
	s.Step(`^(?:I )?pick response json ([a-zA-Z1-9_-]+) as ([a-zA-Z0-9]+)$`, client.PickFromResponseJSONBody)
	// Pick first matching tag attribute as key from html document.
	// Add attributes conditions as a data table to make a more precise selection (will always pick the first
	// matching value in html response)
	// Select format: attribute | matcher | value
	s.Step(
		`(?:I )?pick response html value from tag ([a-z]+[1-9]?) attribute ([a-z]+) as ([A-Za-z0-9]+)(?: with attributes conditions:)?`,
		client.PickResponseHTMLTag,
	)
	s.Step(`^(?:I )?pick response cookie ([a-zA-Z1-9_-]+) as ([a-zA-Z0-9]+)$`, client.PickResponseCookie)

	// s.Step(`^(?:I )?set request cookie from ([a-zA-Z0-9]+)$`, client.SetRequestCookie)

	// ASSERTION ------------------

	// Check HTTP status has correct code
	s.Step(`^response status code should be (\d+)$`, client.ResponseHasStatus)
	// Check is response has empty body
	s.Step(`^response body should (not )?be empty$`, interfaces.AsNot2(client.EmptyResponseBody))
	// Check if client has|do not have a cookie X
	s.Step(`^client should (not )?have an? (.+) cookie$`, func(not, name string) error {
		return interfaces.AsNot(client.ClientHasCookies)(not, name)
	})
	// Check if response has|do not have a cookie X
	s.Step(`^response should (not )?have an? (.+) cookie$`, func(not, name string) error {
		return interfaces.AsNot(client.ResponseHasCookies)(not, name)
	})
	// Check if response cookie X is|is not secure
	s.Step(`^response cookie (.+) should (not )?be secure$`, func(name, not string) error {
		return interfaces.AsNot(client.ResponseCookiesShouldOrShouldNotBeSecure)(not, name)
	})
	// Check if response cookie X is|is http only
	s.Step(`^response cookie (.+) should (not )?be http only$`, func(name, not string) error {
		return interfaces.AsNot(client.ResponseCookiesShouldOrShouldNotBeHTTPOnly)(not, name)
	})
	// Check if response cookie X domain to equal|not equal provided domain
	s.Step(`^response cookie (.+) domain should (not )?be (.+)$`, func(name, not, domain string) error {
		return interfaces.AsNot(client.ResponseCookieDomainShouldOrShouldNotMatch)(not, name, domain)
	})

	// Check if json response object equal provided json (pass as gherkin.DocString)
	s.Step(`^json response should resemble$`, client.ResponseJSONShouldBeEquivalent)
	// Check if json response object contain key/val (pass as gherkin.DataTable). Match only first level key
	// Key work not null check if key exist and contain data
	// If defined as fully contain, all key from the json has to be consumed
	s.Step(
		`^json response should (fully )?contain:$`,
		func(fully string, matchPaths *godog.Table) error {
			return client.ResponseJSONShouldContain(fully != "", matchPaths)
		},
	)

	// Try to match html body with provided html code (as gherkin.DocString)
	s.Step(`^html response should resemble:$`, client.ResponseHTMLShouldBeEquivalent)
	// Look into html body to see if contain substrings (as gherkin.DataTable using a single column)
	s.Step(`^html response should contain:$`, client.ResponseHTMLShouldContains)

	// Check response header X equal|contain|match value Y
	s.Step(
		`^response header (.+) should (not )?(equal|contain|match) (.+)$`,
		func(name, not, matcher, value string) error {
			return interfaces.AsNot(client.ResponseHeaderShouldOrShouldNotMatch)(not, name, matcher, value)
		},
	)

	// OTHERS ------------------
	// Allow trace debug on client.
	s.Step(`^trace client$`, client.Trace)
	// Stop trace debug on client.
	s.Step(`^stop trace client$`, client.DisableTrace)

	s.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		client.Reset()
		return ctx, nil
	})
}
