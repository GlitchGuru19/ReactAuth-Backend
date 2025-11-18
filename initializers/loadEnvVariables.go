package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvVariables loads .env in non-production environments.
// This makes it easy to keep secrets locally in a .env file during development.
func LoadEnvVariables() {
	// Load .env only if ENV != "production"
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Failed to load .env file")
		}
	}
}

// GetPort returns PORT env var or default "8000".
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	return port
}
