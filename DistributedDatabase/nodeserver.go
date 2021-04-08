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
	//Ring           *lib.Ring
}

func newNode(id int, portNo string) *Node {
	ip, _ := lib.ExternalIP()
	return &Node{id, ip, portNo, "", false, lib.RINGSERVER_IP, lib.RINGSERVER_NODES_PORT}
}

func (n *Node) addNodeToRing() {
	nodeData := lib.NodeData{n.Id, n.Ip, n.Port, ""}
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
	if resp.StatusCode == 200 {
		var nodeData2 lib.NodeData
		json.Unmarshal(responseBody, &nodeData2)
		n.Id = nodeData2.Id
		n.Hash = nodeData2.Hash
		n.ConnectedToRing = true
		go n.listenToRing(n.Port)

		// Create folder (unique to node) for storing data (if folder doesnt already exist)
		// folderName := "node" + strconv.Itoa(n.Id)
		folderName := "node" + strconv.Itoa(n.Id)

		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			os.Mkdir(folderName, 0755)
		}
		fmt.Println("Successfully registered. Response:", string(responseBody))
	} else {
		fmt.Println("Failed to register. Response:", string(responseBody))
	}
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
	log.Print(fmt.Sprintf("[NodeServer] Started and Listening at %s:%s.", n.Ip, n.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", n.Port), nil))
}

func (n *Node) TransferHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[NodeServer] Received Transfer Request from RingServer")
	body, _ := ioutil.ReadAll(r.Body)
	var trfMessage lib.TransferMessage
	json.Unmarshal(body, &trfMessage)
	fmt.Println(trfMessage)
	foldername := fmt.Sprintf("./node%d/", n.Id)
	items, _ := ioutil.ReadDir(foldername)

	fmt.Printf("files: \n")
	fmt.Println(items)

	fmt.Printf("going to iterate through the items")

	for _, item := range items {
		fmt.Printf("yay i am looping....")
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

		fmt.Printf("this is the filenamekey: \n")
		fmt.Println(fileNameKey)

		fmt.Printf("this is the newnodekey: \n")
		fmt.Println(newNodeKey)

		ownHash := -1

		if correct, err := strconv.Atoi(n.Hash); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		} else {
			ownHash = correct
		}

		if newNodeKey >= fileNameKey || fileNameKey > ownHash {
			fmt.Printf("trying to read file \n")
			filename := fmt.Sprintf("./node%d/%s", n.Id, item.Name())
			data, err := ioutil.ReadFile(filename)
			fmt.Printf("managed to read file \n")
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			lines := strings.Split(string(data), "\n")

			fmt.Printf("this is the lines list: \n")
			fmt.Println(lines)

			for _, line := range lines {
				fmt.Printf("trying to send lines")
				interim := strings.Split(line, " ")

				fmt.Printf("interim list: \n")
				fmt.Println(interim)

				fmt.Printf("courseid: \n")
				fmt.Println(interim[0])
				courseId := interim[0]

				fmt.Printf("count: \n")
				fmt.Println(interim[1])
				count := interim[1]

				message := lib.Message{lib.Put, courseId, count, item.Name()}
				fmt.Printf("message to be sent over: \n")
				fmt.Println(message)
				fmt.Println(n.Port)
				time.Sleep(time.Second * 5)
				requestBody, _ := json.Marshal(message)
				postURL := fmt.Sprintf("http://%s:%s/write", trfMessage.Ip, trfMessage.Port)
				resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
				if err != nil {
					fmt.Printf("there is an error")
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

			e := os.Remove(filename)
			if e != nil {
				log.Fatal(e)
			} else {
				fmt.Printf("Transferred the data successfully and deleted the file locally")
			}

		}

	}

	fmt.Println("Successfully updated new node!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successfully wrote to node!"))
}

func (n *Node) ReadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[NodeServer] Received Read Request from RingServer")

	courseIdArray, ok := r.URL.Query()["courseid"]
	keyHashArray, ok := r.URL.Query()["keyhash"]
	if !ok || len(courseIdArray) < 1 {
		problem := "Query parameter 'courseid' is missing"
		fmt.Println(problem)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(problem))
		return
	}

	if !ok || len(keyHashArray) < 1 {
		problem := "Query parameter 'keyhash' is missing"
		fmt.Println(problem)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(problem))
		return
	}

	keyHash := keyHashArray[0]
	courseId := courseIdArray[0]

	filename := fmt.Sprintf("./node%d/%s", n.Id, keyHash)
	data, err := ioutil.ReadFile(filename)

	// Check if keyfile exists in node at all
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	lines := strings.Split(string(data), "\n")

	// Check if courseId is in keyfile
	exist := false
	checkCourseId := "-1"
	count := "-1"
	for _, line := range lines {
		//if strings.Contains(line, courseId) {
		//	interim := strings.Split(line, " ")
		//	count = interim[1]
		//}
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(count))
	} else {
		noCourseIdMsg := fmt.Sprintf("CourseID #%s does not exist.", courseId)
		fmt.Println(noCourseIdMsg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(noCourseIdMsg))
	}
}

func (n *Node) WriteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[NodeServer] Received Write Request from RingServer")
	body, _ := ioutil.ReadAll(r.Body)
	var message lib.Message
	json.Unmarshal(body, &message)
	fmt.Println(message)
	filename := fmt.Sprintf("./node%d/%s", n.Id, message.Hash)
	dataToWrite := message.CourseId + " " + message.Count

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
	// filename := fmt.Sprintf("./node%d/%s", n.Id, message.CourseId)
	// data := []byte(message.Count)
	fmt.Println("Successfully wrote to node!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK -- Successfully wrote to node!"))
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
