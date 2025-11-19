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

func ApplyCommand(w *Wal, e *Engine, cmd *Command) error {
	// encode the command in json
	b, err := encode(cmd)
	if err != nil {
		return fmt.Errorf("error encoding command: %w", err)
	}
	// will append the bytes into the wal file
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
	// set in the memory
	e.set(c.Key, c.Value)
	return nil
}

func encode(cmd *Command) ([]byte, error) {
	jsonByte, err := json.Marshal(*cmd)
	if err != nil {
		return nil, fmt.Errorf("error marshaling in the json byte %w", err)
	}
	return jsonByte, nil
}
