// Package config handles application configuration loading.
package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	DatabaseURL  string `mapstructure:"DATABASE_URL" validate:"required"`
	Port         string `mapstructure:"PORT" validate:"required"`
	OTLPEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	Environment  string `mapstructure:"ENV" validate:"required"`
}

// Validate checks the configuration for required fields.
func (c *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// Load reads configuration from environment and .env file.
func Load() (*Config, error) {
	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")

	// Load .env file if exists
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}
