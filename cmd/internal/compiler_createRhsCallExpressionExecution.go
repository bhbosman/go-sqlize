package internal

import (
	"fmt"
	"go/ast"
	"reflect"
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
	compileArguments := func(state State, argss []Node[ast.Node], typeParams map[string]ITypeMapper) []Node[ast.Node] {
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
			node ast.Expr
		},
		funcTypeNode Node[*ast.FuncType],
		args []Node[ast.Node],
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
			nameAndParams := findAllParamNameAndTypes(funcTypeNode.Node.Params)
			if len(nameAndParams) > 0 && len(nameAndParams) == len(node.Node.Args) {

				sss := map[string]ITypeMapper{}
				for _, andParam := range nameAndParams {
					param := ChangeParamNode[*ast.FuncType, ast.Node](funcTypeNode, andParam.node)
					sss, _ = compiler.calculateTypeParams(state, requiredTypeParams, param, sss, funcTypeNode.Node.Params, args)
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
		println("\t knownTypeParams:")
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
			nameAndTypeParams := findAllParamNameAndTypes(funcTypeNode.Node.TypeParams)
			requiredTypeParams := map[string]bool{}
			for _, ss := range nameAndTypeParams {
				if _, ok := knownTypeParams[ss.name]; !ok {
					requiredTypeParams[ss.name] = true
				}
			}

			if mappers, b := createTypeMapperFn(state, requiredTypeParams, node, nameAndTypeParams, funcTypeNode, args); b {
				if len(mappers) >= len(nameAndTypeParams) {
					for key, value := range mappers {
						newContext.TypeParams[key] = value
						knownTypeParams[key] = value
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

// Todo: remove the return value
func (compiler *Compiler) calculateTypeParams(
	state State,
	requiredTypeParams map[string]bool,
	funcDecl Node[ast.Node],
	s map[string]ITypeMapper,
	Params ast.Node,
	args []Node[ast.Node],
) (map[string]ITypeMapper, bool) {

	switch funcDeclItem := funcDecl.Node.(type) {
	default:
		panic(funcDeclItem)
		panic("unreachable")
	case *ast.MapType:
		switch paramItem := Params.(type) {
		default:
			panic("fff")
		case *ast.Ident:
			var b bool
			param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Key)
			if s, b = compiler.calculateTypeParams(state, requiredTypeParams, param, s, paramItem, args); !b {
				return s, len(requiredTypeParams) > 0
			}
			param = ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Value)
			return compiler.calculateTypeParams(state, requiredTypeParams, param, s, paramItem, args)
		case *ast.FieldList:
			for idx, field := range paramItem.List {
				if findTypeMapperForMap, ok := args[idx].Node.(IFindTypeMapper); ok {
					if mapper, b := findTypeMapperForMap.GetTypeMapper(""); b {
						switch mapper[0].Kind() {
						default:
							panic(mapper)
						case reflect.Map:
							keyRt := mapper[0].ActualType().Key()
							mapperKey := &WrapReflectTypeInMapper{keyRt}
							mapperKeyNode := ChangeParamNode[ast.Node, ast.Node](args[idx], mapperKey)

							funcDeclParam := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Key)
							s, b = compiler.calculateTypeParams(state, requiredTypeParams, funcDeclParam, s, field.Type, []Node[ast.Node]{mapperKeyNode})
							if !b {
								return s, len(requiredTypeParams) > 0
							}

							valueRt := mapper[0].ActualType().Elem()
							mapperValue := &WrapReflectTypeInMapper{valueRt}
							mapperValueNode := ChangeParamNode[ast.Node, ast.Node](args[idx], mapperValue)
							funcDeclParam = ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Value)
							s, b = compiler.calculateTypeParams(state, requiredTypeParams, funcDeclParam, s, field.Type, []Node[ast.Node]{mapperValueNode})
							if !b {
								return s, len(requiredTypeParams) > 0
							}
						}
					}
				} else {
					var b bool
					s, b = compiler.calculateTypeParams(state, requiredTypeParams, funcDecl, s, field.Type, []Node[ast.Node]{args[idx]})
					if !b {
						return s, len(requiredTypeParams) > 0
					}
				}
			}
			return s, len(requiredTypeParams) > 0
		}

	case *ast.IndexListExpr:
		for _, index := range funcDeclItem.Indices {
			var b bool
			indexParam := ChangeParamNode[ast.Node, ast.Node](funcDecl, index)
			if s, b = compiler.calculateTypeParams(state, requiredTypeParams, indexParam, s, Params, args); !b {
				return s, len(requiredTypeParams) > 0
			}
		}
		return s, len(requiredTypeParams) > 0
	case *ast.IndexExpr:
		indexParam := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Index)
		return compiler.calculateTypeParams(state, requiredTypeParams, indexParam, s, Params, args)
	case *ast.FuncType:
		switch paramItem := Params.(type) {
		default:
			panic(paramItem)
		case *ast.Ident:
			return s, len(requiredTypeParams) > 0
		case *ast.FuncType:
			a, b := funcDeclItem.Results.List[0].Type.(*ast.Ident), paramItem.Results.List[0].Type.(*ast.Ident)
			if _, ok := requiredTypeParams[a.Name]; ok && a.Name == b.Name {
				typeMapper := compiler.findType(state, args[0], Default|TypeParamType)
				param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Results.List[0].Type)
				return compiler.calculateTypeParams(state, requiredTypeParams, param, s, typeMapper, nil)
			}
			return s, len(requiredTypeParams) > 0
		case *ast.FieldList:
			for idx, field := range paramItem.List {
				var b bool
				param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem)
				s, b = compiler.calculateTypeParams(state, requiredTypeParams, param, s, field.Type, []Node[ast.Node]{args[idx]})
				if !b {
					return s, len(requiredTypeParams) > 0
				}
			}
			return s, len(requiredTypeParams) > 0
		case *ast.ArrayType:
			param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem)
			return compiler.calculateTypeParams(state, requiredTypeParams, param, s, paramItem.Elt, args)
		}
	case *ast.ArrayType:
		param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem.Elt)
		return compiler.calculateTypeParams(state, requiredTypeParams, param, s, Params, args)
	case *ast.Ident:
		if _, ok := requiredTypeParams[funcDeclItem.Name]; ok {
			switch paramItem := Params.(type) {
			default:
				panic(paramItem)
			case *ast.MapType:
				var b bool
				if s, b = compiler.calculateTypeParams(state, requiredTypeParams, funcDecl, s, paramItem.Key, args); !b {
					return s, len(requiredTypeParams) > 0
				}
				return compiler.calculateTypeParams(state, requiredTypeParams, funcDecl, s, paramItem.Value, args)

			case *ast.IndexListExpr:
				for _, index := range paramItem.Indices {
					var b bool
					param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem)
					if s, b = compiler.calculateTypeParams(state, requiredTypeParams, param, s, index, args); !b {
						return s, len(requiredTypeParams) > 0
					}
				}
				return s, len(requiredTypeParams) > 0

			case *ast.IndexExpr:
				param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem)
				return compiler.calculateTypeParams(state, requiredTypeParams, param, s, paramItem.Index, args)
			case *ast.ArrayType:
				param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem)
				return compiler.calculateTypeParams(state, requiredTypeParams, param, s, paramItem.Elt, args)
			case ITypeMapper:
				if _, ok := requiredTypeParams[funcDeclItem.Name]; ok {
					delete(requiredTypeParams, funcDeclItem.Name)
					s[funcDeclItem.Name] = paramItem
				}
				return s, len(requiredTypeParams) > 0
			case *ast.Ident:
				if len(args) == 1 && paramItem.Name == funcDeclItem.Name {
					if findTypeMapper, ok := args[0].Node.(IFindTypeMapper); ok {
						if arr, ok := findTypeMapper.GetTypeMapper(funcDeclItem.Name); ok {
							param := ChangeParamNode[ast.Node, ast.Node](funcDecl, funcDeclItem)
							return compiler.calculateTypeParams(state, requiredTypeParams, param, s, arr[0], nil)
						}
					} else {
						panic("implement me")
					}
				}
				return s, len(requiredTypeParams) > 0
			case *ast.FieldList:
				nameAndParams := findAllParamNameAndTypes(paramItem)
				for idx, arg := range args {
					nameAndParam := nameAndParams[idx]
					var b bool
					s, b = compiler.calculateTypeParams(state, requiredTypeParams, funcDecl, s, nameAndParam.node, []Node[ast.Node]{arg})
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
