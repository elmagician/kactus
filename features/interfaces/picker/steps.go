package picker

import (
	"errors"
	"strings"

	internalPicker "github.com/elmagician/kactus/internal/picker"
	"github.com/google/uuid"
)

var ErrNotPicked = errors.New("unpicked value")

// CreateAndPickUUIDS generates and stores an UUID for each provided keys (`, ` separated) as disposable value.
// It will be forgotten on reset.
func (picker *Picker) CreateAndPickUUIDS(idNames string) error {
	ids := strings.Split(idNames, ",")
	for _, id := range ids {
		picker.store.Pick(strings.TrimSpace(id), uuid.New(), internalPicker.DisposableValue)
	}

	return nil // return nil to be usable in godog Steps
}

// PickString stores provided string as disposable value.
func (picker *Picker) PickString(varName, val string) error {
	picker.store.Pick(varName, val, internalPicker.DisposableValue)
	return nil // return nil to be usable in godog Steps
}

// Forget disposable value.
func (picker *Picker) Forget(key string) error {
	picker.store.Del(key, internalPicker.DisposableValue)
	return nil
}
