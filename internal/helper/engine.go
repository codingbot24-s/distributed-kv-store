package helper

type Engine struct {
	Data map[string]string
}

func NewEngine() *Engine {
	return &Engine{
		Data: make(map[string]string),
	}
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
