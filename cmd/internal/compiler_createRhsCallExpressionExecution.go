package internal

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

func (compiler *Compiler) compileArguments(state State, argss []Node[ast.Node], typeParams map[string]ITypeMapper) []Node[ast.Node] {
	var result []Node[ast.Node]
	for _, arg := range argss {
		tempState := state.setCurrentNode(ChangeParamNode[ast.Node, ast.Node](arg, arg.Node))
		param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, arg.Node)
		fn := compiler.findRhsExpression(tempState, param)
		nodeArg, _ := compiler.executeAndExpandStatement(state, typeParams, argss, fn)
		result = append(result, nodeArg...)
	}
	return result
}

func (compiler *Compiler) createRhsCallExpressionExecution(node Node[*ast.CallExpr]) ExecuteStatement {
	return func(state State, _ map[string]ITypeMapper, _ []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		fmt.Printf("begin %v\n", compiler.Fileset.Position(node.Node.Lparen).String())
		newContext := &CurrentContext{
			ValueInformationMap{},
			map[string]ITypeMapper{},
			LocalTypesMap{},
			false,
			GetCompilerState[*CurrentContext](state),
		}
		state = SetCompilerState(newContext, state)

		knownTypeParams := newContext.flattenTypeParams()
		fmt.Printf("\t knownTypeParams:\n")
		for key, value := range knownTypeParams {
			typ, _ := value.ActualType()
			fmt.Printf("\t\t %s -> %s\n", key, typ.String())
		}
		param := ChangeParamNode[*ast.CallExpr, ast.Node](node, node.Node.Fun)
		tempState02 := state.setCurrentNode(ChangeParamNode[ast.Node, ast.Node](state.currentNode, node.Node.Fun))
		execFn, funcTypeNode := compiler.findFunction(tempState02, param)
		var args []Node[ast.Node]
		for _, arg := range node.Node.Args {
			paramArg := ChangeParamNode[*ast.CallExpr, ast.Node](node, arg)
			args = append(args, paramArg)
		}
		args = compiler.compileArguments(state, args, knownTypeParams)
		if !funcTypeNode.Valid {
			fn, resultType := execFn(tempState02, knownTypeParams, args)
			fmt.Printf("end %v\n", compiler.Fileset.Position(node.Node.Rparen).String())
			return fn, resultType
		} else {
			typeParamsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(funcTypeNode, funcTypeNode.Node.TypeParams))
			requiredTypeParams := map[string]bool{}
			for _, ss := range typeParamsNameAndTypeExpressions.arr {
				if _, ok := knownTypeParams[ss.name]; !ok {
					requiredTypeParams[ss.name] = true
				}
			}
			if len(requiredTypeParams) > 0 {
				var argumentArr []CalculateTypeArgumentType
				for _, arg := range args {
					argumentArr = append(argumentArr, CalculateTypeArgumentType{
						//idx,
						//ChangeParamNode[*ast.CallExpr, ast.Node](node, node.Node.Args[idx]),
						//nameAndParams[idx].node,
						arg,
					})
				}

				if mappers, b := compiler.calculateTypeMapperForCallExpression(
					state,
					requiredTypeParams,
					node,
					typeParamsNameAndTypeExpressions,
					funcTypeNode,
					argumentArr,
				); b {
					if len(mappers) >= len(typeParamsNameAndTypeExpressions.arr) {
						for key, value := range mappers {
							newContext.TypeParams[key] = value
							knownTypeParams[key] = value
						}
						fmt.Printf("\t knownTypeParams after calculation:\n")
						for key, value := range knownTypeParams {
							typ, _ := value.ActualType()
							fmt.Printf("\t\t %s -> %s\n", key, typ.String())
						}
					} else {
						//createTypeMapperFn(node)
						panic("sdfds")
					}
				} else {
					panic("unreachable")
				}
			}

			deltaTypeParams := map[string]ITypeMapper{}
			for _, typeParam := range typeParamsNameAndTypeExpressions.arr {
				deltaTypeParams[typeParam.name] = knownTypeParams[typeParam.name]
			}

			fn, resultType := execFn(tempState02, deltaTypeParams, args)
			fmt.Printf("end %v\n", compiler.Fileset.Position(node.Node.Rparen).String())
			return fn, resultType
		}
	}
}

func (compiler *Compiler) calculateTypeMapperForCallExpression(
	state State,
	requiredTypeParams map[string]bool,
	node Node[*ast.CallExpr],
	nameAndTypeParams findAllParamNameAndTypesResult,
	funcTypeNode Node[*ast.FuncType],
	args []CalculateTypeArgumentType,
) (map[string]ITypeMapper, bool) {
	if len(nameAndTypeParams.arr) == 0 {
		return map[string]ITypeMapper{}, true
	}

	switch nodeItem := node.Node.Fun.(type) {
	case *ast.IndexExpr:
		paramType := ChangeParamNode[*ast.CallExpr, ast.Node](node, nodeItem.Index)
		typeMapper := compiler.findType(state, paramType, TypeParamType|Default)
		return map[string]ITypeMapper{
			nameAndTypeParams.arr[0].name: typeMapper,
		}, true
	case *ast.IndexListExpr:
		results := map[string]ITypeMapper{}
		for idx, index := range nodeItem.Indices {
			paramType := ChangeParamNode[*ast.CallExpr, ast.Node](node, index)
			typeMapper := compiler.findType(state, paramType, TypeParamType|Default)
			results[nameAndTypeParams.arr[idx].name] = typeMapper
		}
		return results, true
	default:
		if !funcTypeNode.Valid {
			return nil, false
		}
		paramsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(funcTypeNode, funcTypeNode.Node.Params))
		if func(paramsNameAndTypeExpressions findAllParamNameAndTypesResult) bool {
			if paramsNameAndTypeExpressions.isVariadic {
				return len(paramsNameAndTypeExpressions.arr) > 0 && len(paramsNameAndTypeExpressions.arr)-1 <= len(node.Node.Args)
			} else {
				return len(paramsNameAndTypeExpressions.arr) > 0 && len(paramsNameAndTypeExpressions.arr) == len(node.Node.Args)
			}
		}(paramsNameAndTypeExpressions) {
			return compiler.CalculateTypeParams(state, requiredTypeParams, funcTypeNode, args, paramsNameAndTypeExpressions)
		}
		return nil, true
	}
}

type CalculateTypeFuncDeclType struct {
	node Node[ast.Node]
}

type CalculateTypeParamType struct {
	node Node[ast.Node]
}

type CalculateTypeArgumentType struct {
	//index int
	//inputArgumentNode Node[ast.Node]
	//paramStruct      Node[ast.Node]
	compiledArgument Node[ast.Node]
}

func NodeStringValue(node Node[ast.Node]) string {
	return internalNodeStringValue(node)
}

func internalNodeStringValue(node Node[ast.Node]) string {
	switch expr := node.Node.(type) {
	default:
		panic("unreachable")
	case *ast.Ident:
		return expr.Name
	case *ast.SelectorExpr:
		x := internalNodeStringValue(ChangeParamNode[ast.Node, ast.Node](node, expr.X))
		return fmt.Sprintf("%v.%v", x, expr.Sel.Name)
	case *ast.IndexExpr:
		x := internalNodeStringValue(ChangeParamNode[ast.Node, ast.Node](node, expr.X))
		index := internalNodeStringValue(ChangeParamNode[ast.Node, ast.Node](node, expr.Index))
		return fmt.Sprintf("%v[%v]", x, index)
	case *ast.IndexListExpr:
		x := internalNodeStringValue(ChangeParamNode[ast.Node, ast.Node](node, expr.X))
		var ss []string
		for _, index := range expr.Indices {
			ss = append(ss, internalNodeStringValue(ChangeParamNode[ast.Node, ast.Node](node, index)))
		}
		return fmt.Sprintf("%v[%v]", x, strings.Join(ss, ","))
	}
}

func (compiler *Compiler) CalculateTypeParams(
	state State,
	requiredTypeParams map[string]bool,
	funcTypeNode Node[*ast.FuncType],
	args []CalculateTypeArgumentType,
	nameAndParams findAllParamNameAndTypesResult,
) (map[string]ITypeMapper, bool) {
	sss := map[string]ITypeMapper{}
	for _, andParam := range nameAndParams.arr {
		sss, _ = compiler.internalCalculateTypeParams(
			state,
			requiredTypeParams,
			CalculateTypeFuncDeclType{ChangeParamNode[*ast.FuncType, ast.Node](funcTypeNode, andParam.node.Node)},
			sss,
			CalculateTypeParamType{ChangeParamNode[*ast.FuncType, ast.Node](funcTypeNode, funcTypeNode.Node.Params)},
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
	funcDecl CalculateTypeFuncDeclType,
	s map[string]ITypeMapper,
	Params CalculateTypeParamType,
	args []CalculateTypeArgumentType,
) (map[string]ITypeMapper, bool) {
	switch funcDeclItem := funcDecl.node.Node.(type) {
	default:
		panic(funcDeclItem)
		panic("unreachable")
	case *ast.MapType:
		switch paramItem := Params.node.Node.(type) {
		default:
			panic("fff")
		case *ast.Ident:
			var b bool
			param := ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Key)
			if s, b = compiler.internalCalculateTypeParams(
				state,
				requiredTypeParams,
				CalculateTypeFuncDeclType{param},
				s,
				CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem)},
				args,
			); !b {
				return s, len(requiredTypeParams) > 0
			}
			param = ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Value)
			return compiler.internalCalculateTypeParams(
				state,
				requiredTypeParams,
				CalculateTypeFuncDeclType{param},
				s,
				CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem)},
				args,
			)
		case *ast.FieldList:
			for idx, field := range paramItem.List {
				if findTypeMapperForMap, ok := args[idx].compiledArgument.Node.(IFindTypeMapper); ok {
					if mapper, b := findTypeMapperForMap.GetTypeMapper(state); b {
						defaultMapper := mapper[0]
						switch defaultMapper.Kind() {
						default:
							panic(mapper)
						case reflect.Map:
							typ, vk := defaultMapper.ActualType()
							keyRt := typ.Key()
							mapperKey := &WrapReflectTypeInMapper{keyRt, vk}
							mapperKeyNode := ChangeParamNode[ast.Node, ast.Node](args[idx].compiledArgument, mapperKey)
							funcDeclParam := ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Key)
							s, b = compiler.internalCalculateTypeParams(
								state,
								requiredTypeParams,
								CalculateTypeFuncDeclType{funcDeclParam},
								s,
								CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
								[]CalculateTypeArgumentType{
									{
										//args[idx].index,
										//args[idx].inputArgumentNode,
										//args[idx].paramStruct,
										mapperKeyNode,
									},
								})
							if !b {
								return s, len(requiredTypeParams) > 0
							}

							valueRt := typ.Elem()
							mapperValue := &WrapReflectTypeInMapper{valueRt, vk}
							mapperValueNode := ChangeParamNode[ast.Node, ast.Node](args[idx].compiledArgument, mapperValue)
							funcDeclParam = ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Value)
							s, b = compiler.internalCalculateTypeParams(
								state,
								requiredTypeParams,
								CalculateTypeFuncDeclType{funcDeclParam},
								s,
								CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
								[]CalculateTypeArgumentType{
									{
										//args[idx].index,
										//args[idx].inputArgumentNode,
										//args[idx].paramStruct,
										mapperValueNode,
									},
								})
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
						CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
						[]CalculateTypeArgumentType{args[idx]},
					)
					if !b {
						return s, len(requiredTypeParams) > 0
					}
				}
			}
			return s, len(requiredTypeParams) > 0
		}
	case *ast.IndexListExpr:
		switch paramItem := Params.node.Node.(type) {
		default:
			panic(paramItem)
		case *ast.IndexListExpr:
			mappers, _ := args[0].compiledArgument.Node.(IFindTypeMapper).GetTypeMapper(state)
			if len(funcDeclItem.Indices) == len(mappers) {
				for idx := 0; idx < len(mappers); idx++ {
					var b bool
					if s, b = compiler.internalCalculateTypeParams(
						state,
						requiredTypeParams,
						CalculateTypeFuncDeclType{ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Indices[idx])},
						s,
						CalculateTypeParamType{Node[ast.Node]{Node: mappers[idx], Valid: true}},
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
					CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
					[]CalculateTypeArgumentType{args[idx]},
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
				CalculateTypeFuncDeclType{ChangeParamNode[ast.Node, ast.Node](funcDecl.node, index)},
				s,
				Params,
				args,
			); !b {
				return s, len(requiredTypeParams) > 0
			}
		}
		return s, len(requiredTypeParams) > 0
	case *ast.IndexExpr:
		switch paramItem := Params.node.Node.(type) {
		default:
			panic(paramItem)
		case *ast.IndexExpr:
			mappers, _ := args[0].compiledArgument.Node.(IFindTypeMapper).GetTypeMapper(state)
			return compiler.internalCalculateTypeParams(
				state,
				requiredTypeParams,
				CalculateTypeFuncDeclType{ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Index)},
				s,
				CalculateTypeParamType{Node[ast.Node]{Node: mappers[0], Valid: true}},
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
					CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
					[]CalculateTypeArgumentType{args[idx]},
				); !b {
					return s, len(requiredTypeParams) > 0
				}
			}
		}
		indexParam := ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Index)
		return compiler.internalCalculateTypeParams(
			state,
			requiredTypeParams,
			CalculateTypeFuncDeclType{indexParam},
			s,
			Params,
			args,
		)
	case *ast.FuncType:
		switch paramItem := Params.node.Node.(type) {
		default:
			panic(paramItem)
		case *ast.Ident:
			return s, len(requiredTypeParams) > 0
		case *ast.FuncType:
			a := NodeStringValue(ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Results.List[0].Type))
			b := NodeStringValue(ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Results.List[0].Type))
			if _, ok := requiredTypeParams[a]; ok && a == b {
				typeMapper := compiler.findType(state, args[0].compiledArgument, Default|TypeParamType)
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
								CalculateTypeFuncDeclType{ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Results.List[0].Type)},
								s,
								CalculateTypeParamType{Node[ast.Node]{Node: outMapper00, Valid: true}},
								nil,
							)
						}
						panic(typeMapper)
					}
				}
			} else if funcLit, ok := args[0].compiledArgument.Node.(FuncLit); ok {
				funcLitTypeParamsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(args[0].compiledArgument, funcLit.Type.Params))
				funcDeclItemParamsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(funcDecl.node, funcDeclItem.Params))
				paramItemParamList := findAllParamNameAndTypes(ChangeParamNode(Params.node, paramItem.Params))
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
								CalculateTypeFuncDeclType{funcDeclItemParamsNameAndTypeExpressions.arr[idx].node},
								s,
								CalculateTypeParamType{Node[ast.Node]{Node: typeMapper, Valid: true}},
								nil,
							); !b {
								return s, len(requiredTypeParams) > 0
							}
						}
					}
				}
			} else if trailArray, ok := args[0].compiledArgument.Node.(TrailArray); ok {
				switch typeMapper := trailArray.typeMapper.(type) {
				default:
					panic("sggdgdf")
				case ITypeMapperForFuncType:
					if funcDeclItem.Params.NumFields() == typeMapper.InCount() && funcDeclItem.Results.NumFields() == typeMapper.OutCount() {
						{
							funcDeclItemParamsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(funcDecl.node, funcDeclItem.Params))
							for idx := 0; idx < typeMapper.InCount(); idx++ {
								itemTypeMapper := typeMapper.In(idx)
								var b bool
								if s, b = compiler.internalCalculateTypeParams(
									state,
									requiredTypeParams,
									CalculateTypeFuncDeclType{funcDeclItemParamsNameAndTypeExpressions.arr[idx].node},
									s,
									CalculateTypeParamType{Node[ast.Node]{Node: itemTypeMapper, Valid: true}},
									nil,
								); !b {
									return s, len(requiredTypeParams) > 0
								}
							}
						}
						{
							funcDeclItemResultsNameAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(funcDecl.node, funcDeclItem.Results))
							for idx := 0; idx < typeMapper.OutCount(); idx++ {
								itemTypeMapper := typeMapper.Out(idx)
								var b bool
								if s, b = compiler.internalCalculateTypeParams(
									state,
									requiredTypeParams,
									CalculateTypeFuncDeclType{funcDeclItemResultsNameAndTypeExpressions.arr[idx].node},
									s,
									CalculateTypeParamType{Node[ast.Node]{Node: itemTypeMapper, Valid: true}},
									nil,
								); !b {
									return s, len(requiredTypeParams) > 0
								}
							}
						}
					}
				}
			} else {
				panic("unreachable")
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
					CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
					[]CalculateTypeArgumentType{args[idx]})
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
				CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Elt)},
				args,
			)
		}
	case *ast.ArrayType:
		param := ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Elt)
		return compiler.internalCalculateTypeParams(
			state,
			requiredTypeParams,
			CalculateTypeFuncDeclType{param},
			s,
			Params,
			args,
		)
	case *ast.Ident:
		if _, ok := requiredTypeParams[funcDeclItem.Name]; ok {
			switch paramItem := Params.node.Node.(type) {
			default:
				panic(paramItem)
			case *ast.MapType:
				var b bool
				if s, b = compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Key)},
					args,
				); !b {
					return s, len(requiredTypeParams) > 0
				}
				return compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Value)},
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
						CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, index)},
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
					CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Index)},
					args,
				)
			case *ast.ArrayType:

				return compiler.internalCalculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					CalculateTypeParamType{ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Elt)},
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
					if findTypeMapper, ok := args[0].compiledArgument.Node.(IFindTypeMapper); ok {
						if arr, ok := findTypeMapper.GetTypeMapper(state); ok {
							return compiler.internalCalculateTypeParams(
								state,
								requiredTypeParams,
								funcDecl,
								s,
								CalculateTypeParamType{Node[ast.Node]{Node: arr[0], Valid: true}},
								nil,
							)
						}
					} else {
						panic("implement me")
					}
				}
				return s, len(requiredTypeParams) > 0
			case *ast.FieldList:
				nameAndParams := findAllParamNameAndTypes(ChangeParamNode(Params.node, paramItem))
				for idx, arg := range args {

					nameAndParam := nameAndParams.arr[idx]
					var b bool
					s, b = compiler.internalCalculateTypeParams(
						state,
						requiredTypeParams,
						funcDecl,
						s,
						CalculateTypeParamType{nameAndParam.node},
						[]CalculateTypeArgumentType{arg})
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
