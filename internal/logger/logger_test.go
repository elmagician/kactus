package logger_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/elmagician/kactus/internal/logger"
	"github.com/elmagician/kactus/internal/matchers"
	. "github.com/elmagician/kactus/internal/test"
	"github.com/elmagician/kactus/internal/types"
)

func init() {
	matchers.NoLog()
	types.NoLog()
}

func TestUnit_Logger(t *testing.T) {
	Convey("Given I wish to log from Kactus, ", t, func() {
		Convey("I should be able to set Default logger", func() {
			So(logger.SetDefault(), ShouldBeNil)

			Convey("Then I should be able to log debug", func() {
				debugLogger := logger.Logger(true)

				So(debugLogger.Core().Enabled(zapcore.DebugLevel), ShouldBeTrue)
			})

			Convey("Then I should be able to log without debug", func() {
				defaultLogger := logger.Logger(false)

				So(defaultLogger.Core().Enabled(zapcore.DebugLevel), ShouldBeFalse)
			})

			Convey("Then I should be able to retrieve logger for internal packages", func() {
				internalLogger := logger.InternalLogger(false)

				So(internalLogger.Core().Enabled(zapcore.DebugLevel), ShouldBeFalse)
			})

		})

		Convey("I should be able to set a custom zap instance", func() {
			core, logs := observer.New(zapcore.DebugLevel)
			logger.Set(zap.New(core))
			So(logs.Len(), ShouldEqual, 0)

			Convey("Then I should be able to log debug", func() {
				debugLogger := logger.Logger(true)

				debugLogger.Debug("test")

				So(logs.Len(), ShouldEqual, 1)
				entry := logs.TakeAll()[0]
				So(entry.Level, ShouldBeEquivalent, zap.DebugLevel)
				So(entry.LoggerName, ShouldEqual, "kactus")
				So(entry.Message, ShouldEqual, "test")
			})

			Convey("Then I should be able to log without debug", func() {
				defaultLogger := logger.Logger(false)

				defaultLogger.Debug("test")

				So(logs.Len(), ShouldEqual, 0)
			})

			Convey("Then I should be able to retrieve logger for internal packages", func() {
				internalLogger := logger.InternalLogger(false)

				internalLogger.Info("test")

				So(logs.Len(), ShouldEqual, 1)
				entry := logs.TakeAll()[0]
				So(entry.Level, ShouldBeEquivalent, zap.InfoLevel)
				So(entry.LoggerName, ShouldEqual, "kactus.internal")
				So(entry.Message, ShouldEqual, "test")
			})

		})

		Convey("I should have a warning log if setting logger instance who doesn't support debug level logging", func() {
			core, logs := observer.New(zapcore.InfoLevel)
			logger.Set(zap.New(core))
			So(logs.Len(), ShouldEqual, 1)
			entry := logs.TakeAll()[0]
			So(entry.Message, ShouldEqual, "Provided logger does not enabled Debug logging. Debug feature will not work as expected.")

		})
	})
}
