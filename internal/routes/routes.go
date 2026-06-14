package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/example/go-user-api/internal/handler"
	"github.com/example/go-user-api/internal/middleware"
)

// Register attaches all routes and middleware to the Fiber app.
func Register(app *fiber.App, userHandler *handler.UserHandler) {
	// Global middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())

	// Health-check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// User routes
	users := app.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.ListUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
