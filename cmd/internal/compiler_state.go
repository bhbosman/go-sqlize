package internal

import (
	"go/ast"
	"reflect"
)

type IABC interface{}

type State struct {
	arr         []IABC
	currentNode Node[ast.Node]
}

func (s State) setCurrentNode(node Node[ast.Node]) State {
	return State{s.arr, node}
}

func RemoveCompilerState[a IABC](state State) State {
	var result []IABC
	for _, compileState := range state.arr {
		if reflect.TypeFor[a]() == reflect.TypeOf(compileState) {
			continue
		}
		result = append(result, compileState)
	}

	return State{result, state.currentNode}
}

func SetCompilerState[a IABC](data a, state State) State {
	result := State{[]IABC{data}, state.currentNode}
	for _, compileState := range state.arr {
		if reflect.TypeFor[a]() == reflect.TypeOf(compileState) {
			continue
		}
		result.arr = append(result.arr, compileState)
	}
	return result
}

func GetCompilerState[a IABC](state State) a {
	for _, compileState := range state.arr {
		if reflect.TypeFor[a]() == reflect.TypeOf(compileState) {
			return compileState.(a)
		}
		if vv, ok := compileState.(a); ok {
			return vv
		}
	}
	var unk IABC = nil
	vv, _ := unk.(a)
	return vv
}
