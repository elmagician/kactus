package picker

import (
	"go.uber.org/zap"
)

// SeedManifest represents the YAML structure to
// load data into store.
type SeedManifest struct {
	Persistent map[string]interface{} `yaml:"persistent"`
	Disposable map[string]interface{} `yaml:"disposable"`
}

// Load parse SeedManifest and stores provided data.
func (s *SeedManifest) Load(store *Store) error {
	for key, val := range s.Persistent {
		log.Debug("picking in persistent", zap.String("key", key), zap.Reflect("val", val))
		store.Pick(key, val, PersistentValue)
	}

	for key, val := range s.Disposable {
		log.Debug("picking in disposable", zap.String("key", key), zap.Reflect("val", val))
		store.Pick(key, val, DisposableValue)
	}

	return nil
}
