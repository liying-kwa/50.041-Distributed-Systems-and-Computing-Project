package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
)

type RingServer struct {
	Ip           string
	NodesPort    string
	FrontendPort string
	Ring         lib.Ring
}

func newRingServer() RingServer {
	ip, _ := lib.ExternalIP()
	return RingServer{
		ip,
		lib.RINGSERVER_NODES_PORT,
		lib.RINGSERVER_FRONTEND_PORT,
		lib.Ring{
			make(map[int]lib.NodeData),
		},
	}
}

// md5 hashing
func HashMD5(text string, max int) int {
	byteArray := md5.Sum([]byte(text))
	var output int
	for _, num := range byteArray {
		output += int(num)
	}
	return output % max
}

//   function to allocate the given CourseId to a node and return that node's ip:port
//func (ringServer *RingServer) AllocateKey(key string) string {
func (ringServer *RingServer) AllocateKey(key string) lib.NodeData {
	nodeMap := ringServer.Ring.RingNodeDataMap
	keyHash := HashMD5(key, lib.MAX_KEYS)
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
			//nodeURL := fmt.Sprintf("%s:%s", nodeMap[key].Ip, nodeMap[key].Port)
			//return nodeURL
			return nodeMap[key]
		}
	}

	//nodeURL := fmt.Sprintf("%s:%s", nodeMap[lowest].Ip, nodeMap[lowest].Port)
	//return nodeURL
	return nodeMap[lowest]
}

func (ringServer *RingServer) ReadFromNode(courseId string) string {
	nodeData := ringServer.AllocateKey(courseId)
	getURL := fmt.Sprintf("http://%s:%s/read?courseid=%s", nodeData.Ip, nodeData.Port, courseId)
	resp, err := http.Get(getURL)
	if err != nil {
		fmt.Println(err)
		return "-1"
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Println("Successfully read from node. Response:", string(body))
		return string(body)
	} else {
		fmt.Println("Failed to read from node. Reason:", string(body))
		return "-1"
	}
}

func (ringServer *RingServer) WriteToNode(courseId string, count string) {
	countInt, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println("Invalid count, must be an integer.")
		return
	}
	if countInt < 0 {
		fmt.Println("Invalid count, must be 0 or more")
		return
	}
	message := lib.Message{lib.Get, courseId, count}
	requestBody, _ := json.Marshal(message)
	nodeData := ringServer.AllocateKey(courseId)
	postURL := fmt.Sprintf("http://%s:%s/write", nodeData.Ip, nodeData.Port)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//Checks response from node
	if resp.StatusCode == 200 {
		fmt.Println("Successfully wrote to node. Response:", string(body))
	} else {
		fmt.Println("Failed to write to node. Reason:", string(body))
	}

}

// Listening on port 5001 for communication with Nodes
func (ringServer RingServer) listenToNodes() {
	// http.HandleFunc("/test", ringServer.test)
	http.HandleFunc("/add-node", ringServer.AddNodeHandler)
	http.HandleFunc("/remove-node", ringServer.RemoveNodeHandler)
	//http.HandleFunc("/revive-node", ringServer.ReviveNodeHandler)
	//http.HandleFunc("/get-node", ringServer.GetNodeHandler)
	//http.HandleFunc("/get-ring", ringServer.GetRingHandler)
	log.Print(fmt.Sprintf("[RingServer] Started and Listening at %s:%s.", ringServer.Ip, ringServer.NodesPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ringServer.NodesPort), nil))
}

func (ringServer *RingServer) AddNodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RingServer] Receiving Registration Request from a Node")
	body, _ := ioutil.ReadAll(r.Body)
	var nodeData lib.NodeData
	json.Unmarshal(body, &nodeData)
	ringNodeDataMap := ringServer.Ring.RingNodeDataMap

	// Check if node is already in ring structure
	for _, nd := range ringNodeDataMap {
		if nd.Ip == nodeData.Ip && nd.Port == nodeData.Port {
			fmt.Printf("Node %s:%s tries to connect but already registered previously. \n", nodeData.Ip, nodeData.Port)
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("409 Conflict -- Already registered"))
			return
		}
	}

	// Assign a random (but unique) key to the node and add to ring
	randomKey := rand.Intn(lib.MAX_KEYS)
	_, taken := ringNodeDataMap[randomKey]
	for taken == true {
		randomKey = rand.Intn(lib.MAX_KEYS)
		_, taken = ringNodeDataMap[randomKey]
	}

	// Add node to ring
	fmt.Printf("Adding node %s:%s to the ring... \n", nodeData.Ip, nodeData.Port)
	ringNodeDataMap[randomKey] = nodeData
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successfully added node to ring!"))

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

	//fmt.Fprintf(w, "Successlly added node to ring! ")

}

func (ringServer *RingServer) RemoveNodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RingServer] Receiving De-registration Request from a Node")
	body, _ := ioutil.ReadAll(r.Body)
	var nodeData lib.NodeData
	json.Unmarshal(body, &nodeData)
	ringNodeDataMap := ringServer.Ring.RingNodeDataMap

	// Check if node is already NOT in ring structure
	notInside := true
	assignedKey := -1
	for key, nd := range ringNodeDataMap {
		if nd.Ip == nodeData.Ip && nd.Port == nodeData.Port {
			notInside = false
			assignedKey = key
			break
		}
	}
	if notInside == true {
		fmt.Printf("Node %s:%s tries to de-register but is already NOT in ring. \n", nodeData.Ip, nodeData.Port)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("409 Conflict -- Already NOT in ring"))
		return
	}

	// Remove node from ring
	fmt.Printf("Removing node %s:%s from the ring... \n", nodeData.Ip, nodeData.Port)
	delete(ringNodeDataMap, assignedKey)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successlly removed node from ring!"))
}

func main() {

	theRingServer := newRingServer()
	go theRingServer.listenToNodes()
	time.Sleep(time.Second * 3)

	for {
		fmt.Printf("RingServer> ")
		cmdString, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		cmdString = strings.TrimSpace(cmdString)
		//fmt.Printf("Command given by RingServer: %s \n", cmdString)
		tokens := strings.Fields(cmdString)
		if len(tokens) == 0 {
			fmt.Println("Please enter a command. Use 'help' to see available commands.")
			continue
		}

		cmd := tokens[0]
		switch cmd {

		case "help":
			fmt.Println("Commands accepted: help, info, ring")

		case "info":
			ringserverJson, _ := json.Marshal(theRingServer)
			fmt.Println("RingServer information:")
			fmt.Println(string(ringserverJson))

		case "ring":
			if len(theRingServer.Ring.RingNodeDataMap) == 0 {
				fmt.Println("Ring is empty at the moment.")
			} else {
				for key, nd := range theRingServer.Ring.RingNodeDataMap {
					nodeDataJson, _ := json.Marshal(nd)
					fmt.Printf("key=%d, %s \n", key, string(nodeDataJson))
				}
			}

		// testing read
		case "read":
			courseId := tokens[1]
			count := theRingServer.ReadFromNode(courseId)
			fmt.Println("Returned count:", count)

		// testing write
		case "write":
			courseId := tokens[1]
			count := tokens[2]
			theRingServer.WriteToNode(courseId, count)

		default:
			fmt.Println("Unknown command. Use 'help' to see available commands.")

		}

		fmt.Println()
	}

}
