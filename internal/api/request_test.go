package api_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/cucumber/messages-go/v10"
	"github.com/elmagician/godog"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/elmagician/kactus/internal/api"
	. "github.com/elmagician/kactus/internal/test"
)

func TestUnit_PrepareRequest(t *testing.T) {
	Convey("When I try to prepare request", t, func() {
		var allowCookie bool

		Convey("should success with cookie allowed", func() {
			allowCookie = true
		})

		Convey("should success with cookie not allowed", func() {
			allowCookie = false
		})

		r := api.PrepareRequest(allowCookie)

		So(r.AllowCookie, ShouldEqual, allowCookie)
	})
}

func TestUnit_RequestPreparation_Empty(t *testing.T) {
	Convey("When I try to check if request if empty", t, func() {
		somebody := "somebody"
		r := api.RequestPreparation{
			AllowCookie: true,
			JSONBody:    &somebody,
		}

		Convey("should get false if body not empty", func() {
			So(r.Empty(), ShouldBeFalse)
		})

		Convey("should get true if body not empty", func() {
			r := api.RequestPreparation{}

			So(r.Empty(), ShouldBeTrue)
		})
	})
}

func TestUnit_RequestPreparation_SetMethod(t *testing.T) {
	Convey("When I try to set method should success", t, func() {
		method := "GET"
		r := api.RequestPreparation{}.SetMethod(method)

		So(r.Method, ShouldEqual, method)
	})
}

func TestUnit_RequestPreparation_SetEndpoint(t *testing.T) {
	Convey("When I try to set endpoint then should success", t, func() {
		endpoint := "someEndpoint"
		r := api.RequestPreparation{}.SetEndpoint(endpoint)

		So(r.Endpoint, ShouldEqual, endpoint)
	})
}

func TestUnit_RequestPreparation_SetJSONBody(t *testing.T) {
	Convey("WHen I try to set json body should success", t, func() {
		content := "{\"foo\":\"bar\"}"
		r := api.RequestPreparation{}.SetJSONBody(&godog.DocString{Content: content})

		So(*r.JSONBody, ShouldEqual, content)
	})
}

func TestUnit_RequestPreparation_SetFORMBody(t *testing.T) {
	Convey("When I try to set form body should success", t, func() {
		key1 := "key1"
		kind1 := "string"
		value1 := "value1"

		key2 := "key2"
		kind2 := "file"
		value2 := "value2"

		body := &godog.Table{
			Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
				{
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: "key"},
						{Value: "value"},
						{Value: "kind"},
					},
				}, {
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: key1},
						{Value: value1},
						{Value: kind1},
					},
				}, {
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: key2},
						{Value: value2},
						{Value: kind2},
					},
				},
			},
		}

		r := api.RequestPreparation{}.SetFORMBody(body)

		So(len(r.FORMBody), ShouldEqual, 2)
	})
}

func TestUnit_RequestPreparation_ResetBody(t *testing.T) {
	Convey("WHen I try to reset body then should success", t, func() {
		val := "someBody"
		r := api.RequestPreparation{
			JSONBody: &val,
		}

		r = r.ResetBody()

		So(r.JSONBody, ShouldBeNil)
		So(r.FORMBody, ShouldBeNil)
	})
}

func TestUnit_RequestPreparation_AddHeader(t *testing.T) {
	Convey("When I try to add header then should success", t, func() {
		headerKey := "key"
		headerValue := "value"

		r := api.PrepareRequest(true).AddHeader(headerKey, headerValue)

		So(r.Headers.Get(headerKey), ShouldEqual, headerValue)
	})
}

func TestUnit_RequestPreparation_ResetHeader(t *testing.T) {
	Convey("When I try to reset header should sucess", t, func() {
		r := api.RequestPreparation{Headers: &http.Header{"header1": []string{"value1"}}}

		r = r.ResetHeader()

		So(r.Headers, ShouldBeEquivalent, &http.Header{})
	})
}

func TestUnit_RequestPreparation_AddArgument(t *testing.T) {
	Convey("When I try to add argument should success", t, func() {
		argKey := "key"
		argValue := "value"

		r := api.PrepareRequest(true).AddArgument(argKey, argValue)

		So(r.Arguments[argKey], ShouldEqual, argValue)
	})
}

func TestUnit_RequestPreparation_ResetArguments(t *testing.T) {
	Convey("When I try to reset arguments should success", t, func() {
		r := api.RequestPreparation{
			Arguments: map[string]string{"key": "value"},
		}

		r = r.ResetArguments()

		So(r.Arguments, ShouldBeEmpty)
	})
}

func TestUnit_RequestPreparation_AddCookie(t *testing.T) {
	Convey("When I try to add cookie", t, func() {
		r := api.PrepareRequest(true)

		key := "key"
		value := "value"
		path := "somePath"
		domain := "someDomain"
		expires := fake.Int64()
		rawExpires := "expires"
		maxAge := 5
		secure := fake.Bool()
		raw := "someRaw"
		sameSite := 15
		httpOnly := fake.Bool()
		unparsed := "unparsed"
		options := &messages.PickleStepArgument_PickleTable{
			Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
				{
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: "path"},
						{Value: "domain"},
						{Value: "expires"},
						{Value: "raw_expires"},
						{Value: "maxAge"},
						{Value: "secure"},
						{Value: "raw"},
						{Value: "sameSite"},
						{Value: "httpOnly"},
						{Value: "unparsed"},
						{Value: "unknown"},
					},
				}, {
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: path},
						{Value: domain},
						{Value: strconv.Itoa(int(expires))},
						{Value: rawExpires},
						{Value: strconv.Itoa(maxAge)},
						{Value: strconv.FormatBool(secure)},
						{Value: raw},
						{Value: strconv.Itoa(sameSite)},
						{Value: strconv.FormatBool(httpOnly)},
						{Value: unparsed},
						{Value: "unknown value"},
					},
				},
			},
		}

		expectedCookie := &http.Cookie{
			Name:       key,
			Value:      value,
			Path:       path,
			Domain:     domain,
			Expires:    time.Unix(expires, 0),
			RawExpires: rawExpires,
			MaxAge:     maxAge,
			Secure:     secure,
			HttpOnly:   httpOnly,
			SameSite:   http.SameSite(sameSite),
			Raw:        raw,
			Unparsed:   strings.Split(unparsed, ";"),
		}

		Convey("should success", func() {
			r, err := r.AddCookie(key, value, options)

			So(err, ShouldBeNil)
			So(len(r.Cookies), ShouldEqual, 1)
			So(r.Cookies[0], ShouldBeEquivalent, expectedCookie)
		})

		Convey("should fail if cannot parse expires", func() {
			options.Rows[1].Cells[2].Value = "someBadValue"

			_, err := r.AddCookie(key, value, options)

			So(err, ShouldBeError)
		})
	})
}

func TestUnit_RequestPreparation_ResetCookies(t *testing.T) {
	Convey("When I try to reset cookies should success", t, func() {
		r := api.PrepareRequest(true)
		r.Cookies = append(r.Cookies, &http.Cookie{Name: "someName"})

		So(r.ResetCookies().Cookies, ShouldBeEmpty)
	})

}

func TestUnit_RequestPreparation_AllowCookies(t *testing.T) {
	Convey("When I try to allow cookies should success", t, func() {
		allowCookie := fake.Bool()

		So(api.RequestPreparation{}.AllowCookies(allowCookie).AllowCookie, ShouldEqual, allowCookie)
	})
}
