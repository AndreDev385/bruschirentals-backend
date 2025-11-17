package models

import "errors"

// Promotion represents a promotional offer for rentals.
type Promotion struct {
	MonthsFree uint8    `json:"months_free"`
	Conditions []string `json:"conditions"`
}

// NewPromotion creates a new Promotion with validation.
func NewPromotion(monthsFree uint8, conditions []string) (Promotion, error) {
	p := Promotion{MonthsFree: monthsFree, Conditions: conditions}
	return p, p.Validate()
}

// Validate checks if the promotion is valid.
func (p Promotion) Validate() error {
	if p.MonthsFree == 0 {
		return errors.New("months free must be greater than 0")
	}
	if len(p.Conditions) == 0 {
		return errors.New("conditions must not be empty")
	}
	for _, cond := range p.Conditions {
		if cond == "" {
			return errors.New("conditions must not contain empty strings")
		}
	}
	return nil
}
