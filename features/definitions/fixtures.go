package definitions

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/elmagician/kactus/features/interfaces/fixtures"
)

// InstallFixtures initialize fixtures.Fixtures features steps.
// It initialize the auto loading process for manifest
// from @manifest tags (@manifest:path/to/manifest.yml).
//
// Provided steps:
// - (?:I )?load(?:ing)? fixture ([a-zA-Z0-9_./-]+)(?: into )?(.+) => execute fixture file in provided instance.
// Fixture kind will be determined from instance
//
//		I load fixture fix/test.sql into pg.example
//	  - (?:I )?load(?:ing)? manifest ([a-zA-Z0-9_./-]+) => parse and load fixture manifest.
//	    I load manifest fix/example_manifest.yml
func InstallFixtures(s *godog.ScenarioContext, fix *fixtures.Fixtures) {
	fix.Reset()

	s.Step("^(?:I )?load(?:ing)? fixture ([a-zA-Z0-9_./-]+)(?: into )?(.+)?$", fix.Load)
	s.Step("^(?:I )?load(?:ing)? manifest ([a-zA-Z0-9_./-]+)$", fix.LoadManifest)

	s.Before(func(ctx context.Context, s *godog.Scenario) (context.Context, error) {
		fix.BeforeScenarioTagLoader(s)
		return ctx, nil
	})
}
