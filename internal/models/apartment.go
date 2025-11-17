// Package models contains domain models for the Bruschi Rentals application.
package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ApartmentType represents the type of apartment (e.g., Studio, OneBed).
type ApartmentType string

// Apartment type constants
const (
	Studio          ApartmentType = "Studio"
	OneBed          ApartmentType = "OneBed"
	TwoBeds         ApartmentType = "TwoBeds"
	ThreeOrMoreBeds ApartmentType = "ThreeOrMoreBeds"
)

// String returns the string representation of ApartmentType
func (a ApartmentType) String() string {
	return string(a)
}

// Apartment represents an apartment listing.
type Apartment struct {
	ID               uuid.UUID     `json:"id"`
	BuildingID       uuid.UUID     `json:"building_id"`
	Type             ApartmentType `json:"type"`
	Price            PriceRange    `json:"price"`
	PromotionalPrice *int64        `json:"promotional_price,omitempty"`
	Images           []string      `json:"images"`
	Videos           []string      `json:"videos"`
	LastUpdate       time.Time     `json:"last_update"`
}

// NewApartment creates a new Apartment with validation.
func NewApartment(id, buildingID uuid.UUID, aptType ApartmentType, price PriceRange, promoPrice *int64, images, videos []string, lastUpdate time.Time) (Apartment, error) {
	a := Apartment{
		ID:               id,
		BuildingID:       buildingID,
		Type:             aptType,
		Price:            price,
		PromotionalPrice: promoPrice,
		Images:           images,
		Videos:           videos,
		LastUpdate:       lastUpdate,
	}
	return a, a.Validate()
}

// Validate checks if the apartment is valid.
func (a Apartment) Validate() error {
	if a.ID == uuid.Nil {
		return errors.New("ID must not be nil")
	}
	if a.BuildingID == uuid.Nil {
		return errors.New("BuildingID must not be nil")
	}
	if err := a.Price.Validate(); err != nil {
		return err
	}
	if a.PromotionalPrice != nil && *a.PromotionalPrice < 0 {
		return errors.New("promotional price must be >= 0")
	}
	// Valid ApartmentType is enforced by type, but could add check
	return nil
}
