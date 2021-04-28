package pubsub

import (
	"github.com/elmagician/kactus/internal/pubsub/google"
)

// Debug start debug logs.
// It will be removed when calling Reset.
func (Google) Debug() error {
	return google.Debug()
}

// DisableDebug stops debugging.
func (Google) DisableDebug() error {
	google.ResetLog()
	return nil
}
