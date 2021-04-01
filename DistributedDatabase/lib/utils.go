package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"sort"
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

//   function to allocate the given CourseId to a node and return that node's ip:port
func (ringServer *RingServer) AllocateKey(key string) string {
	nodeMap := ringServer.ring.RingNodeDataMap
	keyHash := HashMD5(key, MAX_KEYS)
	var lowest int
	lowest = math.MaxInt32

	for key := range nodeMap {
		if key < lowest {
			lowest = key
		}
	}

	keys := make([]int, len(nodeMap))
	i := 0
	for k := range nodeMap {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	for _, key := range keys {
		if keyHash <= key {
			nodeURL := fmt.Sprintf("%s:%s", nodeMap[key].Ip, nodeMap[key].Port)
			return nodeURL
		}
	}

	nodeURL := fmt.Sprintf("%s:%s", nodeMap[lowest].Ip, nodeMap[lowest].Port)
	return nodeURL
}

func SendMessage(message string, nodeData NodeData) {
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
