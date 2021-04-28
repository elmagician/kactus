package picker_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/elmagician/kactus/internal/matchers"
	"github.com/elmagician/kactus/internal/picker"
	"github.com/elmagician/kactus/internal/types"
)

func init() {
	picker.NoLog()
	matchers.NoLog()
	types.NoLog()
}

func TestUnit_Store_Inject(t *testing.T) {
	Convey("Given I wish to inject stored values, ", t, func() {
		store := picker.NewStore()
		store.Pick("test", 124, picker.PersistentValue)
		store.Pick("pi", 3.14, picker.PersistentValue)
		store.Pick("testName", "storeInject", picker.DisposableValue)

		Convey("It should replace provided value if injection can proceed", func() {
			So(store.Inject("{{test}}"), ShouldEqual, 124)
			So(store.Inject("{{pi}}"), ShouldEqual, 3.14)
			So(store.Inject("{{testName}}"), ShouldEqual, "storeInject")
		})

		Convey("It should not inject into non string values", func() {
			So(store.Inject(124), ShouldEqual, 124)
			So(store.Inject(true), ShouldBeTrue)
		})

		Convey("It should not replace if string does not contain value to replace", func() {
			So(store.Inject("no value to replace here"), ShouldEqual, "no value to replace here")
			So(store.Inject("something"), ShouldEqual, "something")
		})

		Convey("It should not inject if asked value is not known in store", func() {
			So(store.Inject("{{unknown}}"), ShouldEqual, "{{unknown}}")
		})
	})
}

func TestUnit_Store_InjectAll(t *testing.T) {
	Convey("Given I wish to inject stored values, ", t, func() {
		store := picker.NewStore()
		store.Pick("test", 124, picker.PersistentValue)
		store.Pick("pi", 3.14, picker.PersistentValue)
		store.Pick("testName", "storeInject", picker.DisposableValue)

		Convey("It should replace provided value if injection can proceed", func() {
			So(store.InjectAll(
				"This test is called {{testName}}. It contains Pi {{pi}} value and a test {{test}} to 124."), ShouldEqual,
				"This test is called storeInject. It contains Pi 3.14 value and a test 124 to 124.")
		})

		Convey("It should not replace if string does not contain value to replace", func() {
			So(store.InjectAll("no value to replace here"), ShouldEqual, "no value to replace here")
			So(store.InjectAll("something"), ShouldEqual, "something")
		})

		Convey("It should not inject if asked value is not known in store", func() {
			So(store.InjectAll(
				"This should stay {{unknown}} cause it is {{notDeclared}}"), ShouldEqual,
				"This should stay {{unknown}} cause it is {{notDeclared}}")
		})
	})
}
