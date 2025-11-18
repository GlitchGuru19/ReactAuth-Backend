package main

import (
	"ReactAuthBackend/initializers"
	"ReactAuthBackend/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	// Auto-migrate your User model - creates/upgrades table automatically
	initializers.DB.AutoMigrate(&models.User{})
}
