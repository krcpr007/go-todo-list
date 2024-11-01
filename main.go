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
	"go.mongodb.org/mongo-driver/bson"

)
type Todo struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}

var collection *mongo.Collection 

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

	collection = client.Database("golang_db").Collection("todos")

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
	// create new todo from request body
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Title == "" { 
		return c.Status(400).JSON(fiber.Map{
			"error": "Title is required",
		})
	}
	// insert into db
	insertRes, err := collection.InsertOne(context.Background(), todo) 

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong",
			"errorDetail": err.Error(),
		})
	}
	todo.ID = insertRes.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Todo created successfully",
		"data": todo,
	})

}

func deleteTodo(c *fiber.Ctx) error {

	// get id from url params
	id := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	filter := bson.M{"_id": objID}
	fmt.Printf("filter: %v\n", filter)
	res, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong",
			"errorDetail": err.Error(),
		})
	}

	fmt.Print(res.DeletedCount, res)
	if res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}

	return c.JSON(fiber.Map{ "success": true, "message": "Todo deleted successfully" })

}

func updateTodo(c *fiber.Ctx) error {
	return c.SendString("Update Todo")
}
