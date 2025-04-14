package internal

import (
	"go/ast"
	"path/filepath"
	"reflect"
)

func (compiler *Compiler) pathFilepathJoinImplementation(state State) ExecuteStatement {
	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		arguments := compiler.compileArguments(state, unprocessedArgs, typeParams)
		rv := reflect.ValueOf(filepath.Join)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic("implementation for this")
	}
}

func (compiler *Compiler) pathFilepathDirImplementation(state State) ExecuteStatement {

	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		arguments := compiler.compileArguments(state, unprocessedArgs, typeParams)
		rv := reflect.ValueOf(filepath.Dir)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic("implementation for this")
	}
}
