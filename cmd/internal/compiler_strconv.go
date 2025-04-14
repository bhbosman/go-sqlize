package internal

import (
	"go/ast"
	"reflect"
	"strconv"
)

func (compiler *Compiler) addStrconvFunctions() {
	compiler.GlobalFunctions[ValueKey{"strconv", "Itoa"}] = functionInformation{compiler.strconvItoaImplementation, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"strconv", "Atoi"}] = functionInformation{compiler.strconvAtoiImplementation, Node[*ast.FuncType]{}, false}
}

func (compiler *Compiler) strconvItoaImplementation(state State) ExecuteStatement {

	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		arguments := compiler.compileArguments(state, unprocessedArgs, typeParams)
		rv := reflect.ValueOf(strconv.Itoa)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		return compiler.coercionString(state)(state, typeParams, unprocessedArgs)
	}
}

func (compiler *Compiler) strconvAtoiImplementation(state State) ExecuteStatement {
	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		arguments := compiler.compileArguments(state, unprocessedArgs, typeParams)
		rv := reflect.ValueOf(strconv.Atoi)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		return compiler.coercionInt(state)(state, typeParams, unprocessedArgs)
	}
}
