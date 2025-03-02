package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig holds standard test configurations
type TestConfig struct {
	name    string
	envVars map[string]string
}

// getBaseConfig returns a base configuration for testing
func getBaseConfig() map[string]string {
	return map[string]string{
		// Server
		"APP_ENV":              "development",
		"SERVER_PORT":          "8080",
		"SERVER_HOST":          "localhost",
		"SERVER_READ_TIMEOUT":  "60s",
		"SERVER_WRITE_TIMEOUT": "60s",
		"SERVER_BASE_URL":      "http://localhost:8080",

		// Database
		"DB_HOST":              "localhost",
		"DB_PORT":              "5432",
		"DB_USER":              "dev_user",
		"DB_PASSWORD":          "dev_password",
		"DB_NAME":              "taskdb_dev",
		"DB_SSL_MODE":          "disable",
		"DB_MAX_OPEN_CONNS":    "10",
		"DB_MAX_IDLE_CONNS":    "5",
		"DB_CONN_MAX_LIFETIME": "5m",

		// API
		"API_BASE_PATH":        "/api/v1",
		"API_KEY":              "dev_12345678901234567890123456789012",
		"CORS_ALLOWED_ORIGINS": "*",

		// Logging
		"LOG_LEVEL":  "debug",
		"LOG_FORMAT": "text",

		// Features
		"ENABLE_SWAGGER":  "true",
		"ENABLE_METRICS":  "true",
		"REPOSITORY_TYPE": "memory",
	}
}

// getEnvironmentConfigs returns environment-specific configurations
func getEnvironmentConfigs() []TestConfig {
	return []TestConfig{
		{
			name:    "development",
			envVars: getBaseConfig(), // Development uses base config
		},
		{
			name: "staging",
			envVars: copyAndModify(getBaseConfig(), map[string]string{
				"APP_ENV":              "staging",
				"SERVER_HOST":          "0.0.0.0",
				"SERVER_READ_TIMEOUT":  "30s",
				"SERVER_WRITE_TIMEOUT": "30s",
				"SERVER_BASE_URL":      "https://staging-api.example.com",
				"DB_HOST":              "staging-db.example.com",
				"DB_USER":              "staging_user",
				"DB_PASSWORD":          "staging_password",
				"DB_NAME":              "taskdb_staging",
				"DB_SSL_MODE":          "verify-full",
				"DB_MAX_OPEN_CONNS":    "25",
				"DB_MAX_IDLE_CONNS":    "10",
				"DB_CONN_MAX_LIFETIME": "10m",
				"API_KEY":              "staging_8901234567890123456789012345678",
				"CORS_ALLOWED_ORIGINS": "https://*.example.com",
				"LOG_LEVEL":            "info",
				"LOG_FORMAT":           "json",
				"REPOSITORY_TYPE":      "postgres",
			}),
		},
		{
			name: "production",
			envVars: copyAndModify(getBaseConfig(), map[string]string{
				"APP_ENV":              "production",
				"SERVER_HOST":          "0.0.0.0",
				"SERVER_READ_TIMEOUT":  "15s",
				"SERVER_WRITE_TIMEOUT": "15s",
				"SERVER_BASE_URL":      "https://api.example.com",
				"DB_HOST":              "prod-db.example.com",
				"DB_USER":              "prod_user",
				"DB_PASSWORD":          "prod_password",
				"DB_NAME":              "taskdb_prod",
				"DB_SSL_MODE":          "verify-full",
				"DB_MAX_OPEN_CONNS":    "100",
				"DB_MAX_IDLE_CONNS":    "25",
				"API_KEY":              "prod_23456789012345678901234567890123",
				"CORS_ALLOWED_ORIGINS": "https://example.com",
				"LOG_LEVEL":            "warn",
				"LOG_FORMAT":           "json",
				"ENABLE_SWAGGER":       "false",
				"REPOSITORY_TYPE":      "postgres",
			}),
		},
	}
}

// copyAndModify creates a deep copy of the base config and applies modifications
func copyAndModify(base map[string]string, modifications map[string]string) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range modifications {
		result[k] = v
	}
	return result
}

// setupTestEnvironment prepares the environment for a test
func setupTestEnvironment(t *testing.T, envVars map[string]string) {
	t.Helper()

	// Save GO_ENV value
	goEnv := os.Getenv("GO_ENV")

	// Clear other environment variables
	os.Clearenv()

	// Restore GO_ENV
	if err := os.Setenv("GO_ENV", goEnv); err != nil {
		t.Fatalf("Failed to restore GO_ENV: %v", err)
	}

	// Set test environment variables
	for k, v := range envVars {
		err := os.Setenv(k, v)
		require.NoError(t, err)
	}
}

func TestEnvironmentSpecificConfigurations(t *testing.T) {
	envConfigs := getEnvironmentConfigs()
	for _, env := range envConfigs {
		t.Run(env.name, func(t *testing.T) {
			setupTestEnvironment(t, env.envVars)

			cfg, err := Load()
			require.NoError(t, err)
			require.NotNil(t, cfg)

			// Assert environment-specific configurations
			assert.Equal(t, env.envVars["SERVER_PORT"], cfg.Server.Port)
			assert.Equal(t, env.envVars["SERVER_HOST"], cfg.Server.Host)
			assert.Equal(t, env.envVars["LOG_LEVEL"], cfg.Logging.Level)
			assert.Equal(t, env.envVars["LOG_FORMAT"], cfg.Logging.Format)
			assert.Equal(t, env.envVars["DB_SSL_MODE"], cfg.Database.SSLMode)
			assert.Equal(t, env.envVars["REPOSITORY_TYPE"], cfg.Repository.Type)
			assert.Equal(t, env.envVars["CORS_ALLOWED_ORIGINS"], cfg.API.AllowedOrigins)
			assert.Equal(t, "[REDACTED]", cfg.Database.Password.String())
			assert.Equal(t, "[REDACTED]", cfg.API.APIKey.String())
		})
	}
}

func TestEnvironmentSpecificValidation(t *testing.T) {
	tests := []struct {
		name          string
		baseEnv       string
		modifications map[string]string
		expectedError bool
		errorMessage  string
	}{
		{
			name:    "development with localhost DB",
			baseEnv: "development",
			modifications: map[string]string{
				"DB_HOST": "localhost",
			},
			expectedError: false,
		},
		{
			name:    "production requires SSL",
			baseEnv: "production",
			modifications: map[string]string{
				"DB_SSL_MODE": "disable",
			},
			expectedError: true,
			errorMessage:  "SSL mode must be verify-full in production",
		},
		{
			name:    "production requires strict CORS",
			baseEnv: "production",
			modifications: map[string]string{
				"CORS_ALLOWED_ORIGINS": "*",
			},
			expectedError: true,
			errorMessage:  "wildcard CORS not allowed in production",
		},
	}

	envConfigs := getEnvironmentConfigs()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Find base environment configuration
			var baseEnv TestConfig
			for _, env := range envConfigs {
				if env.name == tt.baseEnv {
					baseEnv = env
					break
				}
			}
			require.NotEmpty(t, baseEnv.envVars, "Test environment not found")

			// Apply modifications to base config
			testConfig := copyAndModify(baseEnv.envVars, tt.modifications)
			setupTestEnvironment(t, testConfig)

			// Load configuration
			cfg, err := Load()

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name          string
		modifications map[string]string
		expectedError bool
		errorMessage  string
	}{
		{
			name:          "valid configuration",
			modifications: nil,
			expectedError: false,
		},
		{
			name: "invalid URL",
			modifications: map[string]string{
				"SERVER_BASE_URL": "not-a-url",
			},
			expectedError: true,
			errorMessage:  "failed on the 'url' tag",
		},
		{
			name: "short API key",
			modifications: map[string]string{
				"API_KEY": "short",
			},
			expectedError: true,
			errorMessage:  "min",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testConfig := getBaseConfig()
			if tt.modifications != nil {
				testConfig = copyAndModify(testConfig, tt.modifications)
			}

			setupTestEnvironment(t, testConfig)

			cfg, err := Load()
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg)
		})
	}
}

func TestSensitiveValue(t *testing.T) {
	password := SensitiveValue("secret")
	assert.Equal(t, "[REDACTED]", password.String())
}

func TestMain(m *testing.M) {
	// Ensure we're in test mode
	err := os.Setenv("GO_ENV", "test")
	if err != nil {
		fmt.Printf("Failed to set GO_ENV: %v\n", err)
		os.Exit(1)
	}

	// Run the tests
	code := m.Run()

	// Clean up
	os.Unsetenv("GO_ENV")
	os.Exit(code)
}

// Let's also add a test to verify GO_ENV is set correctly
func TestGoEnvIsSet(t *testing.T) {
	env := os.Getenv("GO_ENV")
	assert.Equal(t, "test", env, "GO_ENV should be set to 'test'")
}
