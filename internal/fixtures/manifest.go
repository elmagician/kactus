package fixtures

import (
	"go.uber.org/zap"
)

type (
	// WithInstanceManifest represents a manifest part
	// witch link a fixture file to a specific instance.
	WithInstanceManifest struct {
		Path     string `yaml:"path"`
		Instance string `yaml:"pickedInstanceKey"`
	}

	// WithKindManifest represents manifest part
	// providing fixtures requiring an instance kind.
	WithKindManifest struct {
		Kind     string `yaml:"kind"`
		Path     string `yaml:"path"`
		Instance string `yaml:"pickedInstanceKey"`
	}

	// Manifest represents a YAML fixture manifest.
	Manifest struct {
		Variables []string                `yaml:"variables"`
		API       []*WithInstanceManifest `yaml:"api"`
		Pubsub    []*WithKindManifest     `yaml:"pubsub"`
		Database  []*WithKindManifest     `yaml:"database"`
	}
)

// Load loads manifest using provided fixture instance.
func (m *Manifest) Load(fixtures Fixtures) error {
	log.Debug("loading manifest", zap.Reflect("manifest", m))

	// Loading variables first
	for _, dataManifestPath := range m.Variables {
		if err := fixtures.LoadData(dataManifestPath); err != nil {
			return err
		}
	}

	for _, dbManifest := range m.Database {
		switch dbManifest.Kind {
		case "postgres", "pg":
			log.Debug("loading postgres manifest")

			if err := fixtures.LoadPostgres(dbManifest.Instance, dbManifest.Path); err != nil {
				return err
			}
		}
	}

	for _, pubsubManifest := range m.Pubsub {
		switch pubsubManifest.Kind {
		case "gcp", "google":
			log.Debug("loading gcp pubsub manifest")

			if err := fixtures.LoadGooglePubsub(pubsubManifest.Instance, pubsubManifest.Path); err != nil {
				return err
			}
		}
	}

	return nil
}
