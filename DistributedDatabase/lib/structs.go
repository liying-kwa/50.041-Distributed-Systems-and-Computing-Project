package lib

type Ring struct {
	MaxID           int // maxID in ring. if -1, means no node in ring
	RingNodeDataMap map[int]NodeData
}

type NodeData struct {
	Id              int
	Ip              string
	Port            string
	Hash            string
	PredecessorIP   string
	PredecessorPort string
	SuccessorIP     string
	SuccessorPort   string
	PredecessorIP2   string 
	PredecessorPort2 string
	SuccessorIP2     string
	SuccessorPort2   string
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
