package internal

import "go/ast"

func (compiler *Compiler) compileArguments(state State, argss []Node[ast.Node], typeParams map[string]ITypeMapper, paramInformation *findAllParamNameAndTypesResult) []Node[ast.Node] {
	var result []Node[ast.Node]
	for idx, arg := range argss {
		paramType := func(idx int) Node[ast.Node] {
			if paramInformation == nil {
				return Node[ast.Node]{}
			}
			if paramInformation.isVariadic {
				if idx <= len(paramInformation.arr)-1 {
					return paramInformation.arr[idx].node
				}
				return paramInformation.arr[len(paramInformation.arr)-1].node

			} else {
				return paramInformation.arr[idx].node
			}
		}(idx)

		nodeArg := compiler.internalCompileArguments(state, arg, paramType, argss, typeParams)
		result = append(result, nodeArg...)
	}
	return result
}

func (compiler *Compiler) internalCompileArguments(state State, arg Node[ast.Node], paramType Node[ast.Node], argss []Node[ast.Node], typeParams map[string]ITypeMapper) []Node[ast.Node] {
	if paramType.Valid {
		switch ptItem := paramType.Node.(type) {
		default:
			panic(ptItem)
			panic("unknown type")
		case *ast.FuncType:
			return []Node[ast.Node]{arg}
		case *ast.ArrayType:
			p := ChangeParamNode[ast.Node, ast.Node](paramType, ptItem.Elt)
			return compiler.internalCompileArguments(state, arg, p, argss, typeParams)
		case *ast.Ellipsis:
			p := ChangeParamNode[ast.Node, ast.Node](paramType, ptItem.Elt)
			return compiler.internalCompileArguments(state, arg, p, argss, typeParams)
		case *ast.Ident, *ast.IndexExpr, *ast.MapType, *ast.IndexListExpr, *ast.SelectorExpr:
			break
		}
	}

	tempState := state.setCurrentNode(ChangeParamNode[ast.Node, ast.Node](arg, arg.Node))
	param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, arg.Node)
	fn := compiler.findRhsExpression(tempState, param)
	nodeArg, _ := compiler.executeAndExpandStatement(state, typeParams, nil, fn)
	return nodeArg
}
