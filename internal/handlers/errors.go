package handlers

import (
	"errors"
	"net/http"

	apperrors "github.com/Andre385/bruschirentals-backend/internal/errors"
	"github.com/labstack/echo/v4"
)

// ErrorResponse represents a standardized error response.
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// SendError sends a standardized JSON error response.
func SendError(c echo.Context, code int, message string) error {
	return c.JSON(code, ErrorResponse{
		Error: message,
		Code:  code,
	})
}

// mapErrorToResponse maps service errors to HTTP status codes and sanitized messages.
func mapErrorToResponse(err error) (int, string) {
	if errors.Is(err, apperrors.ErrInvalidID) || errors.Is(err, apperrors.ErrInvalidInput) || errors.Is(err, apperrors.ErrInvalidPriceRange) || errors.Is(err, apperrors.ErrInvalidPromotion) || errors.Is(err, apperrors.ErrInvalidApartment) {
		return http.StatusBadRequest, "invalid request"
	}
	if errors.Is(err, apperrors.ErrNotFound) {
		return http.StatusNotFound, "not found"
	}
	return http.StatusInternalServerError, "internal server error"
}
