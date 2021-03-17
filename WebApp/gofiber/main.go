package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/WebApp/gofiber/api"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/WebApp/gofiber/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func helloWorld(c *fiber.Ctx) error{
	return c.SendString("Hello, World!")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)
	app.Get("/api/v1/student", api.GetStudents)
	app.Get("/api/v1/student/:id", api.GetStudent)
	app.Post("/api/v1/student", api.NewStudent)
	app.Delete("/api/v1/student/:id", api.DelStudent)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("students.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connection successfully opened!")

	database.DBConn.AutoMigrate(&api.Student{})
	fmt.Println("Database Migrated")
}

// go run main.go helper.go
func main(){
	app:= fiber.New()
	initDatabase()
	setupRoutes(app)

	app.Listen(":3000")
}