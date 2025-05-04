package internal

import (
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

func (compiler *Compiler) runBasicLitNode(state State, node Node[*ast.BasicLit]) ([]Node[ast.Node], CallArrayResultType) {
	switch node.Node.Kind {
	case token.INT:
		intValue, _ := strconv.ParseInt(node.Node.Value, 10, 64)
		rv := reflect.ValueOf(intValue).Convert(reflect.TypeFor[int]())
		param := ChangeParamNode[*ast.BasicLit, ast.Node](node, &ReflectValueExpression{rv, intValueKey})
		return []Node[ast.Node]{param}, artValue
	case token.FLOAT:
		floatValue, _ := strconv.ParseFloat(node.Node.Value, 64)
		param := ChangeParamNode[*ast.BasicLit, ast.Node](node, &ReflectValueExpression{reflect.ValueOf(floatValue), float64ValueKey})
		return []Node[ast.Node]{param}, artValue
	case token.IMAG:
		panic("ssfds")
	case token.CHAR:
		panic("ssfds")
	case token.STRING:
		stringValue, _ := strconv.Unquote(node.Node.Value)
		param := ChangeParamNode[*ast.BasicLit, ast.Node](node, &ReflectValueExpression{reflect.ValueOf(stringValue), stringValueKey})
		return []Node[ast.Node]{param}, artValue
	default:
		panic(notFound(node.Node.Kind.String(), "createRhsBasicLitExecution"))
	}
}

func (compiler *Compiler) createRhsBasicLitExecution(node Node[*ast.BasicLit]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.runBasicLitNode(state, node)
	}
}
