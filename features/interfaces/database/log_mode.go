package database

import (
	"github.com/elmagician/kactus/internal/databases/postgres"
)

// Debug start debug logs.
// It will be removed when calling Reset.
func (Postgres) Debug() error {
	return postgres.Debug()
}

// DisableDebug stops debugging.
func (Postgres) DisableDebug() error {
	postgres.ResetLog()
	return nil
}
