package errors

import "errors"

// Common errors used across the application.
var (
	ErrInvalidID    = errors.New("invalid ID format")
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)
