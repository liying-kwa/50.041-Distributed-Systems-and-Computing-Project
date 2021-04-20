package lib

type Ring struct {
	MaxID           int // maxID in ring. if -1, means no node in ring
	RingNodeDataMap map[int]NodeData
}

type NodeData struct {
	Id           int
	Ip           string
	Port         string
	Hash         string
	Predecessors map[int]SimpleNodeData
	Successors   map[int]SimpleNodeData
}

// NodeData that exclude Predecessors and Successors to prevent infinite recursion
type SimpleNodeData struct {
	Id   int
	Ip   string
	Port string
	Hash string
}

type Message struct {
	Type     MessageType
	CourseId string
	Count    string
	Hash     string
	Replica  bool
}

type TransferMessage struct {
	Ip      string
	Port    string
	Hash    string
	Replica bool
}

type MessageType int

const (
	Get MessageType = iota
	Put
)

type PortNo struct {
	PortNo string
}
