// Package repositories provides data access layer implementations.
package repositories

import (
	"context"
	"database/sql"
	"errors"

	apperrors "github.com/Andre385/bruschirentals-backend/internal/errors"
	"github.com/Andre385/bruschirentals-backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// BuildingRepository defines the interface for building data operations.
type BuildingRepository interface {
	Save(ctx context.Context, building models.Building) error
	GetByID(ctx context.Context, id string) (models.Building, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]models.Building, error)
}

// buildingRepository implements BuildingRepository.
type buildingRepository struct {
	db *sqlx.DB
}

// NewBuildingRepository creates a new building repository.
func NewBuildingRepository(db *sqlx.DB) BuildingRepository {
	return &buildingRepository{db: db}
}

// Save inserts or updates a building in the database.
func (r *buildingRepository) Save(ctx context.Context, building models.Building) error {
	query := `INSERT INTO buildings (id, name, neighborhood_id, address) VALUES ($1, $2, $3, $4)
	          ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, neighborhood_id = EXCLUDED.neighborhood_id, address = EXCLUDED.address`
	_, err := r.db.ExecContext(ctx, query, building.ID, building.Name, building.NeighborhoodID, building.Address)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" { // foreign_key_violation
			return apperrors.ErrInvalidInput
		}
		return err
	}
	return nil
}

// GetByID retrieves a building by ID.
func (r *buildingRepository) GetByID(ctx context.Context, id string) (models.Building, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Building{}, apperrors.ErrInvalidID
	}

	var building models.Building
	query := `SELECT id, name, neighborhood_id, address FROM buildings WHERE id = $1`
	err = r.db.GetContext(ctx, &building, query, parsedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Building{}, apperrors.ErrNotFound
		}
		return models.Building{}, err
	}
	return building, nil
}

// Delete removes a building by ID.
func (r *buildingRepository) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return apperrors.ErrInvalidID
	}

	query := `DELETE FROM buildings WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, parsedID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// List retrieves all buildings.
func (r *buildingRepository) List(ctx context.Context) ([]models.Building, error) {
	var buildings []models.Building
	query := `SELECT id, name, neighborhood_id, address FROM buildings ORDER BY name`
	err := r.db.SelectContext(ctx, &buildings, query)
	return buildings, err
}
