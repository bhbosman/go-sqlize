package internal

import (
	"go/ast"
	"reflect"
	"strconv"
)

func (compiler *Compiler) addStrconvFunctions() {
	compiler.GlobalFunctions[ValueKey{"strconv", "Itoa"}] = compiler.strconvItoaImplementation
	compiler.GlobalFunctions[ValueKey{"strconv", "Atoi"}] = compiler.strconvAtoiImplementation
}

func (compiler *Compiler) strconvItoaImplementation(state State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(strconv.Itoa)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		return compiler.coercionString(state, nil, arguments)(state)
	}
}

func (compiler *Compiler) strconvAtoiImplementation(state State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(strconv.Atoi)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		return compiler.coercionInt(state, nil, arguments)(state)
	}
}
