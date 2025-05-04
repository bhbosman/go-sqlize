package internal

import "go/ast"

func (compiler *Compiler) createRhsFuncLitExprExecution(state State, node Node[*ast.FuncLit]) ExecuteStatement {
	currentContext := GetCompilerState[*CurrentContext](state)
	flattenValues := currentContext.flattenVariables()
	fl := FuncLit{node.Node.Type, node.Node.Body, flattenValues}
	return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode[*ast.FuncLit, ast.Node](node, fl)
		return []Node[ast.Node]{param}, artValue
	}
}
