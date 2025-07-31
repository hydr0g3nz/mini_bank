package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hydr0g3nz/mini_bank/internal/infrastructure"
	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	Database infrastructure.DBConfig
	Cache    CacheConfig
	API      APIConfig
	LogLevel string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Host         string
	Environment  string
	ReadTimeout  int // in seconds
	WriteTimeout int // in seconds
	IdleTimeout  int // in seconds
}

// CacheConfig holds Redis cache configuration
type CacheConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// APIConfig holds API configuration
type APIConfig struct {
	Key string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "localhost"),
			Port:         getEnv("PORT", "8080"),
			Environment:  getEnv("GIN_MODE", "debug"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 30),  // 30 seconds
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 30), // 30 seconds
			IdleTimeout:  getEnvAsInt("SERVER_IDLE_TIMEOUT", 60),  // 60 seconds
		},
		Database: infrastructure.DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "mini_bank"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Cache: CacheConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		API: APIConfig{
			Key: getEnv("API_KEY", "your-secret-api-key-change-in-production"),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "release"
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "debug"
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	if c.Server.Host == "" || c.Server.Host == "localhost" {
		return ":" + c.Server.Port
	}
	return c.Server.Host + ":" + c.Server.Port
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.API.Key == "" || c.API.Key == "your-secret-api-key-change-in-production" {
		if c.IsProduction() {
			return fmt.Errorf("API_KEY must be set in production environment")
		}
	}

	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	return nil
}

// getEnvAsInt gets an environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnv gets an environment variable as a string
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
