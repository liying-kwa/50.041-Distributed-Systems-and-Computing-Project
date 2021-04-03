package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
)

type Node struct {
	Id   int
	Ip   string
	Port string

	ConnectedToRing bool
	RingServerIp    string
	RingServerPort  string
	//Ring           *lib.Ring
}

func newNode(id int, portNo string) *Node {
	ip, _ := lib.ExternalIP()
	return &Node{id, ip, portNo, false, lib.RINGSERVER_IP, lib.RINGSERVER_NODES_PORT}
}

func (n *Node) addNodeToRing() {
	nodeData := lib.NodeData{n.Id, n.Ip, n.Port}
	requestBody, _ := json.Marshal(nodeData)
	// Send to ring server
	postURL := fmt.Sprintf("http://%s:%s/add-node", n.RingServerIp, n.RingServerPort)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//Checks response from ring server
	if resp.StatusCode == 200 {
		n.ConnectedToRing = true
		go n.listenToRing(n.Port)
		fmt.Println("Successfully registered. Response:", string(body))
	} else {
		fmt.Println("Failed to register. Response:", string(body))
	}
}

func (n *Node) removeNodeFromRing() {
	nodeData := lib.NodeData{n.Id, n.Ip, n.Port}
	requestBody, _ := json.Marshal(nodeData)
	postURL := fmt.Sprintf("http://%s:%s/remove-node", n.RingServerIp, n.RingServerPort)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//Checks response from ring server
	if resp.StatusCode == 200 {
		n.ConnectedToRing = false
		fmt.Println("Successfully de-registered. Response:", string(body))
	} else {
		fmt.Println("Failed to de-register. Reason:", string(body))
	}
}

func (n *Node) listenToRing(portNo string) {
	http.HandleFunc("/read", n.ReadHandler)
	http.HandleFunc("/write", n.WriteHandler)
	log.Print(fmt.Sprintf("[NodeServer] Started and Listening at %s:%s.", n.Ip, n.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", n.Port), nil))
}

func (n *Node) ReadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[NodeServer] Received Read Request from RingServer")
	body, _ := ioutil.ReadAll(r.Body)
	var message lib.Message
	json.Unmarshal(body, &message)
	filename := "./" + message.CourseId
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println("Returning count:", string(data))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(data)))
}

func (n *Node) WriteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[NodeServer] Received Write Request from RingServer")
	body, _ := ioutil.ReadAll(r.Body)
	var message lib.Message
	json.Unmarshal(body, &message)
	filename := "./" + message.CourseId
	data := []byte(message.Count)
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println("Successfully wrote to node!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successfully wrote to node!"))
}

func main() {

	thisNode := newNode(0, "-1")
	//thisNode.addNodeToRing()

	for {
		fmt.Printf("NodeServer> ")
		cmdString, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		//fmt.Printf("Command given by NodeServer: %s \n", cmdString)
		tokens := strings.Fields(cmdString)
		if len(tokens) == 0 {
			fmt.Println("Please enter a command. Use 'help' to see available commands.")
			continue
		}

		cmd := tokens[0]
		switch cmd {

		case "help":
			fmt.Println("Commands accepted: help, info, register, deregister")

		case "info":
			nodeJson, _ := json.Marshal(thisNode)
			fmt.Println("Node information:")
			fmt.Println(string(nodeJson))

		case "register":
			if len(tokens) < 2 {
				fmt.Println("Please specify a port number.")
				continue
			}
			portNo, err := strconv.Atoi(tokens[1])
			if err != nil {
				fmt.Println("Invalid port number, must be an integer.")
				break
			}
			if portNo >= 0 && portNo <= 65353 {
				thisNode.Port = strconv.Itoa(portNo)
				thisNode.addNodeToRing()
			} else {
				fmt.Println("Invalid port number, must be between 0 and 65353")
			}

		case "deregister":
			thisNode.removeNodeFromRing()

		default:
			fmt.Println("Unknown command. Use 'help' to see available commands.")

		}

		fmt.Println()
	}

}
