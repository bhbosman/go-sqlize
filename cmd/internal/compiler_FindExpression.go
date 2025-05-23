package internal

import (
	"fmt"
	"go/ast"
)

func (compiler *Compiler) findRhsExpression(state State, node Node[ast.Node]) ExecuteStatement {
	return compiler.internalFindRhsExpression(0, state, node).(ExecuteStatement)
}

func (compiler *Compiler) internalFindRhsExpression(stackIndex int, state State, node Node[ast.Node]) interface{} {
	switch item := node.Node.(type) {
	default:
		panic(node.Node)
	case IIsLiterateValue:
		var es func(node Node[ast.Node]) ExecuteStatement = func(node Node[ast.Node]) ExecuteStatement {
			return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
				return []Node[ast.Node]{node}, artValue
			}
		}
		return es(node)
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
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		unk := compiler.internalFindRhsExpression(stackIndex+1, state, param)
		switch vv := unk.(type) {
		case ImportMapEntry:
			vk := ValueKey{vv.Path, item.Sel.Name}
			if globalFunction, ok := compiler.GlobalFunctions[vk]; ok {
				return globalFunction.fn(state, globalFunction.funcType)
			}
			panic(notFound(fmt.Sprintf("%v", vk), "internalFindRhsExpression"))
		case Node[ast.Node]:
			switch vvv := vv.Node.(type) {
			case *TrailRecord:
				return func(trailRecord *TrailRecord, sel *ast.Ident) ExecuteStatement {
					return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
						return []Node[ast.Node]{trailRecord.Value.FieldByName(sel.Name).Interface().(Node[ast.Node])}, artValue
					}
				}(vvv, item.Sel)
			case IfThenElseSingleValueCondition:
				return func(node Node[ast.Node], sel *ast.Ident) ExecuteStatement {
					return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {

						if selector, ok := compiler.expandNodeWithSelector(node, sel); ok {
							return []Node[ast.Node]{selector}, artValue
						}
						panic("fsdfdsfd")
					}
				}(vv, item.Sel)
			case TrailSource:
				var es ExecuteStatement = func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
					typeMapperForStruct := vvv.typeMapper.(*TypeMapperForStruct)
					typeMapper := typeMapperForStruct.typeMapperInstance.FieldByName(item.Sel.Name).Interface().(ITypeMapper)
					result := ChangeParamNode[ast.Node, ast.Node](
						node,
						EntityField{
							vvv.Alias,
							typeMapper,
							item.Sel.Name,
						},
					)
					return []Node[ast.Node]{result}, artValue
				}
				return es
			default:
				panic("implement me")
			}
		case ExecuteStatement:
			return func(es ExecuteStatement, sel *ast.Ident) ExecuteStatement {
				return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
					arr, art := compiler.executeAndExpandStatement(state, typeParams, arguments, es)
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
		return compiler.createRhsFuncLitExprExecution(state, param)
	case *ast.ParenExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		return compiler.findRhsExpression(state, param)
	case *ast.Ident:
		currentContext := GetCompilerState[*CurrentContext](state)
		if value, b := currentContext.FindValueByString(item.Name); b {
			if stackIndex == 0 {
				var es ExecuteStatement = func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
					switch nodeItem := value.Node.(type) {
					case VaradicArgument:
						var result []Node[ast.Node]
						for _, data := range nodeItem.data {
							nn, _ := compiler.findRhsExpression(state, data)(state, nil, nil)
							result = append(result, nn...)

						}
						return result, artValue
					default:
						return []Node[ast.Node]{value}, artValue
					}
				}
				return es
			}
			return value
		}
		if globalFunction, ok := compiler.GlobalFunctions[ValueKey{node.RelPath, item.Name}]; ok {
			return globalFunction.fn(state, globalFunction.funcType)
		}
		if globalFunction, ok := compiler.GlobalFunctions[ValueKey{"", item.Name}]; ok {
			return globalFunction.fn(state, globalFunction.funcType)
		}
		if path, ok := node.ImportMap[item.Name]; ok {
			return path
		}
		panic("unhandled default case")
	}
}

func (compiler *Compiler) onFuncLitExecutionStatement(node Node[FuncLit]) OnCreateExecuteStatement {
	return func(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
		return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
			fmt.Printf("begin %v\n", compiler.Fileset.Position(node.Node.Pos()).String())
			paramNames := findAllParamNameAndTypes(ChangeParamNode[FuncLit, *ast.FieldList](node, node.Node.Type.Params))
			m := ValueInformationMap{}
			if paramNames.isVariadic {
				for idx := 0; idx < len(paramNames.arr)-1; idx++ {
					m[paramNames.arr[idx].name] = ValueInformation{arguments[idx]}
				}
				var varadicArr []Node[ast.Node]
				for idx := len(paramNames.arr) - 1; idx < len(arguments); idx++ {
					varadicArr = append(varadicArr, arguments[idx])
				}
				v := Node[ast.Node]{Node: VaradicArgument{varadicArr}, Valid: true}
				m[paramNames.arr[len(paramNames.arr)-1].name] = ValueInformation{v}

			} else {
				for idx, name := range paramNames.arr {
					m[name.name] = ValueInformation{arguments[idx]}
				}
			}

			newContext := &CurrentContext{m, map[string]ITypeMapper{}, LocalTypesMap{}, false, GetCompilerState[*CurrentContext](state)}
			state = SetCompilerState(newContext, state)
			param := ChangeParamNode[ast.Node, *ast.BlockStmt](state.currentNode, node.Node.Body)
			values, art := compiler.executeBlockStmt(state, param)
			state = SetCompilerState(newContext.Parent, state)
			fmt.Printf("end %v\n", compiler.Fileset.Position(node.Node.End()).String())
			return values, art
		}
	}
}
