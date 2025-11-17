package models

import (
	apperrors "github.com/Andre385/bruschirentals-backend/internal/errors"
	"github.com/google/uuid"
)

// Neighborhood represents a neighborhood location.
type Neighborhood struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

// NewNeighborhood creates a new Neighborhood instance with validation.
func NewNeighborhood(id uuid.UUID, name string) (Neighborhood, error) {
	n := Neighborhood{ID: id, Name: name}
	return n, n.Validate()
}

// Validate checks if the neighborhood is valid.
func (n Neighborhood) Validate() error {
	if n.ID == uuid.Nil {
		return apperrors.ErrInvalidInput
	}
	if n.Name == "" {
		return apperrors.ErrInvalidInput
	}
	return nil
}
