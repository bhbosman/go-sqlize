package internal

import (
	"go/ast"
	"go/token"
	"reflect"
)

func (compiler *Compiler) createRhsBinaryExprExecution(node Node[*ast.BinaryExpr]) ExecuteStatement {
	return func(state State, typeParams ITypeMapperArray, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.X)
		tempState := state.setCurrentNode(ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.X))
		esX := compiler.findRhsExpression(tempState, param)
		x, _ := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, esX)
		rvX, isXLiteral := isLiterateValue(x[0])

		param = ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.Y)
		tempState = state.setCurrentNode(ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.Y))
		esY := compiler.findRhsExpression(tempState, param)
		y, _ := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, esY)
		rvY, isYLiteral := isLiterateValue(y[0])

		switch {
		case !isXLiteral && !isYLiteral:
			newOp := &BinaryExpr{node.Node.OpPos, node.Node.Op, x[0], y[0]}
			return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
		case !isXLiteral && isYLiteral:
			newNode := ChangeParamNode[*ast.BinaryExpr, ast.Node](node, &ReflectValueExpression{rvY})
			newOp := &BinaryExpr{node.Node.OpPos, node.Node.Op, x[0], newNode}
			return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
		case isXLiteral && !isYLiteral:
			newNode := ChangeParamNode[*ast.BinaryExpr, ast.Node](node, &ReflectValueExpression{rvX})
			newOp := &BinaryExpr{node.Node.OpPos, node.Node.Op, newNode, y[0]}
			return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
		case isXLiteral && isYLiteral:
			calculateKind := func(rv reflect.Value) reflect.Kind {
				switch {
				case rv.CanInt():
					return reflect.Int
				case rv.CanUint():
					return reflect.Uint
				case rv.CanFloat():
					return reflect.Float64
				case rv.Kind() == reflect.String:
					return reflect.String
				case rv.Kind() == reflect.Invalid:
					return reflect.Invalid
				default:
					panic("unhandled default case")
				}
			}
			kindX := calculateKind(rvX)
			KindY := calculateKind(rvY)
			switch {
			case kindX == KindY && kindX == reflect.Int:
				switch node.Node.Op {
				case token.ADD:
					newOp := &ReflectValueExpression{reflect.ValueOf(rvX.Int() + rvY.Int())}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				case token.SUB:
					newOp := &ReflectValueExpression{reflect.ValueOf(rvX.Int() - rvY.Int())}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				case token.MUL:
					newOp := &ReflectValueExpression{reflect.ValueOf(rvX.Int() * rvY.Int())}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				case token.QUO:
					newOp := &ReflectValueExpression{reflect.ValueOf(rvX.Int() / rvY.Int())}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				default:
					panic("unhandled default case")
				}
			case kindX == KindY && kindX == reflect.Uint:
				panic("unhandled default case")
			case kindX == KindY && kindX == reflect.Float64:
				panic("unhandled default case")
			case kindX == KindY && kindX == reflect.String:
				switch node.Node.Op {
				case token.ADD:
					newOp := &ReflectValueExpression{reflect.ValueOf(rvX.String() + rvY.String())}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				default:
					panic("unhandled default case")
				}

			case kindX == reflect.Invalid && KindY != reflect.Invalid:
				panic("unhandled default case")
			case kindX != reflect.Invalid && KindY == reflect.Invalid:
				panic("unhandled default case")
			case kindX == reflect.Invalid && KindY == reflect.Invalid:
				switch node.Node.Op {
				case token.NEQ:
					newOp := &ReflectValueExpression{reflect.ValueOf(false)}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				case token.EQL:
					newOp := &ReflectValueExpression{reflect.ValueOf(true)}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				default:
					panic("unhandled default case")
				}

			default:
				panic("unhandled default case")
			}
		default:
			panic("implement me")
		}
	}
}
