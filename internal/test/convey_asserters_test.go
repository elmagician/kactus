package test_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/elmagician/kactus/internal/test"
)

func TestUnit_ErrorShouldBeLike(t *testing.T) {
	Convey("Given some errors", t, func() {
		mainError := errors.New("main error")
		error1 := errors.New("wrong_1")
		error2 := errors.New("wrong_2")
		error1Wrap := fmt.Errorf("wrap: %w", error1)
		error2Wrap := fmt.Errorf("wrap: %w", error2)

		Convey("should success", func() {
			So(mainError, test.ShouldBeLikeError, mainError)
			So(error1, test.ShouldBeLikeError, error1)
			So(error1Wrap, test.ShouldBeLikeError, error1)
		})

		Convey("should fail", func() {
			So(test.ShouldBeLikeError(mainError, error1), ShouldNotBeZeroValue)
			So(test.ShouldBeLikeError(mainError, error1Wrap), ShouldNotBeZeroValue)
			So(test.ShouldBeLikeError(error1, error2), ShouldNotBeZeroValue)
			So(test.ShouldBeLikeError(error1Wrap, error2Wrap), ShouldNotBeZeroValue)
			So(test.ShouldBeLikeError(error1Wrap), ShouldNotBeZeroValue)
			So(test.ShouldBeLikeError("test", error2Wrap), ShouldNotBeZeroValue)
			So(test.ShouldBeLikeError(mainError, "test"), ShouldNotBeZeroValue)
		})
	})
}

func TestUnit_ShouldBeFullFilled(t *testing.T) {
	dbMock1, db1 := test.NewSQLXMocked()
	dbMock2, db2 := test.NewSQLXMocked()
	m1 := &Mock{}
	m2 := &Mock2{}

	dbMockList := []sqlmock.Sqlmock{dbMock1, dbMock2}
	testifyMocks := []*mock.Mock{&m1.Mock, &m2.Mock}

	Convey("Given mocks: ", t, func() {
		Convey("single instance on DB", func() {
			Convey("should success if full filled ", func() {
				dbMock1.ExpectBegin()
				dbMock1.ExpectCommit()
				tx, _ := db1.Beginx()
				_ = tx.Commit()
				So(dbMock1, test.ShouldBeFullFilled)
			})

			Convey("should fail if not full filled", func() {
				dbMock1.ExpectBegin()
				dbMock1.ExpectCommit()
				tx, _ := db1.Beginx()
				So(test.ShouldBeFullFilled(dbMock1), ShouldNotBeZeroValue)
				_ = tx.Commit()
			})
		})

		Convey("multiple instance on DB", func() {
			Convey("should success if full filled ", func() {
				dbMock1.ExpectBegin()
				dbMock1.ExpectCommit()

				dbMock2.ExpectBegin()
				dbMock2.ExpectCommit()

				tx, err := db1.Beginx()
				So(err, ShouldBeNil)
				_ = tx.Commit()

				tx2, _ := db2.Beginx()
				_ = tx2.Commit()

				So(dbMockList, test.ShouldBeFullFilled)
			})

			Convey("should fail if not full filled", func() {
				dbMock1.ExpectBegin()
				dbMock1.ExpectCommit()
				dbMock2.ExpectBegin()
				tx, _ := db1.Beginx()
				_ = tx.Commit()
				So(test.ShouldBeFullFilled(dbMockList), ShouldNotBeZeroValue)
			})
		})

		Convey("single instance on testify", func() {
			Convey("should success if full filled ", func() {
				m1.On("Test", "some").Return(nil).Twice()
				_ = m1.Test("some")
				_ = m1.Test("some")
				So(&m1.Mock, test.ShouldBeFullFilled)
			})

			Convey("should fail if not full filled", func() {
				m1.On("Test", "some").Return(nil).Twice()
				_ = m1.Test("some")
				So(test.ShouldBeFullFilled(&m1.Mock), ShouldNotBeZeroValue)
			})
		})

		Convey("multiple instance on testify", func() {
			Convey("should success if full filled ", func() {
				m1.On("Test", "some").Return(nil).Twice()
				_ = m1.Test("some")
				_ = m1.Test("some")
				m2.On("Do", "some").Return(nil).Once()
				_ = m2.Do("some")
				So(testifyMocks, test.ShouldBeFullFilled)
			})

			Convey("should fail if not full filled", func() {
				m1.On("Test", "some").Return(nil).Twice()
				_ = m1.Test("some")
				m2.On("Do", "some").Return(nil).Once()
				_ = m2.Do("some")
				So(test.ShouldBeFullFilled(testifyMocks), ShouldNotBeZeroValue)
			})
		})

		Convey("unknown should fail", func() {
			So(test.ShouldBeFullFilled("test"), ShouldNotBeZeroValue)
		})

		test.ResetMocks(testifyMocks...)
	})
}
