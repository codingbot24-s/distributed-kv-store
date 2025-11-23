package helper

import (
	"encoding/json"
	"fmt"
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
    w, err := GetWal()
    if err != nil {
        return fmt.Errorf("error getting wal: %w", err)
    }

    l := NewLogEntry()
    logEntry := l.CreateLogEntry(w.Index, w.Term, cmd)
	fmt.Println(logEntry)
    // Encode to JSON bytes
    byteLog, err := EncodeLog(*logEntry)
    if err != nil {
        return fmt.Errorf("error encoding log: %w", err)
    }

    // Append to WAL file
    err = w.Append(byteLog)
    if err != nil {
        return fmt.Errorf("error appending into the file: %w", err)
    }

    // Apply to engine
    e, err := GetEngine()
    if err != nil {
        return fmt.Errorf("error getting engine: %w", err)
    }
	// fix thiss we dont need array of cmd in one logentry just one cmd
    
	// apply to engine	
	switch logEntry.Command.OP {
	case "set":
		e.set(logEntry.Command.Key, logEntry.Command.Value)
	case "delete":
		e.Delete(logEntry.Command.Key)
	case "get":
		// Get operations typically don't modify state
		// but you might want to handle them differently
	default:
		return fmt.Errorf("unknown command: %s", logEntry.Command.OP)
	}
	
    return nil
}

func Encode(cmd *Command) ([]byte, error) {
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

// LOG ENTRY STRUCTER AND ITS METHOD


type LogEntry struct {
	Index   int64  `json:"Index"`
	Term    int64  `json:"Term"`
	Command Command `json:"Command"` 
}

// CreateLogEntry helper to create a log entry
func (l *LogEntry) CreateLogEntry(index, term int64, cmd *Command) *LogEntry {
	return &LogEntry{
		Index:   index,
		Term:    term,
		Command: *cmd, 
	}
}
func NewLogEntry() *LogEntry {
	return &LogEntry{}
}

// LogEntry represents a single WAL entry with one command


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
