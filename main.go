package main

import (
	"fiber-crud/config"
	"fiber-crud/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	config.LoadEnv()
	db := config.InitDB()
	defer db.Close()

	app := fiber.New()
	routes.SetupRoutes(app, db)
	log.Fatal(app.Listen(":3000"))
}
