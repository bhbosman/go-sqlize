package internal

import (
	"go/ast"
	"go/token"
)

func (compiler *Compiler) findStatement(state State, node Node[ast.Stmt]) (ExecuteStatement, Node[ast.Node]) {
	switch item := node.Node.(type) {
	case *ast.FolderContextInformation:
		return func(state State) ([]Node[ast.Node], CallArrayResultType) {
			value := ChangeParamNode[ast.Stmt, ast.Node](node, item)
			return []Node[ast.Node]{value}, artFCI
		}, ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)
	case *ast.AssignStmt:
		value := ChangeParamNode(node, item)
		return compiler.createAssignStatementExecution(value), ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)
	case *ast.ExprStmt:
		param := ChangeParamNode(node, item.X)
		tempState := state.setCurrentNode(ChangeParamNode[ast.Stmt, ast.Node](node, item.X))
		return compiler.findRhsExpression(tempState, param), ChangeParamNode[ast.Stmt, ast.Node](node, item.X)
	case *ast.ReturnStmt:
		value := ChangeParamNode(node, item)
		return compiler.createReturnStmtExecution(value), ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)

	default:
		panic("dddd")
	}
}

func (compiler *Compiler) createReturnStmtExecution(node Node[*ast.ReturnStmt]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var result []Node[ast.Node]
		for _, expr := range node.Node.Results {
			param := ChangeParamNode(node, expr)
			tempState := state.setCurrentNode(ChangeParamNode[*ast.ReturnStmt, ast.Node](node, expr))
			fn := compiler.findRhsExpression(tempState, param)

			param01 := ChangeParamNode[*ast.ReturnStmt, ast.Node](node, expr)
			state = state.setCurrentNode(param01)
			v, _ := fn(state)
			result = append(result, v...)
		}
		return result, artReturn
	}
}

func (compiler *Compiler) createAssignStatementExecution(node Node[*ast.AssignStmt]) ExecuteStatement {
	switch node.Node.Tok {
	case token.DEFINE, token.ASSIGN:
		return func(state State) ([]Node[ast.Node], CallArrayResultType) {
			var rhsArray []Node[ast.Node]

			for _, rhsExpression := range node.Node.Rhs {
				param := ChangeParamNode(node, rhsExpression)
				tempState := state.setCurrentNode(ChangeParamNode[*ast.AssignStmt, ast.Node](node, rhsExpression))
				fn := compiler.findRhsExpression(tempState, param)
				state = state.setCurrentNode(ChangeParamNode[*ast.AssignStmt, ast.Node](node, rhsExpression))
				arr, _ := fn(state)
				rhsArray = append(rhsArray, arr...)
			}
			for idx, lhsExpression := range node.Node.Lhs {
				param := ChangeParamNode(node, lhsExpression)
				state = state.setCurrentNode(ChangeParamNode[*ast.AssignStmt, ast.Node](node, lhsExpression))
				assignStatement := compiler.findLhsExpression(state, param, node.Node.Tok)
				assignStatement(state, rhsArray[idx])
			}

			return nil, artNone
		}
	default:
		panic("dddd")
	}
}
