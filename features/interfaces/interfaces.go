// package interfaces exposes kactus features to enrich godog experience, allowing
// you to define your own steps with it.
//
// To have a more basic experiences around kactus, look into definitions package.
package interfaces

import (
	"github.com/cucumber/messages-go/v10"
)

// ComposeBeforeStep allow to aggregate multiple BeforeStep helper into a single one to pass to
// godog.BeforeStep.
//
// /!\ order is important. BeforeSteps function will be executed in provided order.
func ComposeBeforeStep(beforeSteps ...func(steps *messages.Pickle_PickleStep)) func(steps *messages.Pickle_PickleStep) {
	return func(steps *messages.Pickle_PickleStep) {
		for _, beforeStep := range beforeSteps {
			beforeStep(steps)
		}
	}
}

// AsNot wrap a simple expect/unexpect Step assertion from string to bool.
//
// It is mostly a syntaxic helper allowing to convert a NOT string to boolean.
//
// It matches steps like:
//   I (don't )expect something (p1) to happen when (p2)
func AsNot(fn NotFunction) func(string, ...string) error {
	return func(not string, params ...string) error {
		if not == "" {
			return fn(false, params...)
		}

		return fn(true, params...)
	}
}

func AsNot2(fn NotFunction) func(string) error {
	return func(not string) error {
		if not == "" {
			return fn(false)
		}

		return fn(true)
	}
}

type NotFunction func(bool, ...string) error
