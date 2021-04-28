// package picker provided structures and method
// to manage injected variable system for Godog.
//
// It allows an user to pick value into a store and
// inject them in steps through a variable replacement.
package interfaces

import (
	"go.uber.org/zap"

	"github.com/elmagician/kactus/internal/logger"
)

const localLogName = "interfaces"

var log *zap.Logger

// Reset matcher instance.
func Reset() error {
	ResetLog()
	return nil
}

// Debug activate debug logs.
func Debug() error {
	log = logger.InternalLogger(true).Named(localLogName)
	return nil
}

// ResetLog activate debug logs.
func ResetLog() {
	log = logger.InternalLogger(false).Named(localLogName)
}

// NoLog disable logging under Fatal level.
func NoLog() {
	log = zap.NewNop()
}
