package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	ip, err := ExternalIP()
	newRingServer := RingServer{ip, 8000}
	fmt.Printf("RingServer is serving on %s:%s", newRingServer.ip, newRingServer.port)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Ringserver> ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Printf("Command given: %s \n", cmdString)
	}

}
