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
		var l []Node[ast.Node]
		for _, argument := range arguments {
			l = compiler.internalCombinePredFunctionsWithAndImplementation(state, l, argument)
		}

		mbe := ChangeParamNode[*ast.FuncType, ast.Node](node, MultiBinaryExpr{token.LAND, l, typeMapper})
		return []Node[ast.Node]{mbe}, artValue
	}
}

func (compiler *Compiler) internalCombinePredFunctionsWithAndImplementation(state State, list []Node[ast.Node], argument Node[ast.Node]) []Node[ast.Node] {
	switch arg := argument.Node.(type) {
	default:
		panic(arg)
	case *ast.CallExpr:
		return append(list, argument)
	case *ast.CompositeLit:
		for _, elt := range arg.Elts {
			list = append(list, ChangeParamNode[ast.Node, ast.Node](argument, elt))
		}
		return list
	}
}
