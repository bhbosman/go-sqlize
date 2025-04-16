package internal

import "go/ast"

func (compiler *Compiler) createRhsFuncLitExprExecution(node Node[*ast.FuncLit]) ExecuteStatement {
	return func(state State, typeParams ITypeMapperArray, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode[*ast.FuncLit, ast.Node](node, node.Node)
		return []Node[ast.Node]{param}, artValue
	}
}
