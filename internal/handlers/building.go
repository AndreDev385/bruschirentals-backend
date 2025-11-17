// Package handlers provides HTTP handlers for the API.
package handlers

import (
	"net/http"

	"github.com/Andre385/bruschirentals-backend/internal/services"
	"github.com/labstack/echo/v4"
)

// Building represents a building in the API.
type Building struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	NeighborhoodID string `json:"neighborhood_id"`
	Address        string `json:"address"`
}

// BuildingHandler handles building-related HTTP requests.
type BuildingHandler struct {
	service *services.BuildingService
}

// NewBuildingHandler creates a new building handler.
func NewBuildingHandler(service *services.BuildingService) *BuildingHandler {
	return &BuildingHandler{service: service}
}

// Create handles POST /api/v1/buildings
// @Summary Create a new building
// @Description Create a new building with the given details
// @Tags buildings
// @Accept json
// @Produce json
// @Param request body map[string]string true "Building details"
// @Success 201 {object} Building
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/buildings [post]
func (h *BuildingHandler) Create(c echo.Context) error {
	var req struct {
		Name           string `json:"name"`
		NeighborhoodID string `json:"neighborhood_id"`
		Address        string `json:"address"`
	}
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid request")
	}

	building, err := h.service.CreateBuilding(c.Request().Context(), req.Name, req.NeighborhoodID, req.Address)
	if err != nil {
		status, message := mapErrorToResponse(err)
		return SendError(c, status, message)
	}

	return c.JSON(http.StatusCreated, building)
}

// Get handles GET /api/v1/buildings/:id
// @Summary Get a building by ID
// @Description Retrieve a building by its ID
// @Tags buildings
// @Produce json
// @Param id path string true "Building ID"
// @Success 200 {object} Building
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/buildings/{id} [get]
func (h *BuildingHandler) Get(c echo.Context) error {
	id := c.Param("id")

	building, err := h.service.GetBuilding(c.Request().Context(), id)
	if err != nil {
		status, message := mapErrorToResponse(err)
		return SendError(c, status, message)
	}

	return c.JSON(http.StatusOK, building)
}

// Update handles PUT /api/v1/buildings/:id
// @Summary Update a building
// @Description Update an existing building's details
// @Tags buildings
// @Accept json
// @Produce json
// @Param id path string true "Building ID"
// @Param request body map[string]string true "Updated building details"
// @Success 200 {object} Building
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/buildings/{id} [put]
func (h *BuildingHandler) Update(c echo.Context) error {
	id := c.Param("id")

	var req struct {
		Name           string `json:"name"`
		NeighborhoodID string `json:"neighborhood_id"`
		Address        string `json:"address"`
	}
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid request")
	}

	building, err := h.service.UpdateBuilding(c.Request().Context(), id, req.Name, req.NeighborhoodID, req.Address)
	if err != nil {
		status, message := mapErrorToResponse(err)
		return SendError(c, status, message)
	}

	return c.JSON(http.StatusOK, building)
}

// Delete handles DELETE /api/v1/buildings/:id
// @Summary Delete a building
// @Description Delete a building by its ID
// @Tags buildings
// @Param id path string true "Building ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/buildings/{id} [delete]
func (h *BuildingHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	err := h.service.DeleteBuilding(c.Request().Context(), id)
	if err != nil {
		status, message := mapErrorToResponse(err)
		return SendError(c, status, message)
	}

	return c.NoContent(http.StatusNoContent)
}

// List handles GET /api/v1/buildings
// @Summary List all buildings
// @Description Retrieve a list of all buildings
// @Tags buildings
// @Produce json
// @Success 200 {array} Building
// @Failure 500 {object} map[string]string
// @Router /api/v1/buildings [get]
func (h *BuildingHandler) List(c echo.Context) error {
	buildings, err := h.service.ListBuildings(c.Request().Context())
	if err != nil {
		status, message := mapErrorToResponse(err)
		return SendError(c, status, message)
	}

	return c.JSON(http.StatusOK, buildings)
}
