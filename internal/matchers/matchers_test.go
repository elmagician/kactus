package matchers_test

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/elmagician/kactus/internal/matchers"
	. "github.com/elmagician/kactus/internal/test"
	"github.com/elmagician/kactus/internal/types"
)

func init() {
	matchers.NoLog()
	types.NoLog()
}

func TestUnit_Match(t *testing.T) {
	Convey("Given a regex", t, func() {
		re := "^[a-z]*$"

		Convey("should validate value matching regex", func() {
			ok, err := matchers.Match("abcdefg", re)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			var test interface{}
			test = "test"

			ok, err = matchers.Match(reflect.ValueOf(test), re)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("should not validate unmatching values", func() {
			ok, err := matchers.Match("azr1gfdg23456aaazeds", re)
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("should not validate interface that cannot be resolve as string", func() {
			ok, err := matchers.Match(1234, re)
			So(err, ShouldBeLikeError, matchers.ErrInvalidArgument)
			So(ok, ShouldBeFalse)
		})

		Convey("should not validate uncompilable regex", func() {
			re = "[a-z*"
			ok, err := matchers.Match("1234", re)
			So(err, ShouldBeLikeError, matchers.ErrInvalidArgument)
			So(ok, ShouldBeFalse)
		})
	})
}

func TestUnit_Contain(t *testing.T) {
	Convey("Given a substring", t, func() {
		shouldContain := "test"

		Convey("should validate if contain substring", func() {
			ok, err := matchers.Contain("Ceci est un test", shouldContain)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			var test interface{}
			test = "This is a test"
			ok, err = matchers.Contain(reflect.ValueOf(test), shouldContain)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("should not validate if does not contain substring", func() {
			ok, err := matchers.Contain("Ceci en'est pas", shouldContain)
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("should not validate interface that cannot be resolve as string", func() {
			ok, err := matchers.Contain(415683, shouldContain)
			So(err, ShouldBeLikeError, matchers.ErrInvalidArgument)
			So(ok, ShouldBeFalse)
		})
	})
}

func TestUnit_Equal(t *testing.T) {
	Convey("Given an interface and a string", t, func() {
		Convey("should equal when interface == value parsed from string", func() {
			ok, err := matchers.Equal(int64(123), "123((int))")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			ok, err = matchers.Equal(float64(123.666), "123.666((number))")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			ok, err = matchers.Equal("123.666", "123.666")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			test := "a"
			ok, err = matchers.Equal(reflect.ValueOf(test), "a")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("should not validate when interface != value parsed from string", func() {
			ok, err := matchers.Equal(int64(123), "985((int))")
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("should not validate when value cannot be resolved to interface", func() {
			ok, err := matchers.Equal(123, "123")
			So(err, ShouldBeLikeError, types.ErrUnmatchedType)
			So(ok, ShouldBeFalse)

			ok, err = matchers.Equal(123, "azzertty((int))")
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})
	})
}

func TestUnit_In(t *testing.T) {
	Convey("Given an interface and a string", t, func() {
		Convey("should be true when interface is in slice described by string `,` sep", func() {
			So(matchers.In(int64(123), "1((int)),4((int32)),123((int)),66((int64))"), ShouldBeTrue)
			So(matchers.In("test", "test,4,alpha,roméo"), ShouldBeTrue)
		})

		Convey("should not validate when value is not in list", func() {
			So(matchers.In(123, "so,fun"), ShouldBeFalse)
		})
	})
}

func TestUnit_LenEqual(t *testing.T) {
	Convey("Given an interface and a string", t, func() {
		Convey("should be true when interface len equals int parse from string", func() {
			ok, err := matchers.LenEqual(GetStringOfLength(String, 101), "101")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			ok, err = matchers.LenEqual(GetStringOfLength(String, 1001), "1001")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			ok, err = matchers.LenEqual([]int{0, 5, 7, 8, 9, 10, 4, 6, 7}, "9")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("should panic if element is not of Len able type", func() {
			So(func() { _, _ = matchers.LenEqual(123, "9") }, ShouldPanic)
		})

		Convey("should not validate if provided expectation is not an int", func() {
			ok, err := matchers.LenEqual([]int{0, 5, 7}, "1,2,3")
			So(err, ShouldBeLikeError, matchers.ErrInvalidArgument)
			So(ok, ShouldBeFalse)
		})

		Convey("should not validate when length of interface is not equal to provided one", func() {
			ok, err := matchers.LenEqual(GetStringOfLength(String, 999), "666")
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})
	})
}

func TestUnit_OfType(t *testing.T) {
	Convey("Given an interface and a type", t, func() {
		Convey("should be true if interface match provided type", func() {
			So(matchers.OfType(int64(123), "int64"), ShouldBeTrue)
			So(matchers.OfType(float64(123.666), "float64"), ShouldBeTrue)
			So(matchers.OfType("123.666", "string"), ShouldBeTrue)
			test := "a"
			So(matchers.OfType(reflect.ValueOf(test), "string"), ShouldBeTrue)
		})

		Convey("should not validate when interface does not match type", func() {
			So(matchers.OfType(123, "string"), ShouldBeFalse)
		})

	})
}

func TestUnit_IsDefined(t *testing.T) {
	Convey("Given an interface", t, func() {
		Convey("should validate defined value (string != undefined, non null ptr,...)", func() {
			So(matchers.IsDefined(int64(123)), ShouldBeTrue)
			So(matchers.IsDefined(float64(123.666)), ShouldBeTrue)
			So(matchers.IsDefined("123.666"), ShouldBeTrue)
			test := "a"
			So(matchers.IsDefined(reflect.ValueOf(test)), ShouldBeTrue)
		})

		Convey("should not validate when interface is not defined", func() {
			So(matchers.IsDefined("undefined"), ShouldBeFalse)
			So(matchers.IsDefined(nil), ShouldBeFalse)
			var test *string = nil
			var test2 *int = nil
			So(matchers.IsDefined(reflect.ValueOf(test)), ShouldBeFalse)
			So(matchers.IsDefined(reflect.ValueOf(test2)), ShouldBeFalse)
		})

	})
}

func TestUnit_NotZero(t *testing.T) {
	Convey("Given an interface", t, func() {
		Convey("should validate non zero value (string != undefined, non null ptr,...)", func() {
			So(matchers.NotZero(int64(123)), ShouldBeTrue)
			So(matchers.NotZero(float64(123.666)), ShouldBeTrue)
			So(matchers.NotZero("123.666"), ShouldBeTrue)
			test := "a"
			So(matchers.NotZero(reflect.ValueOf(test)), ShouldBeTrue)
		})

		Convey("should not validate when interface is zero value", func() {
			So(matchers.NotZero("undefined"), ShouldBeFalse)
			So(matchers.NotZero(0), ShouldBeFalse)
			So(matchers.NotZero(""), ShouldBeFalse)
			So(matchers.NotZero(false), ShouldBeFalse)
			So(matchers.NotZero(nil), ShouldBeFalse)
			var test *string = nil
			So(matchers.NotZero(reflect.ValueOf(test)), ShouldBeFalse)
		})

	})
}

func TestUnit_Assert(t *testing.T) {
	Convey("given items", t, func() {
		Convey("should be able to call equal", func() {
			So(matchers.Assert("", int64(10), "10((int))"), ShouldBeNil)
			So(matchers.Assert("==", int64(10), "10((int))"), ShouldBeNil)
			So(matchers.Assert("eq", int64(10), "10((int))"), ShouldBeNil)
			So(matchers.Assert("equal", int64(10), "10((int))"), ShouldBeNil)
		})

		Convey("should be able to call match", func() {
			So(matchers.Assert("=~", "12azer", ".*"), ShouldBeNil)
			So(matchers.Assert("match", "12azer", ".*"), ShouldBeNil)
		})

		Convey("should be able to call contain", func() {
			So(matchers.Assert("contain", "12azer", "aze"), ShouldBeNil)
		})

		Convey("should be able to call defined", func() {
			So(matchers.Assert("defined", "12azer", ""), ShouldBeNil)
		})

		Convey("should be able to call not zero", func() {
			So(matchers.Assert("not zero", 0, ""), ShouldBeLikeError, matchers.ErrUnmatched)
		})

		Convey("should be able to call not of type", func() {
			So(matchers.Assert("type", 0, "int"), ShouldBeNil)
		})

		Convey("should be able to call length equals", func() {
			So(matchers.Assert("length equals", "012345678", "9"), ShouldBeNil)
		})

		Convey("should be able to call in", func() {
			So(matchers.Assert("in", "9", "0,1,5,6,7,9"), ShouldBeNil)
		})

		Convey("should have UndefinedMethod if asserter does not exists", func() {
			So(matchers.Assert("not implemented", 0, "int"), ShouldBeLikeError, matchers.ErrUndefinedMethod)
		})

		Convey("should be able to retrieve method error if exists", func() {
			So(matchers.Assert("equals", 0, "int"), ShouldBeError)
		})
	})
}
