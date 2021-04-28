package api

import (
	"net/url"

	"github.com/elmagician/godog"

	internalPicker "github.com/elmagician/kactus/internal/picker"
)

// PickFromResponseJSONBody picks paths value from a response JSON body.
func (cli Client) PickFromResponseJSONBody(path, pickAs string) error {
	value, err := cli.cli.Response.RetrieveJSON(path)
	if err != nil {
		return err
	}

	cli.store.Pick(pickAs, value, internalPicker.DisposableValue)

	return nil
}

// PickResponseHTMLTag picks tag value from a response HTML body.
func (cli Client) PickResponseHTMLTag(tag, attribute, pickAs string, filters *godog.Table) error {
	value, err := cli.cli.Response.RetrieveHTMLAttribute(tag, attribute, filters)
	if err != nil {
		return err
	}

	cli.store.Pick(pickAs, value, internalPicker.DisposableValue)

	return nil
}

// PickResponseCookie picks cookie from response.
func (cli Client) PickResponseCookie(name, pickAs string) {
	cli.store.Pick(pickAs, cli.cli.Response.GetCookie(name), internalPicker.DisposableValue)
}

// PickResponseHeader picks header from response.
func (cli Client) PickResponseHeader(name, pickAs string) {
	cli.store.Pick(pickAs, cli.cli.Response.RetrieveHeader(name), internalPicker.DisposableValue)
}

// PickArgumentFromURLArg picks header from response.
func (cli Client) PickArgumentFromURLArg(argument, urlCandidate, pickAs string) error {
	parsedURL, err := url.Parse(urlCandidate)
	if err != nil {
		return err
	}

	cli.store.Pick(pickAs, parsedURL.Query().Get(argument), internalPicker.DisposableValue)

	return nil
}
