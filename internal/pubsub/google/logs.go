package google

import (
	"go.uber.org/zap"

	"github.com/elmagician/kactus/internal/logger"
)

const (
	packageLogName = "pubsub"
	localLogName   = "google"
)

var log *zap.Logger

// Reset matcher instance.
func Reset() error {
	ResetLog()
	return nil
}

// Debug activate debug logs.
func Debug() error {
	log = logger.InternalLogger(true).Named(packageLogName).Named(localLogName)
	return nil
}

// ResetLog activate debug logs.
func ResetLog() {
	log = logger.InternalLogger(false).Named(packageLogName).Named(localLogName)
}

// NoLog disable logging under Fatal level.
func NoLog() {
	log = zap.NewNop()
}
