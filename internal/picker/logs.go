// package picker provided structures and method
// to manage injected variable system for Godog.
//
// It allows an user to pick value into a store and
// inject them in steps through a variable replacement.
package picker

import (
	"go.uber.org/zap"

	"github.com/elmagician/kactus/internal/logger"
)

const localLogName = "picker"

var log *zap.Logger

// Reset matcher instance.
func Reset() {
	ResetLog()
}

// Debug activate debug logs.
func Debug() {
	log = logger.InternalLogger(true).Named(localLogName)
}

// ResetLog activate debug logs.
func ResetLog() {
	log = logger.InternalLogger(false).Named(localLogName)
}

// NoLog disable logging under Fatal level.
func NoLog() {
	log = zap.NewNop()
}
