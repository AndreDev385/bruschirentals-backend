package models

import apperrors "github.com/Andre385/bruschirentals-backend/internal/errors"

// PriceRange represents a range of prices with From and To values (in cents).
type PriceRange struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

// NewPriceRange creates a new PriceRange with validation.
func NewPriceRange(from, to int64) (PriceRange, error) {
	pr := PriceRange{From: from, To: to}
	return pr, pr.Validate()
}

// Validate checks if the price range is valid (From < To and non-negative).
func (p PriceRange) Validate() error {
	if p.From < 0 || p.To <= 0 {
		return apperrors.ErrInvalidPriceRange
	}
	if p.From >= p.To {
		return apperrors.ErrInvalidPriceRange
	}
	return nil
}
