package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"50.041-Distributed-Systems-and-Computing-Project/DistributedDatabase/lib"
)

type RingServer struct {
	ip   string
	port string
}

type Ring struct {
	MaxID             int // 0 to maxID inclusive
	RingNodeDataArray []NodeData
}

func newRingServer() RingServer {
	ip, err := lib.ExternalIP()
	if err == nil {
		return RingServer{ip, "5000"}
	} else {
		fmt.Println(err)
		log.Fatalln(err)
		return RingServer{}
	}
}

func main() {

	theRingServer := newRingServer()
	fmt.Printf("RingServer is serving on %s:%s... \n", theRingServer.ip, theRingServer.port)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("RingServer> ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Printf("Command given: %s \n", cmdString)
	}

}
