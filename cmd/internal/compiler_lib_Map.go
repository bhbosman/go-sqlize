package internal

import (
	"fmt"
	"go/ast"
)

func (compiler *Compiler) libMapImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 2 {
			panic(fmt.Errorf("Lib.Map implementation requires 2 arguments, got %d", len(arguments)))
		}
		return compiler.internalLibMapImplementation(state, typeParams, arguments[0], arguments[1])
	}
}

func (compiler *Compiler) internalLibMapImplementation(state State, typeParams map[string]ITypeMapper, arg0, arg1 Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	switch argItem := arg1.Node.(type) {
	default:
		panic(argItem)
	case *ast.Ident, *ast.FuncLit:
		fn := compiler.findRhsExpression(state, arg1)
		v, _ := compiler.executeAndExpandStatement(state, typeParams, nil, fn)
		return compiler.internalLibMapImplementation(state, typeParams, arg0, v[0])
	case FuncLit:
		param := ChangeParamNode[ast.Node, FuncLit](arg1, argItem)
		return compiler.executeFuncLit(state, param, []Node[ast.Node]{arg0}, typeParams)
	}
}
