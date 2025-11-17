package models

import (
	apperrors "github.com/Andre385/bruschirentals-backend/internal/errors"
	"github.com/google/uuid"
)

// Building represents a building in a neighborhood.
type Building struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	NeighborhoodID uuid.UUID `json:"neighborhood_id"`
	Address        string    `json:"address"`
}

// NewBuilding creates a new Building instance with validation.
func NewBuilding(id uuid.UUID, name string, neighborhoodID uuid.UUID, address string) (Building, error) {
	b := Building{
		ID:             id,
		Name:           name,
		NeighborhoodID: neighborhoodID,
		Address:        address,
	}
	return b, b.Validate()
}

// Validate checks if the building is valid.
func (b Building) Validate() error {
	if b.ID == uuid.Nil {
		return apperrors.ErrInvalidInput
	}
	if b.Name == "" {
		return apperrors.ErrInvalidInput
	}
	if b.NeighborhoodID == uuid.Nil {
		return apperrors.ErrInvalidInput
	}
	if b.Address == "" {
		return apperrors.ErrInvalidInput
	}
	return nil
}
