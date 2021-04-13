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
	Id   int
	Ip   string
	Port string
	Hash string

	ConnectedToRing bool
	RingServerIp    string
	RingServerPort  string
	// For replication during writes
	SuccessorIP   string
	SuccessorPort string
	SuccessorIP2   string
	SuccessorPort2 string
}

func newNode(id int, portNo string) *Node {
	ip, _ := lib.ExternalIP()
	return &Node{id, ip, portNo, "", false, lib.RINGSERVER_IP, lib.RINGSERVER_NODES_PORT, "", "", "", ""}
}

func (n *Node) addNodeToRing() {
	nodeData := lib.NodeData{Id: n.Id, Ip: n.Ip, Port: n.Port, Hash: "", PredecessorIP: "", PredecessorPort: "", PredecessorIP2: "", PredecessorPort2: ""}
	requestBody, _ := json.Marshal(nodeData)
	// Send to ring server
	postURL := fmt.Sprintf("http://%s:%s/add-node", n.RingServerIp, n.RingServerPort)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	predecessorIP := ""
	predecessorPort := ""
	predecessorIP2 := ""
	predecessorPort2 := ""
	if resp.StatusCode == 200 {
		var nodeData2 lib.NodeData
		json.Unmarshal(responseBody, &nodeData2)
		n.Id = nodeData2.Id
		n.Hash = nodeData2.Hash
		n.ConnectedToRing = true
		// So that it can send replicated data upon write requests
		n.SuccessorIP = nodeData2.SuccessorIP
		n.SuccessorPort = nodeData2.SuccessorPort
		// To request replicas
		predecessorIP = nodeData2.PredecessorIP
		predecessorPort = nodeData2.PredecessorPort
		predecessorIP2 = nodeData2.PredecessorIP2
		predecessorPort2 = nodeData2.PredecessorPort2
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

	// Request for replica (when more than 1 node in ring)
	fmt.Printf("current ip %s\n", n.Ip)
	fmt.Printf("predcessor IP: %s\n", predecessorIP)
	fmt.Printf("current port: %s\n", n.Port)
	fmt.Printf("predcessor port: %s\n", predecessorPort)

	if n.Ip == predecessorIP && n.Port == predecessorPort {
		return
	} else {
		// Buffer time to allow receiving of supposed data before receiving replica
		time.Sleep(time.Second * 2)
		fmt.Printf("Requesting predecessor %s:%s for replica\n", predecessorIP, predecessorPort)
		go lib.RequestTransfer(n.Ip, n.Port, predecessorIP, predecessorPort, -1, true)
		go lib.RequestTransfer(n.Ip, n.Port, predecessorIP2, predecessorPort2, -1, true)
	}
	// So that the command line can print correctly
	time.Sleep(time.Second)
}

func (n *Node) removeNodeFromRing() {
	// nodeData := lib.NodeData{n.Id, n.Ip, n.Port}
	// requestBody, _ := json.Marshal(nodeData)
	// postURL := fmt.Sprintf("http://%s:%s/remove-node", n.RingServerIp, n.RingServerPort)
	// resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	// if resp.StatusCode == 200 {
	// 	n.ConnectedToRing = false
	// 	fmt.Println("Successfully de-registered. Response:", string(body))
	// } else {
	// 	fmt.Println("Failed to de-register. Reason:", string(body))
	// }
}

func (n *Node) listenToRing(portNo string) {
	http.HandleFunc("/read", n.ReadHandler)
	http.HandleFunc("/write", n.WriteHandler)
	http.HandleFunc("/transfer", n.TransferHandler)
	http.HandleFunc("/loadReplica", n.LoadRepHandler)
	log.Print(fmt.Sprintf("[NodeServer] Started and Listening at %s:%s.", n.Ip, n.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", n.Port), nil))
}

func (n *Node) TransferHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TRANSFER REQUESTED FROM %s", n.Id)
	body, _ := ioutil.ReadAll(r.Body)
	var trfMessage lib.TransferMessage
	json.Unmarshal(body, &trfMessage)
	fmt.Printf("Transfer Message: %v\n", trfMessage)

	if trfMessage.Replica {
		fmt.Print("[NodeServer] Received Transfer Request for data to be replicated")
		n.SuccessorIP = trfMessage.Ip
		n.SuccessorPort = trfMessage.Port
	} else {
		fmt.Print("[NodeServer] Received Transfer Request for Data")

		// New node added, transfer all its data and delete its replica (so that it can re-request for the latest replica)
		folderName := fmt.Sprintf("./node%d/replica/", n.Id)
		err := os.RemoveAll(folderName)
		if err != nil {
			log.Fatal(err)
		}
	}

	foldername := fmt.Sprintf("./node%d/", n.Id)
	items, _ := ioutil.ReadDir(foldername)

	for _, item := range items {
		// Skip transferring the replica folder
		if item.Name() == "replica" {
			continue
		}
		fileNameKey := -1
		newNodeKey := -2
		if correct, err := strconv.Atoi(item.Name()); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		} else {
			fileNameKey = correct
		}

		if correct, err := strconv.Atoi(trfMessage.Hash); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		} else {
			newNodeKey = correct
		}

		// fmt.Printf("this is the filenamekey: \n")
		// fmt.Println(fileNameKey)

		// fmt.Printf("this is the newnodekey: \n")
		// fmt.Println(newNodeKey)

		ownHash := -1

		if correct, err := strconv.Atoi(n.Hash); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		} else {
			ownHash = correct
		}

		filename := fmt.Sprintf("./node%d/%s", n.Id, item.Name())
		lines := []string{}
		print(newNodeKey, fileNameKey, ownHash)
		// Sending all files to be replica or selected files to new node
		if trfMessage.Replica || (newNodeKey > ownHash && (newNodeKey >= fileNameKey && fileNameKey > ownHash)) || (newNodeKey < ownHash && (newNodeKey >= fileNameKey || fileNameKey > ownHash)) {
			print(ownHash)
			fmt.Printf("trying to read file \n")
			data, err := ioutil.ReadFile(filename)
			fmt.Printf("managed to read file \n")
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			lines = strings.Split(string(data), "\n")
		} else {
			continue
		}

		for _, line := range lines {
			fmt.Printf("trying to send lines")
			interim := strings.Split(line, " ")

			// fmt.Printf("interim list: \n")
			// fmt.Println(interim)

			// fmt.Printf("courseid: \n")
			// fmt.Println(interim[0])
			courseId := interim[0]

			// fmt.Printf("count: \n")
			// fmt.Println(interim[1])
			count := interim[1]

			message := lib.Message{}
			if trfMessage.Replica {
				message = lib.Message{Type: lib.Put, CourseId: courseId, Count: count, Hash: item.Name(), Replica: true}
			} else {
				message = lib.Message{Type: lib.Put, CourseId: courseId, Count: count, Hash: item.Name(), Replica: false}
			}

			// Not implemented as async because you want to be sure the message is sent before deleting it
			lib.WriteMessage(message, trfMessage.Ip, trfMessage.Port)
		}
		// Delete the file if not replica
		if !trfMessage.Replica && ((newNodeKey > ownHash && (newNodeKey >= fileNameKey && fileNameKey > ownHash)) || (newNodeKey < ownHash && (newNodeKey >= fileNameKey || fileNameKey > ownHash))) {
			fmt.Printf("[DELETING FILE] %s\n", filename)
			e := os.Remove(filename)
			if e != nil {
				log.Fatal(e)
			} else {
				fmt.Printf("Transferred the data successfully and deleted the file locally")
			}
		}
	}

	if !trfMessage.Replica {
		// Inform successor to refresh its replication set because you have deleted some of your data
		time.Sleep(time.Second * 5)
		// nodeData := lib.NodeData{Id: n.Id, Ip: n.Ip, Port: n.Port, Hash: "", PredecessorIP: "", PredecessorPort: "", PredecessorIP2: "", PredecessorPort2: ""}
		nodeData := lib.NodeData{Id: n.Id, Ip: n.Ip, Port: n.Port, Hash: "", PredecessorIP: "", PredecessorPort: ""}

		requestBody, _ := json.Marshal(nodeData)
		postURL := fmt.Sprintf("http://%s:%s/loadReplica", n.SuccessorIP, n.SuccessorPort)
		resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			fmt.Println("Requested for Replica Refresh!")
		}

		// postURL2 := fmt.Sprintf("http://%s:%s/loadReplica", n.SuccessorIP2, n.SuccessorPort2)
		// resp2, err2 := http.Post(postURL2, "application/json", bytes.NewReader(requestBody))
		// if err2 != nil {
		// 	fmt.Println(err2)
		// 	return
		// }
		// defer resp.Body.Close()

		// if resp2.StatusCode == 200 {
		// 	fmt.Println("Requested for Replica Refresh!")
		// }

	}

	fmt.Println("Successfully updated new node!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successfully wrote to node!"))
}

func (n *Node) ReadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("READ ")

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
		return
		w.Write([]byte(problem))
	}

	courseId := courseIdArray[0]
	keyHash := keyHashArray[0]

	filename := fmt.Sprintf("./node%d/%s", n.Id, keyHash)
	data, err := ioutil.ReadFile(filename)
	// Check if keyfile exists in node at all

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
		w.Write([]byte(err.Error()))
	}
	lines := strings.Split(string(data), "\n")
	exist := false
	// Check if courseId is in keyfile

	count := "-1"
	checkCourseId := "-1"
	for _, line := range lines {
		//if strings.Contains(line, courseId) {
		//	interim := strings.Split(line, " ")
		//	count = interim[1]
		//}
		interim := strings.Split(line, " ")
		checkCourseId = interim[0]
		if checkCourseId == courseId {
			break
			exist = true
			count = interim[1]
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
	fmt.Printf("WRITE")
	log.Print("[NodeServer] Received Write Request")
	body, _ := ioutil.ReadAll(r.Body)
	var message lib.Message
	json.Unmarshal(body, &message)
	fmt.Println(message)
	filename := fmt.Sprintf("./node%d/%s", n.Id, message.Hash)
	dataToWrite := message.CourseId + " " + message.Count

	if message.Replica {
		// Create folder (unique to node) for storing data (if folder doesnt already exist)
		folderName := "./node" + strconv.Itoa(n.Id) + "/replica"

		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			os.MkdirAll(folderName, os.ModePerm)
		}

		filename = fmt.Sprintf("./node%d/replica/%s", n.Id, message.Hash)
	} else {
		// Send to successors to replicate
		print("FORWARDING MESSAGE TO SUCCESSOR TO REPLICATE")
		print(n.SuccessorIP, n.SuccessorPort)
		message.Replica = true
		go lib.WriteMessage(message, n.SuccessorIP, n.SuccessorPort)
	}

	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("File already exists, proceeding to update file... \n")
		isAlreadyInside := false

		// to check if the courseId is alr inside, if it is, then update with this latest value
		if !isAlreadyInside {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatalln(err)
			}
			lines := strings.Split(string(data), "\n")
			for i, line := range lines {
				if strings.Contains(line, message.CourseId) {
					lines[i] = dataToWrite
					isAlreadyInside = true
				}
			}
			output := strings.Join(lines, "\n")
			err = ioutil.WriteFile(filename, []byte(output), 0644)
			if err != nil {
				log.Fatalln(err)
			}
		}

		// if course id is not inside, then append it
		if !isAlreadyInside {
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

func (n *Node) LoadRepHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOADREPHANDLER")
	log.Print("[NodeServer] Received Request to Reload Replica")

	// Delete all its replica before requesting for replica
	folderName := fmt.Sprintf("./node%d/replica/", n.Id)
	err := os.RemoveAll(folderName)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(r.Body)
	var nodeData lib.NodeData
	json.Unmarshal(body, &nodeData)

	print("REQUEST DATA FROM")
	print(n.Port)
	print("TO")
	print(nodeData.Port)
	print("TO BE REPLICA")
	go lib.RequestTransfer(n.Ip, n.Port, nodeData.Ip, nodeData.Port, -1, true)

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
