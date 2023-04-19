package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoResponse struct {
	Data []Todo `json:"data"`
}

func main() {
	app := fiber.New()
	todos := []Todo{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running...")
	})

	// get all todos
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.JSON(TodoResponse{Data: todos})
	})

	// get a todo by Id
	app.Get("/api/todos/:id", func(c *fiber.Ctx) error {
		// Parse the Id from the URL parameter
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}

		// Find the todo with the matching Id
		var todo *Todo
		for _, t := range todos {
			if t.Id == id {
				todo = &t
				break
			}
		}

		// If no matching todo is found, return a 404 error
		if todo == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("Todo with id %d not found", id),
			})
		}

		// Otherwise, return the todo
		return c.JSON(todo)
	})

	// Route to create a new todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		// Parse the request body into a Todo struct
		var newTodo Todo
		if err := c.BodyParser(&newTodo); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bad request",
			})
		}

		// Assign a unique Id to the new todo
		newTodo.Id = len(todos) + 1

		// Add the new todo to the list
		todos = append(todos, newTodo)

		// Return the new todo
		return c.JSON(newTodo)
	})

	// Route to delete a todo by Id
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		// Parse the Id from the URL parameter
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}

		// Find the index of the todo with the matching Id
		var index int = -1
		for i, t := range todos {
			if t.Id == id {
				index = i
				break
			}
		}

		// If no matching todo is found, return a 404 error with the Id in the message
		if index == -1 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("Todo with ID %d not found", id),
			})
		}

		// Remove the todo from the list
		todos = append(todos[:index], todos[index+1:]...)

		// Return a success message
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Todo with ID %d has been deleted", id),
		})
	})

	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		// Get the todo id from the request params
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid todo Id")
		}

		// Find the todo with the given Id and update its fields
		for i, todo := range todos {
			if todo.Id == id {
				var body map[string]interface{}
				if err := json.Unmarshal(c.Body(), &body); err != nil {
					return c.Status(fiber.StatusBadRequest).SendString("Invalid Request Body")
				}

				if completed, ok := body["completed"].(bool); ok {
					todo.Completed = completed
				}

				if title, ok := body["title"].(string); ok {
					todo.Title = title
				}

				todos[i] = todo
				return c.JSON(todo)
			}
		}

		// If no todo is found with the given Id, return a 404 response
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Todo with Id %d not found", id),
		})
	})

	log.Fatal(app.Listen(":1337"))
}
