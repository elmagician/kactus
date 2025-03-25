package fixtures

import (
	"errors"
	"strings"

	"github.com/cucumber/godog"
	"go.uber.org/zap"

	"github.com/elmagician/kactus/internal/picker"
)

const (
	// Variables indicates fixtures to load data
	// using store.
	Variables Kind = iota + 1

	// Postgres indicates fixtures to load SQL
	// data for a picked postgres db instance.
	Postgres

	// Pubsub indicates fixtures to load PUBSUB
	// events for a picked pubsub instance.
	Pubsub

	// API indicates fixtures to load API
	// definitions background data.
	API

	manifestTag = "@manifest"
)

var (
	// ErrInvalidInstance indicates that an object instance was expected but
	// did not matched expected types.
	ErrInvalidInstance = errors.New("invalid picked instance")

	// ErrUnknown indicates that on object should be available but
	// was not stored.
	ErrUnknown = errors.New("unknown instance")

	// ErrUnsupportedFixture indicates provided fixture kind
	// is not supported.
	ErrUnsupportedFixture = errors.New("unsupported fixture kind")
)

type (
	// Fixtures manages fixtures files to
	// help data initialisation
	// for test.
	//
	// It relies on picker.Picker manager.
	Fixtures struct {
		store    *picker.Store
		basePath string
	}

	// Kind of the loaded fixtures.
	Kind int
)

// New initiates a fixtures manager instance using provided store.
func New(store *picker.Store) *Fixtures {
	return &Fixtures{store: store}
}

// WithBasePath add a path element as prefix for any fixture file path.
func (fix *Fixtures) WithBasePath(path string) *Fixtures {
	fix.basePath = strings.TrimSuffix(path, "/") + "/"
	return fix
}

// NewTagLoader initializes tag parsing process to load
// Manifest from `@manifest` tag.
func (fix Fixtures) NewTagLoader(s *godog.Scenario) {
	log.Debug("parsing tags to find fixture manifest", zap.Reflect("tags", s.Tags))

	/*if err := fix.LoadFromTags(s.Tags); err != nil {
		log.Error("could not load fixtures from manifest tag", zap.Error(err))
	}/*/
}

// Reset resets fixture instance.
func (fix Fixtures) Reset() {
	ResetLog()
}

// getPath construct a fixture file loading path using basePath and provided path.
func (fix Fixtures) getPath(path string) string {
	return fix.basePath + strings.TrimPrefix(path, "/")
}
