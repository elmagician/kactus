package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/cucumber/godog"

	"github.com/elmagician/kactus/features/interfaces/picker"
	"github.com/elmagician/kactus/internal/databases/postgres"
	internalPicker "github.com/elmagician/kactus/internal/picker"
)

var (
	// ErrUnknown is raised when trying to operate on
	// an unpicked database.
	ErrUnknown = errors.New("unknown database")

	// ErrInvalidInstance is raised when trying to
	// operate on a non postgresql instance.
	ErrInvalidInstance = errors.New("expected Postgres instance")

	ErrUnexpectedData = errors.New("data should not exists in instance")
)

var (
	instance *Postgres
	pgOnce   sync.Once
)

type (
	// Postgres manages instance for kactus Postgres and
	// provides tools to manipulate and asserts on known
	// postgres databases.
	Postgres struct {
		store *internalPicker.Store
	}

	// PostgresInfo provides a structure to
	// register a new postgres instance to kactus.
	//
	// Key will be used to to pick instance.
	// `pg.` will be append to provided key when picking.
	//
	// DB indicates *sql.DB instance to store under
	// proved key.
	PostgresInfo struct {
		Key string
		DB  *sql.DB
	}
)

// NewPostgres initializes a postgres kactus manager for provided postgres
// database instances. It relies on a singleton pattern so successive call to
// NewPostgres will return the same instance
func NewPostgres(pickerInstance *picker.Picker, databases ...PostgresInfo) *Postgres {
	pgOnce.Do(func() {
		instance = &Postgres{
			store: pickerInstance.This(),
		}

		// Register provided databases
		for _, db := range databases {
			postgres.Database(db.Key, db.DB, pickerInstance.This())
		}
	})

	return instance
}

// AssertData asserts data exists in provided table for known
// db instance using provided where clauses.
func (pg *Postgres) AssertData(instance string, data *godog.Table) error {
	pgInstance, err := pg.getInstance(instance)
	if err != nil {
		return err
	}

	return pgInstance.AssertData(data)
}

// AssertDataDoesNotExists asserts data does not exists in provided table for
// known db instance using provided where clauses.
func (pg *Postgres) AssertDataDoesNotExists(instance string, data *godog.Table) error {
	err := pg.AssertData(instance, data)
	if errors.Is(err, postgres.ErrUnmatched) {
		return nil
	}

	if err == nil {
		return fmt.Errorf("%w %s", ErrUnexpectedData, instance)
	}

	return err
}

// StoreField retrieve existing field from postgres table and store them in
// disposable store under provided key.
func (pg *Postgres) StoreField(field, instance, table, key string, data *godog.Table) error {
	pgInstance, err := pg.getInstance(instance)
	if err != nil {
		return err
	}

	return pgInstance.PickField(pg.store, table, field, key, data)
}

// Reset resets postgres instance.
func (pg *Postgres) Reset() {
	postgres.ResetLog()
}

func (pg *Postgres) getInstance(instance string) (postgres.DB, error) {
	kind, localInstance, exists := pg.store.GetInstance(instance)
	if !exists {
		return postgres.DB{}, fmt.Errorf("%w: %s", ErrUnknown, instance)
	}

	if kind != internalPicker.Postgres {
		return postgres.DB{}, fmt.Errorf("%w", ErrInvalidInstance)
	}

	pgInstance, ok := localInstance.(postgres.DB)
	if !ok {
		return postgres.DB{}, fmt.Errorf("%w, got: %T", ErrInvalidInstance, pg)
	}

	return pgInstance, nil
}
