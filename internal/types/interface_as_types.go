package types

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

// defining const on true && false as ok && nok
// to have a cleaner reading of the function.
// It helps to explicit the fact that
// a true value returns means the conversion
// was correct while a false value stand for
// invalid cast.
const (
	ok  = true
	nok = false
)

// AsInt64 converts int && string interfaces values to int64.
func AsInt64(val interface{}) (int64, bool) {
	switch castedVal := val.(type) {
	case int64:
		return castedVal, ok
	case int:
		return int64(castedVal), ok
	case int8:
		return int64(castedVal), ok
	case int16:
		return int64(castedVal), ok
	case int32:
		return int64(castedVal), ok
	case string:
		intVal, err := strconv.ParseInt(castedVal, 10, 64)
		if err != nil {
			log.Error("could not convert string to int", zap.String("val", castedVal), zap.Error(err))
			return 0, nok
		}

		return intVal, ok
	}

	return 0, nok
}

// AsFloat64 converts int, string && float interfaces values to float64.
func AsFloat64(val interface{}) (float64, bool) {
	switch castedVal := val.(type) {
	case float64:
		return castedVal, true
	case float32:
		return float64(castedVal), true
	case int64:
		return float64(castedVal), true
	case int:
		return float64(castedVal), true
	case int8:
		return float64(castedVal), true
	case int16:
		return float64(castedVal), true
	case int32:
		return float64(castedVal), true
	case string:
		floatVal, err := strconv.ParseFloat(castedVal, 64)
		if err != nil {
			log.Error("could not convert string to float", zap.String("val", castedVal), zap.Error(err))
			return 0, false
		}

		return floatVal, true
	}

	return 0, false
}

// AsBool converts str, bool and int interfaces to boolean.
// For integer: 1 stands for true && 0 for false.
func AsBool(val interface{}) (bool, bool) {
	switch castedVal := val.(type) {
	case bool:
		return castedVal, ok
	case int:
		return castedVal == 1, ok
	case int64:
		return castedVal == 1, ok
	case int8:
		return castedVal == 1, ok
	case int16:
		return castedVal == 1, ok
	case int32:
		return castedVal == 1, ok
	case string:
		intVal, err := strconv.ParseBool(castedVal)
		if err != nil {
			log.Error("could not convert string to bool", zap.String("val", castedVal), zap.Error(err))
			return false, nok
		}

		return intVal, ok
	}

	return false, nok
}

// AsString converts any interface that as a string representation to string
// relying on %v format.
func AsString(val interface{}) (string, bool) {
	return fmt.Sprintf("%v", val), ok
}
