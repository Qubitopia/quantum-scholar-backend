package database

import (
	"log"
	"os"
)

// Global variables for environment variables
var (
	// Email
	SMTP_HOST     string
	SMTP_PORT     string
	SMTP_USERNAME string
	SMTP_PASSWORD string
	FROM_EMAIL    string

	// PostgreSQL for Docker and go
	PGSQL_USER     string
	PGSQL_PASSWORD string
	PGSQL_NAME     string
	PGSQL_HOST     string
	PGSQL_PORT     string
	PGSQL_SSLMODE  string

	// Redis for Docker and go
	REDIS_HOST string
	REDIS_PORT string

	// JWT
	JWT_SECRET string

	// Application Configuration
	MAGIC_LINK_EXPIRY         string
	BASE_URL                  string
	API_RATE_LIMIT_PER_MINUTE string
	EMAIL_RATE_LIMIT          string

	// Razorpay Key
	RZP_KEY_ID         string
	RZP_KEY_SECRET     string
	RZP_WEBHOOK_SECRET string
)

// LoadEnvVariables loads all required environment variables into global variables
func LoadEnvVariables() {
	// Helper to load and check env var
	getEnv := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			log.Fatalf("Warning: Environment variable %s is not set", key)
		}
		return val
	}

	// Email
	SMTP_HOST = getEnv("SMTP_HOST")
	SMTP_PORT = getEnv("SMTP_PORT")
	SMTP_USERNAME = getEnv("SMTP_USERNAME")
	SMTP_PASSWORD = getEnv("SMTP_PASSWORD")
	FROM_EMAIL = getEnv("FROM_EMAIL")

	// PostgreSQL for Docker and go
	PGSQL_USER = getEnv("PGSQL_USER")
	PGSQL_PASSWORD = getEnv("PGSQL_PASSWORD")
	PGSQL_NAME = getEnv("PGSQL_NAME")
	PGSQL_HOST = getEnv("PGSQL_HOST")
	PGSQL_PORT = getEnv("PGSQL_PORT")
	PGSQL_SSLMODE = getEnv("PGSQL_SSLMODE")

	// Redis for Docker and go
	REDIS_HOST = getEnv("REDIS_HOST")
	REDIS_PORT = getEnv("REDIS_PORT")

	// JWT
	JWT_SECRET = getEnv("JWT_SECRET")

	// Application Configuration
	MAGIC_LINK_EXPIRY = getEnv("MAGIC_LINK_EXPIRY")
	BASE_URL = getEnv("BASE_URL")
	API_RATE_LIMIT_PER_MINUTE = getEnv("API_RATE_LIMIT_PER_MINUTE")
	EMAIL_RATE_LIMIT = getEnv("EMAIL_RATE_LIMIT")

	// Razorpay Key
	RZP_KEY_ID = getEnv("RZP_KEY_ID")
	RZP_KEY_SECRET = getEnv("RZP_KEY_SECRET")
	RZP_WEBHOOK_SECRET = getEnv("RZP_WEBHOOK_SECRET")
}
