package object

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), outer: nil}
}
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok {
		if e.outer != nil {
			obj, ok = e.outer.Get(name)
		}
	}
	return obj, ok
}
func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}
