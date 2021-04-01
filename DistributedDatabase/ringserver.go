package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/gofiber/api"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/gofiber/database"
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var wg = &sync.WaitGroup{}

// TODO: Confirm key-value names
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

// Initiate socket of ring on port 5001 (for communication with node server)
func newRingServer() RingServer {
	ip, _ := lib.ExternalIP()
	return RingServer{
		ip,
		"5001",
		lib.Ring{
			make(map[int]lib.NodeData),
		},
	}
}

// Listening on port 3001 (for communication with front-end)
func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Hello, World!") })
	app.Get("/api/v1/student", api.GetStudents)
	app.Get("/api/v1/student/:id", api.GetStudent)
	app.Put("/api/v1/student/:id", api.PutStudent)
	app.Post("/api/v1/student", api.NewStudent)
	app.Delete("/api/v1/student/:id", api.DelStudent)
}

// Listening on port 5001 (for communication with node servers)
func (ringServer RingServer) start() {
	http.HandleFunc("/add-node", ringServer.addNodeHandler)
	//http.HandleFunc("/faint-node", ringServer.FaintNodeHandler)
	//http.HandleFunc("/remove-node", ringServer.RemoveNodeHandler)
	//http.HandleFunc("/revive-node", ringServer.ReviveNodeHandler)
	//http.HandleFunc("/get-node", ringServer.GetNodeHandler)
	//http.HandleFunc("/hb", ringServer.HeartBeatHandler)
	//http.HandleFunc("/get-ring", ringServer.GetRingHandler)
	log.Print(fmt.Sprintf("[RingServer] Started and Listening at %s:%s.", ringServer.ip, ringServer.port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ringServer.port), nil))
}

// Initialise temporary SQL databse (see gofiber/api)
func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("students.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connection successfully opened!")

	database.DBConn.AutoMigrate(&api.Student{})
	fmt.Println("Database Migrated")

	students := []Student{
		{Studentid: 1001234, Course: "DS"},
		{Studentid: 1000000, Course: "DB"},
	}
	for _, c := range students {
		database.DBConn.Create(&c)
	}
}

//   function to allocate the given CourseId to a node and return that node's ip:port
func (ringServer *RingServer) AllocateKey(key string) string {
	nodeMap := ringServer.ring.RingNodeDataMap
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
			nodeURL := fmt.Sprintf("%s:%s", nodeMap[key].Ip, nodeMap[key].Port)
			return nodeURL
		}
	}

	nodeURL := fmt.Sprintf("%s:%s", nodeMap[lowest].Ip, nodeMap[lowest].Port)
	return nodeURL
}

// Receive POST request from :5001/add-node
func (ringServer *RingServer) addNodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RingServer] Receiving Registration from Node %s", r.RemoteAddr)
	body, _ := ioutil.ReadAll(r.Body)
	var nodeData lib.NodeData
	json.Unmarshal(body, &nodeData)

	nodeMap := &ringServer.ring.RingNodeDataMap

	// creating a random key between 0 and 100
	var random int
	random = rand.Intn(lib.MAX_KEYS)

	// interim array to iterate through the keys easier
	keys := make([]int, len(*nodeMap))
	i := 0
	for k := range *nodeMap {
		keys[i] = k
		i++
	}

	// making sure that the assigned key has not alr been assigned before
	idx := 0
	for idx < len(keys) {
		if random == keys[idx] {
			random = rand.Intn(lib.MAX_KEYS)
			idx = 0
		}
		idx++
	}

	// Add node to ring
	(*nodeMap)[random] = nodeData
	fmt.Printf("Ring Structure: %v\n", *nodeMap)

	//---------------------- uncomment block below to just test the hashing function----------------//
	// var CourseID string
	// CourseID = "50005"
	// nodeURL := ringServer.AllocateKey(CourseID)
	// fmt.Println(nodeURL)

	// var CourseIDTwo string
	// CourseIDTwo = "500115"
	// nodeURL2 := ringServer.AllocateKey(CourseIDTwo)
	// fmt.Println(nodeURL2)
	//---------------------- uncomment block above to just test the hashing function----------------//

	// HTTP response
	fmt.Fprintf(w, "Successlly added node to ring! ")
}

// Will take awhile for first run as code imports from Github
func main() {
	app := fiber.New()
	initDatabase()
	setupRoutes(app)

	ip, _ := lib.ExternalIP()

	theRingServer := newRingServer()
	go theRingServer.start()

	api.GetRingStructure(&theRingServer.ring)

	log.Print(fmt.Sprintf("[RingServer] To test, visit %s:%s/api/v1/student", ip, "3001"))

	app.Listen(ip + ":3001")
}
