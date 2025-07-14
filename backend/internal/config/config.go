package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	API      APIConfig
	JWT      JWTConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// APIConfig holds API configuration
type APIConfig struct {
	Key string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "seo_user"),
			Password: getEnv("DB_PASSWORD", "seo_password"),
			Database: getEnv("DB_NAME", "seo_crawler"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		API: APIConfig{
			Key: getEnv("API_KEY", "seo-crawler-api-key-2025"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		},
	}
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
	)
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 