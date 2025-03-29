package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
)

type CurrentContext struct {
	m      map[string]Node[ast.Node]
	parent *CurrentContext
}

func (self *CurrentContext) FindValue(value string) (Node[ast.Node], bool) {
	if v, ok := self.m[value]; ok {
		return v, true
	}
	if self.parent == nil {
		return Node[ast.Node]{}, false
	}
	return self.parent.FindValue(value)
}

func (compiler *Compiler) findLhsExpression(state State, node Node[ast.Expr], tok token.Token) AssignStatement {
	return compiler.internalFindLhsExpression(state, node, tok).(AssignStatement)
}

func (compiler *Compiler) internalFindLhsExpression(state State, node Node[ast.Expr], tok token.Token) interface{} {
	switch item := node.Node.(type) {
	case *ast.Ident:
		if item.Name == "_" {
			fn := func() AssignStatement {
				return func(state State, value Node[ast.Node]) {

				}
			}
			return fn()
		}
		currentContext := GetCompilerState[*CurrentContext](state)
		switch tok {
		case token.DEFINE:
			fn := func(currentContext *CurrentContext, key string) AssignStatement {
				return func(state State, value Node[ast.Node]) {
					currentContext.m[key] = value
				}
			}
			return fn(currentContext, item.Name)
		case token.ASSIGN:
			panic("implement me")
			// walk current context
			//for currentContext != nil {
			//
			//	currentContext = currentContext.parent
			//}
		default:
			panic("unhandled default case")
		}
	default:
		panic(item)
	}
}

func (compiler *Compiler) findRhsExpression(state State, node Node[ast.Expr]) ExecuteStatement {
	return compiler.internalFindRhsExpression(0, state, node).(ExecuteStatement)
}

func (compiler *Compiler) internalFindRhsExpression(stackIndex int, state State, node Node[ast.Expr]) interface{} {
	switch item := node.Node.(type) {
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
		return compiler.createRhsFuncLitExecution(param)

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

		if globalfunction, ok := compiler.GlobalFunctions[ValueKey{node.RelPath, item.Name}]; ok {
			return globalfunction(state, nil, nil)
		}
		panic("unhandled default case")

	default:
		panic(node.Node)
	}
}

func (compiler *Compiler) createRhsBinaryExprExecution(node Node[*ast.BinaryExpr]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var arr []Node[ast.Node]
		param := ChangeParamNode(node, node.Node.X)
		tempState := state.setCurrentNode(ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.X))
		x, _ := compiler.findRhsExpression(tempState, param)(tempState)
		arr = append(arr, x...)

		param = ChangeParamNode(node, node.Node.Y)
		tempState = state.setCurrentNode(ChangeParamNode[*ast.BinaryExpr, ast.Node](node, node.Node.Y))
		y, _ := compiler.findRhsExpression(tempState, param)(tempState)
		arr = append(arr, y...)

		newOp := &BinaryExpr{node.Node.OpPos, node.Node.Op, arr}
		return []Node[ast.Node]{ChangeParamNode[*ast.BinaryExpr, ast.Node](node, newOp)}, artValue
	}
}

func (compiler *Compiler) createRhsFuncLitExecution(node Node[*ast.FuncLit]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode[*ast.FuncLit, ast.Node](node, node.Node)
		return []Node[ast.Node]{param}, artValue
	}
}

func (compiler *Compiler) createRhsBasicLitExecution(node Node[*ast.BasicLit]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		param := ChangeParamNode[*ast.BasicLit, ast.Node](node, node.Node)
		return []Node[ast.Node]{param}, artValue
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
			values, _ := compiler.executeBlockStmt(state, param)
			state = SetCompilerState(newContext.parent, state)
			return values, artValue
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
		panic(item.Name)
	default:
		panic(node.Node)
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
