package main

import (
	"bufio"
	"bytes"
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
	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
)

type RingServer struct {
	Ip           string
	NodesPort    string
	FrontendPort string
	Ring         lib.Ring
	// RingServerChannel lib.Ring
	// quit 	      chan int
}

// Initiate socket of ring on port 5001 (for communication with node server)
func newRingServer() (RingServer, RingServer) {
	ip, _ := lib.ExternalIP()
	

	return RingServer{
		ip,
		lib.RINGSERVER_NODES_PORT,
		lib.RINGSERVER_FRONTEND_PORT,
		lib.Ring{
			-1,
			make(map[int]lib.NodeData),
			true, 
		},
	}, RingServer{
		ip,
		lib.RINGSERVER_SECOND_NODES_PORT,
		lib.RINGSERVER_FRONTEND_PORT,
		lib.Ring{
			-1,
			make(map[int]lib.NodeData),
			true,
		},
	}
}

// Function to allocate the given CourseId to a node and return that nodeData and keyHash
func (ringServer *RingServer) AllocateKey(key string) (lib.NodeData, string) {
	nodeMap := ringServer.Ring.RingNodeDataMap
	keyHash := lib.HashMD5(key, lib.MAX_KEYS)
	fmt.Printf("this is the hash below: \n")
	fmt.Println(keyHash)

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
			return nodeMap[key], strconv.Itoa(keyHash)
		}
	}

	return nodeMap[lowest], strconv.Itoa(keyHash)
}

// Listening on port 3001 for communication with Frontend
func (ringServer RingServer) listenToFrontend() {
	http.HandleFunc("/read-from-node", ringServer.ReadFromNodeHandler)
	http.HandleFunc("/write-to-node", ringServer.WriteToNodeHandler)
	log.Print(fmt.Sprintf("[RingServer] Started and Listening at %s:%s for Frontend.", ringServer.Ip, ringServer.FrontendPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ringServer.FrontendPort), nil))
}

func (ringServer RingServer) ReadFromNodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RingServer] Received Read Request from Frontend")
	courseIdArray, ok := r.URL.Query()["courseid"]
	if !ok || len(courseIdArray) < 1 {
		problem := "Query parameter 'courseid' is missing"
		fmt.Println(problem)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(problem))
		return
	}
	courseId := courseIdArray[0]

	// Create HTTP GET request and send to Node
	nodeData, keyHash := ringServer.AllocateKey(courseId)
	getURL := fmt.Sprintf("http://%s:%s/read?courseid=%s&keyhash=%s", nodeData.Ip, nodeData.Port, courseId, keyHash)
	resp, err := http.Get(getURL)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Echo response back to Frontend
	if resp.StatusCode == 200 {
		fmt.Println("Successfully read from node. Response:", string(body))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(body)))
	} else {
		fmt.Println("Failed to read from node. Reason:", string(body))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(body)))
	}

}

func (ringServer *RingServer) ReadFromNode(courseId string) string {
	fmt.Println(ringServer.Ring.RingNodeDataMap)
	nodeData, keyHash := ringServer.AllocateKey(courseId)
	getURL := fmt.Sprintf("http://%s:%s/read?courseid=%s&keyhash=%s", nodeData.Ip, nodeData.Port, courseId, keyHash)
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

func (ringServer RingServer) WriteToNodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RingServer] Received Write Request from Frontend")
	body, _ := ioutil.ReadAll(r.Body)
	var message lib.Message
	json.Unmarshal(body, &message)
	courseId := message.CourseId
	count := message.Count

	// Create HTTP POST request and send to Node
	countInt, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println("Invalid count, must be an integer.")
		return
	}
	if countInt < 0 {
		fmt.Println("Invalid count, must be 0 or more")
		return
	}
	nodeData, keyHash := ringServer.AllocateKey(courseId)
	message2 := lib.Message{Type: lib.Put, CourseId: courseId, Count: count, Hash: keyHash, Replica: false}
	requestBody, _ := json.Marshal(message2)
	postURL := fmt.Sprintf("http://%s:%s/write", nodeData.Ip, nodeData.Port)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()
	body2, _ := ioutil.ReadAll(resp.Body)

	// Echo response back to Frontend
	if resp.StatusCode == 200 {
		fmt.Println("Successfully wrote to node. Response:", string(body2))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(body2)))
	} else {
		fmt.Println("Failed to write to node. Reason:", string(body2))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(body2)))
	}
}

// [KIV] Similar to WriteToNodeHandler (probably can use this function inside so less duplicates)
func (ringServer *RingServer) WriteToNode(courseId string, count string) {
	fmt.Println(ringServer.Ring.RingNodeDataMap)
	countInt, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println("Invalid count, must be an integer.")
		return
	}
	if countInt < 0 {
		fmt.Println("Invalid count, must be 0 or more")
		return
	}
	nodeData, keyHash := ringServer.AllocateKey(courseId)
	message := lib.Message{Type: lib.Put, CourseId: courseId, Count: count, Hash: keyHash, Replica: false}
	requestBody, _ := json.Marshal(message)
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
func (ringServer *RingServer) listenToNodes() {
	// http.HandleFunc("/test", ringServer.test)
	http.HandleFunc("/add-node", ringServer.AddNodeHandler)
	http.HandleFunc("/remove-node", ringServer.RemoveNodeHandler)
	http.HandleFunc("/send-res", ringServer.SendResponseHandler)
	//http.HandleFunc("/revive-node", ringServer.ReviveNodeHandler)
	//http.HandleFunc("/get-node", ringServer.GetNodeHandler)
	//http.HandleFunc("/get-ring", ringServer.GetRingHandler)
	log.Print(fmt.Sprintf("[RingServer] Started and Listening at %s:%s for Nodes.", ringServer.Ip, ringServer.NodesPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ringServer.NodesPort), nil))
}

func (ringServer *RingServer) AddNodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RingServer] Receiving Registration Request from a Node")
	body, _ := ioutil.ReadAll(r.Body)
	var newNodeData lib.NodeData
	json.Unmarshal(body, &newNodeData)
	ringNodeDataMap := ringServer.Ring.RingNodeDataMap

	// Check if node is already in ring structure
	for _, nd := range ringNodeDataMap {
		if nd.Ip == newNodeData.Ip && nd.Port == newNodeData.Port {
			fmt.Printf("Node %s:%s tries to connect but already registered previously. \n", newNodeData.Ip, newNodeData.Port)
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("409 Conflict -- Already registered"))
			return
		}
	}

	// Assign a random (but unique) key and unique ID to node
	// Add to ring first, to determine successors and predecessors. Update information abt successors and predecessors later.
	newNodeKey := rand.Intn(lib.MAX_KEYS)
	_, taken := ringNodeDataMap[newNodeKey]
	for taken == true {
		newNodeKey = rand.Intn(lib.MAX_KEYS)
		_, taken = ringNodeDataMap[newNodeKey]
	}
	newNodeData.Hash = strconv.Itoa(newNodeKey)
	newNodeData.Id = ringServer.Ring.MaxID + 1
	ringServer.Ring.MaxID += 1
	fmt.Printf("Adding Node #%d with %s:%s to the ring... \n", newNodeData.Id, newNodeData.Ip, newNodeData.Port)
	ringNodeDataMap[newNodeKey] = newNodeData

	// Find the node's RF successors and RF predecessors and update newNodeData
	// SKIP ALL THIS NONSENSE IF NEWNODE IS THE ONLY ONE IN THE RING. Send newNodeData back to new node
	if len(ringNodeDataMap) == 1 {
		responseBody, _ := json.Marshal(newNodeData)
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
		return
	}
	successors := lib.FindSuccessors(newNodeKey, ringNodeDataMap)
	newNodeData.Successors = successors
	predecessors := lib.FindPredecessors(newNodeKey, ringNodeDataMap)
	newNodeData.Predecessors = predecessors

	// Edit newNodeData in ringstructure to have the right predecessors and successors. Send newNodeData back to new node.
	fmt.Printf("Editing Node #%d with %s:%s in the ring... \n", newNodeData.Id, newNodeData.Ip, newNodeData.Port)
	ringNodeDataMap[newNodeKey] = newNodeData
	responseBody, _ := json.Marshal(newNodeData)
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)

	// Update newNode's successors about its changed precedessors. Update in ring as well as send them individually.
	for _, successorSimpleNodeData := range newNodeData.Successors {
		successorKey, _ := strconv.Atoi(successorSimpleNodeData.Hash)
		successorNodeData := ringNodeDataMap[successorKey]
		predecessors := lib.FindPredecessors(successorKey, ringNodeDataMap)
		successorNodeData.Predecessors = predecessors
		// Update in ring
		ringNodeDataMap[successorKey] = successorNodeData
		// Send them individually
		lib.UpdatePredecessors(predecessors, successorNodeData)
	}
	// Update newNode's predecessors about its changed successors. Update in ring as well as send them individually.
	for _, predecessorSimpleNodeData := range newNodeData.Predecessors {
		predecessorKey, _ := strconv.Atoi(predecessorSimpleNodeData.Hash)
		predecessorNodeData := ringNodeDataMap[predecessorKey]
		successors := lib.FindSuccessors(predecessorKey, ringNodeDataMap)
		predecessorNodeData.Successors = successors
		ringNodeDataMap[predecessorKey] = predecessorNodeData
		lib.UpdateSuccessors(successors, predecessorNodeData)
	}

	// Finally, migrate and replicate
	go ringServer.migrateAndReplicate(newNodeKey, ringNodeDataMap)

}

func (ringServer *RingServer) migrateAndReplicate(newNodeKey int, ringNodeDataMap map[int]lib.NodeData) {
	// Tell newNode's immediate successor to migrate part of data to the newNode
	immediateSuccessorNodeData := lib.FindImmediateSuccessor(newNodeKey, ringNodeDataMap)
	newNodeData := ringNodeDataMap[newNodeKey]
	lib.RequestDataMigration(immediateSuccessorNodeData, newNodeData)
	// When data migration is done, broadcast reloadReplica for all affected nodes, i.e. newNode, newNode's successor, newNode's successor's successors
	newSimpleNodeData := lib.SimpleNodeData{
		Id:   newNodeData.Id,
		Ip:   newNodeData.Ip,
		Port: newNodeData.Port,
		Hash: newNodeData.Hash,
	}
	immediateSuccessorSimpleNodeData := lib.SimpleNodeData{
		Id:   immediateSuccessorNodeData.Id,
		Ip:   immediateSuccessorNodeData.Ip,
		Port: immediateSuccessorNodeData.Port,
		Hash: immediateSuccessorNodeData.Hash,
	}
	go lib.InformReloadReplica(newSimpleNodeData)
	go lib.InformReloadReplica(immediateSuccessorSimpleNodeData)
	for _, successor := range immediateSuccessorNodeData.Successors {
		go lib.InformReloadReplica(successor)
	}
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

// func (ringServer RingServer) pingPrimaryRingServer() {
// 	http.HandleFunc("/checkAlive", ringServer.checkAlive) 
// 	log.Print(fmt.Sprintf("[SecondRingServer] Started pinging Primary Server %s:%s", ringServer.Ip, lib.RINGSERVER_NODES_PORT))
// }

func (ringServer *RingServer) SendResponseHandler(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
	// print(lib.PrettyPrintStruct(ringServer))
	if (ringServer.Ring.IsAlive) {
		fmt.Println(ringServer.Ring.IsAlive)
=======
	// ringServer.Ring.IsAlive = false
	print(lib.PrettyPrintStruct(ringServer)) 
	// ringServer.Ring.IsAlive = false
	// time.Sleep(time.Second * 5)
	if (ringServer.Ring.IsAlive) {
>>>>>>> fa77a07805706a04f9c223e1d69229f01c3c07b3
		log.Printf("[RingServer] Receiving Ping Request from Secondary Node")
		ringNodeDataMap := ringServer.Ring.RingNodeDataMap
		numMilliSeconds := rand.Intn(1000) + 3000
		time.Sleep(time.Duration(numMilliSeconds) * time.Millisecond)
		responseBody, _ := json.Marshal(ringNodeDataMap)
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
<<<<<<< HEAD
	} 
	return
	// w.WriteHeader(http.StatusBadRequest)
=======
	} else {
		time.Sleep(time.Second * 10)
	}
>>>>>>> fa77a07805706a04f9c223e1d69229f01c3c07b3
}

func (ringServer RingServer) checkAlive() {
	for {
		log.Print(fmt.Sprintf("[SecondRingServer] Started pinging Primary Server %s:%s\n", ringServer.Ip, lib.RINGSERVER_NODES_PORT))
		numMilliSeconds := rand.Intn(1000) + 3000
		time.Sleep(time.Duration(numMilliSeconds) * time.Millisecond)
<<<<<<< HEAD


		fmt.Println("Primary node is down. Secondary node taking over.")
		requestBody, _ := json.Marshal(lib.RINGSERVER_SECOND_NODES_PORT)
		postURL := fmt.Sprintf("http://%s:%s/update-nodes", ringServer.Ip, )
		resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		responseBody, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			fmt.Println("Successfully updated ring server. Response:", string(responseBody))
		} else {
			fmt.Println("Failed to update ring server. Response:", string(responseBody))
		}




		// pingMessage := lib.Message{Type: lib.Put, CourseId: "", Count: "", Hash: "", Replica: false}
		// // requestBody, _ := json.Marshal(pingMessage)
		// getURL := fmt.Sprintf("http://%s:%s/send-res", ringServer.Ip, lib.RINGSERVER_NODES_PORT)
		

		// // counter := 1
		// // gotResponse := false
		// // for counter < 8 {
		// c := &http.Client {
		// 	Timeout: 5 * time.Second,
		// }
		// resp, err := c.Get(getURL)
		// // if err != nil {
		// // 	fmt.Println(err)
		// 	// w.WriteHeader(http.StatusBadRequest)
		// 	// w.Write([]byte(err.Error()))
		// 	// return
		// // }
		// if err != nil {
		// 	fmt.Println(err)
		// 	fmt.Println("Primary node is down. Secondary node taking over.")
		// }
		// defer resp.Body.Close()
		// body2, _ := ioutil.ReadAll(resp.Body)
		// if resp.StatusCode == 200 {
		// 	fmt.Println("Successfully got back Ring Structure. Response:", string(body2))
		// 	// w.Header().Set("Access-Control-Allow-Origin", "*")
		// 	// w.WriteHeader(http.StatusOK)
		// 	// w.Write([]byte(string(body2)))
		// 	// gotResponse = true
		// 	// break
		// } 
		// 	// time.Sleep(time.Duration(1000) * time.Millisecond)
		// 	// counter += 1
		// 	// fmt.Println(counter)
		// // }
=======
		getURL := fmt.Sprintf("http://%s:%s/send-res", ringServer.Ip, lib.RINGSERVER_NODES_PORT)
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(getURL)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Primary node is down. Secondary node taking over.")
			return
		}
		defer resp.Body.Close()
		body2, _ := ioutil.ReadAll(resp.Body)
>>>>>>> fa77a07805706a04f9c223e1d69229f01c3c07b3

		if resp.StatusCode == 200 {
			fmt.Println("Successfully got back Ring Structure. Response:", string(body2))
		} 
	}
}

func informAllNodes() {

}

func main() {

	// Set a different seed everytime so consistent hashing doesnt hash same keys
	rand.Seed(time.Now().UTC().UnixNano())

	// Initialise ringserver
	theRingServer, theSecondRingServer := newRingServer()
	go theRingServer.listenToFrontend()
	// fmt.Println(lib.PrettyPrintStruct(theSecondRingServer))
	//time.Sleep(time.Second * 3)
	go theRingServer.listenToNodes()
	go theSecondRingServer.checkAlive()
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
			fmt.Println("RingServer information:")
			fmt.Println(lib.PrettyPrintStruct(theRingServer))

		case "ring":
			if len(theRingServer.Ring.RingNodeDataMap) == 0 {
				fmt.Println("Ring is empty at the moment.")
			} else {
				// Print ring pointer? TODO: Fix MaxID display, it always shows -1 for some reason
				fmt.Println(lib.PrettyPrintStruct(&theRingServer.Ring))
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
		
		case "kill": 
			// fmt.Println("Killing primary node")
			theRingServer.Ring.IsAlive = false
			// fmt.Println(theSecondRingServer.Ring.IsAlive)

		default:
			fmt.Println("Unknown command. Use 'help' to see available commands.")
		}

		fmt.Println()
	}

}
