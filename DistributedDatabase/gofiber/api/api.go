package api

import (
	"fmt"
	"math"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/gofiber/database"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
	"gorm.io/gorm"
)

// You would need to make your function exportable with an uppercase for its name
// The database stores key:value of Student-ID: [mod1, mod2, ...]?

var localRing lib.Ring

// Need to be uppercase for first letter and lowercase for the rest
type Student struct {
	gorm.Model
	Studentid int    `json:"key"`
	Course    string `json:"value"`
}

//   function to allocate the given CourseId to a node and return that node's ip:port
func AllocateKey(key string, ring lib.Ring) lib.NodeData {
	nodeMap := ring.RingNodeDataMap
	keyHash := lib.HashMD5(key, lib.MAX_KEYS)
	var lowest int
	lowest = math.MaxInt32

	for key := range nodeMap {
		if key < lowest {
			lowest = key
		}
	}

	keys := make([]int, len(nodeMap))
	i := 0
	for k := range nodeMap {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	for _, key := range keys {
		if keyHash <= key {
			return nodeMap[key]
		}
	}

	return nodeMap[lowest]
}

func GetRingStructure(ip string, port string, ring lib.Ring) {
	localRing = ring
}

// GET all student carts
func GetStudents(c *fiber.Ctx) error {
	db := database.DBConn
	var students []Student
	db.Find(&students)
	// Convert to JSON and set content header to application/json
	return c.JSON(students)
}

// GET student cart
func GetStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	node := AllocateKey(id, localRing)
	fmt.Printf("Received GET request from Frontend, forwarding request to Node %d at %s:%s\n", node.Id, node.Ip, node.Port)
	resp := lib.SendMessage("Testing", node)

	// TODO: Remove temporary SQL DB below
	db := database.DBConn
	var student Student
	db.Where("Studentid = ?", id).First(&student)
	return c.SendString(resp)
}

// POST
func NewStudent(c *fiber.Ctx) error {
	db := database.DBConn
	student := new(Student)

	if err := c.BodyParser(student); err != nil {
		return c.SendStatus(503)
	}

	db.Create(&student)
	return c.JSON(student)
}

// PUT
func PutStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var student Student

	studentNew := new(Student)
	c.BodyParser(studentNew)
	db.Where("Studentid = ?", id).First(&student).Update("Course", studentNew.Course)

	return c.JSON(student)
}

// Clear cart
func DelStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var student Student
	db.Where("Studentid <> ?", id).Find(&student)
	if student.Studentid == 0 {
		return c.Status(500).SendString("No student found with given ID")
	}
	db.Delete(&student)
	return c.SendString("Book Successfully Deleted!")
}
