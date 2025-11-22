package node

import (
	"github.com/codingbot24-s/distributed-kv-store/internal/helper"
)

type RaftNode struct {
	CurrentTerm int64
	VotedFor    int64
	Log         []helper.LogEntry
	CommitIndex int64
	LastApplied int64
	// which follower from which index send
	NextIndex map[string]int64
	// how much node has replicated
	MatchIndex map[string]int64
}
