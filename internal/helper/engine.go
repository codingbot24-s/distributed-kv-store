package helper

import (
	"errors"
	"fmt"
)

type Engine struct {
	Data map[string]string
}

var DefaultEngine *Engine

func NewEngine() {
	DefaultEngine = &Engine{
		Data: make(map[string]string),
	}
}

func GetEngine() (*Engine, error) {
	if DefaultEngine.Data == nil {
		return nil, errors.New("engine data is empty call the new engine")
	}

	return DefaultEngine, nil
}

func (e *Engine) set(key, value string) {
	e.Data[key] = value
}
func (e *Engine) get(key string) (string, bool) {
	value, ok := e.Data[key]
	if !ok {
		return "", false
	}
	return value, true
}
func (e *Engine) delete(key string) {
	delete(e.Data, key)
}

// this apply will create our state machine based on command
func (e *Engine) Apply(cmd *Command) error {
	c := *cmd
	switch c.OP {
	case "set":
		e.set(c.Key, c.Value)
	case "get":
		fmt.Println("key is and value is value", c.Key, c.Value)
	case "delete":
		e.delete(c.Key)
	default:
		return fmt.Errorf("unknown command: %s", c.OP)
	}
	return nil
}

func (e *Engine) Check() {
	fmt.Println("starting the check")
	for k, v := range e.Data {
		fmt.Printf("the key is %s and value is %s\n", k, v)
	}
}
