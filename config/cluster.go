package config

type Peer struct {
	NodeId string `json:"node_id"`
	Addr   string `json:"addr"`
}

type Cluster struct {
	NodeId string `json:"node_id"`
	Addr   string `json:"addr"`
	Peers  []Peer `json:"peers"`
}
