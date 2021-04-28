package types

import (
	"errors"
)

var (
	// ErrUnsupportedType indicates that expected type is not
	// managed by kactus library.
	ErrUnsupportedType = errors.New("unsupported type")

	// ErrUnmatchedType indicates that value interface does not
	// match receptor type.
	ErrUnmatchedType = errors.New("type did not matched expected")
)
