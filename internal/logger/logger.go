package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	debugLogger   *zap.Logger
	defaultLogger *zap.Logger
)

func init() {
	if err := SetDefault(); err != nil {
		fmt.Println("/!\\ Could not initialize default logger.")
	}
}

// Set sets kactus logger using provided zap instance.
// Kactus expect provided logger to be on DEBUG Level.
func Set(logger *zap.Logger) {
	if !logger.Core().Enabled(zapcore.DebugLevel) {
		logger.Warn("Provided logger does not enabled Debug logging. Debug feature will not work as expected.")
	}

	*debugLogger = *logger.Named("kactus")
	*defaultLogger = *debugLogger.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel))
}

// SetDefault sets kactus logger to default values.
func SetDefault() (err error) {
	debugLogger, err = zap.NewDevelopment()
	if err != nil {
		return
	}

	debugLogger = debugLogger.Named("kactus")
	defaultLogger = debugLogger.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel))

	return
}

// Logger copy known loggers to use in kactus logging.
func Logger(isDebug bool) *zap.Logger {
	log := &zap.Logger{}
	if isDebug {
		*log = *debugLogger
	} else {
		*log = *defaultLogger
	}

	return log
}

// InternalLogger get logger for internal packages.
// It is a shortcut for Logger().Named("internal").
func InternalLogger(isDebug bool) *zap.Logger {
	return Logger(isDebug).Named("internal")
}
