package picker

import (
	"github.com/cucumber/messages-go/v10"

	internalPicker "github.com/elmagician/kactus/internal/picker"
)

// BeforeStepReplacer provides a replacer function to use in godog BeforeStep hook.
func (picker *Picker) BeforeStepReplacer(step *messages.Pickle_PickleStep) {
	internalPicker.NewReplacer(picker.store)(step)
}
