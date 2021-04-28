package definitions

import (
	"github.com/cucumber/messages-go/v10"
	"github.com/elmagician/godog"

	"github.com/elmagician/kactus/features/interfaces/picker"
)

// InstallPicker adds picker base definitions to known steps for provided godog ScenarioContext.
// It returns the picker instance if you need to customize it or enrich known steps.
// You can safely ignore return value.
//
// Provided steps:
// - (?:I )?want(?:ing)? to generate uuids (.*) => generate a new uuid for each provided key
//   I want to generate uuids bookID, userID, mouseID
// - (?:I )?set(?:ting)? variable ([a-zA-Z0-9]*) to (.*) => store provided value under variable name
//   I set variable X to 123i
// - (?:I )?want to assert picked variables matches: => ensures provided key value matches condition
// (cf: picker.VerifiesPickedValues)
//   I want to assert picked variables matches:
//   | key | matcher | value         |
//   | X   | =       | 123i((reel))  |
func InstallPicker(s *godog.ScenarioContext, pickerInstance *picker.Picker) {
	pickerInstance.Reset() // nolint: errcheck

	s.BeforeStep(pickerInstance.BeforeStepReplacer)
	s.BeforeScenario(func(_ *messages.Pickle) {
		// nolint: errcheck
		pickerInstance.Reset() // error always nil
	})

	// Pick UUID as id name
	s.Step(`^(?:I )?want(?:ing)? to generate uuids (.*)$`, pickerInstance.CreateAndPickUUIDS)
	s.Step(`^(?:I )?set(?:ting)? variable ([a-zA-Z0-9]*) to (.*)$`, pickerInstance.PickString)
	s.Step(`^(?:I )?want to assert picked variables matches:$`, pickerInstance.VerifiesPickedValues)

	// Debug picker
	s.Step(`^(?:I )?want to debug picker$`, pickerInstance.Debug)
	s.Step(`^(?:I )?want to stop debugging picker$`, pickerInstance.DisableDebug)
}
