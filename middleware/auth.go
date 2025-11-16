package middleware

import (
	"ReactAuthBackend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Protect middleware checks for access token
func Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization") // expects "Bearer <accessToken>"
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization must be Bearer <token>",
			})
		}

		tokenStr := parts[1]
		claims, err := utils.VerifyAccessToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid or expired access token",
			})
		}

		c.Locals("userID", claims.Issuer)
		return c.Next()
	}
}
