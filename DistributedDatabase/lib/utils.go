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
	"strconv"
)

type JsonRequest struct {
	JsonRequestString string `json:"jsonRequestString"`
}

type JsonResponse struct {
	JsonResponseString string `json:"jsonResponseString"`
}

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

func SendMessage(message string, nodeData NodeData) string {
	fmt.Printf("Sending POST request to NodeServer %d at %s:%s\n", nodeData.Id, nodeData.Ip, nodeData.Port)
	msg, _ := json.Marshal(map[string]string{
		"message": message,
	})
	requestBody, _ := json.Marshal(msg)
	// Send to ring server
	postURL := fmt.Sprintf("http://%s:%s/listen", nodeData.Ip, nodeData.Port)
	resp, _ := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	defer resp.Body.Close()

	// Waits for HTTP response
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response from registering w NodeServer %d: %s\n", nodeData.Id, string(body))
	return string(body)
}

func RequestTransfer(requestorIp string, requestorPort string, destinationIp string, destinationPort string, hash int, replica bool) {
	trfMessage := TransferMessage{requestorIp, requestorPort, strconv.Itoa(hash), replica}
	requestBody, _ := json.Marshal(trfMessage)
	postURL := fmt.Sprintf("http://%s:%s/transfer", destinationIp, destinationPort)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body2, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		fmt.Println("Told next node about new node. Response:", string(body2))
	} else {
		fmt.Println("Failed to tell next node about new node. Reason:", string(body2))
	}
}
