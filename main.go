package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

func main() {
	app := fiber.New()

	todos := []Todo{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "Hello World"})
	})

	app.Post("/todos", func(c *fiber.Ctx) error {
		todo := Todo{}
		if err := c.BodyParser(&todo); err != nil {
			return err
		}

		todos = append(todos, todo)

		return c.Status(201).JSON(fiber.Map{"msg": "Todo created", "todos": todos})
	})

	log.Fatal(app.Listen(":3000"))
}