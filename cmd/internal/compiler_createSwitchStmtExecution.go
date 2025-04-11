package internal

import (
	"go/ast"
	"go/token"
)

type CaseClauseNode struct {
	arr   []Node[ast.Node]
	nodes []Node[ast.Node]
}

func (ccn *CaseClauseNode) Pos() token.Pos {
	return token.NoPos
}

func (ccn *CaseClauseNode) End() token.Pos {
	return token.NoPos
}

func (compiler *Compiler) createCaseClauseExecution(node Node[*ast.CaseClause]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var nodes []Node[ast.Node]
		for _, expr := range node.Node.List {
			param := ChangeParamNode[*ast.CaseClause, ast.Expr](node, expr)
			tempState := state.setCurrentNode(ChangeParamNode[*ast.CaseClause, ast.Node](node, expr))
			es := compiler.findRhsExpression(tempState, param)
			nn, _ := compiler.executeAndExpandStatement(tempState, es)
			nodes = append(nodes, nn...)
		}

		for _, stmt := range node.Node.Body {
			param := ChangeParamNode(node, stmt)
			tempState := state.setCurrentNode(ChangeParamNode[*ast.CaseClause, ast.Node](node, stmt))
			statementFn, currentNode := compiler.findStatement(tempState, param)
			tempState = state.setCurrentNode(currentNode)
			arr, art := compiler.executeAndExpandStatement(tempState, statementFn)
			if art == artReturn {
				returnNode := ChangeParamNode[*ast.CaseClause, ast.Node](node, &CaseClauseNode{arr, nodes})
				return []Node[ast.Node]{returnNode}, artReturnAndContinue
			}
		}
		panic("need a return statement")
	}
}

func (compiler *Compiler) createSwitchStmtExecution(node Node[*ast.SwitchStmt]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode(node, node.Node.Body)
		tempState := state.setCurrentNode(ChangeParamNode[*ast.SwitchStmt, ast.Node](node, node.Node.Body))
		stmt, resultType := compiler.executeBlockStmt(tempState, param)
		if resultType != artReturnAndContinue {
			panic("need a return statement")
		}
		for _, n := range stmt {
			println(n.Node)

		}

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
