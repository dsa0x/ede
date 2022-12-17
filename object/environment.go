package object

func NewEnvironment(outer *Environment) *Environment {
	return &Environment{store: make(map[string]Object), outer: outer}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(key string) (Object, bool) {
	if obj, ok := e.store[key]; ok {
		return obj, ok
	}
	if e.outer != nil {
		return e.outer.Get(key)
	}
	return nil, false
}

func (e *Environment) Set(key string, value Object) {
	e.store[key] = value
}
