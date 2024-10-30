package main

import (
	"fmt"
	"log"
	"os"
	"context"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"

)
type Todo struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}
func main() {

	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("Error loading .env file")
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB ATLAS")

	collection := client.Database("golang_db").Collection("todos")
	fmt.Println(collection)

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}


func getTodos(c *fiber.Ctx) error {
	return c.SendString("All Todos")
}

func createTodo(c *fiber.Ctx) error {
	return c.SendString("Create Todo")
}

func deleteTodo(c *fiber.Ctx) error {
	return c.SendString("Delete Todo")
}

func updateTodo(c *fiber.Ctx) error {
	return c.SendString("Update Todo")
}
