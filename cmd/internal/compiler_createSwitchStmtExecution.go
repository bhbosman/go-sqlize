package internal

import (
	"go/ast"
	"go/token"
	"reflect"
	"sort"
)

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
		if node.Node.Tag != nil && node.Node.Body != nil {
			param := ChangeParamNode(node, node.Node.Tag)
			tempState := state.setCurrentNode(ChangeParamNode[*ast.SwitchStmt, ast.Node](node, node.Node.Tag))
			expression, _ := compiler.findRhsExpression(tempState, param)(state)

			paramBody := ChangeParamNode(node, node.Node.Body)
			tempState = state.setCurrentNode(ChangeParamNode[*ast.SwitchStmt, ast.Node](node, node.Node.Body))
			stmt, resultType := compiler.executeBlockStmt(tempState, paramBody)
			if resultType != artReturnAndContinue {
				panic("need a return statement")
			}
			sortNodes := &SortNodes{stmt, func(i, j int) bool {
				ith, iok := stmt[i].Node.(*CaseClauseNode)
				jth, jok := stmt[j].Node.(*CaseClauseNode)
				if jok && iok {
					if len(jth.nodes) == 0 && len(ith.nodes) > 0 {
						return true
					}
				}
				return false

			}}
			sort.Sort(sortNodes)
			var conditionalStatement []MultiValueCondition

			for _, n := range stmt {
				switch item := n.Node.(type) {
				case *CaseClauseNode:
					condition := func(item *CaseClauseNode) Node[ast.Node] {
						if len(item.nodes) == 0 {
							return ChangeParamNode[*ast.SwitchStmt, ast.Node](node, &ReflectValueExpression{reflect.ValueOf(true)})
						}
						return ChangeParamNode[*ast.SwitchStmt, ast.Node](node, &LhsToMultipleRhsOperator{token.EQL, token.LOR, expression[0], item.nodes})
					}(item)
					multiValueCondition := MultiValueCondition{condition: condition, values: item.arr}
					conditionalStatement = append(conditionalStatement, multiValueCondition)
				default:
					panic("need a case statement")
				}
			}
			ite := &IfThenElseMultiValueCondition{conditionalStatement}
			resultValue := ChangeParamNode[*ast.SwitchStmt, ast.Node](node, ite)
			return []Node[ast.Node]{resultValue}, artReturn
		}
		panic("need a return statement")
	}
}
