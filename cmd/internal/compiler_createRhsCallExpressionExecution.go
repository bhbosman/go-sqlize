package internal

import (
	"go/ast"
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
	compileArguments := func(state State, argss []Node[ast.Node], typeParams ITypeMapperArray) []Node[ast.Node] {
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
		node Node[*ast.CallExpr],
		nameAndTypeParams []struct {
			name string
			node ast.Expr
		},
		funcTypeNode Node[*ast.FuncType]) (map[string]ITypeMapper, bool) {
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
				for idx, andParam := range nameAndParams {
					param := ChangeParamNode[*ast.CallExpr, ast.Node](node, node.Node.Args[idx])
					sss = compiler.calculateTypeParams(state, andParam.node, sss, param)
				}
				return sss, true
			}
		}
		return nil, false
	}

	return func(state State, _ ITypeMapperArray, _ []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		println(compiler.Fileset.Position(node.Node.Pos()).String())
		newContext := &CurrentContext{ValueInformationMap{}, map[string]ITypeMapper{}, LocalTypesMap{}, GetCompilerState[*CurrentContext](state)}
		state = SetCompilerState(newContext, state)
		param := ChangeParamNode(node, node.Node.Fun)
		tempState02 := state.setCurrentNode(ChangeParamNode[ast.Node, ast.Node](state.currentNode, node.Node.Fun))
		execFn, funcTypeNode := compiler.findFunction(tempState02, param)
		if !funcTypeNode.Valid {
			var args []Node[ast.Node]
			for _, arg := range node.Node.Args {
				param := ChangeParamNode[*ast.CallExpr, ast.Node](node, arg)
				args = append(args, param)
			}
			args = compileArguments(state, args, nil)
			return execFn(tempState02, nil, args)
		}
		nameAndTypeParams := findAllParamNameAndTypes(funcTypeNode.Node.TypeParams)

		if mappers, b := createTypeMapperFn(state, node, nameAndTypeParams, funcTypeNode); b {
			if len(mappers) >= len(nameAndTypeParams) {
				var calculatedTypeMappers ITypeMapperArray
				for _, typeParam := range nameAndTypeParams {
					calculatedTypeMappers = append(calculatedTypeMappers, mappers[typeParam.name])
				}

				var args []Node[ast.Node]
				for _, arg := range node.Node.Args {
					param := ChangeParamNode[*ast.CallExpr, ast.Node](node, arg)
					args = append(args, param)
				}
				args = compileArguments(state, args, calculatedTypeMappers)
				return execFn(tempState02, calculatedTypeMappers, args)
			} else {
				//createTypeMapperFn(node)
				panic("sdfds")
			}
		} else {
			panic("unreachable")
		}
	}
}

// Todo: remove the return value
func (compiler *Compiler) calculateTypeParams(state State, funcDecl ast.Node, s map[string]ITypeMapper, argument Node[ast.Node]) map[string]ITypeMapper {
	switch funcDeclItem := funcDecl.(type) {
	default:
		panic("unreachable")
	case *ast.SelectorExpr:
		// do nothing ??
		return s
	case *ast.MapType:
		switch itemArgument := argument.Node.(type) {
		default:
			panic("unreachable")
		case *ast.CompositeLit:
			typ := itemArgument.Type
			param := ChangeParamNode[ast.Node, ast.Node](argument, typ)
			return compiler.calculateTypeParams(state, funcDecl, s, param)
		case *ast.MapType:
			paramKey := ChangeParamNode[ast.Node, ast.Node](argument, itemArgument.Key)
			typeMapperKey := compiler.findType(state, paramKey, Default)
			paramKeyArg := ChangeParamNode[ast.Node, ast.Node](argument, typeMapperKey)
			s = compiler.calculateTypeParams(state, funcDeclItem.Key, s, paramKeyArg)

			paramValue := ChangeParamNode[ast.Node, ast.Node](argument, itemArgument.Value)
			typeMapperValue := compiler.findType(state, paramValue, Default)
			paramValueArg := ChangeParamNode[ast.Node, ast.Node](argument, typeMapperValue)
			return compiler.calculateTypeParams(state, funcDeclItem.Value, s, paramValueArg)
		}
	case *ast.Ident:
		switch itemArgument := argument.Node.(type) {
		default:
		case *ast.StructType:
			param := ChangeParamNode[ast.Node, ast.Node](argument, itemArgument)
			typeMapper := compiler.findType(state, param, Default)
			paramValueArg := ChangeParamNode[ast.Node, ast.Node](argument, typeMapper)
			return compiler.calculateTypeParams(state, funcDeclItem, s, paramValueArg)
		case ITypeMapper:
			s[funcDeclItem.Name] = itemArgument
			return s
		case *ast.CompositeLit:
			paramValue := ChangeParamNode[ast.Node, ast.Node](argument, itemArgument.Type)
			typeMapperValue := compiler.findType(state, paramValue, Default)
			paramValueArg := ChangeParamNode[ast.Node, ast.Node](argument, typeMapperValue)
			return compiler.calculateTypeParams(state, funcDeclItem, s, paramValueArg)
		case *BinaryExpr:
			itemArgumentX := itemArgument.left
			paramX := ChangeParamNode[ast.Node, ast.Node](argument, itemArgumentX.Node)
			s = compiler.calculateTypeParams(state, funcDecl, s, paramX)

			itemArgumentY := itemArgument.right
			paramY := ChangeParamNode[ast.Node, ast.Node](argument, itemArgumentY.Node)
			return compiler.calculateTypeParams(state, funcDecl, s, paramY)
		case *ast.Ident:
			if _, ok := s[funcDeclItem.Name]; !ok {
				currentContext := GetCompilerState[*CurrentContext](state)
				if v, ok := currentContext.FindValueByString(itemArgument.Name); ok {
					if findTypeMapper, ok := v.Node.(IFindTypeMapper); ok {
						if typeMapper, ok := findTypeMapper.GetTypeMapper(); ok {
							if len(typeMapper) == 1 {
								s[funcDeclItem.Name] = typeMapper[0]
							} else {
								panic("sfsdfdsfds")
							}
						}
					} else {
						itemArgument := v.Node
						paramX := ChangeParamNode[ast.Node, ast.Node](argument, itemArgument)
						return compiler.calculateTypeParams(state, funcDecl, s, paramX)
					}

				} else if expr, ok := argument.Node.(ast.Expr); ok {
					param := ChangeParamNode[ast.Node, ast.Node](argument, expr)
					s[funcDeclItem.Name] = compiler.findType(state, param, Default)
				}
			}
		case *ast.SelectorExpr:
			if _, ok := s[funcDeclItem.Name]; !ok {
				currentContext := GetCompilerState[*CurrentContext](state)
				if v, ok := currentContext.FindTypeFromNode(argument); ok {
					if len(v) == 1 {
						s[funcDeclItem.Name] = v[0]
					} else {
						panic("dfsdfds")
					}
				} else {
					s[funcDeclItem.Name] = compiler.findType(state, argument, Default)
				}
			}
		case *ast.CallExpr:
			param := ChangeParamNode[ast.Node, ast.Expr](argument, itemArgument.Fun)
			_, funcType := compiler.findFunction(state, param)
			if len(itemArgument.Args) == funcType.Node.Params.NumFields() {
				nameAndParams := findAllParamNameAndTypes(funcType.Node.Params)
				for idx := 0; idx < len(itemArgument.Args); idx++ {
					argItem := itemArgument.Args[idx]
					paramItem := nameAndParams[idx].node
					param := ChangeParamNode[ast.Node, ast.Node](argument, argItem)
					s = compiler.calculateTypeParams(state, paramItem, s, param)
				}
				return s
			}
		case *ast.BinaryExpr:
			itemArgumentX := itemArgument.X
			paramX := ChangeParamNode[ast.Node, ast.Node](argument, itemArgumentX)
			s = compiler.calculateTypeParams(state, funcDecl, s, paramX)

			itemArgumentY := itemArgument.Y
			paramY := ChangeParamNode[ast.Node, ast.Node](argument, itemArgumentY)
			s = compiler.calculateTypeParams(state, funcDecl, s, paramY)

			return s
		case IFindTypeMapper:
			if arr, ok := itemArgument.GetTypeMapper(); ok && len(arr) == 1 {
				s[funcDeclItem.Name] = arr[0]
			} else {
				panic("sfsdfdsfds")
			}
			return s
		case *ast.BasicLit:
			param := ChangeParamNode[ast.Node, *ast.BasicLit](argument, itemArgument)
			arr, _ := compiler.runBasicLitNode(state, param)
			return compiler.calculateTypeParams(state, funcDeclItem, s, arr[0])
		case *ReflectValueExpression:
			s[funcDeclItem.Name] = &WrapReflectTypeInMapper{itemArgument.Rv.Type()}
			return s
		}
		return s
	case *ast.ArrayType:
		return compiler.calculateTypeParams(state, funcDeclItem.Elt, s, argument)
	case *ast.IndexListExpr:
		switch itemArgument := argument.Node.(type) {
		default:
			panic("sfsdfds")
		case *ast.SelectorExpr:
			currentContext := GetCompilerState[*CurrentContext](state)
			param := ChangeParamNode[ast.Node, ast.Node](argument, itemArgument)
			if node, ok := currentContext.FindValueByNode(param); ok {
				if findTypeMapper, ok := node.Node.(IFindTypeMapper); ok {
					if mappers, b := findTypeMapper.GetTypeMapper(); b && len(mappers) == len(funcDeclItem.Indices) {
						for idx, index := range funcDeclItem.Indices {
							param := ChangeParamNode[ast.Node, ast.Node](argument, mappers[idx])
							s = compiler.calculateTypeParams(state, index, s, param)
						}
					}
				}
			}
			return s
		case *ast.Ident:
			currentContext := GetCompilerState[*CurrentContext](state)
			if v, ok := currentContext.FindValueByString(itemArgument.Name); ok {
				if findTypeMapper, ok := v.Node.(IFindTypeMapper); ok {
					if mappers, b := findTypeMapper.GetTypeMapper(); b && len(mappers) == len(funcDeclItem.Indices) {
						for idx, index := range funcDeclItem.Indices {
							param := ChangeParamNode[ast.Node, ast.Node](argument, mappers[idx])
							s = compiler.calculateTypeParams(state, index, s, param)
						}
					} else {
						panic("fix the object count of gettyoemappers")
					}
				} else {
					panic("implement IFindTypeMapper")
				}
			}
			return s
		}
	case *ast.IndexExpr:
		return compiler.calculateTypeParams(state, funcDeclItem.Index, s, argument)
	case *ast.FuncType:
		switch itemA := argument.Node.(type) {
		case *ast.Ident:
			currentContext := GetCompilerState[*CurrentContext](state)
			if v, ok := currentContext.FindValueByString(itemA.Name); ok {
				switch itemV := v.Node.(type) {
				case *ast.FuncLit:
					param := ChangeParamNode[ast.Node, ast.Node](v, itemV)
					return compiler.calculateTypeParams(state, funcDeclItem, s, param)
				default:
					panic("unknown type")
				}
			} else {
				s[itemA.Name] = compiler.findType(state, argument, Default)
			}

			panic("fgfd")
		case *ast.FuncLit:
			if funcDeclItem.Params != nil {
				for idx, field := range funcDeclItem.Params.List {
					param := ChangeParamNode[ast.Node, ast.Node](argument, itemA.Type.Params.List[idx].Type)
					s = compiler.calculateTypeParams(state, field.Type, s, param)
				}
			}
			if funcDeclItem.Results != nil {
				for idx, field := range funcDeclItem.Results.List {
					param := ChangeParamNode[ast.Node, ast.Node](argument, itemA.Type.Params.List[idx].Type)
					s = compiler.calculateTypeParams(state, field.Type, s, param)
				}
			}
			return s
		default:
			panic("unreachable")
		}
	}
}
