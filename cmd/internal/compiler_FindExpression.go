package internal

import (
	"fmt"
	"go/ast"
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
			case *TrailRecord:
				return func(trailRecord *TrailRecord, sel *ast.Ident) ExecuteStatement {
					return func(state State) ([]Node[ast.Node], CallArrayResultType) {
						return []Node[ast.Node]{trailRecord.Value.FieldByName(sel.Name).Interface().(Node[ast.Node])}, artValue
					}
				}(vvv, item.Sel)
			case *IfThenElseSingleValueCondition:
				return func(node Node[ast.Node], sel *ast.Ident) ExecuteStatement {
					return func(state State) ([]Node[ast.Node], CallArrayResultType) {

						if selector, ok := compiler.expandNodeWithSelector(node, sel); ok {
							return []Node[ast.Node]{selector}, artValue
						}
						panic("fsdfdsfd")
					}
				}(vv, item.Sel)
			case *TrailSource:
				var es ExecuteStatement = func(state State) ([]Node[ast.Node], CallArrayResultType) {
					result := ChangeParamNode[ast.Expr, ast.Node](node, &EntityField{node.Node.Pos(), vvv.Alias, item.Sel.Name})
					return []Node[ast.Node]{result}, artValue
				}
				return es
			default:
				panic("implement me")
			}
		case ExecuteStatement:
			return func(es ExecuteStatement, sel *ast.Ident) ExecuteStatement {
				return func(state State) ([]Node[ast.Node], CallArrayResultType) {
					arr, art := compiler.executeAndExpandStatement(state, es)
					if selector, ok := compiler.expandNodeWithSelector(arr[0], sel); ok {
						return []Node[ast.Node]{selector}, art
					}
					return arr, art
				}
			}(vv, item.Sel)
		default:
			panic("implement me")
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

			newContext := &CurrentContext{m, LocalTypesMap{}, GetCompilerState[*CurrentContext](state)}
			state = SetCompilerState(newContext, state)
			param := ChangeParamNode[ast.Node, *ast.BlockStmt](state.currentNode, node.Node.Body)
			values, art := compiler.executeBlockStmt(state, param)
			state = SetCompilerState(newContext.Parent, state)
			return values, art
		}
	}
}
