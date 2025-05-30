package internal

import (
	"go/ast"
	"reflect"
)

func (compiler *Compiler) createRhsCompositeLitExecution(node Node[*ast.CompositeLit]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		typeMapperFn := func(state State, parent Node[*ast.CompositeLit], Type ast.Expr) ITypeMapper {
			if Type != nil {
				param := ChangeParamNode[*ast.CompositeLit, ast.Node](parent, Type)
				return compiler.findType(state, param, Default)
			}
			currentCompositeCreateType := GetCompilerState[*CurrentCompositeCreateType](state)
			return currentCompositeCreateType.typeMapper
		}
		typeMapper := typeMapperFn(state, node, node.Node.Type)
		rtKind := typeMapper.Kind()
		switch rtKind {
		//case reflect.Func:
		//	var arr []Node[ast.Node]
		//	for _, elt := range node.Node.Elts {
		//		currentContext := GetCompilerState[*CurrentContext](state)
		//		flattenValues := currentContext.flattenVariables()
		//		switch etlItem := elt.(type) {
		//		default:
		//			panic("unreachable")
		//		case *ast.CallExpr:
		//			ce := CallExpression{ChangeParamNode[*ast.CompositeLit, *ast.CallExpr](node, etlItem), flattenValues}
		//			fn := func(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
		//				return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		//					onCreateExecuteStatement := compiler.createRhsCallExpressionExecution(ce.CallExpression)
		//					//executeStatement, _ := onCreateExecuteStatement(state, typeParams, arguments)
		//					return onCreateExecuteStatement(state, typeParams, arguments)
		//				}
		//			}
		//			fi := callExpressionInformation{fn, ChangeParamNode[*ast.CompositeLit, *ast.CallExpr](node, etlItem)}
		//			arr = append(arr, ChangeParamNode[*ast.CompositeLit, ast.Node](node, fi))
		//		case *ast.FuncLit:
		//			fl := FuncLit{etlItem.Type, etlItem.Body, flattenValues}
		//			fn := func(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
		//				return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		//					param := ChangeParamNode[*ast.CompositeLit, FuncLit](node, fl)
		//					onCreateExecuteStatement := compiler.onFuncLitExecutionStatement(param)
		//					executeStatement := onCreateExecuteStatement(state, funcTypeNode)
		//					return executeStatement(state, typeParams, unprocessedArgs)
		//				}
		//			}
		//			fi := functionInformation{fn, ChangeParamNode(node, etlItem.Type), true}
		//			arr = append(arr, ChangeParamNode[*ast.CompositeLit, ast.Node](node, fi))
		//		}
		//	}
		//	trailArray := TrailArray{arr, typeMapper}
		//	nodeValue := ChangeParamNode[*ast.CompositeLit, ast.Node](node, trailArray)
		//	return []Node[ast.Node]{nodeValue}, artValue
		case reflect.Struct:
			typeMapperForStruct := typeMapper.(*TypeMapperForStruct)
			if len(node.Node.Elts) == 0 {
				param := ChangeParamNode[*ast.CompositeLit, ast.Node](node, node.Node.Type)
				rv := typeMapperForStruct.createDefaultType(param)
				nodeValue := ChangeParamNode[*ast.CompositeLit, ast.Node](
					node,
					&TrailRecord{
						node.Node.Pos(),
						rv,
						typeMapperForStruct,
					},
				)
				return []Node[ast.Node]{nodeValue}, artValue
			}
			rt := typeMapperForStruct.nodeRt
			rv := reflect.New(rt).Elem()
			for idx, elt := range node.Node.Elts {
				switch expr := elt.(type) {
				case *ast.KeyValueExpr:
					param := ChangeParamNode[*ast.CompositeLit, ast.Node](node, expr.Value)
					tempState := state.setCurrentNode(ChangeParamNode[*ast.CompositeLit, ast.Node](node, expr.Value))
					es := compiler.findRhsExpression(tempState, param)
					vv, _ := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, es)
					switch key := expr.Key.(type) {
					case *ast.Ident:
						itemRv := reflect.ValueOf(vv[0])
						rv.FieldByName(key.Name).Set(itemRv)
					default:
						panic("unhandled key")
					}
				default:
					param := ChangeParamNode[*ast.CompositeLit, ast.Node](node, elt)
					tempState := state.setCurrentNode(ChangeParamNode[*ast.CompositeLit, ast.Node](node, elt))
					es := compiler.findRhsExpression(tempState, param)
					vv, _ := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, es)
					itemRv := reflect.ValueOf(vv[0])
					rv.Field(idx).Set(itemRv)
				}
			}
			nodeValue := ChangeParamNode[*ast.CompositeLit, ast.Node](node, &TrailRecord{node.Node.Pos(), rv, typeMapper})
			return []Node[ast.Node]{nodeValue}, artValue
		case reflect.Map:
			typeMapperForMap := typeMapper.(*TypeMapperForMap)
			rv := reflect.MakeMap(typeMapperForMap.mapRt)
			for _, elt := range node.Node.Elts {
				switch expr := elt.(type) {
				case *ast.KeyValueExpr:
					fn := func(state State, expression ast.Expr, typeMapper ITypeMapper) ([]Node[ast.Node], CallArrayResultType) {
						state = SetCompilerState[*CurrentCompositeCreateType](&CurrentCompositeCreateType{typeMapper}, state)
						param := ChangeParamNode[*ast.CompositeLit, ast.Node](node, expression)
						tempState := state.setCurrentNode(ChangeParamNode[*ast.CompositeLit, ast.Node](node, expression))
						es := compiler.findRhsExpression(tempState, param)
						return compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, es)
					}
					rt, _ := typeMapper.ActualType()
					nodeKey, _ := fn(state, expr.Key, typeMapperForMap.keyTypeMapper)
					rvKey := func(node Node[ast.Node]) reflect.Value {
						if translateNodeValueToReflectValue, ok := typeMapperForMap.keyTypeMapper.(ITranslateNodeValueToReflectValue); ok {
							return translateNodeValueToReflectValue.TranslateNodeValueToReflectValue(node)
						} else if rv, ok := isLiterateValue(node); ok {
							return rv
						}
						// todo: not 100% sold about this,
						panic("unhandled key")
					}(nodeKey[0])
					nodeValue, _ := fn(state, expr.Value, typeMapperForMap.valueTypeMapper)
					rv.SetMapIndex(rvKey.Convert(rt.Key()), reflect.ValueOf(nodeValue[0]))

				default:
					panic("unhandled key")
				}
			}
			nodeValue := ChangeParamNode[*ast.CompositeLit, ast.Node](node, &ReflectValueExpression{rv, ValueKey{"PQR", "PQR"}})
			return []Node[ast.Node]{nodeValue}, artValue
		default:
			panic("dsfsfds")
		}
	}
}
