package fixtures

import (
	"sync"

	"github.com/elmagician/kactus/features/interfaces/picker"
	"github.com/elmagician/kactus/internal/fixtures"
)

var (
	// ErrInvalidInstance indicates that an object instance was expected but
	// did not matched expected types.
	ErrInvalidInstance = fixtures.ErrInvalidInstance

	// ErrUnknown indicates that an error does not exists
	ErrUnknown = fixtures.ErrUnknown

	// ErrUnsupportedFixture indicates provided fixture kind
	// is not supported.
	ErrUnsupportedFixture = fixtures.ErrUnsupportedFixture
)

var (
	instance *Fixtures
	once     sync.Once
)

// Fixtures manages fixtures files to
// help data initialisation
// for test.
//
// It relies on picker.Picker manager.
type Fixtures struct {
	fixtures *fixtures.Fixtures
}

// New initiates a fixtures manager instance using provided picker. It relies on
// a singleton pattern so successive call to New will return the same instance.
func New(pickerInstance *picker.Picker) *Fixtures {
	once.Do(func() {
		instance = &Fixtures{fixtures: fixtures.New(pickerInstance.This())}
	})

	return instance
}

// WithBasePath add a path element as prefix for any fixture file path.
func (fix *Fixtures) WithBasePath(path string) *Fixtures {
	fix.fixtures.WithBasePath(path)
	return fix
}

// Load loads fixtures using provided Kind and path.
// Parameter instancePickedKey is required for some Kind.
func (fix Fixtures) Load(fixturePath string, instancePickedKey string) error {
	if instancePickedKey == "" {
		return fix.fixtures.LoadData(fixturePath)
	}

	return fix.fixtures.LoadFromInstance(fixturePath, instancePickedKey)
}

// LoadManifest loads Manifest from path.
func (fix Fixtures) LoadManifest(manifestPath string) error {
	return fix.fixtures.LoadYAMLManifest(manifestPath)
}

// Reset resets fixture instance.
func (fix Fixtures) Reset() {
	fix.fixtures.Reset()
}
