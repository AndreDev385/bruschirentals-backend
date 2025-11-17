// Package utils have shared domain logic
package utils

import (
	apperrors "github.com/Andre385/bruschirentals-backend/internal/errors"
	"github.com/google/uuid"
)

// ValidateID validates and parses a string ID to UUID.
func ValidateID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, apperrors.ErrInvalidInput
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperrors.ErrInvalidID
	}
	return parsedID, nil
}
