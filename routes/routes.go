package routes

import (
	"ReactAuthBackend/controllers"
	"ReactAuthBackend/middleware"

	"github.com/gin-gonic/gin"
)

// SetUpRoutes configures all routes on the provided Gin engine (router).
// IMPORTANT: variable name is `router` (not `r`) as you requested.
func SetUpRoutes(router *gin.Engine) {
	// Public routes
	router.POST("/api/register", controllers.Register)
	router.POST("/api/login", controllers.Login)
	router.POST("/api/refresh", controllers.RefreshToken)
	router.POST("/api/logout", controllers.Logout)

	// Protected route - uses middleware.Protect() which returns a gin.HandlerFunc
	router.GET("/api/user", middleware.Protect(), controllers.User)
}
