package test

import (
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

// nolint:lll
const (
	LowerLetterCasedString = "abcdefghijklmnopqrstuvwxyz"
	UpperLetterCasedString = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Symbols                = "*/-+$¤£%ù*μ!§~#{'([-|`_\\]°)}´®†‹><²"
	NumberOnlyString       = "0123456789"
	LettersOnlyString      = LowerLetterCasedString + UpperLetterCasedString
	LowerCasedString       = LowerLetterCasedString + Symbols + NumberOnlyString
	UpperCasedString       = UpperLetterCasedString + Symbols + NumberOnlyString
	String                 = LowerLetterCasedString + UpperLetterCasedString + Symbols + NumberOnlyString + LettersOnlyString
)

// init
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Get random string of asked length
// Provided character list to use to generate string
//
// Provided charsets:
// 	LowerLetterCasedString = "abcdefghijklmnopqrstuvwxyz"
//	UpperLetterCasedString = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
//	Symbols                = "*/-+$¤£%ù*μ!§~#{'([-|`_\\]°)}´®†‹><²"
//	NumberOnlyString       = "0123456789"
//	LettersOnlyString      = LowerLetterCasedString + UpperLetterCasedString
//	LowerCasedString       = LowerLetterCasedString + Symbols + NumberOnlyString
//	UpperCasedString       = UpperLetterCasedString + Symbols + NumberOnlyString
//	String                 = LowerLetterCasedString + UpperLetterCasedString + Symbols
//                      	+ NumberOnlyString + LettersOnlyString
func GetStringOfLength(containChars string, length int) string {
	b := make([]byte, length)
	lenProvidedChars := len(containChars)

	for i := range b {
		b[i] = containChars[rand.Intn(lenProvidedChars)] // nolint: gosec
	}

	return string(b)
}

// New mocked sqlx db
func NewSQLXMocked() (sqlmock.Sqlmock, *sqlx.DB) {
	db, mocker, _ := sqlmock.New() // nolint: errcheck
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	return mocker, sqlxDB
}

// AssertMockFullFilled check list of mock to ensure all defined expectations are meet
func AssertMockFullFilled(t *testing.T, mocks ...*mock.Mock) bool {
	res := true
	for _, m := range mocks {
		res = m.AssertExpectations(t) && res
	}

	return res
}

// ExecInBasePath move path up to first matching part
// ex: in path /usr/test/gitlab/platform/goCore/test
// ExecInBasePath(goCore) will move to /usr/test/gitlab/platform/goCore
// ExecInBasePath(test) will move to /usr/test
func ExecInBasePath(baseName string) {
	dir, _ := os.Getwd() // nolint: errcheck
	basePath := ""
	dirList := strings.Split(dir, "/")

	for _, pathElement := range dirList {
		basePath += "/" + pathElement

		if pathElement == baseName {
			break
		}
	}

	_ = os.Chdir(strings.TrimPrefix(basePath, "/")) // nolint: errcheck
}

// ResetMocks reset assertions for provided mock list
func ResetMocks(mocks ...*mock.Mock) {
	for _, m := range mocks {
		ResetMock(m)
	}
}

// ResetMock reset assertion for provided mock
func ResetMock(m *mock.Mock) {
	m.ExpectedCalls = []*mock.Call{}
	m.Calls = []mock.Call{}
}
