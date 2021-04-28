package types_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/elmagician/kactus/internal/test"
	"github.com/elmagician/kactus/internal/types"
)

func init() {
	types.NoLog()
}

func TestUnit_SetTyped(t *testing.T) {
	Convey("Given I want to bind string/string list to typed interface", t, func() {
		Convey("I should be able to retrieve string from single string value", func() {
			var actual string
			expected := "test"
			So(types.SetTyped(expected, &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve int from single string value", func() {
			var actual int
			expected := 666
			So(types.SetTyped("666", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve int8 from single string value", func() {
			var actual int8
			expected := 8
			So(types.SetTyped("8", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve int16 from single string value", func() {
			var actual int16
			expected := 666
			So(types.SetTyped("666", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve int32 from single string value", func() {
			var actual int32
			expected := 666
			So(types.SetTyped("666", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve int64 from single string value", func() {
			var actual int64
			expected := 666
			So(types.SetTyped("666", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve uint from single string value", func() {
			var actual uint
			expected := 666
			So(types.SetTyped("666", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve uint8 from single string value", func() {
			var actual uint8
			expected := 126
			So(types.SetTyped("126", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve uint16 from single string value", func() {
			var actual uint16
			expected := 666
			So(types.SetTyped("666", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve uint32 from single string value", func() {
			var actual uint32
			expected := 666
			So(types.SetTyped("666", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve uint64 from single string value", func() {
			var actual uint64
			expected := 666
			So(types.SetTyped("666", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve float32 from single string value", func() {
			var actual float32
			expected := 666.7
			So(types.SetTyped("666.7", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve float64 from single string value", func() {
			var actual float64
			expected := 666.7
			So(types.SetTyped("666.7", &actual), ShouldBeNil)
			So(actual, ShouldEqual, expected)
		})

		Convey("I should be able to retrieve bool from single string value", func() {
			var actual bool
			So(types.SetTyped("T", &actual), ShouldBeNil)
			So(actual, ShouldBeTrue)

			So(types.SetTyped("F", &actual), ShouldBeNil)
			So(actual, ShouldBeFalse)

			So(types.SetTyped("true", &actual), ShouldBeNil)
			So(actual, ShouldBeTrue)

			So(types.SetTyped("false", &actual), ShouldBeNil)
			So(actual, ShouldBeFalse)

		})

		Convey("I should have an error when provided value cannot be bind to int and expecting int", func() {
			var actual int
			So(types.SetTyped("ceci n'est pas un entier", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to int8 and expecting int8", func() {
			var actual int8
			So(types.SetTyped("ceci n'est pas un entier", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to int16 and expecting int16", func() {
			var actual int16
			So(types.SetTyped("ceci n'est pas un entier", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to int32 and expecting int32", func() {
			var actual int32
			So(types.SetTyped("ceci n'est pas un entier", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to int64 and expecting int64", func() {
			var actual int64
			So(types.SetTyped("ceci n'est pas un entier", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to uint and expecting uint", func() {
			var actual uint
			So(types.SetTyped("ceci n'est pas un entier", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to uint8 and expecting uint8", func() {
			var actual uint8
			So(types.SetTyped("ceci n'est pas un entier", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to uint16 and expecting uint16", func() {
			var actual uint16
			So(types.SetTyped("ceci n'est pas un entier", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to uint32 and expecting uint32", func() {
			var actual uint32
			So(types.SetTyped("ceci n'est pas un entier", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to uint64 and expecting uint64", func() {
			var actual uint64
			So(types.SetTyped("ceci n'est pas un entier", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to float32 and expecting float32", func() {
			var actual float32
			So(types.SetTyped("ceci n'est pas un float", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to float64 and expecting float64", func() {
			var actual float64
			So(types.SetTyped("ceci n'est pas un float", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error when provided value cannot be bind to bool and expecting bool", func() {
			var actual bool
			So(types.SetTyped("ceci n'est pas un bool", &actual), ShouldNotBeNil)
		})

		Convey("I should have an error if type is not supported", func() {
			var actual chan string
			So(types.SetTyped("any", &actual), ShouldBeLikeError, types.ErrUnsupportedType)
		})

	})
}
