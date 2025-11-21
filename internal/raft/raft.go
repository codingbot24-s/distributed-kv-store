package raft

import (
	"encoding/json"
	"fmt"

	"github.com/codingbot24-s/distributed-kv-store/internal/helper"
	"github.com/gofiber/fiber/v2"
)

type LogEntry struct {
	Index   uint64
	Term    uint64
	Command []helper.Command
}


func EncodeLog(l LogEntry) ([]byte, error) {
	jsonByte, err := json.Marshal(l)
	if err != nil {
		return nil, fmt.Errorf("error marshalling log: %v", err)
	}
	return jsonByte, nil
}


func DecodeLog(data []byte) (LogEntry, error) {
	var l LogEntry
	err := json.Unmarshal(data, &l)
	if err != nil {
		return l, fmt.Errorf("error decoding log entry: %w", err)
	}
	return l, nil
}


type AppendRequest struct {
	Term         int
	LeaderID     string
	PrevLogIdx   int
	PrevLogTerm  int
	Entries      LogEntry
	LeaderCommit int
}
type AppendResponse struct {
	Term    int
	Success bool
}


func Append(c *fiber.Ctx) error {
	return nil
}

func Vote(c *fiber.Ctx) error {
	return nil
}
``