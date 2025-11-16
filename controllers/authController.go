package controllers

import (
	"ReactAuthBackend/initializers"
	"ReactAuthBackend/models"
	"ReactAuthBackend/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user
func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Input validation
	if len(strings.TrimSpace(data["name"])) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_name",
			"message": "Name must be at least 3 characters",
		})
	}

	if !strings.Contains(data["email"], "@") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_email",
			"message": "Email is invalid",
		})
	}

	if len(data["password"]) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "weak_password",
			"message": "Password must be at least 8 characters",
		})
	}

	// Check duplicate email
	var existing models.User
	initializers.DB.Where("email = ?", data["email"]).First(&existing)
	if existing.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "email_exists",
			"message": "Email already registered",
		})
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "hash_error",
			"message": "Error hashing password",
		})
	}

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: passwordHash,
	}

	// Save user
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "db_error",
			"message": "Could not create user",
		})
	}

	return c.JSON(fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// Login authenticates a user and sends tokens
func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Find user
	var user models.User
	initializers.DB.Where("email = ?", data["email"]).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "user_not_found",
			"message": "No user with that email",
		})
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "wrong_password",
			"message": "Incorrect password",
		})
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "token_error",
			"message": "Could not generate access token",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "token_error",
			"message": "Could not generate refresh token",
		})
	}

	// Send refresh token in HTTP-only cookie
	cookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Path:     "/api/refresh",
		MaxAge:   60 * 60 * 24 * 30, // 30 days
		// Secure:   true, // uncomment in production
		// SameSite: "Strict", // uncomment in production
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"accessToken": accessToken,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// RefreshToken issues a new access token using the refresh cookie
func RefreshToken(c *fiber.Ctx) error {
	cookie := c.Cookies("refresh_token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "no_token",
			"message": "Refresh token missing",
		})
	}

	claims, err := utils.VerifyRefreshToken(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid_token",
			"message": "Refresh token invalid or expired",
		})
	}

	userID, _ := strconv.Atoi(claims.Issuer)
	accessToken, err := utils.GenerateAccessToken(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "token_error",
			"message": "Could not generate access token",
		})
	}

	return c.JSON(fiber.Map{
		"accessToken": accessToken,
	})
}

// Logout deletes refresh token cookie
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Path:     "/api/refresh",
		MaxAge:   -1,
		// Secure:   true, // uncomment in production
		// SameSite: "Strict", // uncomment in production
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "logged out",
	})
}

// User returns currently authenticated user info
func User(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "unauthenticated",
			"message": "User not authenticated",
		})
	}

	var user models.User
	initializers.DB.Where("id = ?", userID).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "user_not_found",
			"message": "User does not exist",
		})
	}

	return c.JSON(fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
