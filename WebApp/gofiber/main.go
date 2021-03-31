package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/WebApp/gofiber/api"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/WebApp/gofiber/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)
	app.Get("/api/v1/student", api.GetStudents)
	app.Get("/api/v1/student/:id", api.GetStudent)
	app.Put("/api/v1/student/:id", api.PutStudent)
	app.Post("/api/v1/student", api.NewStudent)
	app.Delete("/api/v1/student/:id", api.DelStudent)
	app.Get("/api/v1/courses", api.GetCourses)
}

type Course struct  {
	gorm.Model 
	subjectNum uint 
	courseNum uint 
	courseName string 
	seatsLeft uint 
	class uint 
	section string 
	daysAndTimes string 
	room string 
	instructor string 
	meetingDate string 
	status string 
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("students.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connection successfully opened!")

	database.DBConn.AutoMigrate(&api.Student{}, Course{})
	// database.DBConn.AutoMigrate(&api.Student{}, api)
	seed(database.DBConn)
	// fmt.Println("Database Migrated")
}

func seed(db *gorm.DB) {
	channels := []Course{
		{subjectNum: 50, courseNum: 41, courseName: "Distribute Systems & Computing", seatsLeft: 100, class: 1057, section: "CH01-CLB Regular", daysAndTimes: "Tu 15:00 - 17:00", room: "Think Tank 13 (1.508)", instructor: "Staff", meetingDate: "20/05/2019 - 16/08/2019", status: "Available"},
		{subjectNum: 50, courseNum: 1, courseName: "Information Systems & Programming", seatsLeft: 0, class: 1194, section: "CH01=2-CLB Regular", daysAndTimes: "Th 15:00 - 17:00", room: "Think Tank 15 (1.610)", instructor: "Staff", meetingDate: "20/05/2019 - 16/08/2019", status: "Not Available"},
	}
	for _, c := range channels {
		db.Create(&c)
	}
}

// go run main.go helper.go
func main() {
	app := fiber.New()
	initDatabase()
	setupRoutes(app)

	app.Listen(":3001")
}
