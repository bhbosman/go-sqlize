package internal

import (
	"go/ast"
	"reflect"
)

func (compiler *Compiler) createIfStmtExecution(node Node[*ast.IfStmt]) ExecuteStatement {
	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		var conditionalStatement []MultiValueCondition
		var whatIsReturned CallArrayResultType = 0

		newContext := &CurrentContext{ValueInformationMap{}, map[string]ITypeMapper{}, LocalTypesMap{}, GetCompilerState[*CurrentContext](state)}
		state = SetCompilerState(newContext, state)

		if node.Node.Init != nil {
			param := ChangeParamNode(node, node.Node.Init)
			es, n := compiler.findStatement(state, param)
			tempState := state.setCurrentNode(n)
			_, _ = compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, es)
		}

		fnBoolExpression := func(state State, Cond ast.Expr) ([]Node[ast.Node], CallArrayResultType) {
			param := ChangeParamNode(node, Cond)
			es := compiler.findRhsExpression(state, param)
			return compiler.executeAndExpandStatement(state, typeParams, unprocessedArgs, es)
		}
		boolExpression, _ := fnBoolExpression(state, node.Node.Cond)

		rv, isLogical := isLiterateValue(boolExpression[0])
		if isLogical {
			println(rv.String())
			panic("implement")
		} else {
			if node.Node.Body != nil {
				parent := GetCompilerState[*CurrentContext](state)
				newBodyContext := &CurrentContext{ValueInformationMap{}, map[string]ITypeMapper{}, LocalTypesMap{}, parent}
				state = SetCompilerState(newBodyContext, state)
				{
					param := ChangeParamNode[*ast.IfStmt, ast.Stmt](node, node.Node.Body)
					es, n := compiler.findStatement(state, param)
					tempState := state.setCurrentNode(n)
					bodyValues, resultTypeForBodyPart := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, es)
					whatIsReturned |= resultTypeForBodyPart
					if resultTypeForBodyPart == artReturn {
						conditionalStatementInstance := MultiValueCondition{boolExpression[0], bodyValues}
						conditionalStatement = append(conditionalStatement, conditionalStatementInstance)
					}
				}
				state = SetCompilerState(newBodyContext.Parent, state)
			}
			if node.Node.Else != nil {
				parent := GetCompilerState[*CurrentContext](state)

				newElseContext := &CurrentContext{ValueInformationMap{}, map[string]ITypeMapper{}, LocalTypesMap{}, parent}
				state = SetCompilerState(newElseContext, state)
				{
					param := ChangeParamNode[*ast.IfStmt, ast.Stmt](node, node.Node.Else)
					es, n := compiler.findStatement(state, param)
					tempState := state.setCurrentNode(n)
					bodyElseValues, resultTypeForBodyElsePart := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, es)
					whatIsReturned |= resultTypeForBodyElsePart
					if resultTypeForBodyElsePart == artReturn {
						switch item := bodyElseValues[0].Node.(type) {
						case *IfThenElseSingleValueCondition:
							for idx, stmt := range item.conditionalStatement {
								var nodes []Node[ast.Node]
								for _, value := range bodyElseValues {
									pe := value.Node.(*IfThenElseSingleValueCondition)
									nodes = append(nodes, pe.conditionalStatement[idx].value)
								}
								conditionalStatementInstance := MultiValueCondition{stmt.condition, nodes}
								conditionalStatement = append(conditionalStatement, conditionalStatementInstance)
							}
						default:
							condition := ChangeParamNode[ast.Node, ast.Node](bodyElseValues[0], &ReflectValueExpression{reflect.ValueOf(true)})
							conditionalStatementInstance := MultiValueCondition{condition, bodyElseValues}
							conditionalStatement = append(conditionalStatement, conditionalStatementInstance)
						}
					}
				}
				state = SetCompilerState(newElseContext.Parent, state)
			}
			switch {
			case whatIsReturned == artReturn: // looking for this
				break
			default:
				err := syntaxErrorf(state.currentNode, "an if statement should return on both legs or not on any")
				panic(err)
			}
		}
		state = SetCompilerState(newContext.Parent, state)
		ite := &IfThenElseMultiValueCondition{conditionalStatement}
		resultValue := ChangeParamNode[*ast.IfStmt, ast.Node](node, ite)
		return []Node[ast.Node]{resultValue}, artReturn
	}
}
