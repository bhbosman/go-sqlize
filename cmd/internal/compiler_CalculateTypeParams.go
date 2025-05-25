package internal

import (
	"go/ast"
	"reflect"
)

func (compiler *Compiler) CalculateTypeParams(
	state State,
	requiredTypeParams map[string]bool,
	funcTypeNode Node[*ast.FuncType],
	args []Node[ast.Node],
	nameAndParams findAllParamNameAndTypesResult,
) (map[string]ITypeMapper, bool) {
	sss := map[string]ITypeMapper{}
	for _, andParam := range nameAndParams.arr {
		sss, _ = compiler.internalCalculateTypeParams(
			state,
			requiredTypeParams,
			ChangeParamNode[*ast.FuncType, ast.Node](funcTypeNode, andParam.node.Node),
			sss,
			ChangeParamNode[*ast.FuncType, ast.Node](funcTypeNode, funcTypeNode.Node.Params),
			args,
		)
		if len(requiredTypeParams) == 0 {
			return sss, true
		}
	}
	return sss, true
}

// Todo: remove the return value
func (compiler *Compiler) internalCalculateTypeParams(
	state State,
	requiredTypeParams map[string]bool,
	funcDecl Node[ast.Node],
	s map[string]ITypeMapper,
	Params Node[ast.Node],
	args []Node[ast.Node],
) (map[string]ITypeMapper, bool) {
	switch funcDeclItem := funcDecl.Node.(type) {
	default:
		panic(funcDeclItem)
		panic("unreachable")
	case *ast.Ellipsis:
		param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Elt)
		return compiler.internalCalculateTypeParams(
			state,
			requiredTypeParams,
			param,
			s,
			Params,
			args,
		)
	case *ast.MapType:
		switch paramItem := Params.Node.(type) {
		default:
			panic("fff")
		case *ast.Ident:
			var b bool
			param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Key)
			if s, b = compiler.internalCalculateTypeParams(
				state,
				requiredTypeParams,
				param,
				s,
				ChangeParamNode[ast.Node, ast.Node](Params, paramItem),
				args,
			); !b {
				return s, len(requiredTypeParams) > 0
			}
			param = ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Value)
			return compiler.internalCalculateTypeParams(
				state,
				requiredTypeParams,
				param,
				s,
				ChangeParamNode[ast.Node, ast.Node](Params, paramItem),
				args,
			)
		case *ast.FieldList:
			for idx, field := range paramItem.List {
				if findTypeMapperForMap, ok := args[idx].Node.(IFindTypeMapper); ok {
					if mapper, b := findTypeMapperForMap.GetTypeMapper(state); b {
						defaultMapper := mapper[0]
						switch defaultMapper.Kind() {
						default:
							panic(mapper)
						case reflect.Map:
							typ, vk := defaultMapper.ActualType()
							keyRt := typ.Key()
							mapperKey := &WrapReflectTypeInMapper{keyRt, vk}
							mapperKeyNode := ChangeParamNode[ast.Node, ast.Node](args[idx], mapperKey)
							funcDeclParam := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Key)
							s, b = compiler.internalCalculateTypeParams(
								state,
								requiredTypeParams,
								funcDeclParam,
								s,
								ChangeParamNode[ast.Node, ast.Node](Params, field.Type),
								[]Node[ast.Node]{mapperKeyNode})
							if !b {
								return s, len(requiredTypeParams) > 0
							}

							valueRt := typ.Elem()
							mapperValue := &WrapReflectTypeInMapper{valueRt, vk}
							mapperValueNode := ChangeParamNode[ast.Node, ast.Node](args[idx], mapperValue)
							funcDeclParam = ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Value)
							s, b = compiler.internalCalculateTypeParams(
								state,
								requiredTypeParams,
								funcDeclParam,
								s,
								ChangeParamNode[ast.Node, ast.Node](Params, field.Type),
								[]Node[ast.Node]{mapperValueNode})
							if !b {
								return s, len(requiredTypeParams) > 0
							}
						}
					}
				} else {
					var b bool
					s, b = compiler.internalCalculateTypeParams(
						state,
						requiredTypeParams,
						funcDecl,
						s,
						ChangeParamNode[ast.Node, ast.Node](Params, field.Type),
						[]Node[ast.Node]{args[idx]},
					)
					if !b {
						return s, len(requiredTypeParams) > 0
					}
				}
			}
			return s, len(requiredTypeParams) > 0
		}
	case *ast.IndexListExpr:
		switch paramItem := Params.Node.(type) {
		default:
			panic(paramItem)
		case *ast.IndexListExpr:
			mappers, _ := args[0].Node.(IFindTypeMapper).GetTypeMapper(state)
			if len(funcDeclItem.Indices) == len(mappers) {
				for idx := 0; idx < len(mappers); idx++ {
					var b bool
					if s, b = compiler.internalCalculateTypeParams(
						state,
						requiredTypeParams,
						ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Indices[idx]),
						s,
						Node[ast.Node]{Node: mappers[idx], Valid: true},
						nil,
					); !b {
						return s, len(requiredTypeParams) > 0
					}
				}
				return s, len(requiredTypeParams) > 0
			}
			panic("counts mismatch")
		case *ast.FieldList:
			for idx, field := range paramItem.List {
				var b bool
				if s, b = compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					ChangeParamNode[ast.Node, ast.Node](Params, field.Type),
					[]Node[ast.Node]{args[idx]},
				); !b {
					return s, len(requiredTypeParams) > 0
				}
			}
		}
		for _, index := range funcDeclItem.Indices {
			var b bool
			if s, b = compiler.internalCalculateTypeParams(
				state,
				requiredTypeParams,
				ChangeParamNode[ast.Node, ast.Node](funcDecl, index),
				s,
				Params,
				args,
			); !b {
				return s, len(requiredTypeParams) > 0
			}
		}
		return s, len(requiredTypeParams) > 0
	case *ast.IndexExpr:
		switch paramItem := Params.Node.(type) {
		default:
			panic(paramItem)
		case *ast.IndexExpr:
			mappers, _ := args[0].Node.(IFindTypeMapper).GetTypeMapper(state)
			return compiler.internalCalculateTypeParams(
				state,
				requiredTypeParams,
				ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Index),
				s,
				Node[ast.Node]{Node: mappers[0], Valid: true},
				nil,
			)
		case *ast.FieldList:
			for idx, field := range paramItem.List {
				var b bool
				if s, b = compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					ChangeParamNode[ast.Node, ast.Node](Params, field.Type),
					[]Node[ast.Node]{args[idx]},
				); !b {
					return s, len(requiredTypeParams) > 0
				}
			}
		}
		indexParam := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Index)
		return compiler.internalCalculateTypeParams(
			state,
			requiredTypeParams,
			indexParam,
			s,
			Params,
			args,
		)
	case *ast.FuncType:
		switch paramItem := Params.Node.(type) {
		default:
			panic(paramItem)
		case *ast.Ident:
			return s, len(requiredTypeParams) > 0
		case *ast.FuncType:
			a := NodeStringValue(ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Results.List[0].Type))
			b := NodeStringValue(ChangeParamNode[ast.Node, ast.Node](Params, paramItem.Results.List[0].Type))
			if _, ok := requiredTypeParams[a]; ok && a == b {
				typeMapper := compiler.findType(state, args[0], Default|TypeParamType)
				switch typeMapper.Kind() {
				default:
					panic("sggdgdf")
				case reflect.Func:
					switch tm := typeMapper.(type) {
					case ITypeMapperForFuncType:
						if tm.OutCount() == 1 {
							outMapper00 := tm.Out(0)
							return compiler.internalCalculateTypeParams(
								state,
								requiredTypeParams,
								ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Results.List[0].Type),
								s,
								Node[ast.Node]{Node: outMapper00, Valid: true},
								nil,
							)
						}
						panic(typeMapper)
					}
				}
			} else if _, ok := args[0].Node.(*ast.Ident); ok {
				typeMapper := compiler.findType(state, args[0], Default)
				var b bool
				if s, b = compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					Node[ast.Node]{Node: typeMapper, Valid: true},
					nil,
				); !b {
					return s, len(requiredTypeParams) > 0
				}

			} else if funcLit, ok := args[0].Node.(*ast.FuncLit); ok {
				funcLitTypeParamsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(args[0], funcLit.Type.Params))
				funcDeclItemParamsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(funcDecl, funcDeclItem.Params))
				paramItemParamList := findAllParamNameAndTypes(ChangeParamNode(Params, paramItem.Params))
				if len(funcLitTypeParamsNameAndTypeExpressions.arr) == len(funcDeclItemParamsNameAndTypeExpressions.arr) && len(funcLitTypeParamsNameAndTypeExpressions.arr) == len(paramItemParamList.arr) {
					for idx := 0; idx < len(funcDeclItemParamsNameAndTypeExpressions.arr); idx++ {
						a := NodeStringValue(funcDeclItemParamsNameAndTypeExpressions.arr[idx].node)
						b := NodeStringValue(paramItemParamList.arr[idx].node)
						if _, ok := requiredTypeParams[a]; ok && a == b {
							typeMapper := compiler.findType(state, funcLitTypeParamsNameAndTypeExpressions.arr[idx].node, Default|TypeParamType)
							var b bool
							if s, b = compiler.internalCalculateTypeParams(
								state,
								requiredTypeParams,
								funcDeclItemParamsNameAndTypeExpressions.arr[idx].node,
								s,
								Node[ast.Node]{Node: typeMapper, Valid: true},
								nil,
							); !b {
								return s, len(requiredTypeParams) > 0
							}
						}
					}
				}
			} else if funcLit_, ok := args[0].Node.(FuncLit); ok {
				funcLitTypeParamsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(args[0], funcLit_.Type.Params))
				funcDeclItemParamsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(funcDecl, funcDeclItem.Params))
				paramItemParamList := findAllParamNameAndTypes(ChangeParamNode(Params, paramItem.Params))
				if len(funcLitTypeParamsNameAndTypeExpressions.arr) == len(funcDeclItemParamsNameAndTypeExpressions.arr) && len(funcLitTypeParamsNameAndTypeExpressions.arr) == len(paramItemParamList.arr) {
					for idx := 0; idx < len(funcDeclItemParamsNameAndTypeExpressions.arr); idx++ {
						a := NodeStringValue(funcDeclItemParamsNameAndTypeExpressions.arr[idx].node)
						b := NodeStringValue(paramItemParamList.arr[idx].node)
						if _, ok := requiredTypeParams[a]; ok && a == b {
							typeMapper := compiler.findType(state, funcLitTypeParamsNameAndTypeExpressions.arr[idx].node, Default|TypeParamType)
							var b bool
							if s, b = compiler.internalCalculateTypeParams(
								state,
								requiredTypeParams,
								funcDeclItemParamsNameAndTypeExpressions.arr[idx].node,
								s,
								Node[ast.Node]{Node: typeMapper, Valid: true},
								nil,
							); !b {
								return s, len(requiredTypeParams) > 0
							}
						}
					}
				}
			} else {
				findITypeMapperForFuncType := func() (ITypeMapperForFuncType, bool) {
					if CallExpr, ok := args[0].Node.(*ast.CallExpr); ok {
						p1 := ChangeParamNode[ast.Node, ast.Node](args[0], CallExpr.Fun)
						typeMapper := compiler.findType(state, p1, Default|TypeParamType)
						if typeMapperForFuncType, ok := typeMapper.(ITypeMapperForFuncType); ok {
							return typeMapperForFuncType, true
						}
					}
					if trailArray, ok := args[0].Node.(TrailArray); ok {
						if typeMapperForFuncType, ok := trailArray.typeMapper.(ITypeMapperForFuncType); ok {
							return typeMapperForFuncType, true
						}
					}
					return nil, false
				}

				if typeMapper, ok := findITypeMapperForFuncType(); ok {
					var b bool
					if s, b = compiler.internalCalculateTypeParams(
						state,
						requiredTypeParams,
						funcDecl,
						s,
						Node[ast.Node]{Node: typeMapper, Valid: true},
						nil,
					); !b {
						return s, len(requiredTypeParams) > 0
					}
				}
			}
			return s, len(requiredTypeParams) > 0
		case *ast.FieldList:
			for idx, field := range paramItem.List {
				var b bool
				s, b = compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					ChangeParamNode[ast.Node, ast.Node](Params, field.Type),
					[]Node[ast.Node]{args[idx]})
				if !b {
					return s, len(requiredTypeParams) > 0
				}
			}
			return s, len(requiredTypeParams) > 0
		case *ast.ArrayType:
			return compiler.internalCalculateTypeParams(
				state,
				requiredTypeParams,
				funcDecl,
				s,
				ChangeParamNode[ast.Node, ast.Node](Params, paramItem.Elt),
				args,
			)
		case *ast.Ellipsis:
			return compiler.internalCalculateTypeParams(
				state,
				requiredTypeParams,
				funcDecl,
				s,
				ChangeParamNode[ast.Node, ast.Node](Params, paramItem.Elt),
				args,
			)
		case ITypeMapperForFuncType:
			funcDeclItemParamsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(funcDecl, funcDeclItem.Params))
			for idx := 0; idx < paramItem.InCount(); idx++ {
				itemTypeMapper := paramItem.In(idx)
				var b bool
				if s, b = compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDeclItemParamsNameAndTypeExpressions.arr[idx].node,
					s,
					Node[ast.Node]{Node: itemTypeMapper, Valid: true},
					nil,
				); !b {
					return s, len(requiredTypeParams) > 0
				}
			}

			funcDeclItemResultsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(funcDecl, funcDeclItem.Results))
			for idx := 0; idx < paramItem.OutCount(); idx++ {
				itemTypeMapper := paramItem.Out(idx)
				var b bool
				if s, b = compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDeclItemResultsNameAndTypeExpressions.arr[idx].node,
					s,
					Node[ast.Node]{Node: itemTypeMapper, Valid: true},
					nil,
				); !b {
					return s, len(requiredTypeParams) > 0
				}
			}
			return s, len(requiredTypeParams) > 0
		case ITypeMapper:
			panic(paramItem)
		}
	case *ast.ArrayType:
		param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Elt)
		return compiler.internalCalculateTypeParams(
			state,
			requiredTypeParams,
			param,
			s,
			Params,
			args,
		)
	case *ast.Ident:
		if _, ok := requiredTypeParams[funcDeclItem.Name]; ok {
			switch paramItem := Params.Node.(type) {
			default:
				panic(paramItem)
			case *ast.MapType:
				var b bool
				if s, b = compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					ChangeParamNode[ast.Node, ast.Node](Params, paramItem.Key),
					args,
				); !b {
					return s, len(requiredTypeParams) > 0
				}
				return compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					ChangeParamNode[ast.Node, ast.Node](Params, paramItem.Value),
					args,
				)
			case *ast.IndexListExpr:
				for _, index := range paramItem.Indices {
					var b bool

					if s, b = compiler.internalCalculateTypeParams(
						state,
						requiredTypeParams,
						funcDecl,
						s,
						ChangeParamNode[ast.Node, ast.Node](Params, index),
						args,
					); !b {
						return s, len(requiredTypeParams) > 0
					}
				}
				return s, len(requiredTypeParams) > 0
			case *ast.IndexExpr:
				return compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					ChangeParamNode[ast.Node, ast.Node](Params, paramItem.Index),
					args,
				)
			case *ast.ArrayType:

				return compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					ChangeParamNode[ast.Node, ast.Node](Params, paramItem.Elt),
					args,
				)
			case ITypeMapper:
				if _, ok := requiredTypeParams[funcDeclItem.Name]; ok {
					delete(requiredTypeParams, funcDeclItem.Name)
					s[funcDeclItem.Name] = paramItem
				}
				return s, len(requiredTypeParams) > 0
			case *ast.Ident:
				if len(args) == 1 && paramItem.Name == funcDeclItem.Name {
					if findTypeMapper, ok := args[0].Node.(IFindTypeMapper); ok {
						if arr, ok := findTypeMapper.GetTypeMapper(state); ok {
							return compiler.internalCalculateTypeParams(
								state,
								requiredTypeParams,
								funcDecl,
								s,
								Node[ast.Node]{Node: arr[0], Valid: true},
								nil,
							)
						}
					} else {
						panic("implement me")
					}
				}
				return s, len(requiredTypeParams) > 0
			case *ast.FieldList:
				nameAndParams := findAllParamNameAndTypes(ChangeParamNode(Params, paramItem))
				for idx, arg := range args {

					nameAndParam := nameAndParams.arr[idx]
					var b bool
					s, b = compiler.internalCalculateTypeParams(
						state,
						requiredTypeParams,
						funcDecl,
						s,
						nameAndParam.node,
						[]Node[ast.Node]{arg})
					if !b {
						return s, false
					}
					if _, ok := requiredTypeParams[funcDeclItem.Name]; !ok {
						return s, len(requiredTypeParams) > 0
					}
				}
				return s, len(requiredTypeParams) > 0
			}
		}
		return s, len(requiredTypeParams) > 0
	}
}
