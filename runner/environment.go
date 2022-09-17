package runner

import (
	"errors"
)

type Environment struct {
	Values map[string]interface{}
}

func NewEnvironment() Environment {
	return Environment{
		Values: make(map[string]interface{},0),
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Get(name string) interface{} {
	val, ok := e.Values[name]
	if !ok {
		errors.New("Undefined variable '" + name + "'.")
	}
	return val

}
