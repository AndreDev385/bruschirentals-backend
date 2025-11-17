// Package handlers contains HTTP request handlers.
package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	DB     *sqlx.DB
	Logger *zap.Logger
	Tracer trace.Tracer
}

// NewHealthHandler creates a new HealthHandler with dependencies.
func NewHealthHandler(db *sqlx.DB, logger *zap.Logger, tracer trace.Tracer) *HealthHandler {
	return &HealthHandler{
		DB:     db,
		Logger: logger,
		Tracer: tracer,
	}
}

// CheckHealth performs a health check on the database.
// @Summary Health check
// @Description Check database connectivity
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "status: healthy"
// @Failure 503 {object} ErrorResponse "Database unhealthy"
// @Router /v1/health [get]
func (h *HealthHandler) CheckHealth(c echo.Context) error {
	var span trace.Span
	if h.Tracer != nil {
		_, span = h.Tracer.Start(c.Request().Context(), "health-check")
		defer span.End()
	}

	if err := h.DB.Ping(); err != nil {
		h.Logger.Error("Health check failed", zap.Error(err))
		if span != nil {
			span.RecordError(err)
		}
		return SendError(c, http.StatusServiceUnavailable, "database unhealthy: "+err.Error())
	}
	h.Logger.Info("Health check passed")
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}
