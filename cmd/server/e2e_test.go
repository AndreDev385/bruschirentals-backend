// Package main contains end-to-end tests for the API.
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Andre385/bruschirentals-backend/internal/handlers"
	"github.com/Andre385/bruschirentals-backend/internal/repositories"
	"github.com/Andre385/bruschirentals-backend/internal/services"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	db    *sqlx.DB
	dbURL string
	echo  *echo.Echo
}

func (suite *E2ETestSuite) SetupSuite() {
	// Setup test database
	suite.dbURL = os.Getenv("DATABASE_URL")
	if suite.dbURL == "" {
		suite.dbURL = "postgres://user:password@localhost:5432/bruschi_rentals_test?sslmode=disable"
	}

	var err error
	suite.db, err = sqlx.Connect("postgres", suite.dbURL)
	suite.Require().NoError(err)

	// Setup Echo app
	suite.echo = echo.New()

	// Initialize dependencies
	neighborhoodRepo := repositories.NewNeighborhoodRepository(suite.db)
	neighborhoodService := services.NewNeighborhoodService(neighborhoodRepo)
	neighborhoodHandler := handlers.NewNeighborhoodHandler(neighborhoodService)

	buildingRepo := repositories.NewBuildingRepository(suite.db)
	buildingService := services.NewBuildingService(buildingRepo, neighborhoodRepo)
	buildingHandler := handlers.NewBuildingHandler(buildingService)

	// Setup routes
	suite.echo.POST("/api/v1/neighborhoods", neighborhoodHandler.Create)
	suite.echo.GET("/api/v1/neighborhoods/:id", neighborhoodHandler.Get)
	suite.echo.PUT("/api/v1/neighborhoods/:id", neighborhoodHandler.Update)
	suite.echo.DELETE("/api/v1/neighborhoods/:id", neighborhoodHandler.Delete)
	suite.echo.GET("/api/v1/neighborhoods", neighborhoodHandler.List)

	suite.echo.POST("/api/v1/buildings", buildingHandler.Create)
	suite.echo.GET("/api/v1/buildings/:id", buildingHandler.Get)
	suite.echo.PUT("/api/v1/buildings/:id", buildingHandler.Update)
	suite.echo.DELETE("/api/v1/buildings/:id", buildingHandler.Delete)
	suite.echo.GET("/api/v1/buildings", buildingHandler.List)
}

func (suite *E2ETestSuite) TearDownTest() {
	// Clean up test data after each test
	_, err := suite.db.Exec("TRUNCATE TABLE buildings, neighborhoods RESTART IDENTITY")
	suite.NoError(err)
}

// Helper to create a neighborhood and return its ID
func (suite *E2ETestSuite) createNeighborhood(name string) string {
	createReq := map[string]string{"name": name}
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/neighborhoods", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	var created map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &created) // assume success in helper
	return created["id"]
}

// Helper to create a building and return its ID
func (suite *E2ETestSuite) createBuilding(name string, neighborhoodID string, address string) string {
	createReq := map[string]string{"name": name, "neighborhood_id": neighborhoodID, "address": address}
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	var created map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &created) // assume success in helper
	return created["id"]
}

func (suite *E2ETestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *E2ETestSuite) TestCreateNeighborhood() {
	createReq := map[string]string{"name": "Test Neighborhood"}
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/neighborhoods", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	var created map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &created)
	suite.NoError(err)
	assert.NotEmpty(suite.T(), created["id"])
	assert.Equal(suite.T(), "Test Neighborhood", created["name"])
}

func (suite *E2ETestSuite) TestCreateNeighborhood_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/neighborhoods", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *E2ETestSuite) TestCreateNeighborhood_MissingName() {
	createReq := map[string]string{}
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/neighborhoods", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *E2ETestSuite) TestGetNeighborhood() {
	// Seed: create a neighborhood
	id := suite.createNeighborhood("Seed Neighborhood")

	// Test Get
	req := httptest.NewRequest(http.MethodGet, "/api/v1/neighborhoods/"+id, nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var retrieved map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &retrieved)
	suite.NoError(err)
	assert.Equal(suite.T(), id, retrieved["id"])
	assert.Equal(suite.T(), "Seed Neighborhood", retrieved["name"])
}

func (suite *E2ETestSuite) TestGetNeighborhood_NotFound() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/neighborhoods/11111111-1111-1111-1111-111111111111", nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *E2ETestSuite) TestGetNeighborhood_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/neighborhoods/invalid", nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *E2ETestSuite) TestUpdateNeighborhood() {
	// Seed: create a neighborhood
	id := suite.createNeighborhood("Original Neighborhood")

	// Test Update
	updateReq := map[string]string{"name": "Updated Neighborhood"}
	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/neighborhoods/"+id, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var updated map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &updated)
	suite.NoError(err)
	assert.Equal(suite.T(), id, updated["id"])
	assert.Equal(suite.T(), "Updated Neighborhood", updated["name"])
}

func (suite *E2ETestSuite) TestUpdateNeighborhood_NotFound() {
	updateReq := map[string]string{"name": "Updated Neighborhood"}
	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/neighborhoods/11111111-1111-1111-1111-111111111111", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *E2ETestSuite) TestDeleteNeighborhood() {
	// Seed: create a neighborhood
	id := suite.createNeighborhood("Delete Neighborhood")

	// Test Delete
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/neighborhoods/"+id, nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNoContent, rec.Code)

	// Verify deleted: get should return 404
	req = httptest.NewRequest(http.MethodGet, "/api/v1/neighborhoods/"+id, nil)
	rec = httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *E2ETestSuite) TestDeleteNeighborhood_NotFound() {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/neighborhoods/11111111-1111-1111-1111-111111111111", nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *E2ETestSuite) TestListNeighborhoods() {
	// Seed: create multiple neighborhoods
	names := []string{"Neighborhood 1", "Neighborhood 2", "Neighborhood 3"}
	for _, name := range names {
		suite.createNeighborhood(name)
	}

	// Test List
	req := httptest.NewRequest(http.MethodGet, "/api/v1/neighborhoods", nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var list []map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &list)
	suite.NoError(err)
	assert.Len(suite.T(), list, len(names))
	// Check names are present (order may vary)
	retrievedNames := make(map[string]bool)
	for _, item := range list {
		retrievedNames[item["name"]] = true
	}
	for _, name := range names {
		assert.True(suite.T(), retrievedNames[name])
	}
}

func (suite *E2ETestSuite) TestCreateBuilding() {
	// Seed: create a neighborhood first
	neighborhoodID := suite.createNeighborhood("Test Neighborhood")

	createReq := map[string]string{"name": "Test Building", "neighborhood_id": neighborhoodID, "address": "123 Test St"}
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	var created map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &created)
	suite.NoError(err)
	assert.NotEmpty(suite.T(), created["id"])
	assert.Equal(suite.T(), "Test Building", created["name"])
	assert.Equal(suite.T(), neighborhoodID, created["neighborhood_id"])
	assert.Equal(suite.T(), "123 Test St", created["address"])
}

func (suite *E2ETestSuite) TestCreateBuilding_InvalidNeighborhood() {
	createReq := map[string]string{"name": "Test Building", "neighborhood_id": "11111111-1111-1111-1111-111111111111", "address": "123 Test St"}
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *E2ETestSuite) TestCreateBuilding_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *E2ETestSuite) TestCreateBuilding_MissingFields() {
	createReq := map[string]string{"name": "Test Building"}
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *E2ETestSuite) TestGetBuilding() {
	// Create neighborhood and building
	neighborhoodID := suite.createNeighborhood("Test Neighborhood")
	id := suite.createBuilding("Test Building", neighborhoodID, "123 Test St")

	// Test Get
	req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/"+id, nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var retrieved map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &retrieved)
	suite.NoError(err)
	assert.Equal(suite.T(), id, retrieved["id"])
	assert.Equal(suite.T(), "Test Building", retrieved["name"])
	assert.Equal(suite.T(), neighborhoodID, retrieved["neighborhood_id"])
	assert.Equal(suite.T(), "123 Test St", retrieved["address"])
}

func (suite *E2ETestSuite) TestGetBuilding_NotFound() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/11111111-1111-1111-1111-111111111111", nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *E2ETestSuite) TestGetBuilding_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/invalid", nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *E2ETestSuite) TestUpdateBuilding() {
	// Create neighborhood and building
	neighborhoodID := suite.createNeighborhood("Test Neighborhood")
	id := suite.createBuilding("Original Building", neighborhoodID, "123 Original St")

	// Test Update
	updateReq := map[string]string{"name": "Updated Building", "neighborhood_id": neighborhoodID, "address": "456 Updated St"}
	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/buildings/"+id, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var updated map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &updated)
	suite.NoError(err)
	assert.Equal(suite.T(), id, updated["id"])
	assert.Equal(suite.T(), "Updated Building", updated["name"])
	assert.Equal(suite.T(), neighborhoodID, updated["neighborhood_id"])
	assert.Equal(suite.T(), "456 Updated St", updated["address"])
}

func (suite *E2ETestSuite) TestUpdateBuilding_NotFound() {
	neighborhoodID := suite.createNeighborhood("Test Neighborhood")
	updateReq := map[string]string{"name": "Updated Building", "neighborhood_id": neighborhoodID, "address": "456 Updated St"}
	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/buildings/11111111-1111-1111-1111-111111111111", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *E2ETestSuite) TestDeleteBuilding() {
	// Create neighborhood and building
	neighborhoodID := suite.createNeighborhood("Test Neighborhood")
	id := suite.createBuilding("Delete Building", neighborhoodID, "123 Delete St")

	// Test Delete
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/buildings/"+id, nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNoContent, rec.Code)

	// Verify deleted: get should return 404
	req = httptest.NewRequest(http.MethodGet, "/api/v1/buildings/"+id, nil)
	rec = httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *E2ETestSuite) TestDeleteBuilding_NotFound() {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/buildings/11111111-1111-1111-1111-111111111111", nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *E2ETestSuite) TestListBuildings() {
	// Create neighborhood and multiple buildings
	neighborhoodID := suite.createNeighborhood("Test Neighborhood")
	names := []string{"Building 1", "Building 2", "Building 3"}
	for _, name := range names {
		suite.createBuilding(name, neighborhoodID, "123 Test St")
	}

	// Test List
	req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings", nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var list []map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &list)
	suite.NoError(err)
	assert.Len(suite.T(), list, len(names))
	// Check names are present (order may vary)
	retrievedNames := make(map[string]bool)
	for _, item := range list {
		retrievedNames[item["name"]] = true
	}
	for _, name := range names {
		assert.True(suite.T(), retrievedNames[name])
	}
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
