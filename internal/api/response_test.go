package api_test

import (
	"net/http"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/cucumber/messages-go/v10"
	"github.com/elmagician/godog"
	"github.com/elmagician/kactus/internal/api"
	. "github.com/elmagician/kactus/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnit_NewResponse(t *testing.T) {
	Convey("When I try to get new response should success", t, func() {
		status := int(fake.Int8())
		body := []byte("someBody")
		cookie1Key := "cookie1"
		cookie1Value := "cookie1Value"
		cookie2Key := "cookie2"
		cookie2Value := "cookie2Value"
		cookie1 := &http.Cookie{
			Name:  cookie1Key,
			Value: cookie1Value,
		}
		cookie2 := &http.Cookie{
			Name:  cookie2Key,
			Value: cookie2Value,
		}
		cookies := []*http.Cookie{cookie1, cookie2}
		headers := http.Header{}

		expectedResponse := &api.Response{
			Body:    body,
			Cookies: map[string]*http.Cookie{cookie1Key: cookie1, cookie2Key: cookie2},
			Headers: headers,
			Status:  status,
		}

		r := api.NewResponse(status, body, cookies, headers)

		So(r, ShouldBeEquivalent, expectedResponse)
	})
}

func TestUnit_Response_RetrieveJSON(t *testing.T) {
	Convey("When I try to retrieve JSON", t, func() {
		r := api.Response{
			Body: []byte("{\"foo\":\"bar\"}"),
		}

		key := "foo"
		expectedValue := "bar"

		Convey("should success", func() {
			v, err := r.RetrieveJSON(key)

			So(err, ShouldBeNil)
			So(v, ShouldEqual, expectedValue)
		})

		Convey("Should fail", func() {
			Convey("if empty body", func() {
				r.Body = nil

				v, err := r.RetrieveJSON(key)

				So(v, ShouldBeNil)
				So(err, ShouldBeError, api.ErrNoBody)
			})

			Convey("if unknown key in body", func() {
				v, err := r.RetrieveJSON("someBadKey")

				So(v, ShouldBeNil)
				So(err, ShouldBeLikeError, api.ErrUnknownKey)
			})
		})
	})
}

func TestUnit_Response_RetrieveHeader(t *testing.T) {
	Convey("When I try to retrieve header should success", t, func() {
		headerKey := "Key"
		headerValue := "value"

		r := api.Response{
			Headers: http.Header{headerKey: []string{headerValue}},
		}

		So(r.RetrieveHeader(headerKey), ShouldEqual, headerValue)
	})
}

func TestUnit_Response_RetrieveHTMLAttribute(t *testing.T) {
	Convey("When I try to retrieve HTML attribute", t, func() {
		tag, attribute := "p", "id"
		filters := &godog.Table{
			Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
				{
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: "attribute"},
						{Value: "value"},
						{Value: "match"},
					},
				},
				{
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: "class"},
						{Value: "SomeClass"},
						{Value: "contain"},
					},
				},
			},
		}

		expectedID := "someID"

		r := api.Response{
			Body: []byte("<p id=\"someID\" class=\"foobarSomeClass\">value</p>"),
		}

		Convey("should success", func() {
			Convey("without filters", func() {
				filters.Rows = []*messages.PickleStepArgument_PickleTable_PickleTableRow{
					{
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "attribute"},
							{Value: "value"},
							{Value: "match"},
						},
					}}

				v, err := r.RetrieveHTMLAttribute(tag, attribute, filters)

				So(err, ShouldBeNil)
				So(v, ShouldEqual, expectedID)
			})

			Convey("with filters", func() {
				v, err := r.RetrieveHTMLAttribute(tag, attribute, filters)

				So(err, ShouldBeNil)
				So(v, ShouldEqual, expectedID)
			})
		})

		Convey("should fail if body is empty", func() {
			r.Body = nil

			v, err := r.RetrieveHTMLAttribute(tag, attribute, filters)

			So(v, ShouldBeZeroValue)
			So(err, ShouldBeError, api.ErrNoBody)
		})
	})
}
