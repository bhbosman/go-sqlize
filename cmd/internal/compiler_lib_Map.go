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
		if _, ok := arguments[1].Node.(FuncLit); !ok {
			panic("map implementation requires function literal")
		}
		if funcLit, ok := arguments[1].Node.(FuncLit); ok {
			return compiler.executeFuncLit(state, ChangeParamNode(arguments[1], funcLit), arguments, typeParams)
		}
		panic("map implementation argument 1 is not a function literal")
	}
}
