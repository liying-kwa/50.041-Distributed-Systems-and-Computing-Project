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

type Course struct {
	gorm.Model
	Id           int
	SubjectNum   int
	CourseNum    int
	CourseName   string
	SeatsLeft    int
	Class        int
	Section      string
	DaysAndTimes string
	Room         string
	Instructor   string
	MeetingDate  string
	Status       string
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("students.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connection successfully opened!")

	database.DBConn.AutoMigrate(&api.Student{}, &Course{})
	// database.DBConn.AutoMigrate(&api.Student{}, api)
	seed(database.DBConn)
	// fmt.Println("Database Migrated")
}

func seed(db *gorm.DB) {
	channels := []Course{
		{SubjectNum: 50, CourseNum: 41, CourseName: "Distribute Systems & Computing", SeatsLeft: 100, Class: 1057, Section: "CH01-CLB Regular", DaysAndTimes: "Tu 15:00 - 17:00", Room: "Think Tank 13 (1.508)", Instructor: "Staff", MeetingDate: "20/05/2019 - 16/08/2019", Status: "Available"},
		{SubjectNum: 50, CourseNum: 1, CourseName: "Information Systems & Programming", SeatsLeft: 0, Class: 1194, Section: "CH01=2-CLB Regular", DaysAndTimes: "Th 15:00 - 17:00", Room: "Think Tank 15 (1.610)", Instructor: "Staff", MeetingDate: "20/05/2019 - 16/08/2019", Status: "Not Available"},
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
