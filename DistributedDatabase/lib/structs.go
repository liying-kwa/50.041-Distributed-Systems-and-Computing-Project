package lib

type Ring struct {
	MaxID           int // maxID in ring. if -1, means no node in ring
	RingNodeDataMap map[int]NodeData
}

type NodeData struct {
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
}

type TransferMessage struct {
	Ip   string
	Port string
	Hash string
}

type MessageType int

const (
	Get MessageType = iota
	Put
)
