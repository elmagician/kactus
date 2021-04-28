package api_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/cucumber/messages-go/v10"
	"github.com/elmagician/godog"
	"github.com/elmagician/kactus/internal"
	"github.com/elmagician/kactus/internal/api"
	"github.com/elmagician/kactus/internal/interfaces"
	"github.com/elmagician/kactus/internal/matchers"
	. "github.com/elmagician/kactus/internal/test"
	"github.com/elmagician/kactus/internal/types"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	api.NoLog()
	interfaces.NoLog()
	matchers.NoLog()
	types.NoLog()
}

func TestUnit_Response_JSONContains(t *testing.T) {
	Convey("Given a response", t, func() {
		jsonBody := `[{"name":"george","surname":"weasley","class":{"name":"magician","level":"superior","strong":"funy magic","weak":"serious"},"spells":["willy","shield","fierce"]},{"name":"fred","surname":"weasley","class":{"name":"magician","level":"superior","strong":"funy magic","weak":"serious"},"spells":["link","lang","cards"]},{"name":"hermione","surname":"granger","class":{"name":"magician","level":"great","strong":"learner","weak":"fun"},"spells":["open","great shadow","forget"]}]`
		r := api.Response{
			Body:    []byte(jsonBody),
			Cookies: nil,
			Headers: nil,
			Status:  0,
		}

		Convey("I should be able to match partially response JSON body", func() {
			partialDataTable := &messages.PickleStepArgument_PickleTable{
				Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
					{
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "field"},
							{Value: "matcher"},
							{Value: "value"},
						},
					}, {
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "0"},
							{Value: "not zero"},
							{Value: ""},
						},
					}, {
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "1.name"},
							{Value: "eq"},
							{Value: "fred"},
						},
					}, {
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "2.class.strong"},
							{Value: "eq"},
							{Value: "learner"},
						},
					},
				},
			}

			So(r.JSONContains(false, partialDataTable), ShouldBeNil)
		})
		Convey("I should be able to fully match response JSON body", func() {
			jsonBody := `[
    {
        "name": "fred"
    },
    {
        "class": {
            "strong": "learner"
        }
    }
]`
			r.Body = []byte(jsonBody)

			partialDataTable := &messages.PickleStepArgument_PickleTable{
				Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
					{
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "field"},
							{Value: "matcher"},
							{Value: "value"},
						},
					}, {
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "0.name"},
							{Value: "eq"},
							{Value: "fred"},
						},
					}, {
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "1.class.strong"},
							{Value: "eq"},
							{Value: "learner"},
						},
					},
				},
			}

			So(r.JSONContains(false, partialDataTable), ShouldBeNil)
		})

		Convey("I should have an error if JSON response body does not fully match expected body (too much fields)", func() {
			partialDataTable := &messages.PickleStepArgument_PickleTable{
				Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
					{
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "field"},
							{Value: "matcher"},
							{Value: "value"},
						},
					}, {
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "0"},
							{Value: "not zero"},
							{Value: ""},
						},
					}, {
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "1.name"},
							{Value: "eq"},
							{Value: "fred"},
						},
					}, {
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "2.class.strong"},
							{Value: "eq"},
							{Value: "learner"},
						},
					},
				},
			}

			So(errors.Is(r.JSONContains(true, partialDataTable), api.ErrNotFullyMatch), ShouldBeTrue)
		})
	})
}

func TestUnit_Response_HasStatus(t *testing.T) {
	Convey("When I try to check if response has status", t, func() {

		status := http.StatusOK
		r := api.Response{Status: status}

		Convey("Should get true if response has status", func() {
			So(r.HasStatus(status), ShouldBeTrue)
		})

		Convey("Should get false if response has status", func() {
			So(r.HasStatus(http.StatusNotFound), ShouldBeFalse)
		})
	})
}

func TestUnit_Response_JSONResemble(t *testing.T) {
	Convey("When I try to compare json", t, func() {

		json := "{\"foo\":\"bar\"}"
		r := api.Response{Body: []byte(json)}
		expectedBody := &godog.DocString{
			MediaType: "",
			Content:   json,
		}

		Convey("should success if json resemble", func() {
			So(r.JSONResemble(expectedBody), ShouldBeNil)
		})

		Convey("should fail", func() {
			Convey("if expected body should not contain valid json", func() {
				expectedBody.Content = "[[[}}"

				So(r.JSONResemble(expectedBody), ShouldBeError)
			})

			Convey("if response body do not contain valid json", func() {
				r.Body = []byte("]]}")

				So(r.JSONResemble(expectedBody), ShouldBeError)
			})

			Convey("if response body and expected body are not equals", func() {
				expectedBody.Content = "{\"bar\":\"foo\"}"

				So(r.JSONResemble(expectedBody), ShouldBeLikeError, api.ErrNoMatch)
			})
		})
	})
}

func TestUnit_Response_HTMLContain(t *testing.T) {
	Convey("When I try to check if html contain", t, func() {
		value1 := "someHTML"
		value2 := "another HTML"
		r := api.Response{Body: []byte(value1 + value2)}
		expectedBody := &godog.Table{
			Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
				{
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: value1},
						{Value: value2},
					},
				},
				{
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: value2},
					},
				},
			},
		}

		Convey("should success", func() {
			So(r.HTMLContain(expectedBody), ShouldBeNil)
		})

		Convey("should fail if html do not contain", func() {
			r.Body = []byte("unexpectedBody")

			So(r.HTMLContain(expectedBody), ShouldBeLikeError, api.ErrNoMatch)
		})
	})
}

func TestUnit_Response_HTMLResemble(t *testing.T) {
	Convey("When I try to check if html resemble", t, func() {

		body := "someHTML"
		r := api.Response{Body: []byte(body)}
		expectedBody := &godog.DocString{Content: body}

		Convey("should success", func() {
			So(r.HTMLResemble(expectedBody), ShouldBeNil)
		})

		Convey("should fail if html does not resemble", func() {
			expectedBody.Content = "unexpectedHTML"

			So(r.HTMLResemble(expectedBody), ShouldBeLikeError, api.ErrNoMatch)
		})
	})
}

func TestUnit_Response_HeaderMatches(t *testing.T) {
	Convey("When I try to check if header matches", t, func() {

		header1Key := "Header1"
		header2Key := "Header2"
		header1Value := "header1Value"
		header2Value := "header2Value"

		r := api.Response{
			Headers: http.Header{
				header1Key: []string{header1Value},
				header2Key: []string{header2Value},
			},
		}

		expectedHeaders := &godog.Table{
			Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
				{
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: "key"},
						{Value: "matcher"},
						{Value: "value"},
					},
				}, {
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: header1Key},
						{Value: "eq"},
						{Value: header1Value},
					},
				}, {
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: header2Key},
						{Value: "eq"},
						{Value: header2Value},
					},
				},
			},
		}

		Convey("should success", func() {
			So(r.HeaderMatches(expectedHeaders), ShouldBeNil)
		})

		Convey("should fail", func() {
			Convey("if expected Header does not contain correct head", func() {
				expectedHeaders.Rows[0].Cells[0].Value = "unexpectedValue"

				So(r.HeaderMatches(expectedHeaders), ShouldBeLikeError, internal.ErrUnexpectedColumn)
			})

			Convey("if header does not contain expected value", func() {
				expectedHeaders.Rows[1].Cells[2].Value = "badValue"

				So(r.HeaderMatches(expectedHeaders), ShouldBeError)
			})

			Convey("if header does not contain asked key", func() {
				expectedHeaders.Rows[1].Cells[0].Value = "badKey"

				So(r.HeaderMatches(expectedHeaders), ShouldBeLikeError, api.ErrNoMatch)
			})
		})
	})
}

func TestUnit_Response_HasCookie(t *testing.T) {
	Convey("When I try to check if response has cookie", t, func() {
		cookieName := "someCookieName"
		r := api.Response{Cookies: map[string]*http.Cookie{cookieName: {Name: cookieName}}}

		Convey("should success", func() {
			So(r.HasCookie(cookieName), ShouldBeTrue)
		})

		Convey("should fail if response does not contain cookie", func() {
			So(r.HasCookie("badCookieName"), ShouldBeFalse)
		})
	})
}

func TestUnit_Response_GetCookie(t *testing.T) {
	Convey("When I try to get cookie", t, func() {
		cookieName := "cookieName"
		expectedCookie := &http.Cookie{Name: cookieName}
		r := api.Response{Cookies: map[string]*http.Cookie{cookieName: expectedCookie}}

		Convey("should success", func() {
			So(r.GetCookie(cookieName), ShouldBeEquivalent, expectedCookie)
		})

		Convey("should get nil if cookie does not exist", func() {
			So(r.GetCookie("badCookieName"), ShouldBeNil)
		})
	})
}
