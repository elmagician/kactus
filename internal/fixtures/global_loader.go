package fixtures

import (
	"io/ioutil"
	"strings"

	"github.com/cucumber/godog"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/elmagician/kactus/internal/databases/postgres"
	"github.com/elmagician/kactus/internal/picker"
	"github.com/elmagician/kactus/internal/pubsub/google"
)

// LoadData parses yaml data manifest to stores declared data
func (fix Fixtures) LoadData(fixturePath string) (err error) {
	var (
		manifest  picker.SeedManifest
		bManifest []byte
	)

	bManifest, err = ioutil.ReadFile(fix.getPath(fixturePath))
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(bManifest, &manifest); err != nil {
		return
	}

	return manifest.Load(fix.store)
}

// LoadYAMLManifest loads Manifest from path
func (fix Fixtures) LoadYAMLManifest(manifestPath string) (err error) {
	var (
		manifest  Manifest
		bManifest []byte
	)

	bManifest, err = ioutil.ReadFile(fix.getPath(manifestPath))
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(bManifest, &manifest); err != nil {
		return
	}

	return manifest.Load(fix)
}

func (fix Fixtures) LoadFromInstance(fixturePath string, instanceKey string) error {
	kind, instance, exists := fix.store.GetInstance(instanceKey)
	if !exists {
		return ErrUnknown
	}

	switch kind {
	case picker.Postgres:
		db, ok := instance.(postgres.DB)
		if !ok {
			tmpDB, ok := instance.(*postgres.DB)
			if !ok {
				return ErrInvalidInstance
			}

			db = *tmpDB
		}

		return db.SeedFromFile(fix.getPath(fixturePath), fix.store)
	case picker.GCP:
		gcp, ok := instance.(*google.Client)
		if !ok {
			return ErrInvalidInstance
		}

		return gcp.SeedFromFile(fix.getPath(fixturePath), fix.store)
	case picker.REST:
		return ErrUnsupportedFixture
	case picker.NoInstance, picker.Fixture:
		return ErrInvalidInstance
	default:
		return ErrInvalidInstance
	}
}

func (fix Fixtures) LoadFromTags(tags []*godog.Scenario) error {
	log.Debug("Loading from tags", zap.Reflect("tags", tags))

	for _, tag := range tags {
		log.Debug("analyzing: ", zap.Reflect("current tag", tag))

		splitted := strings.Split(tag.Name, ":")
		if len(splitted) < 2 || splitted[0] != manifestTag {
			log.Debug("tag did not match", zap.Reflect("splitted tag content", splitted))
			continue
		}

		log.Debug("loading manifest", zap.String("manifest path", splitted[1]))

		return fix.LoadYAMLManifest(splitted[1])
	}

	return nil
}
