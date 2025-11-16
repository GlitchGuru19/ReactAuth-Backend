package controllers

import (
	"ReactAuthBackend/initializers"
	"ReactAuthBackend/models"
	"ReactAuthBackend/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// REGISTER: Create a new user
func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Hash password using bcrypt
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)

	// Create user struct
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: passwordHash,
	}

	// Save user to database
	initializers.DB.Create(&user)

	return c.JSON(fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// LOGIN: Authenticate user and return access + refresh tokens
func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	initializers.DB.Where("email = ?", data["email"]).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Incorrect password"})
	}

	// Generate access + refresh tokens
	accessToken, _ := utils.GenerateAccessToken(user.ID)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID)

	return c.JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// REFRESH: Use refresh token to get a new access token
func RefreshToken(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	refreshToken := data["refreshToken"]
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Refresh token required"})
	}

	claims, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid refresh token"})
	}

	userID, _ := strconv.Atoi(claims.Issuer)
	accessToken, _ := utils.GenerateAccessToken(uint(userID))

	return c.JSON(fiber.Map{"accessToken": accessToken})
}

// LOGOUT: Invalidate tokens (frontend can delete them)
func Logout(c *fiber.Ctx) error {
	// Since tokens are stateless JWTs, logout is handled by client deleting tokens.
	// You can optionally implement a blacklist if you want.
	return c.JSON(fiber.Map{"message": "Logout successful. Delete tokens on client side."})
}

// USER: Protected route, returns authenticated user's info
func User(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User not authenticated"})
	}

	var user models.User
	initializers.DB.Where("id = ?", userID).First(&user)
	return c.JSON(user)
}
