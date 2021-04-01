package ConHash

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math"

	"github.com/dgraph-io/badger"
)

type Node struct {
	ID        string
	CName     string
	NumTokens int
	DBPath    string // e.g /tmp/badger
	NodeDB    *badger.DB
	HHQueue   *badger.DB
	IP        string //a.b.c.d:port
	Port      string

	NodeRingPositions []int
	Ring              *Ring

	//{name: str, nodeRingPosition: int}
	NodeDataArray []NodeData
}

type Ring struct {
	MaxID             int // 0 to maxID inclusive
	RingNodeDataArray []NodeData
	NodePrefList      map[int][]NodeData //Map node/virtualNode unique hash to a list of nodeData of virtual/physical nodes belonging to another host
	ReplicationFactor int
	RWFactor          int
	NodeStatuses      map[string]bool
}

type NodeData struct {
	//Node Data contains ip and port and hash This helps in determining which node is responsible for
	//Which request(read/write) and contains relevant info (ip:port) to
	//direct data to that node
	ID    string
	CName string
	Hash  int
	IP    string
	Port  string
}

func ToChar(i int) rune {
	return rune('A' - 1 + i)
}

func NewNodeData(id string, cName string, hash int, ip string, port string) NodeData {
	return NodeData{id, cName, hash, ip, port}
}

func NewNode(numID int, numTokens int, DBPath string, ip string, port string, ring *Ring) Node {
	return Node{
		ID:        string(ToChar(numID)) + "0",
		CName:     string(ToChar(numID)),
		NumTokens: numTokens,
		DBPath:    DBPath,
		IP:        ip,
		Port:      port,
		Ring:      ring,
	}
}

func NewNodeServer(numID int, numTokens int, DBPath string, ip string, port string) Node {
	return Node{
		ID:        string(ToChar(numID)) + "0",
		CName:     string(ToChar(numID)),
		NumTokens: numTokens,
		DBPath:    DBPath,
		IP:        ip,
		Port:      port,
	}
}

func NewRing(maxID int, replicationFactor int, rwFactor int) *Ring {
	nodeDataArray := make([]NodeData, maxID, maxID)
	nodePrefList := make(map[int][]NodeData, maxID)
	fmt.Println(len(nodeDataArray))
	fmt.Println(nodeDataArray[1].ID)
	return &Ring{maxID, nodeDataArray,
		nodePrefList, replicationFactor, rwFactor, map[string]bool{}}
}

//node will create numTokens worth of virtual nodes
func (n *Node) RegisterWithRing(r *Ring) {
	nodeDataArray := []NodeData{}

	for i := 0; i < n.NumTokens+1; i++ {
		id := fmt.Sprintf("%s%d", n.CName, i)
		hash := HashMD5(id, 0, r.MaxID)
		nodeDataArray = append(nodeDataArray, NewNodeData(id, n.CName, hash, n.IP, n.Port))
	}

	fmt.Printf("Node %s registering %s \n", n.ID, ToString(nodeDataArray))
	n.NodeDataArray = r.RegisterNodes(nodeDataArray)
	fmt.Printf("Ring registered for %s: %s  \n", n.ID, ToString(n.NodeDataArray))
}

const RING_URL = config.RING_URL

func (r *Ring) RegisterNodes(nodeDataArray []NodeData) []NodeData {
	fmt.Println("Registering...")
	ret := []NodeData{}
	for _, nd := range nodeDataArray {
		counter := 2
		for {
			//Changed to something like quadratic probing that avoids clusterring better
			if r.RingNodeDataArray[nd.Hash].ID != "" {
				nd.Hash = int(math.Pow(float64(nd.Hash+counter), 2)) % r.MaxID
			} else {
				r.RingNodeDataArray[nd.Hash] = nd
				ret = append(ret, nd)
				break
			}

			counter += 1

		}
	}
	fmt.Println("Done Registering...")
	return ret
}

func ToString(nodeDataArray []NodeData) []string {
	ret := []string{}
	for _, nd := range nodeDataArray {
		ret = append(ret, fmt.Sprintf("(%s, %d) ", nd.ID, nd.Hash))
	}
	return ret
}

//string must end with an integer
func (r *Ring) GetNode(id string) (string, error) {
	var NodeNotFound = errors.New("Node not found")
	hash := HashMD5(id, 0, r.MaxID)

	//Impose an upper bound for probe times
	for i := 0; i < len(r.RingNodeDataArray); i++ {

		if r.RingNodeDataArray[hash].ID == id {
			ip_port := fmt.Sprintf("%s:%s", r.RingNodeDataArray[hash].IP, r.RingNodeDataArray[hash].Port)

			return ip_port, nil
		}
		hash = (hash + 1) % len(r.RingNodeDataArray)
	}

	return id, NodeNotFound
}

func HashMD5(text string, min int, max int) int {
	byteArray := md5.Sum([]byte(text))
	var output int
	for _, num := range byteArray {
		output += int(num)
	}
	return output%(max-min+1) + min
}

// Inputs : key of data, machine's local ring
// Outputs : int (Node hash of the node that's supposed to be responsible for the data)
//           url of Node that's responsible
func (ring *Ring) AllocateKey(key string) (int, string, error) {
	var NodeNotFound = errors.New("Node not found")

	nodeArray := ring.RingNodeDataArray
	keyHash := HashMD5(key, 0, len(nodeArray)-1)
	var firstNodeAddress int //Keep a pointer to the first node address encountered just in case
	firstNodeAddress = -1    // -1 is an impossible number in context of node array,
	//using it as a benchmark to see if it has not been set
	for i := 0; i < len(nodeArray); i++ {
		if nodeArray[i].ID != "" {
			if firstNodeAddress == -1 {
				firstNodeAddress = i
			}
			if keyHash <= i {
				nodeUrl := fmt.Sprintf("%s:%s", nodeArray[i].IP, nodeArray[i].Port)
				return i, nodeUrl, nil
			}
		}
		if i == len(nodeArray)-1 {
			//Reached end of node array and a coordinator node for key still hasn't been allocated
			nodeUrl := fmt.Sprintf("%s:%s", nodeArray[firstNodeAddress].IP, nodeArray[firstNodeAddress].Port)
			return firstNodeAddress, nodeUrl, nil
		}
	}
	return -1, "", NodeNotFound
}

func (ring *Ring) GenPrefList() {
	nodeArray := ring.RingNodeDataArray
	for i := 0; i < len(nodeArray); i++ {
		if nodeArray[i].ID != "" {
			// if node not empty, assign preference list
			ring.NodePrefList[i] = func(i int) []NodeData {
				ret := []NodeData{}
				j := (i + 1) % ring.MaxID
				//Need to track Nodes that have already become a replica by their CNames
				currentReplicas := map[string]struct{}{}
				currentReplicas[nodeArray[i].CName] = struct{}{}
				for j != i {
					_, alrReplica := currentReplicas[nodeArray[j].CName]
					if nodeArray[j].ID != "" && (alrReplica == false) {
						ret = append(ret, nodeArray[j])
						currentReplicas[nodeArray[j].CName] = struct{}{}
						if len(ret) == ring.ReplicationFactor {
							return ret
						}
					}
					j = (j + 1) % ring.MaxID
				}
				return ret
			}(i) // finish assigning preference list to 1 node
		}
	}
}
