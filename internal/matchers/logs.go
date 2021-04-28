package matchers

import (
	"go.uber.org/zap"

	"github.com/elmagician/kactus/internal/logger"
)

const localLogName = "matchers"

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
