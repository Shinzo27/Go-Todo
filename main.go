package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

func main() {
	app := fiber.New()

	todos := []Todo{}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"todos": todos})
	})

	app.Post("/todos", func(c *fiber.Ctx) error {
		todo := Todo{}
		if err := c.BodyParser(&todo); err != nil {
			return err
		}

		todos = append(todos, todo)

		return c.Status(201).JSON(fiber.Map{"msg": "Todo created", "todos": todos})
	})

	app.Patch(("/todos/:id"), func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, t := range todos {
			if fmt.Sprint(t.ID) == id {
				todos[i].Done = true
				return c.Status(200).JSON(fiber.Map{"msg": "Todo updated", "todos": todos})
			}
		}
		return c.Status(404).JSON(fiber.Map{"msg": "Todo not updated"})
	})

	app.Delete("/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, t := range todos {
			if fmt.Sprint(t.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				break
			}
		}
		return c.Status(200).JSON(fiber.Map{"msg": "Todo deleted", "todos": todos})
	})

	log.Fatal(app.Listen(":" + PORT))
}
