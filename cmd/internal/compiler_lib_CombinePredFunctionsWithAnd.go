package internal

import (
	"fmt"
	"go/ast"
	"go/token"
)

func (compiler *Compiler) libCombinePredFunctionsWithAndImplementation(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) == 0 {
			panic(fmt.Errorf("Lib.CombinePredFunctionsWithAnd implementation requires at least 1 argument, got %d", len(arguments)))
		}
		typeMapper := compiler.findType(state, arguments[0], Default)
		mbe := ChangeParamNode[*ast.FuncType, ast.Node](node, MultiBinaryExpr{token.AND, arguments, typeMapper})
		return []Node[ast.Node]{mbe}, artValue
	}
}
