package internal

import (
	"go/ast"
	"go/token"
	"reflect"
)

func (compiler *Compiler) createRhsUnaryExprExecution(node Node[*ast.UnaryExpr]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		switch node.Node.Op {
		case token.SUB:
			tempState := state.setCurrentNode(ChangeParamNode[*ast.UnaryExpr, ast.Node](node, node.Node.X))
			param := ChangeParamNode[*ast.UnaryExpr, ast.Node](node, node.Node.X)
			es := compiler.findRhsExpression(tempState, param)
			arr, _ := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, es)

			if rv, b := isLiterateValue(arr[0]); b {
				switch {
				case rv.CanInt():
					rv = reflect.ValueOf(-1 * rv.Int())
				case rv.CanUint():
					panic("invalid unary Uint()")
				case rv.CanFloat():
					rv = reflect.ValueOf(-1.0 * rv.Float())
				}
				rve := &ReflectValueExpression{rv}
				rveNode := ChangeParamNode[*ast.UnaryExpr, ast.Node](node, rve)
				return []Node[ast.Node]{rveNode}, artValue
			} else {
				rve := &ReflectValueExpression{reflect.ValueOf(-1)}
				rveNode := ChangeParamNode[*ast.UnaryExpr, ast.Node](node, rve)

				be := &BinaryExpr{arr[0].Node.Pos(), token.MUL, arr[0], rveNode}
				beNode := ChangeParamNode[*ast.UnaryExpr, ast.Node](node, be)
				return []Node[ast.Node]{beNode}, artValue
			}

		case token.NOT:
			tempState := state.setCurrentNode(ChangeParamNode[*ast.UnaryExpr, ast.Node](node, node.Node.X))
			param := ChangeParamNode[*ast.UnaryExpr, ast.Node](node, node.Node.X)
			es := compiler.findRhsExpression(tempState, param)
			arr, _ := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, es)
			switch nodeItem := arr[0].Node.(type) {
			case *BinaryExpr:
				switch nodeItem.Op {
				case token.NEQ:
					be := &BinaryExpr{arr[0].Node.Pos(), token.EQL, nodeItem.left, nodeItem.right}
					beNode := ChangeParamNode[*ast.UnaryExpr, ast.Node](node, be)
					return []Node[ast.Node]{beNode}, artValue
				default:
					panic("implement me")
				}
			case *EntityField:
				right := ChangeParamNode[*ast.UnaryExpr, ast.Node](node, &ReflectValueExpression{reflect.ValueOf(true)})
				be := &BinaryExpr{arr[0].Node.Pos(), token.NEQ, arr[0], right}
				beNode := ChangeParamNode[*ast.UnaryExpr, ast.Node](node, be)
				return []Node[ast.Node]{beNode}, artValue
			default:
				panic("implement me")
			}
		default:
			panic("implement me")
		}
	}
}
