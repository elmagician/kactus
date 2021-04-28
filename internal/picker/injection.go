package picker

import (
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

var (
	injectRegex           = regexp.MustCompile(`^{{.+}}$`)
	findAllInjectionRegex = regexp.MustCompile(`{{[^{}]+}}`)
)

// Inject returns value corresponding to element if element matches injectRegex.
// Else, returns provided value.
func (store *Store) Inject(inTo interface{}) interface{} {
	log.Debug("Injecting into element", zap.Reflect("src", inTo))

	t, ok := inTo.(string)
	if ok && injectRegex.MatchString(t) {
		log.Debug("Found variable")

		t = strings.ReplaceAll(t, "{", "")
		t = strings.ReplaceAll(t, "}", "")

		if val, exists := store.Get(t); exists {
			log.Debug("Replacing with value", zap.Reflect("", val))
			return val
		}
	}

	return inTo
}

// InjectAll parses provided string using findAllInjectionRegex and replace
// injected values using Inject method.
func (store *Store) InjectAll(inTo string) string {
	log.Debug("Injecting into", zap.Reflect("src", inTo))

	return findAllInjectionRegex.ReplaceAllStringFunc(
		inTo,
		func(key string) string {
			return fmt.Sprintf("%v", store.Inject(key))
		},
	)
}
