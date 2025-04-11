package internal

import (
	"go/ast"
	"reflect"
)

func (compiler *Compiler) findSources(item *TrailRecord) []string {
	m := map[string]bool{}
	for idx := 0; idx < item.Value.NumField(); idx++ {
		if node, ok := item.Value.Field(idx).Interface().(Node[ast.Node]); ok {
			compiler.internalFindSources(node, m)
		}
	}
	ss := make([]string, 0, len(m))
	for key, _ := range m {
		ss = append(ss, key)
	}
	return ss
}

func (compiler *Compiler) internalFindSources(node Node[ast.Node], m map[string]bool) {
	if !node.Valid {
		return
	}
	switch nodeItem := node.Node.(type) {
	case *EntityField:
		m[nodeItem.alias] = true
		break
	case *ast.BasicLit:
		break
	case *coercion:
		compiler.internalFindSources(nodeItem.Node, m)
		break
	case *BinaryExpr:
		compiler.internalFindSources(nodeItem.left, m)
		compiler.internalFindSources(nodeItem.right, m)
		break
	case *ReflectValueExpression:
		// nothing to do
		break
	case *SupportedFunction:
		for _, param := range nodeItem.params {
			compiler.internalFindSources(param, m)
		}
		break
	case *NilValueExpression:
		// nothing to do
		break
	case *MultiBinaryExpr:
		for _, expression := range nodeItem.expressions {
			compiler.internalFindSources(expression, m)
		}
	case *IfThenElseSingleValueCondition:
		for _, conditionalStatement := range nodeItem.conditionalStatement {
			compiler.internalFindSources(conditionalStatement.condition, m)
			compiler.internalFindSources(conditionalStatement.value, m)
		}
		break
	case *LhsToMultipleRhsOperator:
		compiler.internalFindSources(nodeItem.Lhs, m)
		for _, rhs := range nodeItem.Rhs {
			compiler.internalFindSources(rhs, m)
		}
	default:
		panic(reflect.TypeOf(node.Node).String())
	}
}
