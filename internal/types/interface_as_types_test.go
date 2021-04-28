package types_test

import (
	"strconv"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/elmagician/kactus/internal/types"
)

func TestUnit_AsInt64(t *testing.T) {
	Convey("When I try to get interface as int 64", t, func() {
		var expectedInt64, result int64
		var ok bool
		Convey("should success", func() {
			Convey("with int64", func() {
				expectedInt64 = fake.Int64()
				result, ok = types.AsInt64(expectedInt64)
			})

			Convey("with int", func() {
				val := int(fake.Int16())
				expectedInt64 = int64(val)
				result, ok = types.AsInt64(val)
			})

			Convey("with int8", func() {
				val := fake.Int8()
				expectedInt64 = int64(val)
				result, ok = types.AsInt64(val)
			})

			Convey("with int16", func() {
				val := fake.Int16()
				expectedInt64 = int64(val)
				result, ok = types.AsInt64(val)
			})

			Convey("with int32", func() {
				val := fake.Int32()
				expectedInt64 = int64(val)
				result, ok = types.AsInt64(val)
			})

			Convey("with string", func() {
				val := strconv.Itoa(int(fake.Int16()))
				expectedInt64, _ = strconv.ParseInt(val, 10, 64)
				result, ok = types.AsInt64(val)
			})

			So(ok, ShouldBeTrue)
			So(result, ShouldEqual, expectedInt64)
		})

		Convey("should fail", func() {
			Convey("if invalid string", func() {
				result, ok = types.AsInt64("wrong value")
			})

			Convey("if invalid given type", func() {
				result, ok = types.AsInt64(true)
			})

			So(ok, ShouldBeFalse)
			So(result, ShouldEqual, 0)
		})
	})
}

func TestUnit_AsFloat64(t *testing.T) {
	Convey("When I try to get interface as float 64", t, func() {
		var expectedFloat64, result float64
		var ok bool
		Convey("should success", func() {
			Convey("with float64", func() {
				expectedFloat64 = fake.Float64()
				result, ok = types.AsFloat64(expectedFloat64)
			})

			Convey("with float32", func() {
				val := fake.Float32()
				expectedFloat64 = float64(val)
				result, ok = types.AsFloat64(val)
			})

			Convey("with int64", func() {
				val := int(fake.Int16())
				expectedFloat64 = float64(val)
				result, ok = types.AsFloat64(val)
			})

			Convey("with int", func() {
				val := int(fake.Int16())
				expectedFloat64 = float64(val)
				result, ok = types.AsFloat64(val)
			})

			Convey("with int8", func() {
				val := fake.Int8()
				expectedFloat64 = float64(val)
				result, ok = types.AsFloat64(val)
			})

			Convey("with int16", func() {
				val := fake.Int16()
				expectedFloat64 = float64(val)
				result, ok = types.AsFloat64(val)
			})

			Convey("with int32", func() {
				val := fake.Int32()
				expectedFloat64 = float64(val)
				result, ok = types.AsFloat64(val)
			})

			Convey("with string", func() {
				val := strconv.Itoa(int(fake.Int16()))
				expectedFloat64, _ = strconv.ParseFloat(val, 64)
				result, ok = types.AsFloat64(val)
			})

			So(ok, ShouldBeTrue)
			So(result, ShouldEqual, expectedFloat64)
		})

		Convey("should fail", func() {
			Convey("if invalid string", func() {
				result, ok = types.AsFloat64("wrong value")
			})

			Convey("if invalid given type", func() {
				result, ok = types.AsFloat64(true)
			})

			So(ok, ShouldBeFalse)
			So(result, ShouldEqual, 0)
		})
	})
}

func TestUnit_AsBool(t *testing.T) {

	Convey("When I try to get val as bool", t, func() {
		Convey("should success", func() {
			Convey("with bool", func() {
				result, ok := types.AsBool(true)

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, true)

				result, ok = types.AsBool(false)

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, false)
			})

			Convey("with int64", func() {
				result, ok := types.AsBool(int64(0))

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, false)

				result, ok = types.AsBool(int64(1))

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, true)
			})

			Convey("with int", func() {
				result, ok := types.AsBool(0)

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, false)

				result, ok = types.AsBool(1)

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, true)
			})

			Convey("with int8", func() {
				result, ok := types.AsBool(int8(0))

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, false)

				result, ok = types.AsBool(int8(1))

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, true)
			})

			Convey("with int16", func() {
				result, ok := types.AsBool(int16(0))

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, false)

				result, ok = types.AsBool(int16(1))

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, true)
			})

			Convey("with int32", func() {
				result, ok := types.AsBool(int32(0))

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, false)

				result, ok = types.AsBool(int32(1))

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, true)
			})

			Convey("with string", func() {
				result, ok := types.AsBool("false")

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, false)

				result, ok = types.AsBool("true")

				So(ok, ShouldBeTrue)
				So(result, ShouldEqual, true)
			})
		})

		Convey("should fail", func() {
			var result, ok bool
			Convey("if invalid string", func() {
				result, ok = types.AsBool("hello")
			})

			Convey("if invalid given type", func() {
				result, ok = types.AsBool([]int{})
			})

			So(ok, ShouldBeFalse)
			So(result, ShouldEqual, false)
		})
	})
}

func TestUnit_AsString(t *testing.T) {
	Convey("When I try to get val as string should success", t, func() {
		var expectedString, result string
		var ok bool

		Convey("with string", func() {
			expectedString = fake.Word()

			result, ok = types.AsString(expectedString)
		})

		Convey("with bool", func() {
			val := fake.Bool()
			expectedString = strconv.FormatBool(val)

			result, ok = types.AsString(val)
		})

		Convey("with int32", func() {
			val := 10
			expectedString = strconv.Itoa(val)

			result, ok = types.AsString(val)
		})

		So(result, ShouldEqual, expectedString)
		So(ok, ShouldBeTrue)
	})

}
