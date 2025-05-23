package internal

import (
	"go/ast"
	"go/token"
	"reflect"
)

func (compiler *Compiler) calculateTypeMapperFn(
	state State,
	op token.Token,
	typeMapperX, typeMapperY ITypeMapper,
) ITypeMapper {

	RtCanInt := func(rt reflect.Type) bool {
		switch rt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return true
		default:
			return false
		}
	}

	RtCanUint := func(rt reflect.Type) bool {
		switch rt.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return true
		default:
			return false
		}
	}
	RtCanFloat := func(rt reflect.Type) bool {
		switch rt.Kind() {
		case reflect.Float32, reflect.Float64:
			return true
		default:
			return false
		}
	}

	RtCanString := func(rt reflect.Type) bool {
		switch rt.Kind() {
		case reflect.String:
			return true
		default:
			return false
		}
	}
	RtCanBool := func(rt reflect.Type) bool {
		switch rt.Kind() {
		case reflect.Bool:
			return true
		default:
			return false
		}
	}

	isBoolOperand := func(op token.Token) bool {
		switch op {
		case token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ, token.LAND, token.LOR:
			return true
		default:
			return false
		}
	}

	typX, _ := typeMapperX.ActualType()
	typY, _ := typeMapperY.ActualType()
	switch {
	default:
		panic("unhandled default case")
	case isBoolOperand(op):
		return compiler.registerBool()(state, nil)
	case RtCanInt(typX) && RtCanInt(typY):
		return compiler.registerInt64()(state, nil)
	case RtCanUint(typX) && RtCanUint(typY):
		return compiler.registerUint64()(state, nil)
	case RtCanFloat(typX) && RtCanFloat(typY):
		return compiler.registerFloat64()(state, nil)
	case RtCanString(typX) && RtCanString(typY):
		return compiler.registerString()(state, nil)
	case RtCanBool(typX) && RtCanBool(typY):
		return compiler.registerBool()(state, nil)
	case compiler.isTypeSomeDataType(typX):
		tm, _ := compiler.extractSomeDataTypeMapper(typX)
		return compiler.calculateTypeMapperFn(state, op, tm, typeMapperY)
	case compiler.isTypeSomeDataType(typY):
		tm, _ := compiler.extractSomeDataTypeMapper(typY)
		return compiler.calculateTypeMapperFn(state, op, typeMapperX, tm)
	}
}

func (compiler *Compiler) createRhsBinaryExprExecution(node Node[*ast.BinaryExpr]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.X)
		tempState := state.setCurrentNode(ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.X))
		esX := compiler.findRhsExpression(tempState, param)
		x, _ := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, esX)
		x0 := x[0]
		rvX, isXLiteral := isLiterateValue(x0)

		typeMapperX, _ := x0.Node.(IFindTypeMapper).GetTypeMapper(state)
		typeMapperX0 := typeMapperX[0]

		param = ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.Y)
		tempState = state.setCurrentNode(ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.Y))
		esY := compiler.findRhsExpression(tempState, param)
		y, _ := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, esY)
		y0 := y[0]
		rvY, isYLiteral := isLiterateValue(y0)
		typeMapperY, _ := y0.Node.(IFindTypeMapper).GetTypeMapper(state)
		typeMapperY0 := typeMapperY[0]

		binTypeMapper := compiler.calculateTypeMapperFn(state, node.Node.Op, typeMapperX0, typeMapperY0)

		switch {
		case !isXLiteral && !isYLiteral:
			newOp := BinaryExpr{node.Node.Op, x[0], y[0], binTypeMapper}
			return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
		case !isXLiteral && isYLiteral:
			newNode := ChangeParamNode[*ast.BinaryExpr, ast.Node](node, &ReflectValueExpression{rvY, ValueKey{"JKL", "MNO"}})
			newOp := BinaryExpr{node.Node.Op, x[0], newNode, binTypeMapper}
			return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
		case isXLiteral && !isYLiteral:
			newNode := ChangeParamNode[*ast.BinaryExpr, ast.Node](node, &ReflectValueExpression{rvX, ValueKey{"JKL", "MNO"}})
			newOp := BinaryExpr{node.Node.Op, newNode, y[0], binTypeMapper}
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
					newOp := &ReflectValueExpression{reflect.ValueOf(rvX.Int() + rvY.Int()), intValueKey}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				case token.SUB:
					newOp := &ReflectValueExpression{reflect.ValueOf(rvX.Int() - rvY.Int()), intValueKey}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				case token.MUL:
					newOp := &ReflectValueExpression{reflect.ValueOf(rvX.Int() * rvY.Int()), intValueKey}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				case token.QUO:
					newOp := &ReflectValueExpression{reflect.ValueOf(rvX.Int() / rvY.Int()), intValueKey}
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
					newOp := &ReflectValueExpression{reflect.ValueOf(rvX.String() + rvY.String()), stringValueKey}
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
					newOp := &ReflectValueExpression{reflect.ValueOf(false), boolValueKey}
					return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
				case token.EQL:
					newOp := &ReflectValueExpression{reflect.ValueOf(true), boolValueKey}
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
