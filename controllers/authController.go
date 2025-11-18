package controllers

import (
	"ReactAuthBackend/initializers"
	"ReactAuthBackend/models"
	"ReactAuthBackend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// NOTE:
// This file contains your auth controllers converted from Fiber -> Gin.
// Each handler uses *gin.Context and returns JSON responses with proper HTTP status codes.
// I kept your original logic (bcrypt, tokens, DB calls) but adapted to Gin's APIs.

// Register: Create a new user
func Register(c *gin.Context) {
	// Use a map to receive JSON body (name, email, password)
	var data map[string]string

	// Bind incoming JSON into `data`. If invalid JSON, return 400.
	if err := c.ShouldBindJSON(&data); err != nil {
		// Bad request - invalid JSON
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	// Validate required fields (basic)
	if data["name"] == "" || data["email"] == "" || data["password"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "name, email, and password are required"})
		return
	}

	// Hash password using bcrypt with cost 12
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)
	if err != nil {
		// Internal error during hashing
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password", "error": err.Error()})
		return
	}

	// Build the user model. Password stored as []byte (same as your original model).
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: passwordHash,
	}

	// Save user to DB. GORM will populate user.ID after creation.
	if err := initializers.DB.Create(&user).Error; err != nil {
		// Likely duplicate email or DB error
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		return
	}

	// Return created user info (no password)
	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// Login: Authenticate user and return access + refresh tokens
func Login(c *gin.Context) {
	var data map[string]string

	// Bind incoming JSON to map
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	// Basic validation
	if data["email"] == "" || data["password"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email and password required"})
		return
	}

	// Lookup user by email
	var user models.User
	if err := initializers.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		// User not found
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Compare provided password with stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		// Password mismatch
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect password"})
		return
	}

	// Generate tokens via your utils package (keeps original behavior)
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate access token", "error": err.Error()})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate refresh token", "error": err.Error()})
		return
	}

	// Return tokens and user info (no password)
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// RefreshToken: Use refresh token to get a new access token
func RefreshToken(c *gin.Context) {
	var data map[string]string

	// Bind JSON containing refreshToken
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	refreshToken := data["refreshToken"]
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token required"})
		return
	}

	// Verify the refresh token and extract claims
	claims, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token", "error": err.Error()})
		return
	}

	// claims.Issuer is the user ID as string (kept from original)
	userID, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid token issuer", "error": err.Error()})
		return
	}

	// Issue a fresh access token
	accessToken, err := utils.GenerateAccessToken(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate access token", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}

// Logout: For stateless JWTs, logout is client-side deletion of tokens.
// This endpoint remains for parity with your frontend flow and can be extended
// to support server-side token blacklists if needed.
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful. Delete tokens on client side."})
}

// User: Protected route, returns authenticated user's info
func User(c *gin.Context) {
	// Retrieve userID that middleware stored in context
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	// val is stored as string (claims.Issuer). Convert to integer to query DB.
	userIDStr, ok := val.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user id in context"})
		return
	}

	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user id", "error": err.Error()})
		return
	}

	var user models.User
	// Query user by ID (GORM expects numeric type)
	if err := initializers.DB.Where("id = ?", uint(userIDInt)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Return user (Password has json:"-" so it won't be exposed)
	c.JSON(http.StatusOK, user)
}
