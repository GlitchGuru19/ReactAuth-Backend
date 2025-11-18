package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global gorm DB connection used across the app.
var DB *gorm.DB

// ConnectToDB connects to Postgres using the URI in env variable "URI".
// It sets the global DB variable and logs fatal on failure.
func ConnectToDB() {
	dsn := os.Getenv("URI")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = db
}
