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
	"time"

	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
)

type Node struct {
	Id           int
	Ip           string
	Port         string
	Hash         string
	Predecessors map[int]lib.SimpleNodeData
	Successors   map[int]lib.SimpleNodeData

	ConnectedToRing bool
	RingServerIp    string
	RingServerPort  string
}

func newNode(id int, portNo string) *Node {
	ip, _ := lib.ExternalIP()
	return &Node{
		Id:              id,
		Ip:              ip,
		Port:            portNo,
		Hash:            "",
		Predecessors:    make(map[int]lib.SimpleNodeData),
		Successors:      make(map[int]lib.SimpleNodeData),
		ConnectedToRing: false,
		RingServerIp:    lib.RINGSERVER_IP,
		RingServerPort:  lib.RINGSERVER_NODES_PORT,
	}
}

func (n *Node) addNodeToRing() {

	// Request Ringserver to add this node to the ring
	nodeData := lib.NodeData{
		Id:           n.Id,
		Ip:           n.Ip,
		Port:         n.Port,
		Hash:         n.Hash,
		Predecessors: n.Predecessors,
		Successors:   n.Successors,
	}
	requestBody, _ := json.Marshal(nodeData)
	postURL := fmt.Sprintf("http://%s:%s/add-node", n.RingServerIp, n.RingServerPort)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)

	// Update this node's information locally and create folder for storing data
	if resp.StatusCode == 200 {
		var nodeData2 lib.NodeData
		json.Unmarshal(responseBody, &nodeData2)
		n.Id = nodeData2.Id
		n.Hash = nodeData2.Hash
		n.Predecessors = nodeData2.Predecessors
		n.Successors = nodeData2.Successors
		n.ConnectedToRing = true
		fmt.Println(nodeData2)
		go n.listenToRing(n.Port)
		// Create folder (unique to node) for storing data (if folder doesnt already exist)
		folderName := "node" + strconv.Itoa(n.Id)
		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			os.Mkdir(folderName, 0755)
		}
		fmt.Println("Successfully registered. Response:", string(responseBody))
	} else {
		fmt.Println("Failed to register. Response:", string(responseBody))
	}

	// So that the command line can print correctly
	time.Sleep(time.Second)
}

func (n *Node) removeNodeFromRing() {
	/* nodeData := lib.NodeData{n.Id, n.Ip, n.Port}
	requestBody, _ := json.Marshal(nodeData)
	postURL := fmt.Sprintf("http://%s:%s/remove-node", n.RingServerIp, n.RingServerPort)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		n.ConnectedToRing = false
		fmt.Println("Successfully de-registered. Response:", string(body))
	} else {
		fmt.Println("Failed to de-register. Reason:", string(body))
	} */
}

func (n *Node) listenToRing(portNo string) {
	http.HandleFunc("/read", n.ReadHandler)
	http.HandleFunc("/write", n.WriteHandler)
	http.HandleFunc("/update-predecessors", n.UpdatePredecessorsHandler)
	http.HandleFunc("/update-successors", n.UpdateSuccessorsHandler)
	http.HandleFunc("/migrate-data", n.MigrateDataHandler)
	http.HandleFunc("/reload-replica", n.ReloadReplicaHandler)
	http.HandleFunc("/replicate", n.ReplicateHandler)
	http.HandleFunc("/update-nodes", n.updateRingServer)
	log.Print(fmt.Sprintf("[NodeServer] Started and Listening at %s:%s.", n.Ip, n.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", n.Port), nil))
}

func (n *Node) ReadHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[NodeServer] Received Read Request from RingServer")

	courseIdArray, ok := r.URL.Query()["courseid"]
	keyHashArray, ok := r.URL.Query()["keyhash"]
	if !ok || len(courseIdArray) < 1 {
		problem := "Query parameter 'courseid' is missing"
		fmt.Println(problem)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(problem))
		return
	}
	problem := "Query parameter 'keyhash' is missing"
	fmt.Println(problem)
	if !ok || len(keyHashArray) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(problem))
		return
	}

	courseId := courseIdArray[0]
	keyHash := keyHashArray[0]
	filename := fmt.Sprintf("./node%d/%s", n.Id, keyHash)
	data, err := ioutil.ReadFile(filename)

	// Check if keyfile exists in node at all
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	lines := strings.Split(string(data), "\n")
	exist := false

	// Check if courseId is in keyfile
	count := "-1"
	checkCourseId := "-1"
	for _, line := range lines {
		interim := strings.Split(line, " ")
		checkCourseId = interim[0]
		if checkCourseId == courseId {
			exist = true
			count = interim[1]
			break
		}
	}
	// If courseId exists, return count. Else, return 'error' msg
	if exist == true {
		fmt.Println("Returning count:", count)
		w.Write([]byte(count))
		w.WriteHeader(http.StatusOK)
	} else {
		noCourseIdMsg := fmt.Sprintf("CourseID #%s does not exist.", courseId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(noCourseIdMsg))
		fmt.Println(noCourseIdMsg)
	}
}

func (n *Node) WriteHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[NodeServer] Received Write Request")
	body, _ := ioutil.ReadAll(r.Body)
	var message lib.Message
	json.Unmarshal(body, &message)
	filename := fmt.Sprintf("./node%d/%s", n.Id, message.Hash)
	dataToWrite := message.CourseId + " " + message.Count

	// If data is a replica, then store in ./node/replica folder and stop here.
	// Else, data is original data. Fwd replicas and continue storing as original data here
	if message.Replica {
		// Create folder (unique to node) for storing data (if folder doesnt already exist)
		folderName := "./node" + strconv.Itoa(n.Id) + "/replica"
		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			os.MkdirAll(folderName, os.ModePerm)
		}
		filename = fmt.Sprintf("./node%d/replica/%s", n.Id, message.Hash)
	} else {
		// Send to successors to replicate
		// NOTE: Comment out this branch (remove interference) to test replica migration when node is added
		/* print("FORWARDING MESSAGE TO SUCCESSOR TO REPLICATE")
		for _, successor := range n.Successors {
			print(successor.Ip, successor.Port)
			message.Replica = true
			go lib.WriteMessage(message, successor.Ip, successor.Port)
		} */
	}

	// Write data
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("File already exists, proceeding to update file... \n")
		isAlreadyInside := false

		// Check if the courseId is alr inside
		data, _ := ioutil.ReadFile(filename)
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			interim := strings.Split(line, " ")
			courseId := interim[0]
			if courseId == message.CourseId {
				lines[i] = dataToWrite
				isAlreadyInside = true
				break
			}
		}
		// If courseId is alr inside, then update with latest value. Else, append it.
		if isAlreadyInside {
			output := strings.Join(lines, "\n")
			err = ioutil.WriteFile(filename, []byte(output), 0644)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if _, err = f.WriteString("\n" + dataToWrite); err != nil {
				panic(err)
			}
		}

	} else {
		fmt.Printf("File does not exist \n")
		fmt.Printf("Creating file.... \n")
		err := ioutil.WriteFile(filename, []byte(dataToWrite), 0644)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

	}

	fmt.Println("Successfully wrote to node!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successfully wrote to node!"))
}

func (n *Node) UpdatePredecessorsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RingServer] Received Update Predecessors Request from a Ringserver")
	body, _ := ioutil.ReadAll(r.Body)
	var newPredecessors map[int]lib.SimpleNodeData
	json.Unmarshal(body, &newPredecessors)
	n.Predecessors = newPredecessors
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successfully updated node's predecessors!"))
}

func (n *Node) UpdateSuccessorsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RingServer] Received Update Successors Request from a Ringserver")
	body, _ := ioutil.ReadAll(r.Body)
	var newSuccessors map[int]lib.SimpleNodeData
	json.Unmarshal(body, &newSuccessors)
	n.Successors = newSuccessors
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successfully updated node's successors!"))
}

// This is the newNode's successor. Pass part of this node's data to newNode
func (n *Node) MigrateDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[NodeServer] Received Transfer Request for data to be migrated to new node")
	body, _ := ioutil.ReadAll(r.Body)
	var newNodeData lib.NodeData
	json.Unmarshal(body, &newNodeData)
	newNodeKey, _ := strconv.Atoi(newNodeData.Hash)
	thisNodeKey, _ := strconv.Atoi(n.Hash)
	// Check each keyfile to see if need to send it over to newNode
	foldername := fmt.Sprintf("./node%d/", n.Id)
	files, _ := ioutil.ReadDir(foldername)
	for _, file := range files {
		// Skip transferring the replica folder
		if file.Name() == "replica" {
			continue
		}
		// Get compare fileKey with thisNodeKey and nextNodeKey to see if need to transfer
		filename := fmt.Sprintf("./node%d/%s", n.Id, file.Name())
		fileKey, _ := strconv.Atoi(file.Name())

		if (thisNodeKey > newNodeKey && (fileKey <= newNodeKey || fileKey > thisNodeKey)) ||
			(newNodeKey > thisNodeKey && (fileKey <= newNodeKey && fileKey > thisNodeKey)) {
			data, _ := ioutil.ReadFile(filename)
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				interim := strings.Split(line, " ")
				courseId := interim[0]
				count := interim[1]
				message := lib.Message{
					Type:     lib.Put,
					CourseId: courseId,
					Count:    count,
					Hash:     file.Name(),
					Replica:  false,
				}
				lib.WriteMessage(message, newNodeData.Ip, newNodeData.Port)
			}
			// Done transferring, delete file.
			e := os.Remove(filename)
			if e != nil {
				log.Fatal(e)
			} else {
				fmt.Println("Transferred the data successfully and deleted the file locally")
			}
		}

	}
	fmt.Println("Successfully transferred data to newNode!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successfully transferred data to newNode!"))
}

func (n *Node) ReloadReplicaHandler(w http.ResponseWriter, r *http.Request) {
	// Delete its replica folder (so that it can re-request for the latest replica)
	folderName := fmt.Sprintf("./node%d/replica/", n.Id)
	err := os.RemoveAll(folderName)
	if err != nil {
		log.Fatal(err)
	}
	// Request for its predecessors to send data over here to store as replica here
	for _, predecessor := range n.Predecessors {
		thisNodeData := lib.NodeData{
			Id:           n.Id,
			Ip:           n.Ip,
			Port:         n.Port,
			Hash:         n.Hash,
			Predecessors: n.Predecessors,
			Successors:   n.Successors,
		}
		requestBody, _ := json.Marshal(thisNodeData)
		postURL := fmt.Sprintf("http://%s:%s/replicate", predecessor.Ip, predecessor.Port)
		resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		responseBody, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			fmt.Println("Successfully requested for transfer of data (to replicate) from predecessor. Response:", string(responseBody))
			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)
		} else {
			fmt.Println("Failed to request for transfer of data (to replicate) from predecessor. Response:", string(responseBody))
			w.WriteHeader(resp.StatusCode)
			w.Write(responseBody)
		}
	}
}

// This is one of the predecessors of an affected node. Send data back to sourceNode (where the request came from) to replicate
func (n *Node) ReplicateHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[NodeServer] Received Transfer Request for data to be replicated to sourceNode")
	body, _ := ioutil.ReadAll(r.Body)
	var sourceNodeData lib.NodeData
	json.Unmarshal(body, &sourceNodeData)
	// Transfer all data files in this node's folder
	foldername := fmt.Sprintf("./node%d/", n.Id)
	files, _ := ioutil.ReadDir(foldername)
	for _, file := range files {
		// Skip transferring the replica folder
		if file.Name() == "replica" {
			continue
		}
		// Transfer file contents
		filename := fmt.Sprintf("./node%d/%s", n.Id, file.Name())
		data, _ := ioutil.ReadFile(filename)
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			interim := strings.Split(line, " ")
			courseId := interim[0]
			count := interim[1]
			message := lib.Message{
				Type:     lib.Put,
				CourseId: courseId,
				Count:    count,
				Hash:     file.Name(),
				Replica:  true,
			}
			lib.WriteMessage(message, sourceNodeData.Ip, sourceNodeData.Port)
		}
	}
	fmt.Println("Successfully replicated data to affectedNode!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successfully replicated data to affectedNode!"))
}

func (n *Node) updateRingServer(portNo string) {
	fmt.Println("oi hello")
	n.RingServerPort = portNo
}

func main() {

	thisNode := newNode(0, "-1")

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
			fmt.Println("Node information:")
			fmt.Println(lib.PrettyPrintStruct(thisNode))

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
