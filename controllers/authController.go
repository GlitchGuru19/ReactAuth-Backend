package controllers

import (
	"ReactAuthBackend/initializers"
	"ReactAuthBackend/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	// Get the user off req body
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// encrypt the password
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	// map what the user inputs to the field we have
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	// Add user to the database
	initializers.DB.Create(&user)

	// Return the user
	return c.JSON(user)
}

// Function to log in a user
func Login(c *fiber.Ctx) error {
	var data map[string]string

	// pass the req body
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// use the email for checking the users as the email is unique
	var user models.User

	initializers.DB.Where("email = ?", data["email"]).First(&user)

	// Check for thr user in the db using the email
	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)

		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	// compare the passwords and if its wrong the user shouldn't be allowed access
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	// creating a jwt
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	// create a token
	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not log in",
		})
	}

	// create a cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	// Return the token
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	initializers.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	// remove the cookie by setting the expiry time to an hour ago
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
