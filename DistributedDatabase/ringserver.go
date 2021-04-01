package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
)

type RingServer struct {
	ip   string
	port string
	ring lib.Ring
}

func newRingServer() RingServer {
	ip, _ := lib.ExternalIP()
	return RingServer{
		ip,
		"5001",
		lib.Ring{
			lib.RING_MAX_ID,
			make(map[int]lib.NodeData),
		},
	}
}

func (ringServer RingServer) start() {
	http.HandleFunc("/add-node", ringServer.addNodeHandler)
	//http.HandleFunc("/faint-node", ringServer.FaintNodeHandler)
	//http.HandleFunc("/remove-node", ringServer.RemoveNodeHandler)
	//http.HandleFunc("/revive-node", ringServer.ReviveNodeHandler)
	//http.HandleFunc("/get-node", ringServer.GetNodeHandler)
	//http.HandleFunc("/hb", ringServer.HeartBeatHandler)
	//http.HandleFunc("/get-ring", ringServer.GetRingHandler)
	log.Print(fmt.Sprintf("[RingServer] Started and Listening at %s:%s.", ringServer.ip, ringServer.port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ringServer.port), nil))
}

func (ringServer *RingServer) addNodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RingServer] Receiving Registration from a Node")
	body, _ := ioutil.ReadAll(r.Body)
	var nodeData lib.NodeData
	json.Unmarshal(body, &nodeData)
	// Add node to ring
	ringServer.ring.RingNodeDataMap[nodeData.Id] = nodeData
	fmt.Fprintf(w, "Successlly added node to ring! ")
}

func main() {

	theRingServer := newRingServer()
	go theRingServer.start()

	time.Sleep(time.Second * 20)

	println(theRingServer.ring.RingNodeDataMap[0].Id)
	println(theRingServer.ring.RingNodeDataMap[0].Ip)
	println(theRingServer.ring.RingNodeDataMap[0].Port)

	/* reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("RingServer> ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Printf("Command given: %s \n", cmdString)
	} */

}
