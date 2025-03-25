package picker

import (
	"fmt"

	"github.com/cucumber/godog"

	"github.com/elmagician/kactus/features/interfaces"
	"github.com/elmagician/kactus/internal/matchers"
)

// VerifiesPickedValues asserts variables values matches expected conditions.
func (picker *Picker) VerifiesPickedValues(assertions *godog.Table) error {
	var key, val, matcher string

	head := assertions.Rows[0].Cells

	for i := 1; i < len(assertions.Rows); i++ {
		for n, cell := range assertions.Rows[i].Cells {
			switch head[n].Value {
			case "key":
				key = cell.Value
			case "matcher":
				matcher = cell.Value
			case "value":
				val = cell.Value
			default:
				return fmt.Errorf("%w %s", interfaces.ErrUnexpectedColumn, head[n].Value)
			}
		}

		if actualVal, exists := picker.store.Get(key); exists {
			if err := matchers.Assert(matcher, actualVal, val); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("key %v do not exist", key)
		}

		key = ""
		val = ""
		matcher = ""
	}

	return nil // return nil to be usable in godog Steps
}

// VerifiesPickedValue asserts variable value matches expected condition.
func (picker *Picker) VerifiesPickedValue(pickedKey, matcher, expectedValue string) error {
	actualVal, exists := picker.store.Get(pickedKey)
	if !exists {
		return fmt.Errorf("%w: %s", ErrNotPicked, pickedKey)
	}

	if err := matchers.Assert(matcher, actualVal, expectedValue); err != nil {
		return err
	}

	return nil
}
