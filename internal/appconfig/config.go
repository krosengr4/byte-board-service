package appconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// Server configuration
	Port string `env:"PORT" envDefault:"8080"`

	// Database Configuration
	PostgresHost         string `env:"POSTGRES_HOST"`
	PostgresPort         string `env:"POSTGRES_PORT"`
	PostgresDB           string `env:"POSTGRES_DB"`
	PostgresUser         string `env:"POSTGRES_USER"`
	PostgresPasswordFile string `env:"POSTGRES_PASSWORD_FILE"`
	// PostgresPassword string `env:"POSTGRES_PASSWORD_FILE"`
	PostgresSSLMode string `env:"POSTGRES_SSL_MODE"`

	// FrontendURL string `env:"FRONTEND_URL"`

	// Allowed Origins
	AllowedOrigins string `env:"ALLOWED_ORIGINS"`

	// Secrets Configuration
	SecretsPath string `env:"SECRETS_PATH"`

	// Logging Configuration
	LogLevel  string `env:"LOG_LEVEL"`
	LogFormat string `env:"LOG_FORMAT"`
}

// Load loads the configuration from envrionment variables and .env files
func Load() (*Config, error) {
	// Lodd the env file if one exists
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("No .env file found or failed to load")
	}

	cfg := &Config{}

	// Parse env variables using caarlos0/env
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse the environment variables: %w", err)
	}

	// Validate required configurations
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	log.Info().
		Str("port:", cfg.PostgresPort).
		Str("database: ", cfg.PostgresDB).
		Str("log level: ", cfg.LogLevel).
		Msg("Configuration loaded succesfully")

	return cfg, nil
}

// Validate will validate the configuration
func (c *Config) Validate() error {
	// Check each individual database component
	if c.PostgresHost == "" {
		return fmt.Errorf("POSTGRES_HOST is required")
	}
	if c.PostgresPort == "" {
		return fmt.Errorf("POSTGRES_PORT is required")
	}
	if c.PostgresDB == "" {
		return fmt.Errorf("POSTGRES_DB is required")
	}
	if c.PostgresUser == "" {
		return fmt.Errorf("POSTGRES_USER is required")
	}
	if c.PostgresPasswordFile == "" {
		return fmt.Errorf("POSTGRES_PASSWORD_FILE is required")
	}

	// Check that SECRETS_PATH is set
	if !filepath.IsAbs(c.PostgresPasswordFile) && c.SecretsPath == "" {
		return fmt.Errorf("SECRETS_PATH is required when using relative paths for POSTGRES_PASSWORD_FILE")
	}

	return nil
}

// Constructs the database URL from individual components (No Fallbacks)
func (c *Config) GetDatabaseURL() (string, error) {
	// Always construct the URL from individual components
	password, err := c.GetDatabasePassword()
	if err != nil {
		return "", fmt.Errorf("failed to get the database password: %w", err)
	}

	// Construct the PostgresSQL connection string
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.PostgresUser, password, c.PostgresHost, c.PostgresPort, c.PostgresDB, c.PostgresSSLMode)

	log.Info().
		Str("host", c.PostgresHost).
		Str("port", c.PostgresPort).
		Str("database", c.PostgresDB).
		Str("user", c.PostgresUser).
		Str("ssl_mode", c.PostgresSSLMode).
		Msg("Constructed database URL from individual components")

	return dbURL, nil
}

func (c *Config) GetDatabasePassword() (string, error) {
	// Always require password file
	if c.PostgresPasswordFile == "" {
		return "", fmt.Errorf("POSTGRES_PASSWORD_FILE is required - could not find POSTGRES_PASSWORD_FILE")
	}

	filePath := c.PostgresPasswordFile

	// If file path is not absolute and secrets path is set, use SECRETS_PATH as base directory
	if !filepath.IsAbs(filePath) && c.SecretsPath != "" {
		filePath = filepath.Join(c.SecretsPath, filePath)
		log.Debug().
			Str("relative_path", c.PostgresPasswordFile).
			Str("secrets_path", c.SecretsPath).
			Str("full_path", filePath).
			Msg("Using relative path with SECRETS_PATH")
	} else if !filepath.IsAbs(filePath) && c.SecretsPath == "" {
		return "", fmt.Errorf("relative path provided for POSTGRES_PASSWORD_FILE but SECRETS_PATH is not set")
	}

	passwordBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read database password from file %s: %w", filePath, err)
	}
	password := string(passwordBytes)
	// Remove any white space / new lines
	password = strings.TrimSpace(password)
	log.Info().
		Str("password_file", filePath).
		Msg("Using POSTGRES_PASSWORD from file")

	return password, nil
}

// GetAllowedOrigins returns the list of allowed CORS origins
func (c *Config) GetAllowedOrigins() []string {
	if c.AllowedOrigins == "" {
		// Default to localhost for development
		return []string{"http://localhost:3000"}
	}

	// Split comma-separated origins and trim whitespace
	origins := strings.Split(c.AllowedOrigins, ",")
	result := make([]string, 0, len(origins))
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
