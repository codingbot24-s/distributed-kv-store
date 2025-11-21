package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Peer struct {
	NodeId string `json:"node_id"`
	Addr   string `json:"addr"`
}

type Cluster struct {
	NodeId string `json:"node_id"`
	Addr   string `json:"addr"`
	Peers  []Peer `json:"peers"`
}

func LoadConfig(path string) (*Cluster, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %s", err.Error())
	}
	var c Cluster
	if err := json.Unmarshal(content, &c); err != nil {
		return nil, fmt.Errorf("unmarshal config file error: %s", err.Error())
	}

	return &c, nil
}

func (c *Cluster) CreatePeer() {
	
}
