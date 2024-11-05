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


// func GenerateJwtToken(email string) (string, error) {

// 	token := jwt.New(jwt.SigningMethodHS256)
// 	claims := token.Claims.(jwt.MapClaims)
// 	if claims == nil {
// 		claims = jwt.MapClaims{}
// 		token.Claims = claims
// 	}
// 	claims["email"] = email
// 	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

// 	secret := os.Getenv("JWT_SECRET")
// 	if secret == "" {
// 		return "", fmt.Errorf("JWT_SECRET environment variable not set")
// 	}
// 	t, err := token.SignedString([]byte(secret))
	
// 	if err != nil {
// 		return "", err
// 	}

// 	return t, nil
// }

func GenerateJwtToken(email string) (string, error) {
	// Create a new token object, specifying the signing method and the claims
	
	// Retrieve the secret key from environment variables
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}

	// // Sign the token with the secret key
	// t, err := token.SignedString([]byte(secret))
	// if err != nil {
	// 	return "", err
	// }

	// return t, nil

	// Create a new token object, specifying signing method and the claims
// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))

	fmt.Println(tokenString, err)
	return tokenString, err
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
	fmt.Print(err)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims")
	}

	return claims, nil
}

func ExtractEmail(tokenString string) (string, error) {
	claims, err := ExtractClaims(tokenString)
	
	fmt.Print(err)
	if err != nil {
		return "", err
	}
	// fmt.Print(claims)
	//map[email:krcpr080@gmail.com exp:1.730888215e+09]
	email , ok := claims["email"].(string)

	if !ok {
		return "", fmt.Errorf("Invalid email claim")
	}
	// fmt.Print(email)
	return email, nil
}

