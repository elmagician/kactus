package test_test

import (
	"errors"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/elmagician/kactus/internal/test"
)

func TestUnit_ResetMock(t *testing.T) {
	Convey("Given some mock", t, func() {
		initTestStack()

		Convey("should be able to reset one", func() {
			So(testMock.Test("test"), ShouldBeNil)
			So(testMock.Test("another"), ShouldResemble, errors.New("err2"))

			test.ResetMock(&testMock.Mock)
			testMock.On("Test", "test").Return(errors.New("random"))

			So(testMock.Test("test"), ShouldResemble, errors.New("random"))
			So(testMock2.Do("smtg"), ShouldResemble, errors.New("err"))
		})

		Convey("should be able to reset list", func() {
			So(testMock.Test("test"), ShouldBeNil)
			So(testMock2.Do("smtg"), ShouldResemble, errors.New("err"))
			So(testMock.Test("another"), ShouldResemble, errors.New("err2"))

			test.ResetMocks(&testMock.Mock, &testMock2.Mock)
			testMock.On("Test", "test").Return(errors.New("random"))
			testMock2.On("Do", "smtg").Return(errors.New("some"))

			So(testMock.Test("test"), ShouldResemble, errors.New("random"))
			So(testMock2.Do("smtg"), ShouldResemble, errors.New("some"))
		})

		test.ResetMocks(&testMock.Mock, &testMock2.Mock)
	})
}

func TestUnit_AssertMockFullFilled(t *testing.T) {
	Convey("Given some mock", t, func() {
		initTestStack()
		tester := &testing.T{}
		Convey("should be ok if all expectation were met", func() {
			_ = testMock.Test("test")
			_ = testMock2.Do("smtg")
			_ = testMock.Test("another")

			So(test.AssertMockFullFilled(tester, &testMock.Mock, &testMock2.Mock), ShouldBeTrue)
			So(tester.Failed(), ShouldBeFalse)
		})

		Convey("should be false if some expectations were not met", func() {
			_ = testMock.Test("test")
			So(test.AssertMockFullFilled(tester, &testMock.Mock), ShouldBeFalse)
			So(test.AssertMockFullFilled(tester, &testMock2.Mock), ShouldBeFalse)
			So(tester.Failed(), ShouldBeTrue)

			_ = testMock2.Do("smtg")
			So(test.AssertMockFullFilled(tester, &testMock2.Mock), ShouldBeTrue)
			So(test.AssertMockFullFilled(tester, &testMock2.Mock, &testMock.Mock), ShouldBeFalse)

			_ = testMock.Test("another")
			So(test.AssertMockFullFilled(tester, &testMock2.Mock, &testMock.Mock), ShouldBeTrue)
			So(tester.Failed(), ShouldBeTrue)

			test.ResetMocks(&testMock.Mock, &testMock2.Mock)
			initTestStack()

			_ = testMock.Test("test")
			So(test.AssertMockFullFilled(tester, &testMock2.Mock, &testMock.Mock), ShouldBeFalse)
			_ = testMock.Test("another")
			So(test.AssertMockFullFilled(tester, &testMock.Mock), ShouldBeTrue)
			So(test.AssertMockFullFilled(tester, &testMock.Mock, &testMock2.Mock), ShouldBeFalse)
			So(tester.Failed(), ShouldBeTrue)
		})

		test.ResetMocks(&testMock.Mock, &testMock2.Mock)
	})
}

func TestUnit_ExecInBasePath(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	Convey("Given current path", t, func() {
		wd, _ := os.Getwd()
		basePath1, expectedPath1 := makeExpectedPath(wd)

		var basePath2 string
		var expectedPath2 string
		for {
			basePath2, expectedPath2 = makeExpectedPath(wd)
			if expectedPath2 != expectedPath1 {
				break
			}
		}

		Convey("should be able to move to path", func() {
			test.ExecInBasePath(basePath1)
			nwd, _ := os.Getwd()
			So(nwd, ShouldEqual, expectedPath1)
			_ = os.Chdir(wd)

			test.ExecInBasePath(basePath2)
			nwd, _ = os.Getwd()
			So(nwd, ShouldEqual, expectedPath2)
		})
	})
}

func TestUnit_NewSQLXMocked(t *testing.T) {
	dbMock, conn := test.NewSQLXMocked()

	Convey("I should be able to mock sql", t, func() {
		expectError := errors.New("rollback")
		dbMock.ExpectBegin()
		dbMock.ExpectRollback().WillReturnError(expectError)
		tx, err := conn.Beginx()
		So(err, ShouldBeNil)
		So(tx.Rollback(), ShouldBeError, expectError)
		So(dbMock.ExpectationsWereMet(), ShouldBeNil)
	})
}

func initTestStack() {
	testMock.On("Test", "test").Return(nil)
	testMock2.On("Do", "smtg").Return(errors.New("err"))
	testMock.On("Test", "another").Return(errors.New("err2"))
}

func makeExpectedPath(path string) (baseElem string, expectedPath string) {
	pathElem := strings.Split(path, "/")
	rngStop := (rand.Int() % (len(pathElem) - 1)) + 1
	baseElem = pathElem[rngStop]

	for _, el := range pathElem {
		expectedPath += "/" + el
		if baseElem == el {
			break
		}
	}
	expectedPath = strings.TrimPrefix(expectedPath, "/")
	return baseElem, expectedPath
}

type Mock struct {
	mock.Mock
}

func (m *Mock) Test(s string) error {
	args := m.Called(s)
	return args.Error(0)
}

type Mock2 struct {
	mock.Mock
}

func (m *Mock2) Do(action string) error {
	args := m.Called(action)
	return args.Error(0)
}

var (
	testMock  = &Mock{}
	testMock2 = &Mock2{}
)
