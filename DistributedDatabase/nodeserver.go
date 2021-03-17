package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

func newNode(id int) Node {
	ip, _ := lib.ExternalIP()
	return Node{id, ip, "6001", lib.RING_IP, lib.RING_PORT}
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

	aNode := newNode(0)
	aNode.addNodeToRing()

	/* for {
		fmt.Printf("Ringserver> ")
		query, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Printf("Query given by node: %s \n", query)

	}
	*/
}
