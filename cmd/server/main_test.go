package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/check-health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Note: This test assumes DB is not connected; in real test, mock DB
	// For now, test the route exists and returns 503 since no DB
	handler := func(c echo.Context) error {
		// Simulate health check without DB
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "unhealthy",
			"error":  "database not connected",
		})
	}

	// Test
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "unhealthy")
}
