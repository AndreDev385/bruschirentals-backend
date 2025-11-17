package models

import "errors"

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
		return errors.New("price range From must be >= 0 and To > 0")
	}
	if p.From >= p.To {
		return errors.New("price range From must be less than To")
	}
	return nil
}
