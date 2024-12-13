package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error while loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	//Create Todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		// Parse the request body into a Todo struct
		var todo Todo
		if err := c.BodyParser(&todo); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Validate required fields
		if strings.TrimSpace(todo.Body) == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		// Generate a unique ID and add the todo item to the list
		todo.ID = len(todos) + 1
		todos = append(todos, todo)

		// Respond with the created todo
		return c.Status(201).JSON(todo)
	})

	// Update Todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		// Extract the ID from the request parameters
		id := c.Params("id")

		// Parse the ID into an integer for better validation
		todoID, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid Todo ID"})
		}

		// Find and update the corresponding Todo
		for i := range todos {
			if todos[i].ID == todoID {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		// Return a 404 response if the Todo was not found
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// Delete Todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		// Extract the ID from the request parameters
		id := c.Params("id")

		// Parse the ID into an integer for validation
		todoID, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid Todo ID"})
		}

		// Find the index of the Todo to delete
		for i, todo := range todos {
			if todo.ID == todoID {
				// Remove the Todo by slicing
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"message": "Todo deleted successfully"})
			}
		}

		// Return a 404 response if the Todo was not found
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
