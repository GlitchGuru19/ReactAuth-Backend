package main

import (
	"ReactAuthBackend/initializers"
	"ReactAuthBackend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {

	app := fiber.New()

	// Correct Fiber CORS (not Gin)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	}))

	routes.SetUpRoutes(app)

	app.Listen(":8000")
}
