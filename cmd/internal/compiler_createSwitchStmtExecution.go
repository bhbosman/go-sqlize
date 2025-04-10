package internal

import "go/ast"

func (compiler *Compiler) createCaseClauseExecution(node Node[*ast.CaseClause]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var nodes []Node[ast.Node]
		for _, expr := range node.Node.List {
			param := ChangeParamNode[*ast.CaseClause, ast.Expr](node, expr)
			tempState := state.setCurrentNode(ChangeParamNode[*ast.CaseClause, ast.Node](node, expr))
			nn, _ := compiler.findRhsExpression(tempState, param)(tempState)
			nodes = append(nodes, nn...)
		}
		return nodes, artValue
	}
}

func (compiler *Compiler) createSwitchStmtExecution(node Node[*ast.SwitchStmt]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode(node, node.Node.Body)
		tempState := state.setCurrentNode(ChangeParamNode[*ast.SwitchStmt, ast.Node](node, node.Node.Body))
		compiler.executeBlockStmt(tempState, param)
		panic("fdsfdsfds")

		//if node.Node.Tag != nil && node.Node.Body != nil {
		//	//param := ChangeParamNode(node, node.Node.Tag)
		//	//tempState := state.setCurrentNode(ChangeParamNode[*ast.SwitchStmt, ast.Node](node, node.Node.Tag))
		//	//expression, _ := compiler.findRhsExpression(tempState, param)(state)
		//
		//	for _, stmt := range node.Node.Body.List {
		//		stmtParam := ChangeParamNode[*ast.SwitchStmt, ast.Stmt](node, stmt)
		//		stmtTempState := state.setCurrentNode(ChangeParamNode[*ast.SwitchStmt, ast.Node](node, stmt))
		//		es, n := compiler.findStatement(stmtTempState, stmtParam)
		//		tempState := state.setCurrentNode(n)
		//		es(tempState)
		//		//
		//		//compiler.findStatement(stmtTempState, stmtParam)
		//		//bodyItemParam := ChangeParamNode[*ast.SwitchStmt, ast.Expr](node, bodyItem)
		//		//bodyItemTempState := state.setCurrentNode(ChangeParamNode[*ast.SwitchStmt, ast.Node](node, node.Node.Tag))
		//		//
		//		//compiler.findRhsExpression(bodyItemTempState, bodyItemParam)
		//
		//	}
		//
		//	switch node.Node.Tag.(type) {
		//	case *ast.Ident:
		//	default:
		//		panic(node.Node.Tag)
		//
		//	}
		//
		//	panic("unhandled default case")
		//
		//} else {
		//	panic(node.Node.Tag)
		//}
	}
}
