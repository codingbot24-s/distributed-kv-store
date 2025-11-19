package helper

type Engine struct {
	Data map[string]string
}

func (e *Engine) NewEngine () *Engine{
	return &Engine{
		Data: make(map[string]string),
	}
}

func(e *Engine) set (key,value string) {}
func(e *Engine) get (key string) {}
func(e *Engine) delete (key string) {}