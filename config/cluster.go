package config

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *Cluster) CreatePeer() map[string]*http.Client {
	peerClients := make(map[string]*http.Client)
	for _, peer := range c.Peers {
		peerClients[peer.NodeId] = new(http.Client)
	}

	return peerClients
}
