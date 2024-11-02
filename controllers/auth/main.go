package auth

import (
	"context"
	"fmt"
	"todo-list/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	// "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson"
	// "../../auth-utils"
)

var userCollection *mongo.Collection
type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Email string  `json:"email" bson:"email" unique:"true"` // unique email
	Password string `json:"password" bson:"password"`
}

func LoginUser(c *fiber.Ctx) error {

	// create new user from request body
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email and Password are required",
		})
	}

	// find user in db
	filter := bson.M{"email": user.Email}
	var existingUser User
	userCollection = db.ConnectToCollection("users")
	err := userCollection.FindOne(context.Background(), filter).Decode(&existingUser)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// auth token
	token, err := GenerateJwtToken(user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
			"errorDetail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"authToken": token,
	})
}


func RegisterUser(c *fiber.Ctx) error {
	// create new user from request body
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	if user.Email == "" || user.Password == ""  || user.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			// yeah I know this is not the best way to handle this
			"error": "Name, Email and Password are required",
		})
	}

	// insert into db but first check if user already exists 
	filter := bson.M{"email": user.Email}
	var existingUser User
	userCollection = db.ConnectToCollection("users")
	err := userCollection.FindOne(context.Background(), filter).Decode(&existingUser)

	if err == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "User already exists",
		})
	}

	// hash password with bcrypt 
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to hash password",
			"errorDetail": err.Error(),
		})
	}
	user.Password = string(hashedPassword)
	

	// insert into db with hashed password
	insertRes, err := userCollection.InsertOne(context.Background(), user)
	fmt.Print(insertRes) 
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong",
			"errorDetail": err.Error(),
		})
	}

	// auth token  
	token, err := GenerateJwtToken(user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
			"errorDetail": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "User created successfully",
		"authToken": token,
	})
	
}
