package utils

import (
	"errors"
	"fmt"
	"strings"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var SecretKey = []byte("12345")

// GetUserIDFromToken extracts the user_id from the JWT token in the Authorization header
func GetUserIDFromToken(c echo.Context) (int, error) {
	// Get the Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("authorization header is missing")
	}

	// Check if the token is prefixed with "Bearer "
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return 0, errors.New("invalid authorization header format")
	}

	tokenString := parts[1]

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64) // JWT usually stores numbers as float64
		if !ok {
			return 0, errors.New("user_id not found in token claims")
		}
		return int(userID), nil
	}

	return 0, errors.New("invalid token claims")
}