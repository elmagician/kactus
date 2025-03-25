package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
	"golang.org/x/net/html"

	"github.com/elmagician/kactus/internal/interfaces"
	match "github.com/elmagician/kactus/internal/matchers"
)

var (
	ErrNoBody     = errors.New("impossible to retrieve data from empty response body")
	ErrUnknownKey = errors.New("unknown key")
)

type Response struct {
	Status  int
	Headers http.Header
	Body    []byte
	Cookies map[string]*http.Cookie
}

func NewResponse(status int, body []byte, cookies []*http.Cookie, headers http.Header) *Response {
	cookieMap := make(map[string]*http.Cookie)

	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie
	}

	return &Response{
		Status:  status,
		Body:    body,
		Cookies: cookieMap,
		Headers: headers,
	}
}

func (r Response) RetrieveJSON(key string) (interface{}, error) {
	if r.HasEmptyBody() {
		return nil, ErrNoBody
	}
	var body interface{}

	if err := json.Unmarshal(r.Body, &body); err != nil {
		return nil, err
	}

	val, ok := interfaces.GetFieldFromPath(body, key)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownKey, key)
	}

	return val.Interface(), nil
}

func (r Response) RetrieveHeader(key string) string {
	return r.Headers.Get(key)
}

func (r Response) RetrieveHTMLAttribute(tag, attribute string, filters *godog.Table) (string, error) {
	if r.HasEmptyBody() {
		return "", ErrNoBody
	}

	var key, val, matcher, candidate string

	type matchOptions struct {
		Key     string
		Value   string
		Matcher string
	}

	matched := 0
	matchingOptions := map[string]matchOptions{}

	head := filters.Rows[0].Cells

	for i := 1; i < len(filters.Rows); i++ {
		for n, cell := range filters.Rows[i].Cells {
			switch head[n].Value {
			case "attribute":
				key = cell.Value
			case "value":
				val = cell.Value
			case "match":
				matcher = cell.Value
			}
		}

		matchingOptions[key] = matchOptions{
			Key: key, Value: val, Matcher: matcher,
		}
		key = ""
		val = ""
		matcher = ""
	}

	tokenized := html.NewTokenizer(bytes.NewReader(r.Body))

	for {
		next := tokenized.Next()
		//nolint:exhaustive
		switch next {
		case html.ErrorToken:
			return "", fmt.Errorf("could not find data for tag: %s attribute %s", tag, attribute)
		// nolint: gocritic
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenized.Token()
			if token.Data == tag {
				for _, att := range token.Attr {
					if att.Key == attribute {
						candidate = att.Val
					}

					if mo, ok := matchingOptions[att.Key]; ok {
						if err := match.Assert(mo.Matcher, att.Val, mo.Value); err == nil {
							matched++
						}
					}
				}

				if matched == len(matchingOptions) && candidate != "" {
					return candidate, nil
				}
			}

			matched = 0
		}
	}
}
