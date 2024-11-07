package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID   primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello, world!")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection = client.Database("todo").Collection("todos")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS, PATCH",
	}))

	app.Get("/api/getTodo", getTodos) 
	app.Post("/api/createTodo", createTodo)
	app.Patch("/api/updateTodo/:id", updateTodo)
	app.Delete("/api/deleteTodo/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))	
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Name == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Todo name is required"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(fiber.Map{"msg": "Todo created", "todos": todo})
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": "Invalid ID"})
	}

	_,err = collection.UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": bson.M{"done": true}})
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"msg": "Todo updated"})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": "Invalid ID"})
	}

	_,err = collection.DeleteOne(context.Background(), bson.M{"_id": objectId})
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"msg": "Todo deleted"})
}