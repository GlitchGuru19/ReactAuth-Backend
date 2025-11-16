package routes

import (
	"ReactAuthBackend/controllers"
	"ReactAuthBackend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/refresh", controllers.RefreshToken)
	app.Post("/api/logout", controllers.Logout)

	// Protected route
	app.Get("/api/user", middleware.Protect(), controllers.User)
}
