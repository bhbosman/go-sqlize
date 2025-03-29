package internal

import (
	"go/ast"
	"path/filepath"
	"reflect"
)

func (compiler *Compiler) pathFilepathJoinImplementation(state State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(filepath.Join)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic("implementation for this")
	}
}

func (compiler *Compiler) pathFilepathDirImplementation(state State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(filepath.Dir)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic("implementation for this")
	}
}
