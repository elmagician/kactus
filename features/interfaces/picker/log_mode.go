package picker

import (
	internalPicker "github.com/elmagician/kactus/internal/picker"
)

// Debug start debug logs for picker.
// It will be removed when calling Reset.
func (picker *Picker) Debug() {
	internalPicker.Debug()
}

// DisableDebug stops debugging.
func (picker *Picker) DisableDebug() {
	internalPicker.ResetLog()
}
