package api

import (
	"github.com/gofiber/fiber/v2"
)

// The database stores key:value of Student-ID: [mod1, mod2, ...]?

// Get all student carts
func getStudents(c *fiber.Ctx) error{
	return c.SendString("All Student Carts")
}

// Get student cart
func getStudent(c *fiber.Ctx) error{
	return c.SendString("Student Cart!")
}

// New cart
func newStudent(c *fiber.Ctx) error{
	return c.SendString("New Cart!")
}

// Clear cart
func delStudent(c *fiber.Ctx) error{
	return c.SendString("Clear Cart!")
}