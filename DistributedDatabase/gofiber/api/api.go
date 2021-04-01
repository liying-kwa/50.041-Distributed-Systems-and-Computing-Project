package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/gofiber/database"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
	"gorm.io/gorm"
)

// You would need to make your function exportable with an uppercase for its name
// The database stores key:value of Student-ID: [mod1, mod2, ...]?

// Need to be uppercase for first letter and lowercase for the rest
type Student struct {
	gorm.Model
	Studentid int    `json:"key"`
	Course    string `json:"value"`
}

type RingServer struct {
	ip   string
	port string
	ring lib.Ring
}

// GET all student carts
func GetStudents(c *fiber.Ctx) error {
	// TODO: Replace this with actual Node information. Ring Server should be running when executing this.
	nodeData := lib.NodeData{1, "127.0.0.1", "5001"}
	requestBody, _ := json.Marshal(nodeData)
	// Send to ring server
	postURL := fmt.Sprintf("http://%s:%s/add-node", "192.168.56.1", "5001")
	resp, _ := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	fmt.Printf("Sending POST request to ring server %s:%s\n", "192.168.56.1", "5001")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//Checks response from registering with ring server
	fmt.Println("Response from registering w Ring Server: ", string(body))

	// Temporary return SQL values
	db := database.DBConn
	var students []Student
	db.Find(&students)
	// Convert to JSON and set content header to application/json
	return c.JSON(students)
}

// GET student cart
func GetStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var student Student
	db.Where("Studentid = ?", id).First(&student)
	return c.JSON(student)
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
