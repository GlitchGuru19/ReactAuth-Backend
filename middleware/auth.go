package middleware

import (
	"net/http"
	"strings"

	"ReactAuthBackend/utils"

	"github.com/gin-gonic/gin"
)

// Protect middleware checks for access token and validates it.
// Converted from Fiber middleware to Gin.HandlerFunc.
// Behavior:
// - Expects Authorization header: "Bearer <accessToken>"
// - Verifies token using utils.VerifyAccessToken
// - Stores userID in Gin context with key "userID" (string)
func Protect() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read Authorization header
		authHeader := c.GetHeader("Authorization") // expects "Bearer <accessToken>"
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization header required",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization must be Bearer <token>",
			})
			return
		}

		tokenStr := parts[1]
		claims, err := utils.VerifyAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired access token",
				"error":   err.Error(),
			})
			return
		}

		// Save the userID (issuer) into Gin context for handlers to use
		// We store it as string because claims.Issuer is a string.
		c.Set("userID", claims.Issuer)

		// Continue to next handler
		c.Next()
	}
}
