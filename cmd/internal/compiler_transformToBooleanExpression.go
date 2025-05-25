package internal

import (
	"go/ast"
	"go/token"
	"reflect"
)

func (compiler *Compiler) transformToBooleanExpression(state State, op token.Token, node Node[ast.Node]) Node[BooleanCondition] {
	allConditions := compiler.internalTransformToBooleanExpression(state, nil, node)

	var uniqueConditions []Node[ast.Node]
	m := make(map[uint32]bool)
	for _, condition := range allConditions {
		hashValue := compiler.calculateHash(condition)
		if _, ok := m[hashValue]; !ok {
			m[hashValue] = true
			uniqueConditions = append(uniqueConditions, condition)
		}
	}

	booleanCondition := BooleanCondition{uniqueConditions, op}
	return ChangeParamNode[ast.Node, BooleanCondition](node, booleanCondition)
}

func (compiler *Compiler) internalTransformToBooleanExpression(state State, conditions []Node[ast.Node], node Node[ast.Node]) []Node[ast.Node] {
	switch nodeItem := node.Node.(type) {
	default:
		panic(nodeItem)
	case BooleanCondition:
		conditions = append(conditions, node)
		return conditions
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
				p01 := ChangeParamNode[ast.Node, ast.Node](node, mbe)
				conditions = append(conditions, p01)
			}
		}
		return conditions
	}
}
