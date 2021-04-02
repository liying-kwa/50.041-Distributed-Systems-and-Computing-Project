package lib

type Ring struct {
	// MaxID           int // 0 to maxID inclusive
	RingNodeDataMap map[int]NodeData
}

type NodeData struct {
	Id   int
	Ip   string
	Port string
}

type Message struct {
	Type     MessageType
	CourseId string
	Count    int
}

type MessageType int

const (
	Get MessageType = iota
	Put
)
