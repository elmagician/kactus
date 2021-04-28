package types

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	hasExpectedType = 2 // const for split on matcher with types
)

// ToInterface converts value to asked type.
func ToInterface(value string) (interface{}, error) {
	val, wantType := getValAndType(value)
	log.Debug(
		"Converting value...",
		zap.String("value", val), zap.String("required type", wantType),
	)

	switch wantType {
	case "int":
		var res int64
		if err := SetTyped(val, &res); err != nil {
			return nil, err
		}

		log.Debug("Correctly converted value to int", zap.Reflect("value", res))

		return res, nil
	case "float", "number":
		var res float64
		if err := SetTyped(val, &res); err != nil {
			return nil, err
		}

		log.Debug("Correctly converted value to float", zap.Reflect("value", res))

		return res, nil
	case "boolean", "bool":
		var res bool
		if err := SetTyped(val, &res); err != nil {
			return nil, err
		}

		log.Debug("Correctly converted value to bool", zap.Reflect("value", res))

		return res, nil
	case "uuid":
		var res uuid.UUID
		if err := SetTyped(val, &res); err != nil {
			return nil, err
		}

		log.Debug("Correctly converted value to uuid", zap.Reflect("value", res))

		return res, nil
	case "string":
		log.Debug("No conversion required", zap.String("value", val))
		return val, nil
	case "array":
		if val[0] == '[' {
			val = val[1:]
		}

		if val[len(val)-1] == ']' {
			val = val[:len(val)-1]
		}

		var res []interface{}

		values := strings.Split(val, ",")
		for _, v := range values {
			tmp, err := ToInterface(v)
			if err != nil {
				return nil, err
			}

			res = append(res, tmp)
		}

		log.Debug("Correctly converted value to array", zap.Reflect("value", res))

		return res, nil
	}

	return value, nil
}

// AssertEqual ensures actual interfaced value is equal to expected string typed value. Type can be ignored
// on string assertion only.
func AssertEqual(expected string, actual interface{}) (bool, error) { // nolint: gocyclo
	_, wantType := getValAndType(expected)
	if wantType == "" {
		actualCasted, ok := actual.(string)
		if !ok {
			return false, fmt.Errorf("actual value %s %w string", actual, ErrUnmatchedType)
		}

		log.Debug(
			"Asserting using default String type",
			zap.String("expected", expected), zap.Reflect("actual", actual),
		)

		return expected == actualCasted, nil
	}

	valAsInterface, err := ToInterface(expected)
	if err != nil {
		return false, err
	}

	if strActual, ok := actual.(string); ok {
		actual, err = ToInterface(strActual + "((" + wantType + "))")
		if err != nil {
			return false, err
		}
	}

	// disable errcheck to avoid lint errors when not checking cast on expected value interface
	// this check make no sense here has value conversion is already made and check
	// by ToInterface method.
	// nolint: errcheck
	switch wantType {
	case "int":
		castedExpected := valAsInterface.(int64)
		castedActual, ok2 := AsInt64(actual)

		if !ok2 {
			return false, fmt.Errorf("actual value %v %w %s", actual, ErrUnmatchedType, wantType)
		}

		log.Debug(
			"asserting equality using INT type",
			zap.Int64("expected", castedExpected), zap.Int64("actual", castedActual),
		)

		return castedExpected == castedActual, nil

	case "float", "number":
		castedExpected := valAsInterface.(float64)
		castedActual, ok2 := AsFloat64(actual)

		if !ok2 {
			return false, fmt.Errorf("actual value %v %w %s", actual, ErrUnmatchedType, wantType)
		}

		log.Debug(
			"Asserting equality using NUMBER type",
			zap.Float64("expected", castedExpected), zap.Float64("actual", castedActual),
		)

		return castedExpected == castedActual, nil

	case "string":
		castedExpected := valAsInterface.(string)
		castedActual, ok2 := AsString(actual)

		if !ok2 {
			return false, fmt.Errorf("actual value %v %w %s", actual, ErrUnmatchedType, wantType)
		}

		log.Debug(
			"Asserting equality using STRING type",
			zap.String("expected", castedExpected), zap.String("actual", castedActual),
		)

		return castedExpected == castedActual, nil
	case "uuid":
		castedExpected := valAsInterface.(uuid.UUID)
		castedActual, ok2 := actual.(uuid.UUID)

		if !ok2 {
			return false, fmt.Errorf("actual value %v %w %s", actual, ErrUnmatchedType, wantType)
		}

		log.Debug(
			"Asserting equality using UUID type",
			zap.String("expected", castedExpected.String()), zap.String("actual", castedActual.String()),
		)

		return castedActual.String() == castedExpected.String(), nil
	case "boolean", "bool":
		castedExpected := valAsInterface.(bool)
		castedActual, ok2 := AsBool(actual)

		if !ok2 {
			return false, fmt.Errorf("actual value %v %w %s", actual, ErrUnmatchedType, wantType)
		}

		log.Debug(
			"Asserting equality using BOOL type",
			zap.Bool("expected", castedExpected), zap.Bool("actual", castedActual),
		)

		return castedExpected == castedActual, nil
	case "array":
		castedExpected := valAsInterface.([]interface{})
		castedActual, ok := actual.([]interface{})

		if !ok {
			return false, fmt.Errorf("actual value %v %w %s", actual, ErrUnmatchedType, wantType)
		}

		maxIndex := len(castedActual) - 1

		for ind, val := range castedExpected {
			if ind > maxIndex || val != castedActual[ind] {
				return false, nil
			}
		}

		log.Debug(
			"Asserting equality using ARRAY type",
			zap.Reflect("expected", castedExpected), zap.Reflect("actual", castedActual),
		)

		return true, nil
	}

	return valAsInterface == actual, nil
}

func getValAndType(value string) (val string, wantType string) {
	trimmed := strings.TrimSuffix(value, "))")
	listedVal := strings.Split(trimmed, "((")

	listedLen := len(listedVal)

	switch {
	case listedLen < hasExpectedType:
		val = value
		wantType = ""

	case listedLen > hasExpectedType:
		wantType = listedVal[listedLen-1]
		val = strings.Join(listedVal[:listedLen-1], "((")

	default:
		val = listedVal[0]
		wantType = listedVal[1]
	}

	return
}
