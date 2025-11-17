// Package services provides business logic layer implementations.
package services

import (
	"context"

	"github.com/Andre385/bruschirentals-backend/internal/models"
	"github.com/Andre385/bruschirentals-backend/internal/repositories"
	"github.com/Andre385/bruschirentals-backend/internal/utils"
	"github.com/google/uuid"
)

// BuildingService handles business logic for buildings.
type BuildingService struct {
	repo             repositories.BuildingRepository
	neighborhoodRepo repositories.NeighborhoodRepository
}

// NewBuildingService creates a new building service.
func NewBuildingService(repo repositories.BuildingRepository, neighborhoodRepo repositories.NeighborhoodRepository) *BuildingService {
	return &BuildingService{repo: repo, neighborhoodRepo: neighborhoodRepo}
}

// CreateBuilding creates a new building.
func (s *BuildingService) CreateBuilding(ctx context.Context, name string, neighborhoodID string, address string) (models.Building, error) {
	// Validate neighborhood ID
	neighborhoodUUID, err := utils.ValidateID(neighborhoodID)
	if err != nil {
		return models.Building{}, err
	}

	// Check if neighborhood exists
	_, err = s.neighborhoodRepo.GetByID(ctx, neighborhoodID)
	if err != nil {
		return models.Building{}, err
	}

	id := uuid.New()
	building, err := models.NewBuilding(id, name, neighborhoodUUID, address)
	if err != nil {
		return models.Building{}, err
	}

	err = s.repo.Save(ctx, building)
	if err != nil {
		return models.Building{}, err
	}

	return building, nil
}

// GetBuilding retrieves a building by ID.
func (s *BuildingService) GetBuilding(ctx context.Context, id string) (models.Building, error) {
	_, err := utils.ValidateID(id)
	if err != nil {
		return models.Building{}, err
	}

	building, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.Building{}, err
	}

	return building, nil
}

// UpdateBuilding updates an existing building.
func (s *BuildingService) UpdateBuilding(ctx context.Context, id string, name string, neighborhoodID string, address string) (models.Building, error) {
	// Validate building ID
	buildingUUID, err := utils.ValidateID(id)
	if err != nil {
		return models.Building{}, err
	}

	// Check if building exists
	_, err = s.repo.GetByID(ctx, id)
	if err != nil {
		return models.Building{}, err
	}

	// Validate neighborhood ID
	neighborhoodUUID, err := utils.ValidateID(neighborhoodID)
	if err != nil {
		return models.Building{}, err
	}

	// Check if neighborhood exists
	_, err = s.neighborhoodRepo.GetByID(ctx, neighborhoodID)
	if err != nil {
		return models.Building{}, err
	}

	building, err := models.NewBuilding(buildingUUID, name, neighborhoodUUID, address)
	if err != nil {
		return models.Building{}, err
	}

	err = s.repo.Save(ctx, building)
	if err != nil {
		return models.Building{}, err
	}

	return building, nil
}

// DeleteBuilding removes a building by ID.
func (s *BuildingService) DeleteBuilding(ctx context.Context, id string) error {
	_, err := utils.ValidateID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

// ListBuildings retrieves all buildings.
func (s *BuildingService) ListBuildings(ctx context.Context) ([]models.Building, error) {
	return s.repo.List(ctx)
}
