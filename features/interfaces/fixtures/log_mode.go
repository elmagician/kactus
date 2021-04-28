package fixtures

import (
	"github.com/elmagician/kactus/internal/fixtures"
)

// Debug start debug logs for picker.
// It will be removed when calling Reset.
func (Fixtures) Debug() error {
	return fixtures.Debug()
}

// DisableDebug stops debugging.
func (Fixtures) DisableDebug() error {
	fixtures.ResetLog()
	return nil
}
