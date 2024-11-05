package main

import (
	// "context"
	// "fmt"
	"log"
	"os"
	"todo-list/controllers/auth"
	"todo-list/controllers/todo"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "./controllers/auth"
)

// var todoCollection *mongo.Collection 
// var userCollection *mongo.Collection

func main() {

	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("Error loading .env file")
	}

	// connect to mongodb not the any collection but only a db 
	// MONGODB_URI := os.Getenv("MONGODB_URI")
	// clientOptions := options.Client().ApplyURI(MONGODB_URI)
	// client, err := mongo.Connect(context.Background(), clientOptions)


	if err != nil {
		log.Fatal(err)
	}

	
	// err = client.Ping(context.Background(), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

    // collection := 	client.Database("golang_db")

	app := fiber.New()


	app.Post("/api/auth/login", auth.LoginUser)
	app.Post("/api/auth/register", auth.RegisterUser)

	app.Get("/api/todos", auth.IsAuthenticated, todo.GetTodos)
	app.Get("/api/todos/:id", auth.IsAuthenticated, todo.GetTodoById)
	app.Post("/api/todos", auth.IsAuthenticated, todo.CreateTodo)
	app.Patch("/api/todos/:id", auth.IsAuthenticated, todo.UpdateTodo)
	app.Delete("/api/todos/:id", auth.IsAuthenticated, todo.DeleteTodo)


	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

