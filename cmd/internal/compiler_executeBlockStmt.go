package internal

import (
	"go/ast"
	"reflect"
)

func (compiler *Compiler) executeBlockStmt(state State, node Node[*ast.BlockStmt]) ([]Node[ast.Node], CallArrayResultType) {
	var conditionalStatement []ConditionalStatement
	newContext := &CurrentContext{map[string]Node[ast.Node]{}, GetCompilerState[*CurrentContext](state)}
	state = SetCompilerState(newContext, state)
	for _, item := range node.Node.List {
		param := ChangeParamNode(node, item)
		tempState := state.setCurrentNode(ChangeParamNode[*ast.BlockStmt, ast.Node](node, item))
		statementFn, currentNode := compiler.findStatement(tempState, param)
		tempState = state.setCurrentNode(currentNode)
		arr, rt := statementFn(tempState)
		switch rt {
		case artFCI:
			switch instance := arr[0].Node.(type) {
			case *ast.FolderContextInformation:
				node = Node[*ast.BlockStmt]{node.Key, node.Node, instance.Imports, instance.AbsPath, instance.RelPath, instance.FileName, node.Fs, node.Valid}
			}
		case artValue:
		case artReturn:
			if len(arr) == 0 {
				panic("big error")
			}

			if _, ok := arr[0].Node.(*IfThenElseCondition); !ok && len(conditionalStatement) == 0 {
				return arr, artReturn
			}

			if itec, ok := arr[0].Node.(*IfThenElseCondition); ok {
				conditionalStatement = append(conditionalStatement, itec.conditionalStatement...)
			} else {
				condition := ChangeParamNode[*ast.BlockStmt, ast.Node](node, &ReflectValueExpression{reflect.ValueOf(true)})
				conditionalStatementInstance := ConditionalStatement{condition, arr}
				conditionalStatement = append(conditionalStatement, conditionalStatementInstance)
			}

			var partialExpressionArr []Node[ast.Node]
			for _, partial := range conditionalStatement[0].values {
				partialExpressionArr = append(partialExpressionArr, ChangeParamNode[ast.Node, ast.Node](partial, &PartialExpression{}))
			}

			for _, partial := range conditionalStatement {
				for idx := range partial.values {
					partialExpressionArr[idx].Node.(*PartialExpression).conditionalStatement = append(
						partialExpressionArr[idx].Node.(*PartialExpression).conditionalStatement,
						struct {
							condition Node[ast.Node]
							value     Node[ast.Node]
						}{partial.condition, partial.values[idx]})
				}
			}
			validNodes := isValidNodes(partialExpressionArr)
			if !validNodes {
				errNode := ChangeParamNode[*ast.BlockStmt, ast.Node](node, node.Node)
				err := syntaxErrorf(errNode, "something went wrong with building partialExpression")
				panic(err)
			}
			return partialExpressionArr, artReturn
		default:
			continue
		}
	}
	state = SetCompilerState(newContext.Parent, state)
	if len(conditionalStatement) > 0 {
		panic("implement me")
		//pe := &PartialExpressions{partials}
		//result := ChangeParamNode[*ast.BlockStmt, ast.Node](node, pe)
		//return []Node[ast.Node]{result}, artPartialReturn
	}
	return nil, artNone
}
