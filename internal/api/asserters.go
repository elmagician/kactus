package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/elmagician/godog"

	match "github.com/elmagician/kactus/internal/matchers"

	"github.com/elmagician/kactus/internal"
	"github.com/elmagician/kactus/internal/interfaces"
)

// const for godog.Table default header on assertion
const (
	valueHeader   = "value"
	fieldHeader   = "field"
	matcherHeader = "matcher"
)

var (
	ErrNoMatch       = errors.New("actual value does not match expected")
	ErrNoPath        = errors.New("path does not exists in object")
	ErrNotFullyMatch = errors.New("expected values to be fully matched but it was not")
)

func (r Response) HasStatus(status int) bool {
	return r.Status == status
}

func (r Response) HasEmptyBody() bool {
	return len(r.Body) == 0
}

func (r Response) JSONContains(fully bool, expected *godog.Table) error {
	var (
		actual               interface{}
		actualFullPaths      []string
		expectedPaths        []string
		path, value, matcher string
	)

	if r.HasEmptyBody() {
		return ErrNoBody
	}

	if err := json.Unmarshal(r.Body, &actual); err != nil {
		return err
	}

	if fully {
		// cannot fail unmarshal has all values are the same than previous call
		actualFullPaths = interfaces.GenerateFieldList(actual)
	}

	head := expected.Rows[0].Cells

	for i := 1; i < len(expected.Rows); i++ {
		for n, cell := range expected.Rows[i].Cells {
			switch head[n].Value {
			case fieldHeader:
				path = cell.Value
			case matcherHeader:
				matcher = cell.Value
			case valueHeader:
				value = cell.Value
			default:
				return fmt.Errorf("%w %s", internal.ErrUnexpectedColumn, head[n].Value)
			}
		}

		if actualVal, exists := interfaces.GetFieldFromPath(actual, path); exists {
			if err := match.Assert(matcher, actualVal, value); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%w: %v", ErrUnknownKey, path)
		}

		if fully {
			if !inList(path, expectedPaths) {
				expectedPaths = append(expectedPaths, path)
			}
		}

		path = ""
		value = ""
		matcher = ""
	}

	if fully {
		for _, path = range expectedPaths {
			for i, actualPath := range actualFullPaths {
				if path == actualPath {
					actualFullPaths = remove(actualFullPaths, i)
				}
			}
		}

		if len(actualFullPaths) != 0 {
			return fmt.Errorf("%w: missing %v", ErrNotFullyMatch, actualFullPaths)
		}
	}

	return nil
}

func (r Response) JSONResemble(expectedBody *godog.DocString) error {
	var expected, actual interface{}
	var err error

	// re-encode expected response
	if err = json.Unmarshal([]byte(expectedBody.Content), &expected); err != nil {
		return err
	}

	// re-encode actual response too
	if err = json.Unmarshal(r.Body, &actual); err != nil {
		return err
	}

	// the matching may be adapted per different requirements.
	if !cmp.Equal(expected, actual) {
		return fmt.Errorf("%w, %v vs. %v", ErrNoMatch, expected, actual)
	}

	return nil
}

func (r Response) HTMLContain(expectedBody *godog.Table) error {
	actual := string(r.Body)

	for _, row := range expectedBody.Rows {
		for _, cell := range row.Cells {
			val := cell.Value

			if !strings.Contains(actual, val) {
				return fmt.Errorf("%w, %v vs. %v", ErrNoMatch, val, actual)
			}
		}
	}

	return nil
}

func (r Response) HTMLResemble(expectedBody *godog.DocString) error {
	var expected, actual interface{}

	expected = expectedBody.Content
	actual = string(r.Body)

	// the matching may be adapted per different requirements.
	if expected != actual {
		return fmt.Errorf("%w, %v vs. %v", ErrNoMatch, expected, actual)
	}

	return nil
}

func (r Response) HeaderMatches(expected *godog.Table) error {
	var (
		key, value, matcher string
	)

	head := expected.Rows[0].Cells

	for i := 1; i < len(expected.Rows); i++ {
		for n, cell := range expected.Rows[i].Cells {
			switch head[n].Value {
			case "key":
				key = cell.Value
			case matcherHeader:
				matcher = cell.Value
			case valueHeader:
				value = cell.Value
			default:
				return fmt.Errorf("%w %s", internal.ErrUnexpectedColumn, head[n].Value)
			}
		}

		if actualVal := r.Headers.Get(key); actualVal != "" {
			if err := match.Assert(matcher, actualVal, value); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%w: %v", ErrNoMatch, key)
		}

		key = ""
		value = ""
		matcher = ""
	}

	return nil
}

func (r Response) HasCookie(cookieName string) bool {
	_, has := r.Cookies[cookieName]
	return has
}

func (r Response) GetCookie(cookieName string) *http.Cookie {
	cookie, has := r.Cookies[cookieName]
	if !has {
		return nil
	}

	return cookie
}

func (cli Client) HasCookie(cookieName string) bool {
	if cli.request == nil {
		return false
	}

	for _, cookie := range cli.client.Jar.Cookies(cli.request.URL) {
		if cookie.Name == cookieName {
			return true
		}
	}

	return false
}

func inList(e string, l []string) bool {
	for _, known := range l {
		if e == known {
			return true
		}
	}

	return false
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
