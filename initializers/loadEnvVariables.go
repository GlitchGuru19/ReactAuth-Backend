package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {

	// Load .env file if not in production
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Failed to load .env file")
		}
	}
}
