package routes

import (
	"database/sql"
	"fiber-crud/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	userHandler := handlers.NewUserHandler(db)

	api := app.Group("/api")
	users := api.Group("/users")

	users.Post("/", userHandler.CreateUser)
	// users.Post("/", func(c *fiber.Ctx) error {
	// 	fmt.Println("Hi")
	// 	return nil
	// })
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
