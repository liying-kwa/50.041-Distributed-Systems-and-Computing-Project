package lib

type Ring struct {
	RingNodeDataMap map[int]NodeData
}

type NodeData struct {
	Id   int
	Ip   string
	Port string
}

type RingServer struct {
	ip   string
	port string
	ring Ring
}
