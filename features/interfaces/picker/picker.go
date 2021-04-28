// package picker provides simple function to
// create steps
package picker

import (
	"sync"

	"github.com/elmagician/kactus/internal/interfaces"
	"github.com/elmagician/kactus/internal/matchers"
	internalPicker "github.com/elmagician/kactus/internal/picker"
	"github.com/elmagician/kactus/internal/types"
)

const (
	// Persistent exposes internal picker.PersistentValue to
	// indicate values should be stored in the Persistent scope.
	Persistent = internalPicker.PersistentValue

	// Disposable exposes internal picker.DisposableValue to
	// indicate values should be stored in the Disposable scope.
	Disposable = internalPicker.DisposableValue
)

var (
	instance *Picker
	once     sync.Once
)

// Picker is a key/value store to manage
// variable injection through godog steps.
type Picker struct {
	store *internalPicker.Store
}

// New initialize a Picker instance.
func New() *Picker {
	once.Do(func() {
		instance = &Picker{
			store: internalPicker.NewStore(),
		}
	})

	return instance
}

// RegisterVariables adds variables to known values using provided scopes.
func (picker *Picker) RegisterVariables(variables map[string]interface{}, scope internalPicker.DataScope) {
	for k, v := range variables {
		picker.store.Pick(k, v, scope)
	}
}

// RegisterTemporaryVariable adds variable to known values up to instance reset.
func (picker *Picker) RegisterTemporaryVariable(key string, value interface{}) {
	picker.store.Pick(key, value, internalPicker.DisposableValue)
}

// PersistVariable adds variable to known values until suite is done or value is manually forgotten.
func (picker *Picker) PersistVariable(key string, value interface{}) {
	picker.store.Pick(key, value, internalPicker.PersistentValue)
}

// ForgetPersistent forgets persisted value.
func (picker *Picker) ForgetPersistent(key string) {
	picker.store.Del(key, internalPicker.PersistentValue)
}

// Retrieve value if exists
func (picker *Picker) Retrieve(key string) (interface{}, bool) {
	return picker.store.Get(key)
}

// This return internal store to pass to other features.
func (picker *Picker) This() *internalPicker.Store {
	return picker.store
}

// Reset picker
func (picker *Picker) Reset() {
	types.Reset()
	matchers.Reset()
	internalPicker.Reset()
	interfaces.ResetLog()
	picker.store.Reset()
}
