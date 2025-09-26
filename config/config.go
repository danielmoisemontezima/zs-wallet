package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	APIPort    string
	JWTSecret  string
}

// LoadConfig loads configuration from .env file
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "wallet"),
		APIPort:    getEnv("API_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
