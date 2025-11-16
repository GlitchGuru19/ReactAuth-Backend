package middleware

import (
	"ReactAuthBackend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RequireAuth is a middleware that protects routes using access token
func RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization") // get token from header

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "no_token",
			"message": "Authorization header is missing",
		})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid_header",
			"message": "Authorization header format must be Bearer {token}",
		})
	}

	token := tokenParts[1]
	claims, err := utils.VerifyAccessToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid_token",
			"message": "Access token is invalid or expired",
		})
	}

	// Store user ID in locals for controllers
	c.Locals("userID", claims.Issuer)

	return c.Next()
}
