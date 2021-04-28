package fixtures

import (
	"github.com/elmagician/godog"
)

// BeforeScenarioTagLoader initialize the TagLoader for fixture.
// It allows to pass a `@manifest` tag to load fixture manifest
//   @manifest:fixture/example.yml
func (fix *Fixtures) BeforeScenarioTagLoader(s *godog.Scenario) {
	fix.fixtures.NewTagLoader(s)
}
