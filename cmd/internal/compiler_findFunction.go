package internal

import (
	"fmt"
	"go/ast"
)

func (compiler *Compiler) findFunction(state State, node Node[ast.Expr]) (ExecuteStatement, Node[*ast.FuncType]) {
	unk01, unk02 := compiler.internalFindFunction(0, state, node)
	return unk01.(ExecuteStatement), unk02
}

func (compiler *Compiler) internalFindFunction(stackIndex int, state State, node Node[ast.Expr]) (interface{}, Node[*ast.FuncType]) {
	switch item := node.Node.(type) {
	case *ast.FuncLit:
		param := ChangeParamNode[ast.Expr, *ast.FuncLit](node, item)
		paramFuncType := ChangeParamNode[ast.Expr, *ast.FuncType](node, item.Type)
		es := compiler.onFuncLitExecutionStatement(param)
		return compiler.initExecutionStatement(state, stackIndex, es, paramFuncType, nil)
	case *ast.IndexExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		indexParam := ChangeParamNode[ast.Expr, ast.Node](node, item.Index)
		unk, unk2 := compiler.internalFindFunction(stackIndex+1, state, param)
		return compiler.initExecutionStatement(state, stackIndex, unk, unk2, []Node[ast.Node]{indexParam})
	case *ast.IndexListExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		var arrIndices []Node[ast.Node]
		for _, index := range item.Indices {
			indexParam := ChangeParamNode[ast.Expr, ast.Node](node, index)
			arrIndices = append(arrIndices, indexParam)
		}
		unk, unk2 := compiler.internalFindFunction(stackIndex+1, state, param)
		return compiler.initExecutionStatement(state, stackIndex, unk, unk2, arrIndices)
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		unk, _ := compiler.internalFindFunction(stackIndex+1, state, param)
		switch value := unk.(type) {
		case ast.ImportMapEntry:
			vk := ValueKey{value.Path, item.Sel.Name}
			returnValue, ok := compiler.GlobalFunctions[vk]
			if ok {
				return compiler.initExecutionStatement(state, stackIndex, returnValue.fn, returnValue.funcType, nil)
			}
			panic(fmt.Errorf("can not find function %s", vk))
		case Node[ast.Node]:
			switch nodeItem := value.Node.(type) {
			case *ReflectValueExpression:
				rvFn := nodeItem.Rv.MethodByName(item.Sel.Name)
				return compiler.initExecutionStatement(state, stackIndex, compiler.builtInStructMethods(rvFn), Node[*ast.FuncType]{}, nil)
			default:
				panic(value.Node)
			}
		default:
			panic("sdfdsfds")
		}
	case *ast.Ident:
		if path, ok := node.ImportMap[item.Name]; ok {
			return compiler.initExecutionStatement(state, stackIndex, path, Node[*ast.FuncType]{}, nil)
		}

		currentContext := GetCompilerState[*CurrentContext](state)
		if v, ok := currentContext.FindValueByString(item.Name); ok {
			return compiler.initExecutionStatement(state, stackIndex, v, Node[*ast.FuncType]{}, nil)
		}

		if fn, ok := compiler.GlobalFunctions[ValueKey{"", item.Name}]; ok {
			return compiler.initExecutionStatement(state, stackIndex, fn.fn, fn.funcType, nil)
		}
		if fn, ok := compiler.GlobalFunctions[ValueKey{node.RelPath, item.Name}]; ok {
			return compiler.initExecutionStatement(state, stackIndex, fn.fn, fn.funcType, nil)
		}
		panic(item.Name)
	default:
		panic(node.Node)
	}
}

func (compiler *Compiler) initExecutionStatement(state State, stackIndex int, unk interface{}, unk02 Node[*ast.FuncType], typeParams []Node[ast.Node]) (interface{}, Node[*ast.FuncType]) {
	if stackIndex != 0 {
		return unk, unk02
	}
	switch value := unk.(type) {
	case OnCreateExecuteStatement:
		return value(state), unk02
	case Node[ast.Node]:
		switch value02 := value.Node.(type) {
		case ast.Expr:
			param := ChangeParamNode(value, value02)
			unkValue, unkValue02 := compiler.internalFindFunction(stackIndex+1, state, param)
			return compiler.initExecutionStatement(state, stackIndex, unkValue, unkValue02, typeParams)
		default:
			panic(unk)
		}
	case functionInformation:
		panic("gfdgf")
		return value.fn(state), value.funcType
	default:
		panic(unk)
	}
}
