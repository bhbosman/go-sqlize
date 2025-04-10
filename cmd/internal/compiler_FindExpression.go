package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

func (compiler *Compiler) findRhsExpression(state State, node Node[ast.Expr]) ExecuteStatement {
	return compiler.internalFindRhsExpression(0, state, node).(ExecuteStatement)
}

func (compiler *Compiler) internalFindRhsExpression(stackIndex int, state State, node Node[ast.Expr]) interface{} {
	switch item := node.Node.(type) {
	case *ast.UnaryExpr:
		param := ChangeParamNode(node, item)
		return compiler.createRhsUnaryExprExecution(param)
	case *ast.BinaryExpr:
		param := ChangeParamNode(node, item)
		return compiler.createRhsBinaryExprExecution(param)
	case *ast.BasicLit:
		param := ChangeParamNode(node, item)
		return compiler.createRhsBasicLitExecution(param)
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		unk := compiler.internalFindRhsExpression(stackIndex+1, state, param)
		switch vv := unk.(type) {
		case ast.ImportMapEntry:
			vk := ValueKey{vv.Path, item.Sel.Name}
			if globalFunction, ok := compiler.GlobalFunctions[vk]; ok {
				return globalFunction(state, nil, nil)
			}
			panic(notFound(fmt.Sprintf("%v", vk), "internalFindRhsExpression"))
		case Node[ast.Node]:
			switch vvv := vv.Node.(type) {
			case *TrailSource:
				var es ExecuteStatement = func(state State) ([]Node[ast.Node], CallArrayResultType) {
					result := ChangeParamNode[ast.Expr, ast.Node](node, &EntityField{node.Node.Pos(), vvv.Alias, item.Sel.Name})
					return []Node[ast.Node]{result}, artValue
				}
				return es
			}
			panic("implement me")
		default:
			return unk
		}
	case *ast.CallExpr:
		param := ChangeParamNode(node, item)
		return compiler.createRhsCallExpressionExecution(param)
	case *ast.CompositeLit:
		param := ChangeParamNode(node, item)
		return compiler.createRhsCompositeLitExecution(param)
	case *ast.FuncLit:
		param := ChangeParamNode(node, item)
		return compiler.createRhsFuncLitExprExecution(param)
	case *ast.ParenExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		return compiler.findRhsExpression(state, param)
	case *ast.Ident:
		currentContext := GetCompilerState[*CurrentContext](state)
		if value, b := currentContext.FindValue(item.Name); b {
			if stackIndex == 0 {
				var es ExecuteStatement = func(state State) ([]Node[ast.Node], CallArrayResultType) {
					return []Node[ast.Node]{value}, artValue
				}
				return es
			}
			return value
		}
		if globalFunction, ok := compiler.GlobalFunctions[ValueKey{node.RelPath, item.Name}]; ok {
			return globalFunction(state, nil, nil)
		}
		if globalFunction, ok := compiler.GlobalFunctions[ValueKey{"", item.Name}]; ok {
			return globalFunction(state, nil, nil)
		}
		if path, ok := node.ImportMap[item.Name]; ok {
			return path
		}
		panic("unhandled default case")

	default:
		panic(node.Node)
	}
}

func (compiler *Compiler) createRhsUnaryExprExecution(node Node[*ast.UnaryExpr]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		switch node.Node.Op {
		case token.SUB:
			tempState := state.setCurrentNode(ChangeParamNode[*ast.UnaryExpr, ast.Node](node, node.Node.X))
			param := ChangeParamNode(node, node.Node.X)
			arr, _ := compiler.findRhsExpression(tempState, param)(tempState)

			rve := &ReflectValueExpression{reflect.ValueOf(-1)}
			rveNode := ChangeParamNode[*ast.UnaryExpr, ast.Node](node, rve)

			be := &BinaryExpr{arr[0].Node.Pos(), token.MUL, arr[0], rveNode}
			beNode := ChangeParamNode[*ast.UnaryExpr, ast.Node](node, be)
			return []Node[ast.Node]{beNode}, artValue

		case token.NOT:
			tempState := state.setCurrentNode(ChangeParamNode[*ast.UnaryExpr, ast.Node](node, node.Node.X))
			param := ChangeParamNode(node, node.Node.X)
			arr, _ := compiler.findRhsExpression(tempState, param)(tempState)
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

func (compiler *Compiler) createRhsBinaryExprExecution(node Node[*ast.BinaryExpr]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode(node, node.Node.X)
		tempState := state.setCurrentNode(ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.X))
		x, _ := compiler.findRhsExpression(tempState, param)(tempState)
		rvX, isXLiteral := isLiterateValue(x[0])

		param = ChangeParamNode(node, node.Node.Y)
		tempState = state.setCurrentNode(ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.Y))
		y, _ := compiler.findRhsExpression(tempState, param)(tempState)
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

				panic("unhandled default case")
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
			panic("unhandled default case")

		default:
			panic("implement me")
		}

	}
}

func (compiler *Compiler) createRhsFuncLitExprExecution(node Node[*ast.FuncLit]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode[*ast.FuncLit, ast.Node](node, node.Node)
		return []Node[ast.Node]{param}, artValue
	}
}

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

func (compiler *Compiler) findFunction(state State, node Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return compiler.internalFindFunction(0, state, node, arguments).(ExecuteStatement)
}

func (compiler *Compiler) initExecutionStatement(state State, stackIndex int, unk interface{}, typeParams []Node[ast.Expr], arguments []Node[ast.Node]) interface{} {
	if stackIndex != 0 {
		return unk
	}
	switch value := unk.(type) {
	case OnCreateExecuteStatement:
		return value(state, typeParams, arguments)
	case Node[ast.Node]:
		switch value02 := value.Node.(type) {
		case ast.Expr:
			param := ChangeParamNode(value, value02)
			unkValue := compiler.internalFindFunction(stackIndex+1, state, param, arguments)
			return compiler.initExecutionStatement(state, stackIndex, unkValue, typeParams, arguments)
		default:
			panic(unk)
		}
	default:
		panic(unk)
	}
}

func (compiler *Compiler) onFuncLitExecutionStatement(node Node[*ast.FuncLit]) OnCreateExecuteStatement {
	return func(state State, typeParams []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
		return func(state State) ([]Node[ast.Node], CallArrayResultType) {
			var names []*ast.Ident
			if node.Node.Type.Params != nil {
				for _, field := range node.Node.Type.Params.List {
					names = append(names, field.Names...)
				}
			}
			m := map[string]Node[ast.Node]{}
			for idx, name := range names {
				m[name.Name] = arguments[idx]
			}

			newContext := &CurrentContext{m, GetCompilerState[*CurrentContext](state)}
			state = SetCompilerState(newContext, state)
			param := ChangeParamNode[ast.Node, *ast.BlockStmt](state.currentNode, node.Node.Body)
			values, art := compiler.executeBlockStmt(state, param)
			state = SetCompilerState(newContext.Parent, state)
			return values, art
		}
	}
}

func (compiler *Compiler) internalFindFunction(stackIndex int, state State, node Node[ast.Expr], arguments []Node[ast.Node]) interface{} {
	switch item := node.Node.(type) {
	case *ast.FuncLit:
		param := ChangeParamNode[ast.Expr, *ast.FuncLit](node, item)
		var es OnCreateExecuteStatement = compiler.onFuncLitExecutionStatement(param)
		return compiler.initExecutionStatement(state, stackIndex, es, nil, arguments)
	case *ast.IndexExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		indexParam := ChangeParamNode(node, item.Index)
		unk := compiler.internalFindFunction(stackIndex+1, state, param, arguments)
		return compiler.initExecutionStatement(state, stackIndex, unk, []Node[ast.Expr]{indexParam}, arguments)
	case *ast.IndexListExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		var arrIndices []Node[ast.Expr]
		for _, index := range item.Indices {
			indexParam := ChangeParamNode(node, index)
			arrIndices = append(arrIndices, indexParam)
		}
		unk := compiler.internalFindFunction(stackIndex+1, state, param, arguments)
		return compiler.initExecutionStatement(state, stackIndex, unk, arrIndices, arguments)
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		unk := compiler.internalFindFunction(stackIndex+1, state, param, arguments)
		switch value := unk.(type) {
		case ast.ImportMapEntry:
			vk := ValueKey{value.Path, item.Sel.Name}
			returnValue, ok := compiler.GlobalFunctions[vk]
			if ok {
				return compiler.initExecutionStatement(state, stackIndex, returnValue, nil, arguments)
			}
			panic(fmt.Errorf("can not find function %s", vk))
		case Node[ast.Node]:
			switch nodeItem := value.Node.(type) {
			case *ReflectValueExpression:
				rvFn := nodeItem.Rv.MethodByName(item.Sel.Name)
				return compiler.initExecutionStatement(state, stackIndex, compiler.builtInStructMethods(rvFn), nil, arguments)
			default:
				panic(value.Node)
			}
		default:
			panic("sdfdsfds")
		}
	case *ast.Ident:
		if path, ok := node.ImportMap[item.Name]; ok {
			return compiler.initExecutionStatement(state, stackIndex, path, nil, arguments)
		}

		currentContext := GetCompilerState[*CurrentContext](state)
		if v, ok := currentContext.FindValue(item.Name); ok {
			return compiler.initExecutionStatement(state, stackIndex, v, nil, arguments)
		}

		if fn, ok := compiler.GlobalFunctions[ValueKey{"", item.Name}]; ok {
			return compiler.initExecutionStatement(state, stackIndex, fn, nil, arguments)
		}
		if fn, ok := compiler.GlobalFunctions[ValueKey{node.RelPath, item.Name}]; ok {
			return compiler.initExecutionStatement(state, stackIndex, fn, nil, arguments)
		}

		panic(item.Name)
	default:
		panic(node.Node)
	}
}

func (compiler *Compiler) builtInStructMethods(rv reflect.Value) OnCreateExecuteStatement {
	return func(state State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
		return func(state State) ([]Node[ast.Node], CallArrayResultType) {
			if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
				return outputNodes, art
			}
			panic(fmt.Errorf("builtInStructMethods only accept literal values"))
		}
	}
}

func (compiler *Compiler) createRhsCompositeLitExecution(node Node[*ast.CompositeLit]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode(node, node.Node.Type)
		rt := compiler.findType(state, param)
		rtKind := rt.Kind()
		switch rtKind {
		case reflect.Struct:
			rv := reflect.New(rt).Elem()
			for idx, elt := range node.Node.Elts {
				switch expr := elt.(type) {
				case *ast.KeyValueExpr:
					param = ChangeParamNode(node, expr.Value)
					tempState := state.setCurrentNode(ChangeParamNode[*ast.CompositeLit, ast.Node](node, expr.Value))
					vv, _ := compiler.findRhsExpression(tempState, param)(tempState)
					switch key := expr.Key.(type) {
					case *ast.Ident:
						rv.FieldByName(key.Name).Set(reflect.ValueOf(vv[0]))
					default:
						panic("unhandled key")
					}
				default:
					param = ChangeParamNode(node, elt)
					tempState := state.setCurrentNode(ChangeParamNode[*ast.CompositeLit, ast.Node](node, elt))
					vv, _ := compiler.findRhsExpression(tempState, param)(tempState)
					itemRv := reflect.ValueOf(vv[0])
					rv.Field(idx).Set(itemRv)

				}
			}

			nodeValue := ChangeParamNode[*ast.CompositeLit, ast.Node](
				node,
				&TrailRecord{
					node.Node.Pos(),
					rv,
				},
			)

			return []Node[ast.Node]{nodeValue}, artValue
		default:
			panic("dsfsfds")
		}
	}
}

func (compiler *Compiler) createRhsCallExpressionExecution(node Node[*ast.CallExpr]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var args []Node[ast.Node]
		for _, arg := range node.Node.Args {
			param := ChangeParamNode(state.currentNode, arg)
			tempState := state.setCurrentNode(ChangeParamNode[ast.Node, ast.Node](state.currentNode, arg))
			fn := compiler.findRhsExpression(tempState, param)
			nodeArg, _ := fn(state)
			args = append(args, nodeArg...)
		}
		param := ChangeParamNode(node, node.Node.Fun)
		tempState02 := state.setCurrentNode(ChangeParamNode[ast.Node, ast.Node](state.currentNode, node.Node.Fun))
		execFn := compiler.findFunction(tempState02, param, args)
		return execFn(tempState02)
	}
}
