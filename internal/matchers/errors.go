package matchers

import (
	"errors"
)

var (
	ErrUndefinedMethod = errors.New("undefined matcher method")
	ErrInvalidArgument = errors.New("at least one arguments did not match asserter expectation")
	ErrUnmatched       = errors.New("value does not match")
)
