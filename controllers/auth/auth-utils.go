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
		fmt.Print("No auth header")
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}	

	tokenString := authHeader[len("Bearer "):]
	_, err := ExtractEmail(tokenString)

	if err != nil {
		fmt.Print("Invalid token")
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	return c.Next()
}


func GenerateJwtToken(email string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	if claims == nil {
		claims = jwt.MapClaims{}
		token.Claims = claims
	}
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}
	t, err := token.SignedString([]byte(secret))
	
	if err != nil {
		return "", err
	}

	return t, nil
}

func VerifyJwtToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExtractClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := VerifyJwtToken(tokenString)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims")
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
		return "", fmt.Errorf("Invalid email claim")
	}

	return email, nil
}

