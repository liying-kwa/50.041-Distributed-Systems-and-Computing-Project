package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/WebApp/gofiber/database"
	"gorm.io/gorm"
)

// You would need to make your function exportable with an uppercase for its name
// The database stores key:value of Student-ID: [mod1, mod2, ...]?

type Student struct {
	gorm.Model
	StudentID int    `json:"key"`
	Course    string `json:"value"`
}

// Get all student carts
func GetStudents(c *fiber.Ctx) error {
	db := database.DBConn
	var students []Student
	db.Find(&students)
	// Convert to JSON and set content header to application/json
	return c.JSON(students)
}

// Get student cart
func GetStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var student Student
	db.Find(&student, id)
	return c.JSON(student)
}

// New cart
func NewStudent(c *fiber.Ctx) error {
	db := database.DBConn
	student := new(Student)

	if err := c.BodyParser(student); err != nil {
		return c.SendStatus(503)
	}

	db.Create(&student)
	return c.JSON(student)
}

// Clear cart
func DelStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var student Student
	db.Find(&student, id)
	print(student.StudentID)
	if student.StudentID == 0 {
		return c.Status(500).SendString("No student found with given ID")
	}
	db.Delete(&student)
	return c.SendString("Book Successfully Deleted!")
}
