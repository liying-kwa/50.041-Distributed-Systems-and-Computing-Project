package main

type RingServer struct {
	ip   string
	port string
}

type Ring struct {
	MaxID             int // 0 to maxID inclusive
	RingNodeDataArray []NodeData
}

type NodeData struct {
	ID string
	//CName string
	//Hash  int
	IP   string
	Port string
}
