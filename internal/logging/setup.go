// Package logging provides logger initialization.
package logging

import (
	"log"

	"go.uber.org/zap"
)

// InitLogger initializes and returns a production Zap logger.
func InitLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	return logger
}
