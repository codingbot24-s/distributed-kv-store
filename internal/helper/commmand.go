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

// handler will call with command and this will append to the file
func ApplyCommand(cmd *Command) error {
	// encode the command in jsonByte
	b, err := encode(cmd)
	if err != nil {
		return fmt.Errorf("error encoding command: %w", err)
	}
	// will append the bytes into the wal file
	w, err := GetWal()
	if err != nil {
		return fmt.Errorf("error getting wal: %w", err)
	}
	err = w.Append(b)
	if err != nil {
		return fmt.Errorf("error appending into the file: %w", err)
	}
	var c Command
	// unmarshall json byte
	err = json.Unmarshal(b, &c)
	if err != nil {
		return fmt.Errorf("error decoding command: %w", err)
	}
	e, err := GetEngine()
	if err != nil {
		return fmt.Errorf("error getting engine: %w", err)
	}
	switch c.OP {
	case "set":
		e.set(c.Key, c.Value)
	default:
		return fmt.Errorf("unknown command: %s", c.OP)
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
