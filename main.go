package main

import (
	"log"

	"ReactAuthBackend/initializers"
	"ReactAuthBackend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	// Load environment variables (from .env in development)
	initializers.LoadEnvVariables()

	// Connect to database (sets initializers.DB)
	initializers.ConnectToDB()
}

func main() {
	// Create a Gin engine named `router` (not `r`)
	// gin.Default() gives you logger + recovery middleware out of the box.
	router := gin.Default()

	// Configure CORS to allow your frontend during development.
	// Adjust origins to match your frontend origin(s).
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Register routes (controllers & middleware)
	routes.SetUpRoutes(router)

	// Start server using port from initializers
	port := initializers.GetPort()
	log.Printf("Server is running on port %s\n", port)
	log.Fatal(router.Run(":" + port))
}
