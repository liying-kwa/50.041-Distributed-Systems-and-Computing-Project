package main

import (
	"bufio"
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
	"strings"
	"time"

	"50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
)

type RingServer struct {
	Ip   string
	Port string
	Ring lib.Ring
}

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
func (ringServer *RingServer) AllocateKey(key string) string {
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
			nodeURL := fmt.Sprintf("%s:%s", nodeMap[key].Ip, nodeMap[key].Port)
			return nodeURL
		}
	}

	nodeURL := fmt.Sprintf("%s:%s", nodeMap[lowest].Ip, nodeMap[lowest].Port)
	return nodeURL
}

func (ringServer RingServer) start() {
	http.HandleFunc("/add-node", ringServer.AddNodeHandler)
	// http.HandleFunc("/test", ringServer.test)
	//http.HandleFunc("/faint-node", ringServer.FaintNodeHandler)
	http.HandleFunc("/remove-node", ringServer.RemoveNodeHandler)
	//http.HandleFunc("/revive-node", ringServer.ReviveNodeHandler)
	//http.HandleFunc("/get-node", ringServer.GetNodeHandler)
	//http.HandleFunc("/hb", ringServer.HeartBeatHandler)
	//http.HandleFunc("/get-ring", ringServer.GetRingHandler)
	log.Print(fmt.Sprintf("[RingServer] Started and Listening at %s:%s.", ringServer.Ip, ringServer.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ringServer.Port), nil))
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

	/* // interim array to iterate through the keys easier
	keys := make([]int, len(ringNodeDataMap))
	i := 0
	for k := range ringNodeDataMap {
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
	} */

	// Add node to ring
	fmt.Printf("Adding node %s:%s to the ring... \n", nodeData.Ip, nodeData.Port)
	ringNodeDataMap[randomKey] = nodeData
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successlly added node to ring!"))

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
	go theRingServer.start()
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

		default:
			fmt.Println("Unknown command. Use 'help' to see available commands.")

		}

		fmt.Println()
	}

}
