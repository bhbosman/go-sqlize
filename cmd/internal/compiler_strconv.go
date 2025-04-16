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
	return func(state State, typeParams ITypeMapperArray, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.strconvItoaCompiled(state, arguments)
	}
}

func (compiler *Compiler) strconvItoaCompiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	rv := reflect.ValueOf(strconv.Itoa)
	if outputNodes, art, b := compiler.genericCall(state, rv, compiledArguments); b {
		return outputNodes, art
	}
	return compiler.coercionStringCompiled(state, compiledArguments)
}

func (compiler *Compiler) strconvAtoiImplementation(state State) ExecuteStatement {
	return func(state State, typeParams ITypeMapperArray, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.strconvAtoiCompiled(state, arguments)
	}
}

func (compiler *Compiler) strconvAtoiCompiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	rv := reflect.ValueOf(strconv.Atoi)
	if outputNodes, art, b := compiler.genericCall(state, rv, compiledArguments); b {
		return outputNodes, art
	}
	return compiler.coercionIntCompiled(state, compiledArguments)
}
