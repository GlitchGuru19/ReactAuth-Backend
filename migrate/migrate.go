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
	initializers.DB.AutoMigrate(&models.User{})
}
