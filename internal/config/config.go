package config

import "os"

type Config struct {
	Server struct {
		Port string
		Host string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	JWT struct {
		Secret        string
		ExpiryHours   int
		RefreshSecret string // A separate secret key for signing refresh tokens
	}
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	// Server config
	config.Server.Host = getEnvOrDefault("SERVER_HOST", "0.0.0.0")
	config.Server.Port = getEnvOrDefault("SERVER_PORT", "50051")

	// Database config
	config.Database.Host = getEnvOrDefault("DB_HOST", "localhost")
	config.Database.Port = getEnvOrDefault("DB_PORT", "5432")
	config.Database.User = getEnvOrDefault("DB_USER", "postgres")
	config.Database.Password = getEnvOrDefault("DB_PASSWORD", "password")
	config.Database.DBName = getEnvOrDefault("DB_NAME", "rmshop")
	config.Database.SSLMode = getEnvOrDefault("DB_SSLMODE", "disable")

	// JWT config
	config.JWT.Secret = getEnvOrDefault("JWT_SECRET", "fawas's-secret-key")
	config.JWT.ExpiryHours = 24
	config.JWT.RefreshSecret = getEnvOrDefault("JWT_REFRESH_SECRET", "fawas's-refresh-secret-key")

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
