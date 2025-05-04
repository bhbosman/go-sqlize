package internal

import (
	"go/ast"
	"go/token"
	"reflect"
)

func (compiler *Compiler) transformToBooleanExpression(state State, node Node[ast.Node]) Node[ast.Node] {
	conditions := compiler.internalTransformToBooleanExpression(state, nil, node)
	booleanCondition := BooleanCondition{conditions}
	return ChangeParamNode[ast.Node, ast.Node](node, booleanCondition)
}

func (compiler *Compiler) internalTransformToBooleanExpression(state State, conditions []Node[ast.Node], node Node[ast.Node]) []Node[ast.Node] {
	switch nodeItem := node.Node.(type) {
	default:
		panic(nodeItem)
	case BinaryExpr:
		conditions = append(conditions, node)
		return conditions
	case IfThenElseSingleValueCondition:
		for _, conditionalStatement := range nodeItem.conditionalStatement {
			value, b := isLiterateValue(conditionalStatement.value)
			if b {
				if value.Kind() == reflect.Bool {
					if value.Bool() {
						conditions = append(conditions, conditionalStatement.condition)
					}
				} else {
					panic("unhandled default case")
				}
			} else {
				localConditions := compiler.internalTransformToBooleanExpression(state, nil, conditionalStatement.value)
				expressions := []Node[ast.Node]{conditionalStatement.condition}
				expressions = append(expressions, localConditions...)
				mbe := MultiBinaryExpr{token.LAND, expressions, compiler.registerBool()(state, nil)}
				conditions = append(conditions, ChangeParamNode[ast.Node, ast.Node](node, mbe))
			}
		}
		return conditions
	}
}
