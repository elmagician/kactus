package internal

import (
	"errors"
)

var (
	// ErrUnexpectedColumn is raised when loading a godog table with
	// unexpected column header.
	ErrUnexpectedColumn = errors.New("unexpected column name")
)
