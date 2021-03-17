package main

import (
	"bufio"
	"fmt"
	"os"
)

type NodeData struct {
	ID string
	//CName string
	//Hash  int
	IP   string
	Port string
}

func nodeMain() {
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
