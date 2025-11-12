package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/stevenwijaya/finance-tracker/pkg/log"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	JWTSecret  string
	AppPort    string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Warn(".env file not found, using default environment variables")
	}

	config := Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "admin"),
		DBName:     getEnv("DB_NAME", "finance_tracker"),
		DBPort:     getEnv("DB_PORT", "5432"),
		JWTSecret:  getEnv("JWT_SECRET", "your_jwt_secret_key"),
		AppPort:    getEnv("APP_PORT", "8081"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
