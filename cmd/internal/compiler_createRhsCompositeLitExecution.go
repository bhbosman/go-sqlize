package internal

import (
	"go/ast"
	"reflect"
)

func (compiler *Compiler) createRhsCompositeLitExecution(node Node[*ast.CompositeLit]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		typeMapperFn := func(state State, parent Node[*ast.CompositeLit], Type ast.Expr) ITypeMapper {
			if Type != nil {
				param := ChangeParamNode(parent, Type)
				return compiler.findType(state, param)
			}
			currentCompositeCreateType := GetCompilerState[*CurrentCompositeCreateType](state)
			return &ReflectTypeHolder{
				func(state State, rv reflect.Value) reflect.Value {
					panic("ddddd")
				},
				func(state State) reflect.Type {
					return currentCompositeCreateType.rt
				},
			}
		}
		typeMapper := typeMapperFn(state, node, node.Node.Type)
		rt := typeMapper.Type(state)
		rtKind := rt.Kind()
		switch rtKind {
		case reflect.Struct:

			if len(node.Node.Elts) == 0 {
				param := ChangeParamNode[*ast.CompositeLit, ast.Node](node, node.Node.Type)
				typeMapperForStruct := typeMapper.(*TypeMapperForStruct)

				rv := typeMapperForStruct.createDefaultType(param)
				nodeValue := ChangeParamNode[*ast.CompositeLit, ast.Node](
					node,
					&TrailRecord{
						node.Node.Pos(),
						rv,
					},
				)
				return []Node[ast.Node]{nodeValue}, artValue
			}
			rv := reflect.New(rt).Elem()
			for idx, elt := range node.Node.Elts {
				switch expr := elt.(type) {
				case *ast.KeyValueExpr:
					param := ChangeParamNode(node, expr.Value)
					tempState := state.setCurrentNode(ChangeParamNode[*ast.CompositeLit, ast.Node](node, expr.Value))
					es := compiler.findRhsExpression(tempState, param)
					vv, _ := compiler.executeAndExpandStatement(tempState, es)
					switch key := expr.Key.(type) {
					case *ast.Ident:
						rv.FieldByName(key.Name).Set(reflect.ValueOf(vv[0]))
					default:
						panic("unhandled key")
					}
				default:
					param := ChangeParamNode(node, elt)
					tempState := state.setCurrentNode(ChangeParamNode[*ast.CompositeLit, ast.Node](node, elt))
					es := compiler.findRhsExpression(tempState, param)
					vv, _ := compiler.executeAndExpandStatement(tempState, es)
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
		case reflect.Map:
			typeMapperForMap := typeMapper.(*TypeMapperForMap)
			rv := reflect.MakeMap(rt)
			for _, elt := range node.Node.Elts {
				switch expr := elt.(type) {
				case *ast.KeyValueExpr:
					fn := func(state State, expression ast.Expr, typeMapper ITypeMapper, rt reflect.Type) ([]Node[ast.Node], CallArrayResultType) {
						state = SetCompilerState[*CurrentCompositeCreateType](&CurrentCompositeCreateType{rt, typeMapper}, state)
						param := ChangeParamNode(node, expression)
						tempState := state.setCurrentNode(ChangeParamNode[*ast.CompositeLit, ast.Node](node, expression))
						es := compiler.findRhsExpression(tempState, param)
						return compiler.executeAndExpandStatement(tempState, es)
					}

					nodeKey, _ := fn(state, expr.Key, typeMapperForMap.keyTypeMapper, rt.Key())
					if rvKey, okKey := isLiterateValue(nodeKey[0]); okKey {
						rvKey = typeMapperForMap.keyTypeMapper.Create(state, rvKey)
						nodeValue, _ := fn(state, expr.Value, typeMapperForMap.valueTypeMapper, rt.Elem())
						if rvValue, okValue := isLiterateValue(nodeValue[0]); okValue {
							rvValue = typeMapperForMap.valueTypeMapper.Create(state, rvValue)
							rv.SetMapIndex(rvKey.Convert(rt.Key()), rvValue.Convert(rt.Elem()))
							continue
						}
					}

					panic("must be literal values")
				default:
					panic("unhandled key")
				}
			}

			nodeValue := ChangeParamNode[*ast.CompositeLit, ast.Node](node, &ReflectValueExpression{rv})
			return []Node[ast.Node]{nodeValue}, artValue

		default:
			panic("dsfsfds")
		}
	}
}
