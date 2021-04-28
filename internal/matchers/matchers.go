package matchers

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/elmagician/kactus/internal/types"
)

// Match asserts expected value matched provided regular expression.
// It will try to compile the expression then match value against it.
func Match(actual interface{}, expr string) (bool, error) {
	log.Debug("Asserting using Match function")

	re, err := regexp.Compile(expr)
	if err != nil {
		log.Error("could not compile regex", zap.Error(err))
		return false, fmt.Errorf("%w: expr should be a golang compatible regular expression", ErrInvalidArgument)
	}

	switch matchedVal := actual.(type) {
	case string:
		log.Debug("Actual value is a string")
		return re.MatchString(matchedVal), nil
	case reflect.Value:
		log.Debug("Actual value is a reflect.Value")

		if matchedVal.CanInterface() {
			log.Debug("Actual value can interface")
			return Match(matchedVal.Interface(), expr)
		}
	}

	return false, fmt.Errorf("%w: actual value should be a string", ErrInvalidArgument)
}

// Contain asserts expected value is contain by actual value.
func Contain(actual interface{}, expected string) (bool, error) {
	log.Debug("Asserting using Contain function")

	switch matchedVal := actual.(type) {
	case string:
		log.Debug("Actual value is a string")
		return strings.Contains(matchedVal, expected), nil
	case reflect.Value:
		log.Debug("Actual value is a reflect.Value")

		if matchedVal.CanInterface() {
			log.Debug("Actual value can interface")
			return Contain(matchedVal.Interface(), expected)
		}
	}

	return false, fmt.Errorf("%w: actual value should be a string", ErrInvalidArgument)
}

// Equal asserts expected value is equal to actual value.
func Equal(actual interface{}, expected string) (bool, error) {
	log.Debug("Asserting using Equal function")

	matchedVal, ok := actual.(reflect.Value)
	if ok && matchedVal.CanInterface() {
		return Equal(matchedVal.Interface(), expected)
	}

	res, err := types.AssertEqual(expected, actual)

	if err != nil {
		log.Error("could not assert equality:", zap.Error(err))

		if errors.Is(err, types.ErrUnsupportedType) || errors.Is(err, types.ErrUnmatchedType) {
			return false, err
		}

		return false, nil
	}

	return res, nil
}

// OfType asserts actual value matches expected type
func OfType(actual interface{}, expectedType string) bool {
	log.Debug("Asserting using OfType function")

	matchedVal, ok := actual.(reflect.Value)
	if ok && matchedVal.CanInterface() {
		return OfType(matchedVal.Interface(), expectedType)
	}

	return expectedType == reflect.TypeOf(actual).String()
}

// IsDefined asserts value is a JSON defined field value.
func IsDefined(value interface{}) bool {
	log.Debug("Asserting using IsDefined function")

	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		return false
	}

	switch matchedVal := value.(type) {
	case string:
		log.Debug("Value is a string")
		return matchedVal != "undefined"
	case reflect.Value:
		log.Debug("Value is a reflect.Value")
		return matchedVal.CanInterface() && IsDefined(matchedVal.Interface())
	}

	return true
}

// NotZero asserts value is not a ZeroValue.
// /!\ undefined string value is considered as Zero!
func NotZero(value interface{}) bool {
	log.Debug("Asserting using NotZero function")

	if value == nil {
		log.Debug("Value is a nil pointer")
		return false
	}

	switch matchedVal := value.(type) {
	case string:
		log.Debug("Value is a string")
		return matchedVal != "undefined" && matchedVal != ""
	case reflect.Value:
		log.Debug("Value is a reflect.Value")
		return matchedVal.IsValid() && !matchedVal.IsZero()
	default:
		log.Debug("Value is not yet assertable. Calling back NotZero using reflect.ValueOf.")
		return NotZero(reflect.ValueOf(value))
	}
}

// LenEqual asserts actual value length equals expected length.
//
// /!\ It will panic if actual value Length cannot be asserted.
func LenEqual(actual interface{}, expected string) (bool, error) {
	log.Debug("Asserting using LenEqual function")

	expectedLen, err := strconv.Atoi(expected)
	if err != nil {
		log.Error("expectedVal is not an int", zap.Error(err))
		return false, fmt.Errorf("%w: expected value should be an integer", ErrInvalidArgument)
	}

	switch matchedVal := actual.(type) {
	case string:
		log.Debug("Value is a string")
		return len(matchedVal) == expectedLen, nil
	case reflect.Value:
		log.Debug("Value is a reflect.Value")
		return matchedVal.Len() == expectedLen, nil
	default:
		log.Debug("Value is not yet assertable. Calling back NotZero using reflect.ValueOf.")
		return LenEqual(reflect.ValueOf(actual), expected) // cannot ensure value can have a length. Will panic if it does not :(
	}
}

// In asserts actual value matches one of the expected value.
// Expected values has to be `,` separated.
func In(actual interface{}, expected string) bool {
	log.Debug("Asserting using In function")

	list := strings.Split(expected, ",")
	for _, candidate := range list {
		log.Debug(
			"Matching against candidate",
			zap.Reflect("actual", actual), zap.String("expected", expected),
		)

		if ok, err := Equal(actual, candidate); ok && err == nil {
			return true
		}
	}

	return false
}

// Assert asserts actual value matches expected value using provided method.
// Assert is designed to be used in steps so it will return a single error.
// If assertions produces error, assertion error will be returned.
// If assertion is false, ErrUnmatched will be return.
func Assert(method string, actualValue interface{}, expectedValue string) error {
	log.Debug(
		"Asserting with",
		zap.String("method", method),
		zap.String("expected", expectedValue),
		zap.Reflect("actual", actualValue),
	)

	var matches bool
	var err error

	switch strings.ToLower(method) {
	case "equal", "eq", "equals", "=", "==", "":
		matches, err = Equal(actualValue, expectedValue)
	case "match", "matches", "=~":
		matches, err = Match(actualValue, expectedValue)
	case "contain", "contains":
		matches, err = Contain(actualValue, expectedValue)
	case "present", "defined":
		matches = IsDefined(actualValue)
	case "not zero":
		matches = NotZero(actualValue)
	case "in":
		matches = In(actualValue, expectedValue)
	case "length equals", "l=", "s=", "size is":
		matches, err = LenEqual(actualValue, expectedValue)
	case "type":
		matches = OfType(actualValue, expectedValue)
	default:
		log.Warn("Undefined matcher " + method)
		return ErrUndefinedMethod
	}

	if err != nil {
		log.Debug("Assertions encounters error")
		return err
	}

	return AsError(matches, actualValue, expectedValue)
}

// AsError produces error on assertion results.
func AsError(res bool, actualValue interface{}, expectedValue string) error {
	if !res {
		log.Debug("Failed assertion")
		return fmt.Errorf("actual (%v) %w expected (%s) value", actualValue, ErrUnmatched, expectedValue)
	}

	log.Debug("Values matches !")

	return nil
}
