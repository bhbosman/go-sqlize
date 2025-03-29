package internal

import (
	"go/ast"
	"reflect"
)

func (compiler *Compiler) createIfStmtExecution(node Node[*ast.IfStmt]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {

		var conditionalStatement []ConditionalStatement
		//artResultCount := 0
		var whatIsReturned CallArrayResultType = 0
		var isLogical bool = false

		newContext := &CurrentContext{map[string]Node[ast.Node]{}, GetCompilerState[*CurrentContext](state)}
		state = SetCompilerState(newContext, state)
		{
			if node.Node.Init != nil {
				param := ChangeParamNode(node, node.Node.Init)
				es, n := compiler.findStatement(state, param)
				tempState := state.setCurrentNode(n)
				_, _ = es(tempState)
			}
			param := ChangeParamNode(node, node.Node.Cond)
			es := compiler.findRhsExpression(state, param)
			boolExpression, _ := es(state)

			var rv reflect.Value
			rv, isLogical = isLiterateValue(boolExpression[0])

			{
				// do the body part of the if statement
				if !isLogical || (isLogical && rv.Bool()) {
					if node.Node.Body != nil {
						parent := GetCompilerState[*CurrentContext](state)
						parent = parent.Flatten()
						newBodyContext := &CurrentContext{map[string]Node[ast.Node]{}, parent}
						state = SetCompilerState(newBodyContext, state)
						{
							param := ChangeParamNode[*ast.IfStmt, ast.Stmt](node, node.Node.Body)
							es, n := compiler.findStatement(state, param)
							tempState := state.setCurrentNode(n)
							bodyValues, resultTypeForBodyPart := es(tempState)
							whatIsReturned |= resultTypeForBodyPart
							if resultTypeForBodyPart == artReturn {
								conditionalStatementInstance := ConditionalStatement{boolExpression[0], bodyValues}
								conditionalStatement = append(conditionalStatement, conditionalStatementInstance)
							}
						}
						state = SetCompilerState(newBodyContext.parent, state)
					}
				}
				// do the else part of the if statement
				if !isLogical || (isLogical && !rv.Bool()) {
					if node.Node.Else != nil {
						parent := GetCompilerState[*CurrentContext](state)
						parent = parent.Flatten()
						newElseContext := &CurrentContext{map[string]Node[ast.Node]{}, parent}
						state = SetCompilerState(newElseContext, state)
						{
							addResultToConditions := func(bodyElseValues []Node[ast.Node], resultTypeForBodyElsePart CallArrayResultType) {
								switch nodeItem := bodyElseValues[0].Node.(type) {
								case *IfThenElseCondition:
									conditionalStatement = append(conditionalStatement, nodeItem.conditionalStatement...)
								case *ReflectValueExpression:
									condition := ChangeParamNode[ast.Node, ast.Node](bodyElseValues[0], &ReflectValueExpression{reflect.ValueOf(true)})
									conditionalStatementInstance := ConditionalStatement{condition, bodyElseValues}
									conditionalStatement = append(conditionalStatement, conditionalStatementInstance)
								default:
									panic("ddddd")
								}
							}
							param := ChangeParamNode[*ast.IfStmt, ast.Stmt](node, node.Node.Else)
							es, n := compiler.findStatement(state, param)
							tempState := state.setCurrentNode(n)
							bodyElseValues, resultTypeForBodyElsePart := es(tempState)
							whatIsReturned |= resultTypeForBodyElsePart
							switch {
							case resultTypeForBodyElsePart == artReturn && len(bodyElseValues) > 0:
								addResultToConditions(bodyElseValues, resultTypeForBodyElsePart)
							default:
								panic("ddddd")
							}
						}
						state = SetCompilerState(newElseContext.parent, state)
					}
				}
			}
			switch {
			case (whatIsReturned == artReturn): // looking for this
				break
			default:
				err := syntaxErrorf(state.currentNode, "an if statement should return on both legs or not on any")
				panic(err)
			}
		}
		state = SetCompilerState(newContext.parent, state)
		ite := &IfThenElseCondition{conditionalStatement}
		resultValue := ChangeParamNode[*ast.IfStmt, ast.Node](node, ite)
		return []Node[ast.Node]{resultValue}, artReturn

	}
}
