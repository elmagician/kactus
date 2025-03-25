package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
)

const (
	unknown formElementKind = iota
	file
)

type (
	RequestPreparation struct {
		AllowCookie bool
		JSONBody    *string
		FORMBody    map[string]formElement
		Headers     *http.Header
		Cookies     []*http.Cookie
		Arguments   map[string]string
		Endpoint    string
		Method      string
	}

	formElement struct {
		element string
		kind    formElementKind
	}

	formElementKind int
)

func PrepareRequest(withCookie bool) RequestPreparation {
	return RequestPreparation{
		AllowCookie: withCookie,
		Method:      "GET",
		Headers:     &http.Header{},
		Arguments:   make(map[string]string),
	}
}

func (request RequestPreparation) Empty() bool {
	return cmp.Equal(request, RequestPreparation{})
}

func (request RequestPreparation) SetMethod(method string) RequestPreparation {
	request.Method = method
	return request
}

func (request RequestPreparation) SetEndpoint(endpoint string) RequestPreparation {
	request.Endpoint = endpoint
	return request
}

// Body management

func (request RequestPreparation) SetJSONBody(body *godog.DocString) RequestPreparation {
	content := body.Content
	request.JSONBody = &content

	return request
}

func (request RequestPreparation) SetFORMBody(body *godog.Table) RequestPreparation {
	var key, val, kind string

	if request.FORMBody == nil {
		request.FORMBody = make(map[string]formElement)
	}

	headers := body.Rows[0].Cells

	for i := 1; i < len(body.Rows); i++ {
		for n, cell := range body.Rows[i].Cells {
			switch headers[n].Value {
			case "key":
				key = cell.Value
			case "value", "val":
				val = cell.Value
			case "kind":
				kind = cell.Value
			}
		}

		request.FORMBody[key] = formElement{element: val, kind: formElementKindFromString(kind)}

		log.Debug(
			"adding form raw",
			zap.String("key", key), zap.String("val", val),
			zap.String("kind", kind), zap.Reflect("form", request.FORMBody[key]),
		)

		val = ""
		kind = ""
		key = ""
	}

	return request
}

func (request RequestPreparation) ResetBody() RequestPreparation {
	request.JSONBody = nil
	request.FORMBody = nil

	return request
}

// Headers management

func (request RequestPreparation) AddHeader(key, value string) RequestPreparation {
	request.Headers.Add(key, value)
	return request
}

func (request RequestPreparation) ResetHeader() RequestPreparation {
	request.Headers = &http.Header{}
	return request
}

// Arguments management

func (request RequestPreparation) AddArgument(key, value string) RequestPreparation {
	request.Arguments[key] = value
	return request
}

func (request RequestPreparation) ResetArguments() RequestPreparation {
	request.Arguments = make(map[string]string)
	return request
}

// Cookies management

// nolint: gocyclo
func (request RequestPreparation) AddCookie(key, val string, options *godog.Table) (RequestPreparation, error) {
	var (
		err error

		path, domain, rawExpires, expires, maxAge string
		secure, httpOnly, sameSite, raw, unparsed string
	)

	if len(options.Rows) > 1 {
		head := options.Rows[0].Cells

		for n, cell := range options.Rows[1].Cells {
			switch head[n].Value {
			case "path":
				path = cell.Value
			case "domain":
				domain = cell.Value
			case "expires":
				expires = cell.Value
			case "raw_expires", "rawExpires", "rawexpires":
				rawExpires = cell.Value
			case "maxAge", "max_age", "maxage":
				maxAge = cell.Value
			case "secure":
				secure = cell.Value
			case "raw":
				raw = cell.Value
			case "sameSite", "same_site", "samesite":
				sameSite = cell.Value
			case "httpOnly", "http_only", "httponly":
				httpOnly = cell.Value
			case "unparsed":
				unparsed = cell.Value
			default:
				log.Debug("unexpected column name %s" + head[n].Value)
			}
		}
	}

	timeExpires := time.Time{}

	if expires != "" {
		i, sErr := strconv.ParseInt(expires, 10, 64)
		if sErr != nil {
			return request, sErr
		}

		timeExpires = time.Unix(i, 0)
	}

	intMaxAge := 0
	if maxAge != "" {
		intMaxAge, err = strconv.Atoi(maxAge)
	}

	boolSecure := false
	if secure != "" {
		boolSecure, err = strconv.ParseBool(secure)
	}

	boolHTTPpOnly := false
	if httpOnly != "" {
		boolHTTPpOnly, err = strconv.ParseBool(httpOnly)
	}

	sameSiteInt := 0
	if sameSite != "" {
		sameSiteInt, err = strconv.Atoi(sameSite)
	}

	var unparsedList []string
	if unparsed != "" {
		unparsedList = strings.Split(unparsed, ";")
	}

	request.Cookies = append(request.Cookies, &http.Cookie{
		Name:       key,
		Value:      val,
		Path:       path,
		Domain:     domain,
		Expires:    timeExpires,
		RawExpires: rawExpires,
		MaxAge:     intMaxAge,
		Secure:     boolSecure,
		HttpOnly:   boolHTTPpOnly,
		SameSite:   http.SameSite(sameSiteInt),
		Raw:        raw,
		Unparsed:   unparsedList,
	})

	return request, err
}

func (request RequestPreparation) ResetCookies() RequestPreparation {
	request.Cookies = []*http.Cookie{}
	return request
}

func (request RequestPreparation) AllowCookies(allow bool) RequestPreparation {
	request.AllowCookie = allow
	return request
}

// Generate request
func (request RequestPreparation) GenerateRequest(jar http.CookieJar) (*http.Request, error) { //nolint:gocyclo
	var (
		req *http.Request
		err error
		b   bytes.Buffer

		formWriter  = multipart.NewWriter(&b)
		hasForm     = false
		hasBody     = false
		contentType = ""
		body        = ""
	)

	if request.JSONBody != nil {
		hasBody = true
		body += *request.JSONBody

		contentType = "application/json"
	}

	if body == "" && request.FORMBody != nil {
		hasForm = true
		hasBody = true

		for key, element := range request.FORMBody {
			var fw io.Writer

			switch element.kind {
			case file:
				if fw, err = formWriter.CreateFormField(key); err != nil {
					openFile, sErr := os.Open(element.element)
					if sErr != nil {
						return nil, sErr
					}

					if fw, err = formWriter.CreateFormFile(key, element.element); err != nil {
						return nil, err
					}

					if _, err = io.Copy(fw, openFile); err != nil {
						return nil, err
					}

					if err = openFile.Close(); err != nil {
						return nil, err
					}
				}

				if _, err = io.Copy(fw, strings.NewReader(element.element)); err != nil {
					return nil, err
				}
			case unknown:
				if writeErr := formWriter.WriteField(key, element.element); writeErr != nil {
					return nil, writeErr
				}
			}

			contentType = formWriter.FormDataContentType()
		}
	}

	if hasBody {
		if request.Headers == nil {
			request.Headers = &http.Header{}
		}

		request.Headers.Set("Content-Type", contentType)
	}

	if err = formWriter.Close(); err != nil {
		return nil, err
	}

	if hasForm {
		req, err = http.NewRequest(request.Method, request.Endpoint, &b)
	} else {
		req, err = http.NewRequest(request.Method, request.Endpoint, strings.NewReader(body))
	}

	if err != nil {
		return nil, err
	}

	if request.AllowCookie {
		jar.SetCookies(req.URL, request.Cookies)
	}

	q := req.URL.Query()

	for key, value := range request.Arguments {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	req.Header = *request.Headers

	return req, nil
}

func formElementKindFromString(kind string) formElementKind {
	switch kind {
	case "file":
		return file
	default:
		return unknown
	}
}
