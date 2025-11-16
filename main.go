package main

import (
	"ReactAuthBackend/initializers"
	"ReactAuthBackend/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	}))

	routes.SetUpRoutes(app)

	// Get port from initializer
	port := initializers.GetPort()
	log.Printf("Server is running")
	log.Fatal(app.Listen(":" + port))
}
