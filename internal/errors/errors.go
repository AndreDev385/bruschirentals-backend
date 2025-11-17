// Package errors centralize our app errors in one place
package errors

import "errors"

// Common errors used across the application.
var (
	ErrInvalidID         = errors.New("invalid ID format")
	ErrNotFound          = errors.New("not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidPriceRange = errors.New("invalid price range")
	ErrInvalidPromotion  = errors.New("invalid promotion")
	ErrInvalidApartment  = errors.New("invalid apartment")
)
