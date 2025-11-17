// Package handlers provides HTTP handlers for the API.
package handlers

import (
	"errors"
	"net/http"

	apperrors "github.com/Andre385/bruschirentals-backend/internal/errors"
	"github.com/Andre385/bruschirentals-backend/internal/services"
	"github.com/labstack/echo/v4"
)

// mapErrorToResponse maps service errors to HTTP status codes and sanitized messages.
func mapErrorToResponse(err error) (int, map[string]string) {
	if errors.Is(err, apperrors.ErrInvalidID) || errors.Is(err, apperrors.ErrInvalidInput) {
		return http.StatusBadRequest, map[string]string{"error": "invalid request"}
	}
	if errors.Is(err, apperrors.ErrNotFound) {
		return http.StatusNotFound, map[string]string{"error": "not found"}
	}
	return http.StatusInternalServerError, map[string]string{"error": "internal server error"}
}

// Neighborhood represents a neighborhood location.
// @Description Neighborhood model
// @Success 200 {object} Neighborhood
type Neighborhood struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NeighborhoodHandler handles neighborhood-related HTTP requests.
type NeighborhoodHandler struct {
	service *services.NeighborhoodService
}

// NewNeighborhoodHandler creates a new neighborhood handler.
func NewNeighborhoodHandler(service *services.NeighborhoodService) *NeighborhoodHandler {
	return &NeighborhoodHandler{service: service}
}

// Create handles POST /api/v1/neighborhoods
// @Summary Create a new neighborhood
// @Description Create a new neighborhood with the given name
// @Tags neighborhoods
// @Accept json
// @Produce json
// @Param request body map[string]string true "Neighborhood name"
// @Success 201 {object} Neighborhood
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/neighborhoods [post]
func (h *NeighborhoodHandler) Create(c echo.Context) error {
	var req struct {
		Name string `json:"name" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}

	neighborhood, err := h.service.CreateNeighborhood(c.Request().Context(), req.Name)
	if err != nil {
		status, resp := mapErrorToResponse(err)
		return c.JSON(status, resp)
	}

	return c.JSON(http.StatusCreated, neighborhood)
}

// Get handles GET /api/v1/neighborhoods/:id
// @Summary Get a neighborhood by ID
// @Description Retrieve a neighborhood by its ID
// @Tags neighborhoods
// @Produce json
// @Param id path string true "Neighborhood ID"
// @Success 200 {object} Neighborhood
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/neighborhoods/{id} [get]
func (h *NeighborhoodHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	neighborhood, err := h.service.GetNeighborhood(c.Request().Context(), id)
	if err != nil {
		status, resp := mapErrorToResponse(err)
		return c.JSON(status, resp)
	}

	return c.JSON(http.StatusOK, neighborhood)
}

// Update handles PUT /api/v1/neighborhoods/:id
// @Summary Update a neighborhood
// @Description Update an existing neighborhood's name
// @Tags neighborhoods
// @Accept json
// @Produce json
// @Param id path string true "Neighborhood ID"
// @Param request body map[string]string true "Updated neighborhood name"
// @Success 200 {object} Neighborhood
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/neighborhoods/{id} [put]
func (h *NeighborhoodHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	var req struct {
		Name string `json:"name" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}

	neighborhood, err := h.service.UpdateNeighborhood(c.Request().Context(), id, req.Name)
	if err != nil {
		status, resp := mapErrorToResponse(err)
		return c.JSON(status, resp)
	}

	return c.JSON(http.StatusOK, neighborhood)
}

// Delete handles DELETE /api/v1/neighborhoods/:id
// @Summary Delete a neighborhood
// @Description Delete a neighborhood by its ID
// @Tags neighborhoods
// @Param id path string true "Neighborhood ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/neighborhoods/{id} [delete]
func (h *NeighborhoodHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	err := h.service.DeleteNeighborhood(c.Request().Context(), id)
	if err != nil {
		status, resp := mapErrorToResponse(err)
		return c.JSON(status, resp)
	}

	return c.NoContent(http.StatusNoContent)
}

// List handles GET /api/v1/neighborhoods
// @Summary List all neighborhoods
// @Description Retrieve a list of all neighborhoods
// @Tags neighborhoods
// @Produce json
// @Success 200 {array} Neighborhood
// @Failure 500 {object} map[string]string
// @Router /api/v1/neighborhoods [get]
func (h *NeighborhoodHandler) List(c echo.Context) error {
	neighborhoods, err := h.service.ListNeighborhoods(c.Request().Context())
	if err != nil {
		status, resp := mapErrorToResponse(err)
		return c.JSON(status, resp)
	}

	return c.JSON(http.StatusOK, neighborhoods)
}
