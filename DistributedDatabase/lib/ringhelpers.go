package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

// Local helper function
func GetSortedKeys(ringNodeDataMap map[int]NodeData) []int {
	keys := []int{}
	for k, _ := range ringNodeDataMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// Helper for Ringserver to find a node's precedessors
func FindPredecessors(thisNodeKey int, ringNodeDataMap map[int]NodeData) map[int]SimpleNodeData {
	keys := GetSortedKeys(ringNodeDataMap)
	thisNodeKeyIndex := FindIndexOfArray(thisNodeKey, keys)
	predecessorKeyIndex := thisNodeKeyIndex
	predecessors := make(map[int]SimpleNodeData)
	for i := 0; i < REPLICATION_FACTOR-1; i++ {
		predecessorKeyIndex -= 1
		if predecessorKeyIndex < 0 {
			predecessorKeyIndex += len(keys)
		}
		if predecessorKeyIndex == thisNodeKeyIndex {
			continue
		}
		predecessorKey := keys[predecessorKeyIndex]
		predecessorNodeData := ringNodeDataMap[predecessorKey]
		predecessorSimpleNodeData := SimpleNodeData{
			Id:   predecessorNodeData.Id,
			Ip:   predecessorNodeData.Ip,
			Port: predecessorNodeData.Port,
			Hash: predecessorNodeData.Hash,
		}
		predecessors[predecessorKey] = predecessorSimpleNodeData
	}
	return predecessors
}

// Helper for Ringserver to find a node's successors
func FindSuccessors(thisNodeKey int, ringNodeDataMap map[int]NodeData) map[int]SimpleNodeData {
	keys := GetSortedKeys(ringNodeDataMap)
	thisNodeKeyIndex := FindIndexOfArray(thisNodeKey, keys)
	successorKeyIndex := thisNodeKeyIndex
	successors := make(map[int]SimpleNodeData)
	for i := 0; i < REPLICATION_FACTOR-1; i++ {
		successorKeyIndex += 1
		successorKeyIndex %= len(keys)
		if successorKeyIndex == thisNodeKeyIndex {
			continue
		}
		successorKey := keys[successorKeyIndex]
		successorNodeData := ringNodeDataMap[successorKey]
		successorSimpleNodeData := SimpleNodeData{
			Id:   successorNodeData.Id,
			Ip:   successorNodeData.Ip,
			Port: successorNodeData.Port,
			Hash: successorNodeData.Hash,
		}
		successors[successorKey] = successorSimpleNodeData
	}
	return successors
}

// Helper for Ringserver to update a node about it's predecessors
func UpdatePredecessors(predecessors map[int]SimpleNodeData, destNode NodeData) {
	requestBody, _ := json.Marshal(predecessors)
	postURL := fmt.Sprintf("http://%s:%s/update-predecessors", destNode.Ip, destNode.Port)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Println("Successfully updated Successor's predecessors. Response:", string(responseBody))
	} else {
		fmt.Println("Failed to update Successor's predecessors. Response:", string(responseBody))
	}
}

// Helper for Ringserver to update a node about it's successors
func UpdateSuccessors(successors map[int]SimpleNodeData, destNode NodeData) {
	requestBody, _ := json.Marshal(successors)
	postURL := fmt.Sprintf("http://%s:%s/update-successors", destNode.Ip, destNode.Port)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Println("Successfully updated Predecessor's successors. Response:", string(responseBody))
	} else {
		fmt.Println("Failed to update Predecessor's successors. Response:", string(responseBody))
	}
}

// Helper for Ringserver to find a node's immediate successor
func FindImmediateSuccessor(thisNodeKey int, ringNodeDataMap map[int]NodeData) NodeData {
	keys := GetSortedKeys(ringNodeDataMap)
	thisNodeKeyIndex := FindIndexOfArray(thisNodeKey, keys)
	immediateSuccessorKeyIndex := thisNodeKeyIndex
	immediateSuccessorKeyIndex += 1
	immediateSuccessorKeyIndex %= len(keys)
	immediateSuccessorNodeData := ringNodeDataMap[keys[immediateSuccessorKeyIndex]]
	return immediateSuccessorNodeData
}

// Ringserver requests newNode's immeidate successor to transfer data to it.
func RequestData(successorNodeData NodeData, newNodeData NodeData) {
	requestBody, _ := json.Marshal(newNodeData)
	postURL := fmt.Sprintf("http://%s:%s/transferdata", successorNodeData.Ip, successorNodeData.Port)
	resp, err := http.Post(postURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Println("Successully told newNode's successor about newNode. Response:", string(responseBody))
	} else {
		fmt.Println("Failed to tell newNode's successor about newNode. Reason:", string(responseBody))
	}
}
