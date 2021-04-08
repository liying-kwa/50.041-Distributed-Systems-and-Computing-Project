package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/liying-kwa/50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
)

type Node struct {
	id   int
	ip   string
	port string

<<<<<<< Updated upstream
	ringServerIp   string
	ringServerPort string
	//ring           *lib.Ring
=======
	ConnectedToRing  bool
	RingServerIp     string
	RingServerPort   string
	NodesToReplicate []lib.NodeData
	//Ring           *lib.Ring
>>>>>>> Stashed changes
}

func newNode(id int) Node {
	ip, _ := lib.ExternalIP()
<<<<<<< Updated upstream
	return Node{id, ip, "6001", lib.RING_IP, lib.RING_PORT}
=======
	return &Node{id, ip, portNo, false, lib.RINGSERVER_IP, lib.RINGSERVER_NODES_PORT, []lib.NodeData{}}
>>>>>>> Stashed changes
}

func (n *Node) addNodeToRing() {
	nodeData := lib.NodeData{n.id, n.ip, n.port}
	requestBody, _ := json.Marshal(nodeData)
	// Send to ring server
	postURL := fmt.Sprintf("http://%s:%s/add-node", n.ringServerIp, n.ringServerPort)
	resp, _ := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	fmt.Printf("Sending POST request to ring server %s:%s\n", n.ringServerIp, n.ringServerPort)
	defer resp.Body.Close()

	// Waits for HTTP response
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response from registering w Ring Server: ", string(body))
}

<<<<<<< Updated upstream
func listen(w http.ResponseWriter, r *http.Request) {
	// TODO: Read the message and take necessary action
	fmt.Printf("[NodeServer] Receiving Message from Ring Server\n")
	// HTTP response
	fmt.Fprintf(w, "Value: 100")
=======
func (n *Node) listenToRing(portNo string) {
	http.HandleFunc("/read", n.ReadHandler)
	http.HandleFunc("/write", n.WriteHandler)
	http.HandleFunc("/loadReplica", n.LoadRepHandler)
	// http.HandleFunc("/reqData", n.ReqDataHandler)
	// http.HandleFunc("/recData", n.RecDataHandler)
	log.Print(fmt.Sprintf("[NodeServer] Started and Listening at %s:%s.", n.Ip, n.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", n.Port), nil))
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
>>>>>>> Stashed changes
}

func (n *Node) LoadRepHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[NodeServer] Received Request from RingServer to Reload Replica")

	body, _ := ioutil.ReadAll(r.Body)
	var nodeData lib.NodeData
	json.Unmarshal(body, &nodeData)

	if nodeData.Port != "" {
		fmt.Println(nodeData.Ip, nodeData.Port)

		responseBody, _ := json.Marshal(nodeData)
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
	}
}

func main() {

	aNode := newNode(0)
	aNode.addNodeToRing()

	http.HandleFunc("/listen", listen)
	http.ListenAndServe(":6001", nil)

	/* for {
		fmt.Printf("Ringserver> ")
		query, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Printf("Query given by node: %s \n", query)

	}
	*/
}
