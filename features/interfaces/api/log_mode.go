package api

import (
	"github.com/elmagician/kactus/internal/api"
)

// Debug start debug logs.
// It will be removed when calling Reset.
func (*Client) Debug() {
	api.Debug()
}

// DisableDebug stops debugging.
func (*Client) DisableDebug() {
	api.ResetLog()
}
