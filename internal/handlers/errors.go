package handlers

import (
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
