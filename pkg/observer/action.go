package observer

import (
	"fmt"
	"reflect"
)

type Action interface {
	Execute(event ...interface{})
}

type action struct {
	target interface{}
	method reflect.Method
}

func NewAction(target interface{}, method reflect.Method) Action {
	a := &action{
		target: target,
		method: method,
	}
	return Action(a)
}

func (a action) Execute(event ...interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	if event != nil && a.method.Type.NumIn()-1 == len(event) {
		params := make([]reflect.Value, 0)
		params = append(params, reflect.ValueOf(a.target))
		for _, param := range event {
			params = append(params, reflect.ValueOf(param))
		}
		a.method.Func.Call(params)
	}
}
