package internal

import (
	"go/ast"
)

func (compiler *Compiler) compileArguments(state State, argss []Node[ast.Expr], typeParams []ITypeMapper) []Node[ast.Node] {
	var result []Node[ast.Node]
	for _, arg := range argss {

		tempState := state.setCurrentNode(ChangeParamNode[ast.Expr, ast.Node](arg, arg.Node))

		param := ChangeParamNode[ast.Node, ast.Expr](state.currentNode, arg.Node)
		fn := compiler.findRhsExpression(tempState, param)
		nodeArg, _ := compiler.executeAndExpandStatement(state, typeParams, argss, fn)
		result = append(result, nodeArg...)
	}
	return result
}

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
	return func(state State, _ []ITypeMapper, _ []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		println(compiler.Fileset.Position(node.Node.Pos()).String())
		newContext := &CurrentContext{ValueInformationMap{}, map[string]ITypeMapper{}, LocalTypesMap{}, GetCompilerState[*CurrentContext](state)}

		state = SetCompilerState(newContext, state)

		param := ChangeParamNode(node, node.Node.Fun)
		tempState02 := state.setCurrentNode(ChangeParamNode[ast.Node, ast.Node](state.currentNode, node.Node.Fun))

		execFn, funcTypeNode := compiler.findFunction(tempState02, param)
		if funcTypeNode.Valid {
		}
		nameAndTypeParams := findAllParamNameAndTypes(funcTypeNode.Node.TypeParams)
		//if len(nameAndTypeParams) > 0 {
		{
			createTypeMapperFn := func(node Node[*ast.CallExpr]) (map[string]ITypeMapper, bool) {
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
					var results map[string]ITypeMapper
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
			if mappers, b := createTypeMapperFn(node); b {
				if len(mappers) >= len(nameAndTypeParams) {
					var calculatedTypeMappers []ITypeMapper
					for _, typeParam := range nameAndTypeParams {
						calculatedTypeMappers = append(calculatedTypeMappers, mappers[typeParam.name])
					}

					var args []Node[ast.Expr]
					for _, arg := range node.Node.Args {
						param := ChangeParamNode(node, arg)
						args = append(args, param)
					}
					return execFn(tempState02, calculatedTypeMappers, args)
				} else {
					//createTypeMapperFn(node)
					panic("sdfds")
				}
			} else {
				panic("unreachable")
			}
			panic("unreachable")
		}

		panic("dfgdfgdfg")
	}
}

// Todo: remove the return value
func (compiler *Compiler) calculateTypeParams(state State, funcDecl ast.Expr, s map[string]ITypeMapper, argument Node[ast.Node]) map[string]ITypeMapper {
	switch funcDeclItem := funcDecl.(type) {
	case *ast.Ident:
		switch itemArgument := argument.Node.(type) {
		case *BinaryExpr:
			itemArgumentX := itemArgument.left
			paramX := ChangeParamNode[ast.Node, ast.Node](argument, itemArgumentX.Node)
			s = compiler.calculateTypeParams(state, funcDecl, s, paramX)

			itemArgumentY := itemArgument.right
			paramY := ChangeParamNode[ast.Node, ast.Node](argument, itemArgumentY.Node)
			s = compiler.calculateTypeParams(state, funcDecl, s, paramY)

			return s
		case *ast.Ident:
			if _, ok := s[funcDeclItem.Name]; !ok {
				currentContext := GetCompilerState[*CurrentContext](state)
				if v, ok := currentContext.FindValue(itemArgument.Name); ok {
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
			panic("sfsdfdsfds")
		default:
			panic("unknown type")
		}
		return s
	case *ast.ArrayType:
		return compiler.calculateTypeParams(state, funcDeclItem.Elt, s, argument)
	case *ast.IndexExpr:
		return compiler.calculateTypeParams(state, funcDeclItem.Index, s, argument)
	case *ast.FuncType:
		switch itemA := argument.Node.(type) {
		case *ast.Ident:
			currentContext := GetCompilerState[*CurrentContext](state)

			if v, ok := currentContext.FindValue(itemA.Name); ok {
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
	default:
		panic("unreachable")
	}
}
