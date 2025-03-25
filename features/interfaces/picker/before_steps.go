package picker

import (
	"github.com/cucumber/godog"

	internalPicker "github.com/elmagician/kactus/internal/picker"
)

// BeforeStepReplacer provides a replacer function to use in godog BeforeStep hook.
func (picker *Picker) BeforeStepReplacer(step *godog.Step) {
	internalPicker.NewReplacer(picker.store)(step)
}
