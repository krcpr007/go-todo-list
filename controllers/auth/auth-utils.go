package auth 

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
	"fmt"
)

func IsAuthenticated(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized - No auth header",
		})
	}	

	// Check if the header starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized - Invalid auth header format",
		})
	}

	// Extract token by removing "Bearer " prefix
	tokenString := authHeader[7:]
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized - Empty token",
		})
	}

	// Verify the token and extract email
	_, err := ExtractEmail(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized - Invalid token",
			"details": err.Error(),
		})
	}

	return c.Next()
}

func GenerateJwtToken(email string) (string, error) {
	// Retrieve the secret key from environment variables
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}

	// Create claims with proper expiration time
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"iat":   time.Now().Unix(), // Added issued at time
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func VerifyJwtToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify that the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, fmt.Errorf("JWT_SECRET environment variable not set")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func ExtractClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := VerifyJwtToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims format")
	}

	return claims, nil
}

func ExtractEmail(tokenString string) (string, error) {
	claims, err := ExtractClaims(tokenString)
	if err != nil {
		return "", err
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("invalid or missing email claim")
	}

	return email, nil
}

