package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
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

func sendMessage(message string, nodeData NodeData) {
	fmt.Printf("Sending POST request to node server %d at %s:%s\n", nodeData.Id, nodeData.Ip, nodeData.Port)
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
	fmt.Println("Response from registering w Ring Server: ", string(body))
}
