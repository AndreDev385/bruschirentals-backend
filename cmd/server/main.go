// Package main is the entry point for the Bruschi Rentals backend server.
// @title Bruschi Rentals API
// @version 1.0
// @description API for managing rentals and clients.
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Andre385/bruschirentals-backend/docs"
	"github.com/Andre385/bruschirentals-backend/internal/config"
	"github.com/Andre385/bruschirentals-backend/internal/handlers"
	"github.com/Andre385/bruschirentals-backend/internal/logging"
	"github.com/Andre385/bruschirentals-backend/internal/middleware"
	"github.com/Andre385/bruschirentals-backend/internal/repositories"
	"github.com/Andre385/bruschirentals-backend/internal/services"
	"github.com/Andre385/bruschirentals-backend/internal/tracing"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize tracing if OTLP endpoint is set
	var tp *sdktrace.TracerProvider
	if cfg.OTLPEndpoint != "" {
		tp = tracing.InitTracer()
	}

	// Initialize logger
	logger := logging.InitLogger()
	defer logger.Sync()

	// Database connection
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	defer db.Close()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	e := echo.New()

	// Add middlewares
	e.Use(echomw.CORS())
	e.Use(echomw.Recover())

	// Add tracing middleware if tracing enabled
	if cfg.OTLPEndpoint != "" {
		e.Use(middleware.Tracing())
	}

	// Add logging middleware
	e.Use(middleware.Logging(logger))

	// Initialize repositories
	neighborhoodRepo := repositories.NewNeighborhoodRepository(db)
	buildingRepo := repositories.NewBuildingRepository(db)

	// Initialize services
	neighborhoodService := services.NewNeighborhoodService(neighborhoodRepo)
	buildingService := services.NewBuildingService(buildingRepo, neighborhoodRepo)

	// Initialize handlers
	var tracer trace.Tracer
	if cfg.OTLPEndpoint != "" {
		tracer = otel.Tracer("health-handler")
	}
	healthHandler := handlers.NewHealthHandler(db, logger, tracer)
	neighborhoodHandler := handlers.NewNeighborhoodHandler(neighborhoodService)
	buildingHandler := handlers.NewBuildingHandler(buildingService)

	e.GET("/api/v1/health", healthHandler.CheckHealth)

	// Neighborhood routes
	e.POST("/api/v1/neighborhoods", neighborhoodHandler.Create)
	e.GET("/api/v1/neighborhoods/:id", neighborhoodHandler.Get)
	e.PUT("/api/v1/neighborhoods/:id", neighborhoodHandler.Update)
	e.DELETE("/api/v1/neighborhoods/:id", neighborhoodHandler.Delete)
	e.GET("/api/v1/neighborhoods", neighborhoodHandler.List)

	// Building routes
	e.POST("/api/v1/buildings", buildingHandler.Create)
	e.GET("/api/v1/buildings/:id", buildingHandler.Get)
	e.PUT("/api/v1/buildings/:id", buildingHandler.Update)
	e.DELETE("/api/v1/buildings/:id", buildingHandler.Delete)
	e.GET("/api/v1/buildings", buildingHandler.List)

	// Swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server in goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Port))
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for signal
	<-sigChan
	logger.Info("Shutting down server...")

	// Shutdown server
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error("Failed to shutdown server", zap.Error(err))
	}

	// Shutdown tracer
	if tp != nil {
		if err := tp.Shutdown(shutdownCtx); err != nil {
			logger.Error("Failed to shutdown tracer", zap.Error(err))
		}
	}

	logger.Info("Server stopped")
}
