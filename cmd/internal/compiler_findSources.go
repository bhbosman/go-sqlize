package internal

import (
	"go/ast"
	"reflect"
)

type ISource interface {
	SourceName() string
	Dependencies() map[string]ISource
}

type PrimarySource struct {
	sourceName string
}

func (ps PrimarySource) Dependencies() map[string]ISource {
	return map[string]ISource{}
}

func (ps PrimarySource) SourceName() string {
	return ps.sourceName
}

func (compiler *Compiler) findSourcesFromNode(node Node[ast.Node]) map[string]ISource {
	m := map[string]bool{}
	compiler.internalFindSourcesFromNode(node, m)
	arr := map[string]ISource{}
	for key, _ := range m {
		arr[key] = PrimarySource{key}
	}
	return arr
}

func (compiler *Compiler) internalFindSourcesFromNode(node Node[ast.Node], m map[string]bool) {
	if !node.Valid {
		return
	}
	switch nodeItem := node.Node.(type) {
	case *CheckForNotNullExpression:
		compiler.internalFindSourcesFromNode(nodeItem.node, m)

	case EntityField:
		m[nodeItem.alias] = true
		break
	case *ast.BasicLit:
		break
	case coercion:
		compiler.internalFindSourcesFromNode(nodeItem.Node, m)
		break
	case BinaryExpr:
		compiler.internalFindSourcesFromNode(nodeItem.left, m)
		compiler.internalFindSourcesFromNode(nodeItem.right, m)
		break
	case *ReflectValueExpression:
		// nothing to do
		break
	case *SupportedFunction:
		for _, param := range nodeItem.params {
			compiler.internalFindSourcesFromNode(param, m)
		}
		break
	case MultiBinaryExpr:
		for _, expression := range nodeItem.expressions {
			compiler.internalFindSourcesFromNode(expression, m)
		}
	case IfThenElseSingleValueCondition:
		for _, conditionalStatement := range nodeItem.conditionalStatement {
			compiler.internalFindSourcesFromNode(conditionalStatement.condition, m)
			compiler.internalFindSourcesFromNode(conditionalStatement.value, m)
		}
		break
	case *LhsToMultipleRhsOperator:
		compiler.internalFindSourcesFromNode(nodeItem.Lhs, m)
		for _, rhs := range nodeItem.Rhs {
			compiler.internalFindSourcesFromNode(rhs, m)
		}
	case *TrailRecord:
		rv := nodeItem.Value
		for idx := range rv.NumField() {
			switch rvIdxField := rv.Field(idx).Interface().(type) {
			case Node[ast.Node]:
				if !rvIdxField.Valid {
					continue
				}
				compiler.internalFindSourcesFromNode(rvIdxField, m)
			}
		}
	default:
		panic(reflect.TypeOf(node.Node).String())
	}
}
