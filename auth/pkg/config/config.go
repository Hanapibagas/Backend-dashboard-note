package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  int // in seconds
	RefreshTokenExpiry int // in seconds
}

type AppConfig struct {
	Env string
}

var AppConfigInstance *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file from cmd directory
	err := godotenv.Load("cmd/.env")
	if err != nil {
		// .env file is optional in production
		fmt.Println("Warning: cmd/.env file not found, using system environment variables")
	}

	config := &Config{}

	// Database Config
	config.Database.Host = getEnv("DB_HOST", "")
	config.Database.Port = getEnv("DB_PORT", "")
	config.Database.Name = getEnv("DB_NAME", "")
	config.Database.User = getEnv("DB_USER", "")
	config.Database.Password = getEnv("DB_PASSWORD", "")

	// Server Config
	config.Server.Port = getEnv("SERVER_PORT", "")
	config.Server.GinMode = getEnv("GIN_MODE", "")

	// JWT Config
	config.JWT.Secret = getEnv("JWT_SECRET", "")
	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	config.JWT.AccessTokenExpiry, _ = strconv.Atoi(getEnv("ACCESS_TOKEN_EXPIRY", ""))   // 1 hour
	config.JWT.RefreshTokenExpiry, _ = strconv.Atoi(getEnv("REFRESH_TOKEN_EXPIRY", "")) // 7 days

	// App Config
	config.App.Env = getEnv("APP_ENV", "")

	AppConfigInstance = config
	return config, nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// GetDSN returns MySQL Data Source Name
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}
