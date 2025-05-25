package internal

import (
	"fmt"
	"go/ast"
	"strings"
)

func (compiler *Compiler) createRhsCallExpressionExecution(node Node[*ast.CallExpr]) ExecuteStatement {
	return func(state State, _ map[string]ITypeMapper, _ []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		fmt.Printf("begin %v\n", compiler.Fileset.Position(node.Node.Lparen).String())
		defer fmt.Printf("end %v\n", compiler.Fileset.Position(node.Node.Rparen).String())
		newContext := &CurrentContext{ValueInformationMap{}, map[string]ITypeMapper{}, LocalTypesMap{}, false, GetCompilerState[*CurrentContext](state)}
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
		paramInformation := func(ft Node[*ast.FuncType]) *findAllParamNameAndTypesResult {
			if ft.Node != nil {
				r := findAllParamNameAndTypes(ChangeParamNode(funcTypeNode, funcTypeNode.Node.Params))
				return &r
			}
			return nil
		}(funcTypeNode)
		args = compiler.compileArguments(
			state,
			args,
			knownTypeParams,
			paramInformation)
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
				var argumentArr []Node[ast.Node]
				for _, arg := range args {
					argumentArr = append(argumentArr, arg)
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
			if resultType == artReturn {
				return fn, resultType
			}
			if len(fn) == 1 {
				switch dd := fn[0].Node.(type) {
				case FuncLit:
					param := ChangeParamNode[ast.Node, FuncLit](fn[0], dd)
					onCreateExecuteStatement := compiler.onFuncLitExecutionStatement(param)
					executeStatement := onCreateExecuteStatement(state, funcTypeNode)
					return executeStatement(state, deltaTypeParams, args)
				default:
					return fn, resultType
				}
			}
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
	args []Node[ast.Node],
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
