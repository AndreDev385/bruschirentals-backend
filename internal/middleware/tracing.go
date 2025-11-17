// Package middleware provides custom Echo middlewares.
package middleware

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

// Tracing returns an Echo middleware that adds OpenTelemetry tracing to requests.
func Tracing() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tracer := otel.Tracer("echo-server")
			ctx, span := tracer.Start(c.Request().Context(), c.Request().Method+" "+c.Request().URL.Path)
			defer span.End()

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
