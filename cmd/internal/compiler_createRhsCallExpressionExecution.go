package internal

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

func buildMap(ss []struct {
	name string
	node ast.Expr
}) map[string]ITypeMapper {
	m := map[string]ITypeMapper{}
	for _, vv := range ss {
		m[vv.name] = nil
	}
	return m
}

func (compiler *Compiler) createRhsCallExpressionExecution(node Node[*ast.CallExpr]) ExecuteStatement {
	compileArguments := func(
		state State,
		argss []Node[ast.Node],
		typeParams map[string]ITypeMapper,
	) []Node[ast.Node] {
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

	createTypeMapperFn := func(
		state State,
		requiredTypeParams map[string]bool,
		node Node[*ast.CallExpr],
		nameAndTypeParams []struct {
			name string
			node Node[ast.Node]
		},
		funcTypeNode Node[*ast.FuncType],
		args []CalculateTypeArgumentType,
	) (map[string]ITypeMapper, bool) {
		if len(nameAndTypeParams) == 0 {
			return map[string]ITypeMapper{}, true
		}

		if !funcTypeNode.Valid {
			return nil, false
		}
		switch nodeItem := node.Node.Fun.(type) {
		case *ast.IndexExpr:
			paramType := ChangeParamNode[*ast.CallExpr, ast.Node](node, nodeItem.Index)
			typeMapper := compiler.findType(state, paramType, TypeParamType|Default)
			return map[string]ITypeMapper{
				nameAndTypeParams[0].name: typeMapper,
			}, true
		case *ast.IndexListExpr:
			results := map[string]ITypeMapper{}
			for idx, index := range nodeItem.Indices {
				paramType := ChangeParamNode[*ast.CallExpr, ast.Node](node, index)
				typeMapper := compiler.findType(state, paramType, TypeParamType|Default)
				results[nameAndTypeParams[idx].name] = typeMapper
			}
			return results, true
		default:
			nameAndParams := findAllParamNameAndTypes(ChangeParamNode(funcTypeNode, funcTypeNode.Node.Params))
			if len(nameAndParams) > 0 && len(nameAndParams) == len(node.Node.Args) {

				sss := map[string]ITypeMapper{}
				for _, andParam := range nameAndParams {
					sss, _ = compiler.calculateTypeParams(
						state,
						requiredTypeParams,
						CalculateTypeFuncDeclType{0, ChangeParamNode[*ast.FuncType, ast.Node](funcTypeNode, andParam.node.Node)},
						sss,
						CalculateTypeParamType{0, ChangeParamNode[*ast.FuncType, ast.Node](funcTypeNode, funcTypeNode.Node.Params)},
						args,
					)
					if len(requiredTypeParams) == 0 {
						return sss, true
					}
				}
				return sss, true
			}
		}
		return nil, true
	}

	return func(state State, _ map[string]ITypeMapper, _ []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {

		fmt.Printf("begin %v\n", compiler.Fileset.Position(node.Node.Lparen).String())
		newContext := &CurrentContext{ValueInformationMap{}, map[string]ITypeMapper{}, LocalTypesMap{}, GetCompilerState[*CurrentContext](state)}
		state = SetCompilerState(newContext, state)

		knownTypeParams := newContext.flattenTypeParams()
		fmt.Printf("\t knownTypeParams:")
		for key, value := range knownTypeParams {
			fmt.Printf("\t\t %s -> %s\n", key, value.ActualType().String())
		}
		param := ChangeParamNode(node, node.Node.Fun)
		tempState02 := state.setCurrentNode(ChangeParamNode[ast.Node, ast.Node](state.currentNode, node.Node.Fun))
		execFn, funcTypeNode := compiler.findFunction(tempState02, param)
		var args []Node[ast.Node]
		for _, arg := range node.Node.Args {
			paramArg := ChangeParamNode[*ast.CallExpr, ast.Node](node, arg)
			args = append(args, paramArg)
		}
		args = compileArguments(state, args, knownTypeParams)
		if !funcTypeNode.Valid {
			fn, resultType := execFn(tempState02, knownTypeParams, args)
			fmt.Printf("end %v\n", compiler.Fileset.Position(node.Node.Rparen).String())
			return fn, resultType
		} else {
			nameAndTypeParams := findAllParamNameAndTypes(ChangeParamNode(funcTypeNode, funcTypeNode.Node.TypeParams))
			nameAndParams := findAllParamNameAndTypes(ChangeParamNode(funcTypeNode, funcTypeNode.Node.Params))

			requiredTypeParams := map[string]bool{}
			for _, ss := range nameAndTypeParams {
				if _, ok := knownTypeParams[ss.name]; !ok {
					requiredTypeParams[ss.name] = true
				}
			}

			var argumentArr []CalculateTypeArgumentType
			for idx, arg := range args {
				argumentArr = append(argumentArr, CalculateTypeArgumentType{idx, ChangeParamNode[*ast.CallExpr, ast.Node](node, node.Node.Args[idx]), nameAndParams[idx].node, arg})
			}

			if mappers, b := createTypeMapperFn(state, requiredTypeParams, node, nameAndTypeParams, funcTypeNode, argumentArr); b {
				if len(mappers) >= len(nameAndTypeParams) {
					for key, value := range mappers {
						newContext.TypeParams[key] = value
						knownTypeParams[key] = value
					}

					fmt.Printf("\t knownTypeParams after calculation:\n")
					for key, value := range knownTypeParams {
						fmt.Printf("\t\t %s -> %s\n", key, value.ActualType().String())
					}

					fn, resultType := execFn(tempState02, knownTypeParams, args)
					fmt.Printf("end %v\n", compiler.Fileset.Position(node.Node.Rparen).String())
					return fn, resultType
				} else {
					//createTypeMapperFn(node)
					panic("sdfds")
				}
			} else {
				panic("unreachable")
			}
		}
	}
}

type CalculateTypeFuncDeclType struct {
	index int
	node  Node[ast.Node]
}

type CalculateTypeParamType struct {
	index int
	node  Node[ast.Node]
}

type CalculateTypeArgumentType struct {
	index             int
	inputArgumentNode Node[ast.Node]
	paramStruct       Node[ast.Node]
	compiledArgument  Node[ast.Node]
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

// Todo: remove the return value
func (compiler *Compiler) calculateTypeParams(
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
			if s, b = compiler.calculateTypeParams(
				state,
				requiredTypeParams,
				CalculateTypeFuncDeclType{funcDecl.index, param},
				s,
				CalculateTypeParamType{Params.index, ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem)},
				args,
			); !b {
				return s, len(requiredTypeParams) > 0
			}
			param = ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Value)
			return compiler.calculateTypeParams(
				state,
				requiredTypeParams,
				CalculateTypeFuncDeclType{funcDecl.index, param},
				s,
				CalculateTypeParamType{Params.index, ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem)},
				args,
			)
		case *ast.FieldList:
			for idx, field := range paramItem.List {
				if findTypeMapperForMap, ok := args[idx].compiledArgument.Node.(IFindTypeMapper); ok {
					if mapper, b := findTypeMapperForMap.GetTypeMapper(); b {
						defaultMapper := mapper[0]
						switch defaultMapper.typeMapper.Kind() {
						default:
							panic(mapper)
						case reflect.Map:
							keyRt := defaultMapper.typeMapper.ActualType().Key()
							mapperKey := &WrapReflectTypeInMapper{keyRt}
							mapperKeyNode := ChangeParamNode[ast.Node, ast.Node](args[idx].compiledArgument, mapperKey)
							funcDeclParam := ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Key)
							s, b = compiler.calculateTypeParams(
								state,
								requiredTypeParams,
								CalculateTypeFuncDeclType{funcDecl.index, funcDeclParam},
								s,
								CalculateTypeParamType{Params.index, ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
								[]CalculateTypeArgumentType{
									{
										args[idx].index,
										args[idx].inputArgumentNode,
										args[idx].paramStruct,
										mapperKeyNode,
									},
								})
							if !b {
								return s, len(requiredTypeParams) > 0
							}

							valueRt := defaultMapper.typeMapper.ActualType().Elem()
							mapperValue := &WrapReflectTypeInMapper{valueRt}
							mapperValueNode := ChangeParamNode[ast.Node, ast.Node](args[idx].compiledArgument, mapperValue)
							funcDeclParam = ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Value)
							s, b = compiler.calculateTypeParams(
								state,
								requiredTypeParams,
								CalculateTypeFuncDeclType{funcDecl.index, funcDeclParam},
								s,
								CalculateTypeParamType{Params.index, ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
								[]CalculateTypeArgumentType{
									{
										args[idx].index,
										args[idx].inputArgumentNode,
										args[idx].paramStruct,
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
					s, b = compiler.calculateTypeParams(
						state,
						requiredTypeParams,
						funcDecl,
						s,
						CalculateTypeParamType{Params.index, ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
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
			mappers, _ := args[0].compiledArgument.Node.(IFindTypeMapper).GetTypeMapper()
			if len(funcDeclItem.Indices) == len(mappers) {
				for idx := 0; idx < len(mappers); idx++ {
					var b bool
					if s, b = compiler.calculateTypeParams(
						state,
						requiredTypeParams,
						CalculateTypeFuncDeclType{idx, ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Indices[idx])},
						s,
						CalculateTypeParamType{idx, Node[ast.Node]{Node: mappers[idx].typeMapper, Valid: true}},
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
				if s, b = compiler.calculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					CalculateTypeParamType{idx, ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
					[]CalculateTypeArgumentType{args[idx]},
				); !b {
					return s, len(requiredTypeParams) > 0
				}
			}
		}
		for idx, index := range funcDeclItem.Indices {
			var b bool
			if s, b = compiler.calculateTypeParams(
				state,
				requiredTypeParams,
				CalculateTypeFuncDeclType{idx, ChangeParamNode[ast.Node, ast.Node](funcDecl.node, index)},
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
			mappers, _ := args[0].compiledArgument.Node.(IFindTypeMapper).GetTypeMapper()
			return compiler.calculateTypeParams(
				state,
				requiredTypeParams,
				CalculateTypeFuncDeclType{0, ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Index)},
				s,
				CalculateTypeParamType{0, Node[ast.Node]{Node: mappers[0].typeMapper, Valid: true}},
				nil,
			)
		case *ast.FieldList:
			for idx, field := range paramItem.List {
				var b bool
				if s, b = compiler.calculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					CalculateTypeParamType{idx, ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type)},
					[]CalculateTypeArgumentType{args[idx]},
				); !b {
					return s, len(requiredTypeParams) > 0
				}
			}
		}
		indexParam := ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Index)
		return compiler.calculateTypeParams(
			state,
			requiredTypeParams,
			CalculateTypeFuncDeclType{funcDecl.index, indexParam},
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
				return compiler.calculateTypeParams(
					state,
					requiredTypeParams,
					CalculateTypeFuncDeclType{
						0,
						ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Results.List[0].Type),
					},
					s,
					CalculateTypeParamType{
						0,
						Node[ast.Node]{
							Node:  compiler.findType(state, args[0].compiledArgument, Default|TypeParamType),
							Valid: true,
						},
					},
					nil,
				)
			}
			return s, len(requiredTypeParams) > 0
		case *ast.FieldList:
			for idx, field := range paramItem.List {
				var b bool
				s, b = compiler.calculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					CalculateTypeParamType{
						Params.index,
						ChangeParamNode[ast.Node, ast.Node](Params.node, field.Type),
					},
					[]CalculateTypeArgumentType{args[idx]})
				if !b {
					return s, len(requiredTypeParams) > 0
				}
			}
			return s, len(requiredTypeParams) > 0
		case *ast.ArrayType:
			return compiler.calculateTypeParams(
				state,
				requiredTypeParams,
				funcDecl,
				s,
				CalculateTypeParamType{Params.index, ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Elt)},
				args,
			)
		}
	case *ast.ArrayType:
		param := ChangeParamNode[ast.Node, ast.Node](funcDecl.node, funcDeclItem.Elt)
		return compiler.calculateTypeParams(
			state,
			requiredTypeParams,
			CalculateTypeFuncDeclType{funcDecl.index, param},
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
				if s, b = compiler.calculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					CalculateTypeParamType{Params.index, ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Key)},
					args,
				); !b {
					return s, len(requiredTypeParams) > 0
				}
				return compiler.calculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					CalculateTypeParamType{Params.index, ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Value)},
					args,
				)
			case *ast.IndexListExpr:
				for idx, index := range paramItem.Indices {
					var b bool

					if s, b = compiler.calculateTypeParams(
						state,
						requiredTypeParams,
						funcDecl,
						s,
						CalculateTypeParamType{idx, ChangeParamNode[ast.Node, ast.Node](Params.node, index)},
						args,
					); !b {
						return s, len(requiredTypeParams) > 0
					}
				}
				return s, len(requiredTypeParams) > 0
			case *ast.IndexExpr:

				return compiler.calculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					CalculateTypeParamType{Params.index, ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Index)},
					args,
				)
			case *ast.ArrayType:

				return compiler.calculateTypeParams(
					state,
					requiredTypeParams,
					funcDecl,
					s,
					CalculateTypeParamType{Params.index, ChangeParamNode[ast.Node, ast.Node](Params.node, paramItem.Elt)},
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
						if arr, ok := findTypeMapper.GetTypeMapper(); ok {
							return compiler.calculateTypeParams(
								state,
								requiredTypeParams,
								funcDecl,
								s,
								CalculateTypeParamType{Params.index, Node[ast.Node]{Node: arr[0].typeMapper, Valid: true}},
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

					nameAndParam := nameAndParams[idx]
					var b bool
					s, b = compiler.calculateTypeParams(
						state,
						requiredTypeParams,
						funcDecl,
						s,
						CalculateTypeParamType{idx, nameAndParam.node},
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
