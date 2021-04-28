package fixtures

import (
	"fmt"

	"github.com/elmagician/kactus/internal/databases/postgres"
	"github.com/elmagician/kactus/internal/picker"
	"github.com/elmagician/kactus/internal/pubsub/google"
)

// LoadPostgres applies sql seed fixtures to picked database name.
func (fix Fixtures) LoadPostgres(db, fixturePath string) error {
	kind, pgDb, exists := fix.store.GetInstance(db)
	if !exists {
		return fmt.Errorf("%w: expected database: %s to be known", ErrUnknown, db)
	}

	if kind != picker.Postgres {
		return fmt.Errorf("%w: expected: %T does not match: %T", ErrInvalidInstance, postgres.DB{}, pgDb)
	}

	postgresLoader, ok := pgDb.(postgres.DB)
	if !ok {
		postgresLoaderPointer, ok := pgDb.(*postgres.DB)
		if !ok {
			return fmt.Errorf("%w: expected: %T does not match: %T", ErrInvalidInstance, postgres.DB{}, pgDb)
		}

		postgresLoader = *postgresLoaderPointer
	}

	return postgresLoader.SeedFromFile(fix.getPath(fixturePath), fix.store)
}

// LoadGooglePubsub applies yaml pubsub manifest seed fixtures to picked google instance name.
func (fix Fixtures) LoadGooglePubsub(gcp, fixturePath string) error {
	kind, googlePubsub, exists := fix.store.GetInstance(gcp)
	if !exists {
		return fmt.Errorf("%w: expected google pubsub: %s to be known", ErrUnknown, gcp)
	}

	if kind != picker.GCP {
		return fmt.Errorf("%w: expected: %T does not match: %T", ErrInvalidInstance, google.Client{}, googlePubsub)
	}

	googlePubsubLoader, ok := googlePubsub.(*google.Client)
	if !ok {
		return fmt.Errorf("%w: expected: %T does not match: %T", ErrInvalidInstance, google.Client{}, googlePubsub)
	}

	return googlePubsubLoader.SeedFromFile(fix.getPath(fixturePath), fix.store)
}
