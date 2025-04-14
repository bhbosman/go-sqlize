package internal

import (
	"fmt"
	"go/ast"
	"os"
	"reflect"
)

func (compiler *Compiler) addOsFunctions() {
	compiler.GlobalFunctions[ValueKey{"os", "Getwd"}] = compiler.osGetWdImplementation
	compiler.GlobalFunctions[ValueKey{"os", "ModePerm"}] = compiler.osModePermImplementation
	compiler.GlobalFunctions[ValueKey{"os", "MkdirAll"}] = compiler.osMkdirAllImplementation
	compiler.GlobalFunctions[ValueKey{"os", "Create"}] = compiler.osCreateImplementation
	compiler.GlobalFunctions[ValueKey{"os", "Stdout"}] = compiler.osStdout
}

func (compiler *Compiler) genericValue(rv reflect.Value) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		rvNode := &ReflectValueExpression{rv}
		vv := ChangeParamNode[ast.Node, ast.Node](state.currentNode, rvNode)
		return []Node[ast.Node]{vv}, artValue
	}
}

func (compiler *Compiler) osModePermImplementation(State, []Node[ast.Expr], []Node[ast.Node]) ExecuteStatement {
	rv := reflect.ValueOf(os.ModePerm)
	return compiler.genericValue(rv)
}

func (compiler *Compiler) osStdout(State, []Node[ast.Expr], []Node[ast.Node]) ExecuteStatement {
	rv := reflect.ValueOf(os.Stdout)
	return compiler.genericValue(rv)
}

func (compiler *Compiler) osGetWdImplementation(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(os.Getwd)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic(fmt.Errorf("os.Getwd only accept literal values"))
	}
}

func (compiler *Compiler) osMkdirAllImplementation(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(os.MkdirAll)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic(fmt.Errorf("os.MkdirAll only accept literal values"))
	}
}

func (compiler *Compiler) osCreateImplementation(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(os.Create)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic(fmt.Errorf("os.Create only accept literal values"))
	}
}
