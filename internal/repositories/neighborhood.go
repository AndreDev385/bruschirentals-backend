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

// NeighborhoodRepository defines the interface for neighborhood data operations.
type NeighborhoodRepository interface {
	Save(ctx context.Context, neighborhood models.Neighborhood) error
	GetByID(ctx context.Context, id string) (models.Neighborhood, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]models.Neighborhood, error)
}

// neighborhoodRepository implements NeighborhoodRepository.
type neighborhoodRepository struct {
	db *sqlx.DB
}

// NewNeighborhoodRepository creates a new neighborhood repository.
func NewNeighborhoodRepository(db *sqlx.DB) NeighborhoodRepository {
	return &neighborhoodRepository{db: db}
}

// Save inserts or updates a neighborhood in the database.
func (r *neighborhoodRepository) Save(ctx context.Context, neighborhood models.Neighborhood) error {
	query := `INSERT INTO neighborhoods (id, name) VALUES ($1, $2)
	          ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`
	_, err := r.db.ExecContext(ctx, query, neighborhood.ID, neighborhood.Name)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // unique_violation
			return apperrors.ErrInvalidInput
		}
		return err
	}
	return nil
}

// GetByID retrieves a neighborhood by ID.
func (r *neighborhoodRepository) GetByID(ctx context.Context, id string) (models.Neighborhood, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Neighborhood{}, apperrors.ErrInvalidID
	}

	var neighborhood models.Neighborhood
	query := `SELECT id, name FROM neighborhoods WHERE id = $1`
	err = r.db.GetContext(ctx, &neighborhood, query, parsedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Neighborhood{}, apperrors.ErrNotFound
		}
		return models.Neighborhood{}, err
	}
	return neighborhood, nil
}

// Delete removes a neighborhood by ID.
func (r *neighborhoodRepository) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return apperrors.ErrInvalidID
	}

	query := `DELETE FROM neighborhoods WHERE id = $1`
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

// List retrieves all neighborhoods.
func (r *neighborhoodRepository) List(ctx context.Context) ([]models.Neighborhood, error) {
	var neighborhoods []models.Neighborhood
	query := `SELECT id, name FROM neighborhoods ORDER BY name`
	err := r.db.SelectContext(ctx, &neighborhoods, query)
	return neighborhoods, err
}
