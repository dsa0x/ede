package object

func NewEnvironment(outer *Environment) *Environment {
	return &Environment{store: make(map[string]Object), outer: outer}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

// Get retrieves the object from the environment, if not in the current env,
// it looks deeper into the parent. The environment returned
// is the one where the key was found
func (e *Environment) Get(key string) (Object, bool) {
	if obj, ok := e.store[key]; ok {
		return obj, ok
	}
	if e.outer != nil {
		return e.outer.Get(key)
	}
	return nil, false
}

// Update updates the key
func (e *Environment) Update(key string, val Object) bool {
	if _, ok := e.store[key]; ok {
		e.store[key] = val
		return ok
	}
	if e.outer != nil {
		return e.outer.Update(key, val)
	}
	return false
}

func (e *Environment) Set(key string, value Object) {
	e.store[key] = value
}
