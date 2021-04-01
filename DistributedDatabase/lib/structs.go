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
