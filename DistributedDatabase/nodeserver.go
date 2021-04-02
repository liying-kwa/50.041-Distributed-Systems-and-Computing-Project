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
	id   int
	ip   string
	port string

	ringServerIp   string
	ringServerPort string
	//ring           *lib.Ring
}

func newNode(id int, portNo string) Node {
	ip, _ := lib.ExternalIP()
	return Node{id, ip, portNo, lib.RING_IP, lib.RING_PORT}
}

func (n *Node) addNodeToRing() {
	nodeData := lib.NodeData{n.id, n.ip, n.port}
	requestBody, _ := json.Marshal(nodeData)
	// Send to ring server
	postURL := fmt.Sprintf("http://%s:%s/add-node", n.ringServerIp, n.ringServerPort)
	resp, _ := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//Checks response from registering with ring server
	fmt.Println("Response from registering w Ring Server: ", string(body))
}

func main() {

	//aNode := newNode(0)
	//aNode.addNodeToRing()

	for {
		fmt.Printf("NodeServer> ")
		cmdString, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		//fmt.Printf("Command given by NodeServer: %s \n", cmdString)
		tokens := strings.Fields(cmdString)
		if len(tokens) == 0 {
			fmt.Println("Please enter a command.")
			continue
		}
		cmd := tokens[0]
		switch cmd {

		case "help":
			fmt.Println("Commands accepted: start")

		case "start":
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
				aNode := newNode(0, strconv.Itoa(portNo))
				aNode.addNodeToRing()
			} else {
				fmt.Println("Invalid port number, must be between 0 and 65353")
			}

		default:
			fmt.Println("Unknown command.")

		}
	}

}
