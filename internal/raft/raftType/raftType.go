package raftType

import "github.com/codingbot24-s/distributed-kv-store/internal/helper"

type AppendRequest struct {
	Term         int
	LeaderID     string
	PrevLogIdx   int
	PrevLogTerm  int
	Entries      helper.LogEntry
	LeaderCommit int
}
type AppendResponse struct {
	Term    int
	Success bool
}
