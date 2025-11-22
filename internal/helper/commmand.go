package helper

import (
	"encoding/json"
	"fmt"

	"github.com/codingbot24-s/distributed-kv-store/internal/raft"
)

type Command struct {
	OP    string
	Key   string
	Value string
}

func NewCommand() *Command {
	return &Command{}
}

// handler will call with command and this will append to the wal log file
// return the encoded bytes of command
func ApplyCommand(cmd *Command) error {
	// encode the command in jsonByte
	// encode directly in log entry
	// b, err := encode(cmd)
	//if err != nil {
	//	return fmt.Errorf("error encoding command: %w", err)
	//}
	// will append the bytes into the wal file
	w, err := GetWal()
	if err != nil {
		return fmt.Errorf("error getting wal: %w", err)
	}
	// append to log file
	index := w.Index
	term := w.Term
	l := raft.NewLogEntry()
	logEntry := l.CreateLogEntry(index, term, cmd)
	byteLog, err := raft.EncodeLog(logEntry)
	if err != nil {
		return fmt.Errorf("error encoding log: %w", err)
	}
	err = w.Append(byteLog)
	if err != nil {
		return fmt.Errorf("error appending into the file: %w", err)
	}
	var logEntryStruct raft.LogEntry
	// unmarshall json byte
	err = json.Unmarshal(byteLog, &logEntryStruct)
	if err != nil {
		return fmt.Errorf("error decoding command: %w", err)
	}
	e, err := GetEngine()
	if err != nil {
		return fmt.Errorf("error getting engine: %w", err)
	}
	for _, cmd := range logEntryStruct.Command {
		switch cmd.OP {
		case "set":
			e.set(cmd.Key, cmd.Value)
		default:
			return fmt.Errorf("unknown command: %s", cmd.OP)
		}

	}
	return nil
}

func encode(cmd *Command) ([]byte, error) {
	jsonByte, err := json.Marshal(*cmd)
	if err != nil {
		return nil, fmt.Errorf("error marshaling in the json byte %w", err)
	}
	return jsonByte, nil
}

func DecodeCommand(data []byte) (*Command, error) {
	var cmd Command
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		return nil, fmt.Errorf("error decoding the command: %w", err)
	}
	return &cmd, nil
}
