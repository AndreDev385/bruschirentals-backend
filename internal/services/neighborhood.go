// Package services provides business logic layer implementations.
package services

import (
	"context"

	apperrors "github.com/Andre385/bruschirentals-backend/internal/errors"
	"github.com/Andre385/bruschirentals-backend/internal/models"
	"github.com/Andre385/bruschirentals-backend/internal/repositories"
	"github.com/google/uuid"
)

// NeighborhoodService handles business logic for neighborhoods.
type NeighborhoodService struct {
	repo repositories.NeighborhoodRepository
}

// validateID validates and parses a string ID to UUID.
func (s *NeighborhoodService) validateID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, apperrors.ErrInvalidInput
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperrors.ErrInvalidID
	}
	return parsedID, nil
}

// NewNeighborhoodService creates a new neighborhood service.
func NewNeighborhoodService(repo repositories.NeighborhoodRepository) *NeighborhoodService {
	return &NeighborhoodService{repo: repo}
}

// CreateNeighborhood creates a new neighborhood with validation.
func (s *NeighborhoodService) CreateNeighborhood(ctx context.Context, name string) (models.Neighborhood, error) {
	id := uuid.New()
	neighborhood, err := models.NewNeighborhood(id, name)
	if err != nil {
		return models.Neighborhood{}, err
	}

	err = s.repo.Create(ctx, neighborhood)
	if err != nil {
		return models.Neighborhood{}, err
	}

	return neighborhood, nil
}

// GetNeighborhood retrieves a neighborhood by ID.
func (s *NeighborhoodService) GetNeighborhood(ctx context.Context, id string) (models.Neighborhood, error) {
	_, err := s.validateID(id)
	if err != nil {
		return models.Neighborhood{}, err
	}

	neighborhood, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.Neighborhood{}, err
	}

	return neighborhood, nil
}

// UpdateNeighborhood updates an existing neighborhood with validation.
func (s *NeighborhoodService) UpdateNeighborhood(ctx context.Context, id string, name string) (models.Neighborhood, error) {
	parsedID, err := s.validateID(id)
	if err != nil {
		return models.Neighborhood{}, err
	}

	neighborhood, err := models.NewNeighborhood(parsedID, name)
	if err != nil {
		return models.Neighborhood{}, err
	}

	err = s.repo.Update(ctx, neighborhood)
	if err != nil {
		return models.Neighborhood{}, err
	}

	return neighborhood, nil
}

// DeleteNeighborhood removes a neighborhood by ID.
func (s *NeighborhoodService) DeleteNeighborhood(ctx context.Context, id string) error {
	_, err := s.validateID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

// ListNeighborhoods retrieves all neighborhoods.
func (s *NeighborhoodService) ListNeighborhoods(ctx context.Context) ([]models.Neighborhood, error) {
	return s.repo.List(ctx)
}
