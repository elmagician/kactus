package types_test

import (
	"testing"

	. "github.com/elmagician/kactus/internal/test"
	"github.com/elmagician/kactus/internal/types"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	types.NoLog()
}

func TestUnit_ToInterface(t *testing.T) {
	Convey("Given a string representing a matchable type", t, func() {
		Convey("should get string if type is not provided", func() {
			val, err := types.ToInterface("test")
			So(err, ShouldBeNil)
			matched, ok := val.(string)
			So(ok, ShouldBeTrue)
			So(matched, ShouldEqual, "test")
		})

		Convey("should get string when asked to", func() {
			val, err := types.ToInterface("test((string))")
			So(err, ShouldBeNil)
			matched, ok := val.(string)
			So(ok, ShouldBeTrue)
			So(matched, ShouldEqual, "test")
		})

		Convey("should get int when asked to", func() {
			val, err := types.ToInterface("12((int))")
			So(err, ShouldBeNil)
			matched, ok := val.(int64)
			So(ok, ShouldBeTrue)
			So(matched, ShouldEqual, 12)
		})

		Convey("should get float when asked to", func() {
			val, err := types.ToInterface("12.45((float))")
			So(err, ShouldBeNil)
			matched, ok := val.(float64)
			So(ok, ShouldBeTrue)
			So(matched, ShouldEqual, 12.45)
		})

		Convey("should get bool when asked to", func() {
			val, err := types.ToInterface("t((bool))")
			So(err, ShouldBeNil)
			matched, ok := val.(bool)
			So(ok, ShouldBeTrue)
			So(matched, ShouldBeTrue)
		})

		Convey("should get uuid when asked to", func() {
			test := uuid.New()
			val, err := types.ToInterface(test.String() + "((uuid))")
			So(err, ShouldBeNil)
			matched, ok := val.(uuid.UUID)
			So(ok, ShouldBeTrue)
			So(matched, ShouldEqual, test)
		})

		Convey("should be able to parse an array of interfaces", func() {
			val, err := types.ToInterface("[1((int)),2.45((float)),test]((array))")
			So(err, ShouldBeNil)
			matched, ok := val.([]interface{})
			So(ok, ShouldBeTrue)
			So(matched, ShouldBeEquivalent, []interface{}{int64(1), 2.45, "test"})
		})

		Convey("should have error when type asked could not be casted", func() {
			val, err := types.ToInterface("azetr((int))")
			So(err, ShouldNotBeNil)
			So(val, ShouldBeNil)

			val, err = types.ToInterface("azetr((float))")
			So(err, ShouldNotBeNil)
			So(val, ShouldBeNil)

			val, err = types.ToInterface("azetr((bool))")
			So(err, ShouldNotBeNil)
			So(val, ShouldBeNil)

			val, err = types.ToInterface("azetr((uuid))")
			So(err, ShouldNotBeNil)
			So(val, ShouldBeNil)

			val, err = types.ToInterface("azetr((int)),test((bool))((array))")
			So(err, ShouldNotBeNil)
			So(val, ShouldBeNil)
		})
	})
}

func TestUnit_AssertEqual(t *testing.T) {
	Convey("Given a string representing a matchable type and an interface", t, func() {
		Convey("should be able to check equality as strings", func() {
			ok, err := types.AssertEqual("test", "test")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			ok, err = types.AssertEqual("pas glop((string))", "glop")
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("should be able to check equality as int", func() {
			ok, err := types.AssertEqual("645((int))", int64(645))
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			ok, err = types.AssertEqual("666((int))", int64(645))
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("should be able to check equality as float", func() {
			ok, err := types.AssertEqual("645.448((number))", 645.448)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			ok, err = types.AssertEqual("124.45((float))", "11124")
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)

			ok, err = types.AssertEqual("666((number))", 645.77)
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("should be able to check equality as uuid", func() {
			testUUID := uuid.New()
			ok, err := types.AssertEqual(testUUID.String()+"((uuid))", testUUID)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			ok, err = types.AssertEqual(uuid.New().String()+"((uuid))", testUUID)
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("should be able to check equality as bool", func() {
			ok, err := types.AssertEqual("true((bool))", true)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			ok, err = types.AssertEqual("false((bool))", true)
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("should be able to check equality as array", func() {
			ok, err := types.AssertEqual("true((bool)),12((int)),test((array))", []interface{}{true, int64(12), "test"})
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			ok, err = types.AssertEqual("[true((bool)),24((int)),test]((array))", []interface{}{true, int64(12), "test"})
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("should have error when type asked could not be casted or interface does not match asked type", func() {
			ok, err := types.AssertEqual("azetr((int))", 125)
			So(err, ShouldNotBeNil)
			So(ok, ShouldBeFalse)

			ok, err = types.AssertEqual("12((int))", true)
			So(err, ShouldNotBeNil)
			So(ok, ShouldBeFalse)

			ok, err = types.AssertEqual("124.45((float))", make(chan int))
			So(err, ShouldNotBeNil)
			So(ok, ShouldBeFalse)

			ok, err = types.AssertEqual("true((bool))", 124.45)
			So(err, ShouldNotBeNil)
			So(ok, ShouldBeFalse)

			ok, err = types.AssertEqual("124.45((array))", 12)
			So(err, ShouldNotBeNil)
			So(ok, ShouldBeFalse)

			ok, err = types.AssertEqual(uuid.New().String()+"((uuid))", 12)
			So(err, ShouldNotBeNil)
			So(ok, ShouldBeFalse)

		})
	})
}
