package todo

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"todo-list/controllers/auth"

)
type Todo struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email string `json:"email" bson:"email"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}



var todoCollection *mongo.Collection 

func GetTodos(c *fiber.Ctx) error {

	email, _err := auth.ExtractEmail(c.Get("Authorization"))
	if _err != nil {
		fmt.Print("No email found", _err)
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	cursor, err := todoCollection.Find(context.Background(), bson.M{ "email": email })
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong",
			"errorDetail": err.Error(),
		})
	}

	var todos []Todo = make([]Todo, 0)

	if err = cursor.All(context.Background(), &todos); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong",
			"errorDetail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": todos,
	})

}

func CreateTodo(c *fiber.Ctx) error {
	// create new todo from request body
	email, _err := auth.ExtractEmail(c.Get("Authorization"))
	fmt.Print(email)
	if _err != nil {
		fmt.Print("No email")
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	// create new todo from request body with email 
	// todo := &Todo{Email: email}
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Title == "" { 
		return c.Status(400).JSON(fiber.Map{
			"error": "Title is required",
		})
	}
	if !todo.Completed {
		
		return c.Status(400).JSON(fiber.Map{
			"error": "Completed is required",
		})

	}



	// insert into db
	insertRes, err := todoCollection.InsertOne(context.Background(), todo) 

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

func DeleteTodo(c *fiber.Ctx) error {

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
	res, err := todoCollection.DeleteOne(context.Background(), filter)

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



func GetTodoById(c *fiber.Ctx) error {
	// get id from url params
	id := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	filter := bson.M{"_id": objID}

	var todo Todo

	err = todoCollection.FindOne(context.Background(), filter).Decode(&todo)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": todo,
	})
}

func UpdateTodo(c *fiber.Ctx) error {
	// get id from url params
	id := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	filter := bson.M{"_id": objID}

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

	update := bson.M{
		"$set": bson.M{
			"title": todo.Title,
			"completed": todo.Completed,
		},
	}

	_, err = todoCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong",
			"errorDetail": err.Error(),
		})
	}

	todo.ID = objID

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Todo updated successfully",
		"data": todo,
	})
}


