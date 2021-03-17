package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/WebApp/gofibre/api"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/WebApp/gofibre/database"
)

func helloWorld(c *fiber.Ctx) error{
	return c.SendString("Hello, World!")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)
	app.Get("/api/v1/student", api.getStudents)
	app.Get("/api/v1/student/:id", api.getStudent)
	app.Post("/api/v1/student", api.newStudent)
	app.Delete("/api/v1/student/:id", api.delStudent)
}


// go run main.go helper.go
func main(){
	app:= fiber.New()

	setupRoutes(app)

	app.Listen(":3000")
}