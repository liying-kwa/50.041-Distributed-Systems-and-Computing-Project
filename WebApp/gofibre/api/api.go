package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// You would need to make your function exportable with an uppercase for its name
// The database stores key:value of Student-ID: [mod1, mod2, ...]?

type Student struct{
	gorm.Model
	StudentID int `json:"studentid"`
	Name string `json:"name"`
	Course string `json:"course"`

}

// Get all student carts
func GetStudents(c *fiber.Ctx) error{
	return c.SendString("All Student Carts")
}

// Get student cart
func GetStudent(c *fiber.Ctx) error{
	return c.SendString("Student Cart!")
}

// New cart
func NewStudent(c *fiber.Ctx) error{
	return c.SendString("New Cart!")
}

// Clear cart
func DelStudent(c *fiber.Ctx) error{
	return c.SendString("Clear Cart!")
}