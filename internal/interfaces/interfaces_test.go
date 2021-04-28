package interfaces_test

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/elmagician/kactus/internal/interfaces"
	. "github.com/elmagician/kactus/internal/test"
)

func init() {
	interfaces.NoLog()
}

func TestUnit_GetFieldFromPath(t *testing.T) {
	Convey("Given a json object interface", t, func() {
		ptrExpectedVal := map[string]string{"test": "ok from pointer"}
		var testNil *map[string]interface{}

		test := map[string]interface{}{
			"test": "vale",
			"object": map[string]interface{}{
				"some": "val",
				"to":   "test",
			},
			"array": [8]string{"test", "1", "deux", "666", "allo", "abra", "racourcix", "kkk"},
			"slice": []int{1, 2, 4, 8, 36, 666, 77, 589, 7, 231},
			"item_slice": []map[string]interface{}{
				{"test": "key", "lambda": 666},
				{"test": "val", "lambda": 777},
			},
			"ptr": &ptrExpectedVal,
			"nil": testNil,
			"struct": testStruct{
				Val: "kinder",
				Ok:  true,
			},
		}

		Convey("I should be able to get value for basic keys (non struct, non array/slice, non objects, non ptr", func() {
			val, ok := interfaces.GetFieldFromPath(test, "test")
			So(ok, ShouldBeTrue)
			So(val.Interface(), ShouldEqual, "vale")
		})

		Convey("I should be able to get value from items behind pointers", func() {
			val, ok := interfaces.GetFieldFromPath(test, "ptr.test")
			So(ok, ShouldBeTrue)
			So(val.Interface(), ShouldEqual, ptrExpectedVal["test"])
		})

		Convey("I should be able to get value from value in slice", func() {
			val, ok := interfaces.GetFieldFromPath(test, "slice.5")
			So(ok, ShouldBeTrue)
			So(val.Interface(), ShouldEqual, 666)
		})

		Convey("I should be able to get value behind slice/array", func() {
			val, ok := interfaces.GetFieldFromPath(test, "item_slice.0.test")
			So(ok, ShouldBeTrue)
			So(val.Interface(), ShouldEqual, "key")

			val, ok = interfaces.GetFieldFromPath(test, "item_slice.1.lambda")
			So(ok, ShouldBeTrue)
			So(val.Interface(), ShouldEqual, 777)
		})

		Convey("I should be able to get value from value in arrays", func() {
			val, ok := interfaces.GetFieldFromPath(test, "array.7")
			So(ok, ShouldBeTrue)
			So(val.Interface(), ShouldEqual, "kkk")
		})

		Convey("I should be able to get value from struct", func() {
			val, ok := interfaces.GetFieldFromPath(test, "struct.ok")
			So(ok, ShouldBeTrue)
			So(val.Interface(), ShouldBeTrue)
		})

		// Convey("I should have a panic if I pass a non integer path for slice/array", func() {
		// 	So(func() { interfaces.GetFieldFromPath(test, "array.azerty") }, ShouldPanicWith, interfaces.ErrInvalidPathElement.Error())
		// })

		Convey("I should have a panic if index is out of range on slice/array", func() {
			So(func() { interfaces.GetFieldFromPath(test, "array.152") }, ShouldPanic)
		})

		Convey("I should be able to know that value does not exists", func() {
			_, ok := interfaces.GetFieldFromPath(test, "nofield")
			So(ok, ShouldBeFalse)
			_, ok = interfaces.GetFieldFromPath(test, "test.test")
			So(ok, ShouldBeFalse)
			_, ok = interfaces.GetFieldFromPath(test, "nil.test")
			So(ok, ShouldBeFalse)
		})

	})
}

func TestUnit_GenerateFieldList(t *testing.T) {
	Convey("I should be able to list all path from interface object", t, func() {
		testJSON := []byte(`[{"key": {"path": [ "test","alpha","pas voulue"]}}]`)
		var asInterface interface{}
		So(json.Unmarshal(testJSON, &asInterface), ShouldBeNil)

		expectedPaths := []string{"0.key.path.0", "0.key.path.1", "0.key.path.2"}
		So(interfaces.GenerateFieldList(asInterface), ShouldBeEquivalent, expectedPaths)
	})
}

type testStruct struct {
	Val  string
	Kind *int
	Ok   bool
}
