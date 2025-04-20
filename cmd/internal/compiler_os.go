package internal

import (
	"fmt"
	"go/ast"
	"os"
	"reflect"
)

func (compiler *Compiler) addOsFunctions() {
	compiler.GlobalFunctions[ValueKey{"os", "Getwd"}] = functionInformation{compiler.osGetWdImplementation, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"os", "ModePerm"}] = functionInformation{compiler.osModePermImplementation, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"os", "MkdirAll"}] = functionInformation{compiler.osMkdirAllImplementation, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"os", "Create"}] = functionInformation{compiler.osCreateImplementation, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"os", "Stdout"}] = functionInformation{compiler.osStdout, Node[*ast.FuncType]{}, false}
}

func (compiler *Compiler) genericValue(rv reflect.Value, vk ValueKey) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		rvNode := &ReflectValueExpression{rv, vk}
		vv := ChangeParamNode[ast.Node, ast.Node](state.currentNode, rvNode)
		return []Node[ast.Node]{vv}, artValue
	}
}

func (compiler *Compiler) osModePermImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	rv := reflect.ValueOf(os.ModePerm)
	return compiler.genericValue(rv, ValueKey{})
}

func (compiler *Compiler) osStdout(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	rv := reflect.ValueOf(os.Stdout)
	return compiler.genericValue(rv, ValueKey{})
}

func (compiler *Compiler) osGetWdImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {

	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(os.Getwd)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic(fmt.Errorf("os.Getwd only accept literal values"))
	}
}

func (compiler *Compiler) osMkdirAllImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(os.MkdirAll)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic(fmt.Errorf("os.MkdirAll only accept literal values"))
	}
}

func (compiler *Compiler) osCreateImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(os.Create)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic(fmt.Errorf("os.Create only accept literal values"))
	}
}
