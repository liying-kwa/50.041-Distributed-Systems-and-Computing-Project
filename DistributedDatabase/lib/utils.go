package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

// md5 hashing
func HashMD5(text string, max int) int {
	byteArray := md5.Sum([]byte(text))
	var output int
	for _, num := range byteArray {
		output += int(num)
	}
	return output % max
}

func PrettyPrintStruct(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func FindIndexOfArray(toFind int, array []int) int {
	for idx, element := range array {
		if element == toFind {
			return idx
		}
	}
	return -1
}

func WriteMessage(message Message, destIP string, destPort string) {
	fmt.Printf("Writing message to NodeServer at %s:%s\n", destIP, destPort)

	requestBody, _ := json.Marshal(message)
	postURL := fmt.Sprintf("http://%s:%s/write", destIP, destPort)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// Checks response from node
	if resp.StatusCode == 200 {
		fmt.Println("Successfully wrote to node. Response:", string(body))
	} else {
		fmt.Println("Failed to write to node. Reason:", string(body))
	}
}
