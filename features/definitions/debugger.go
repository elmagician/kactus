package definitions

import (
	"github.com/elmagician/godog"

	"github.com/elmagician/kactus/internal/matchers"
	"github.com/elmagician/kactus/internal/types"
)

// InstallDebug adds debug base definitions to known steps for provided godog ScenarioContext.
// It allows to activate or disable debug logs on specific features offered by kactus.
//
// Provided steps:
// - (?:I )?want to debug types converters => activate debug logs for internal type conversion instance
//   I want to debug types converters
// - (?:I )?want to stop debugging types converters => disable debug logs for internal type conversion instance
//   I want to debug types converters
//
// - (?:I )?want to debug  matchers => activate debug logs for internal matchers operation instance
//   I want to debug matchers
// - (?:I )?want to stop debugging matchers$ => disable debug logs for internal matchers operation instance
//   I want to debug matchers
func InstallDebug(s *godog.ScenarioContext) {
	// Debug types
	s.Step(`^(?:I )?want to debug types converters$`, types.Debug)
	s.Step(`^(?:I )?want to stop debugging types converters$`, func() error {
		types.ResetLog()
		return nil
	})

	// Debug matchers
	s.Step(`^(?:I )?want to debug matchers$`, matchers.Debug)
	s.Step(`^(?:I )?want to stop debugging matchers$`, func() error {
		matchers.ResetLog()
		return nil
	})
}
