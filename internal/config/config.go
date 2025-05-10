package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIConfig        APIConfig
	PostgresConfig   PostgresConfig
	GoogleOIDCConfig GoogleOIDCConfig
}

type APIConfig struct {
	BaseURL string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

type GoogleOIDCConfig struct {
	ClientID     string
	ClientSecret string
}

func (c *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.DbName)
}

func LoadConfig() (Config, error) {
	env := os.Getenv("OIDC_DEMO_ENV")
	if env == "" {
		env = "development"
	}

	// First, load environment-specific config,
	// then load default (existing values loaded
	// take precedence over the default values).
	err := godotenv.Load(fmt.Sprintf(".env.%s", env))
	if err != nil {
		return Config{}, fmt.Errorf("error loading env file: %w", err)
	}

	err = godotenv.Load(".env")
	if err != nil {
		return Config{}, fmt.Errorf("error loading env file: %w", err)
	}

	return Config{
		APIConfig: APIConfig{
			BaseURL: os.Getenv("OIDC_DEMO_API_BASE_URL"),
		},
		PostgresConfig: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DbName:   os.Getenv("POSTGRES_DB"),
		},
		GoogleOIDCConfig: GoogleOIDCConfig{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		},
	}, nil
}
