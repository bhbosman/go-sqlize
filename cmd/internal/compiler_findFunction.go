package internal

import (
	"fmt"
	"go/ast"
)

func (compiler *Compiler) findFunction(state State, node Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return compiler.internalFindFunction(0, state, node, arguments).(ExecuteStatement)
}

func (compiler *Compiler) internalFindFunction(stackIndex int, state State, node Node[ast.Expr], arguments []Node[ast.Node]) interface{} {
	switch item := node.Node.(type) {
	case *ast.FuncLit:
		param := ChangeParamNode[ast.Expr, *ast.FuncLit](node, item)
		var es OnCreateExecuteStatement = compiler.onFuncLitExecutionStatement(param)
		return compiler.initExecutionStatement(state, stackIndex, es, nil, arguments)
	case *ast.IndexExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		indexParam := ChangeParamNode(node, item.Index)
		unk := compiler.internalFindFunction(stackIndex+1, state, param, arguments)
		return compiler.initExecutionStatement(state, stackIndex, unk, []Node[ast.Expr]{indexParam}, arguments)
	case *ast.IndexListExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		var arrIndices []Node[ast.Expr]
		for _, index := range item.Indices {
			indexParam := ChangeParamNode(node, index)
			arrIndices = append(arrIndices, indexParam)
		}
		unk := compiler.internalFindFunction(stackIndex+1, state, param, arguments)
		return compiler.initExecutionStatement(state, stackIndex, unk, arrIndices, arguments)
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		unk := compiler.internalFindFunction(stackIndex+1, state, param, arguments)
		switch value := unk.(type) {
		case ast.ImportMapEntry:
			vk := ValueKey{value.Path, item.Sel.Name}
			returnValue, ok := compiler.GlobalFunctions[vk]
			if ok {
				return compiler.initExecutionStatement(state, stackIndex, returnValue, nil, arguments)
			}
			panic(fmt.Errorf("can not find function %s", vk))
		case Node[ast.Node]:
			switch nodeItem := value.Node.(type) {
			case *ReflectValueExpression:
				rvFn := nodeItem.Rv.MethodByName(item.Sel.Name)
				return compiler.initExecutionStatement(state, stackIndex, compiler.builtInStructMethods(rvFn), nil, arguments)
			default:
				panic(value.Node)
			}
		default:
			panic("sdfdsfds")
		}
	case *ast.Ident:
		if path, ok := node.ImportMap[item.Name]; ok {
			return compiler.initExecutionStatement(state, stackIndex, path, nil, arguments)
		}

		currentContext := GetCompilerState[*CurrentContext](state)
		if v, ok := currentContext.FindValue(item.Name); ok {
			return compiler.initExecutionStatement(state, stackIndex, v, nil, arguments)
		}

		if fn, ok := compiler.GlobalFunctions[ValueKey{"", item.Name}]; ok {
			return compiler.initExecutionStatement(state, stackIndex, fn, nil, arguments)
		}
		if fn, ok := compiler.GlobalFunctions[ValueKey{node.RelPath, item.Name}]; ok {
			return compiler.initExecutionStatement(state, stackIndex, fn, nil, arguments)
		}
		panic(item.Name)
	default:
		panic(node.Node)
	}
}
