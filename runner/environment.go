package runner

import (
	"errors"
)

type Environment struct {
	Enclosing *Environment
	Values map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		Values: make(map[string]interface{},0),
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Assign(name string, value interface{}) {
	if _, ok := e.Values[name]; ok {
		e.Values[name] = value
		return
	}

	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}

	panic(errors.New("Undefined variable '" + name + "'."))
}

func (e *Environment) Get(name string) interface{} {
	val, ok := e.Values[name]
	if !ok {
		if e.Enclosing != nil {
			return e.Enclosing.Get(name)
		}
		panic(errors.New("Undefined variable '" + name + "'."))
	}
	return val

}
