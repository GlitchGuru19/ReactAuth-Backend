package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateAccessToken creates a JWT that expires in 15 minutes
func GenerateAccessToken(userID uint) (string, error) {
	claims := jwt.StandardClaims{
		Issuer:    string(rune(userID)),
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken creates a JWT that expires in 30 days
func GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.StandardClaims{
		Issuer:    string(rune(userID)),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

// VerifyAccessToken verifies access token and returns claims
func VerifyAccessToken(tokenString string) (*jwt.StandardClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*jwt.StandardClaims), nil
}

// VerifyRefreshToken verifies refresh token and returns claims
func VerifyRefreshToken(tokenString string) (*jwt.StandardClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*jwt.StandardClaims), nil
}
