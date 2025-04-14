package internal

import (
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

func (compiler *Compiler) createRhsBasicLitExecution(node Node[*ast.BasicLit]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		switch node.Node.Kind {
		case token.INT:
			intValue, _ := strconv.ParseInt(node.Node.Value, 10, 64)
			param := ChangeParamNode[*ast.BasicLit, ast.Node](node, &ReflectValueExpression{reflect.ValueOf(intValue)})
			return []Node[ast.Node]{param}, artValue
		case token.FLOAT:
			floatValue, _ := strconv.ParseFloat(node.Node.Value, 64)
			param := ChangeParamNode[*ast.BasicLit, ast.Node](node, &ReflectValueExpression{reflect.ValueOf(floatValue)})
			return []Node[ast.Node]{param}, artValue
		case token.IMAG:
			panic("ssfds")
		case token.CHAR:
			panic("ssfds")
		case token.STRING:
			stringValue, _ := strconv.Unquote(node.Node.Value)
			param := ChangeParamNode[*ast.BasicLit, ast.Node](node, &ReflectValueExpression{reflect.ValueOf(stringValue)})
			return []Node[ast.Node]{param}, artValue
		default:
			panic(notFound(node.Node.Kind.String(), "createRhsBasicLitExecution"))
		}
	}
}
