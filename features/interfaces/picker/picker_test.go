package picker_test

import (
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"

	"github.com/elmagician/kactus/features/interfaces/picker"
	"github.com/elmagician/kactus/internal/logger"
)

func init() {
	logger.Set(zap.NewNop())
}

func Example() {
	var scenarioCtx *godog.ScenarioContext

	pickerInstance := picker.New()

	// Create a step registering variables that will be reset on picker.Cancel
	scenarioCtx.Step("^I wish to pick key ([a-z]*) as value .*$", func(key, value string) error {
		pickerInstance.RegisterTemporaryVariable(key, value)
		return nil
	})

	// Create a step to retrieve a picked value.
	scenarioCtx.Step("^I wish to use value [a-z]*$", func(key string) error {
		value, ok := pickerInstance.Retrieve(key)
		if !ok {
			return picker.ErrNotPicked
		}
		fmt.Println(value) // do something with value.
		return nil
	})

	// Ensure picked value matches expected one.
	scenarioCtx.Step("^I wish to ensure value [a-z]* (equals|matches) .*$", pickerInstance.VerifiesPickedValue)

	// Ensure some picked values matches expected data  (using gherkin data table)
	scenarioCtx.Step("^I wish to ensure values are correct:$", pickerInstance.VerifiesPickedValues)

	// Registering persistent information before scenario. Those key will not be reset. To forget them,
	// call picker.ForgetPersistent
	scenarioCtx.BeforeScenario(func(sc *godog.Scenario) {
		pickerInstance.PersistVariable("test", true)
		pickerInstance.RegisterVariables(map[string]interface{}{"name": "sample", "simple": true}, picker.Persistent)
	})

	// Setup reset after steps (should not be a good idea in a normal use case ^^)
	scenarioCtx.AfterStep(func(st *godog.Step, err error) {
		pickerInstance.Reset()
	})

	// Forget test name after scenario
	scenarioCtx.AfterScenario(func(sc *godog.Scenario, err error) {
		pickerInstance.Forget("name")
	})

}

func TestUnit__Picker(t *testing.T) {

	Convey("I should be able to use Picker instance", t, func() {
		var cleanUp func()
		pick := picker.New()
		pick.Reset()

		Convey("to register values into Disposable scope", func() {
			pick.RegisterVariables(map[string]interface{}{"test": true, "vlod": "demort"}, picker.Disposable)
			pick.RegisterTemporaryVariable("tmp", 1)

			Convey("then retrieve it from Retrieve method", func() {
				val, ok := pick.Retrieve("test")
				So(ok, ShouldBeTrue)
				So(val, ShouldBeTrue)

				val, ok = pick.Retrieve("vlod")
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, "demort")

				val, ok = pick.Retrieve("tmp")
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, 1)
			})

			Convey("or forget them", func() {
				So(pick.Forget("tmp"), ShouldBeNil)
				_, ok := pick.Retrieve("tmp")
				So(ok, ShouldBeFalse)
			})

			cleanUp = func() {}
		})

		Convey("to register values into Persistent scope", func() {
			pick.RegisterVariables(map[string]interface{}{"test": true, "vlod": "demort"}, picker.Persistent)
			pick.PersistVariable("persistent", 2)

			Convey("then retrieve it from Retrieve method", func() {
				val, ok := pick.Retrieve("test")
				So(ok, ShouldBeTrue)
				So(val, ShouldBeTrue)

				val, ok = pick.Retrieve("vlod")
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, "demort")

				val, ok = pick.Retrieve("persistent")
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, 2)
			})

			Convey("or forget them", func() {
				pick.ForgetPersistent("persistent")
				_, ok := pick.Retrieve("persistent")
				So(ok, ShouldBeFalse)
			})

			cleanUp = func() {
				pick.ForgetPersistent("persistent")
				pick.ForgetPersistent("vol")
				pick.ForgetPersistent("test")
			}
		})

		pick.Reset()
		cleanUp()
	})
}
