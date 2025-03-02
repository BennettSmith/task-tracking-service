package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// SensitiveValue is used for fields that shouldn't be logged
type SensitiveValue string

func (s SensitiveValue) String() string {
	return "[REDACTED]"
}

// Config struct with proper validation tags
type Config struct {
	Environment string           `validate:"required,oneof=development staging production"`
	Server      ServerConfig     `validate:"required"`
	Database    DatabaseConfig   `validate:"required"`
	API         APIConfig        `validate:"required"`
	Logging     LogConfig        `validate:"required"`
	Features    FeatureConfig    `validate:"required"`
	Repository  RepositoryConfig `validate:"required"`
}

type ServerConfig struct {
	Port         string `validate:"required,numeric"`
	Host         string `validate:"required"`
	ReadTimeout  string `validate:"required"`
	WriteTimeout string `validate:"required"`
	BaseURL      string `validate:"required,url"`
}

type DatabaseConfig struct {
	Host            string         `validate:"required"`
	Port            string         `validate:"required,numeric"`
	User            string         `validate:"required"`
	Password        SensitiveValue `validate:"required"`
	Name            string         `validate:"required"`
	SSLMode         string         `validate:"required,oneof=disable enable verify-full"`
	MaxOpenConns    int            `validate:"required,min=1"`
	MaxIdleConns    int            `validate:"required,min=1"`
	ConnMaxLifetime string         `validate:"required"`
}

type APIConfig struct {
	BasePath       string         `validate:"required"`
	APIKey         SensitiveValue `validate:"required,min=32"`
	AllowedOrigins string         `validate:"required"`
}

type LogConfig struct {
	Level  string `validate:"required,oneof=debug info warn error"`
	Format string `validate:"required,oneof=text json"`
}

type FeatureConfig struct {
	EnableSwagger bool
	EnableMetrics bool
}

type RepositoryConfig struct {
	Type string `validate:"required,oneof=memory postgres"`
}

// ValidationError holds validation error details
type ValidationError struct {
	Field string
	Error string
}

// Validate performs configuration validation
func (c *Config) Validate() error {
	validate := validator.New()

	// Register custom validation for environment-specific rules
	if err := validate.RegisterValidation("ssl_mode_production", validateSSLModeProduction); err != nil {
		return fmt.Errorf("failed to register ssl mode validator: %w", err)
	}

	if err := validate.RegisterValidation("cors_production", validateCORSProduction); err != nil {
		return fmt.Errorf("failed to register CORS validator: %w", err)
	}

	// First validate struct fields based on tags
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Then perform environment-specific validation
	if c.Environment == "production" {
		// Validate SSL mode in production
		if c.Database.SSLMode != "verify-full" {
			return fmt.Errorf("SSL mode must be verify-full in production")
		}

		// Validate CORS in production
		if c.API.AllowedOrigins == "*" {
			return fmt.Errorf("wildcard CORS not allowed in production")
		}
	}

	return nil
}

// Custom validator functions
func validateSSLModeProduction(fl validator.FieldLevel) bool {
	config, ok := fl.Parent().Interface().(Config)
	if !ok {
		return false
	}
	if config.Environment == "production" {
		return fl.Field().String() == "verify-full"
	}
	return true
}

func validateCORSProduction(fl validator.FieldLevel) bool {
	config, ok := fl.Parent().Interface().(Config)
	if !ok {
		return false
	}
	if config.Environment == "production" {
		return fl.Field().String() != "*"
	}
	return true
}

// Change setDefaults to accept a Viper instance
func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("SERVER_HOST", "localhost")
	v.SetDefault("SERVER_READ_TIMEOUT", "60s")
	v.SetDefault("SERVER_WRITE_TIMEOUT", "60s")
	v.SetDefault("SERVER_BASE_URL", "http://localhost:8080")

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNS", 10)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "5m")

	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "text")

	v.SetDefault("ENABLE_SWAGGER", true)
	v.SetDefault("ENABLE_METRICS", true)

	v.SetDefault("REPOSITORY_TYPE", "memory")
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	config := &Config{}

	// Initialize viper
	v := viper.New()

	// Set defaults using the same Viper instance
	setDefaults(v)

	// Only attempt to load .env file if not in test mode
	if os.Getenv("GO_ENV") != "test" {
		// Don't set config file at all in test mode
		v.SetConfigFile(".env")
		if err := v.ReadInConfig(); err != nil {
			// Only return error if file exists but cannot be read
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	// Set up environment variables
	v.AutomaticEnv()

	// Use the local viper instance instead of the global one
	config.Environment = v.GetString("APP_ENV")
	config.Server.Port = v.GetString("SERVER_PORT")
	config.Server.Host = v.GetString("SERVER_HOST")
	config.Server.ReadTimeout = v.GetString("SERVER_READ_TIMEOUT")
	config.Server.WriteTimeout = v.GetString("SERVER_WRITE_TIMEOUT")
	config.Server.BaseURL = v.GetString("SERVER_BASE_URL")

	config.Database.Host = v.GetString("DB_HOST")
	config.Database.Port = v.GetString("DB_PORT")
	config.Database.User = v.GetString("DB_USER")
	config.Database.Password = SensitiveValue(v.GetString("DB_PASSWORD"))
	config.Database.Name = v.GetString("DB_NAME")
	config.Database.SSLMode = v.GetString("DB_SSL_MODE")
	config.Database.MaxOpenConns = v.GetInt("DB_MAX_OPEN_CONNS")
	config.Database.MaxIdleConns = v.GetInt("DB_MAX_IDLE_CONNS")
	config.Database.ConnMaxLifetime = v.GetString("DB_CONN_MAX_LIFETIME")

	config.API.BasePath = v.GetString("API_BASE_PATH")
	config.API.APIKey = SensitiveValue(v.GetString("API_KEY"))
	config.API.AllowedOrigins = v.GetString("CORS_ALLOWED_ORIGINS")

	config.Logging.Level = v.GetString("LOG_LEVEL")
	config.Logging.Format = v.GetString("LOG_FORMAT")

	config.Features.EnableSwagger = v.GetBool("ENABLE_SWAGGER")
	config.Features.EnableMetrics = v.GetBool("ENABLE_METRICS")

	config.Repository.Type = v.GetString("REPOSITORY_TYPE")

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}
