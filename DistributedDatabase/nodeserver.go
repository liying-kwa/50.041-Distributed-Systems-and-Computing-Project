package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	return &Node{id, ip, portNo, false, lib.RING_IP, lib.RING_PORT}
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
	//Checks response from registering with ring server
	if resp.StatusCode == 200 {
		//fmt.Println("Response from registering w Ring Server: ", string(body))
		n.ConnectedToRing = true
		fmt.Println("Successfully registered in ring!")
	} else {
		fmt.Println("Failed to register. Reason: ", string(body))
	}
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
			fmt.Println("Commands accepted: help, info, register")

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

		default:
			fmt.Println("Unknown command. Use 'help' to see available commands.")

		}
	}

}
