package internal

import "go/ast"

func (compiler *Compiler) createRhsCallExpressionExecution(node Node[*ast.CallExpr]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var args []Node[ast.Node]
		for _, arg := range node.Node.Args {
			param := ChangeParamNode(state.currentNode, arg)
			tempState := state.setCurrentNode(ChangeParamNode[ast.Node, ast.Node](state.currentNode, arg))
			fn := compiler.findRhsExpression(tempState, param)
			nodeArg, _ := compiler.executeAndExpandStatement(state, fn)
			args = append(args, nodeArg...)
		}
		param := ChangeParamNode(node, node.Node.Fun)
		tempState02 := state.setCurrentNode(ChangeParamNode[ast.Node, ast.Node](state.currentNode, node.Node.Fun))
		execFn := compiler.findFunction(tempState02, param, args)
		return execFn(tempState02)
	}
}
