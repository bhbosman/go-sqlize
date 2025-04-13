package internal

import (
	"go/ast"
	"reflect"
)

func (compiler *Compiler) findType(state State, node Node[ast.Expr]) ITypeMapper {
	if node.Node == nil {
		panic("node cannot be nil")
	}
	return compiler.internalFindType(0, state, node).(ITypeMapper)
}

func (compiler *Compiler) internalFindType(stackIndex int, state State, node Node[ast.Expr]) interface{} {
	initOnCreateType := func(stackIndex int, unk interface{}, indexes []Node[ast.Expr]) interface{} {
		if stackIndex != 0 {
			return unk
		}
		switch value := unk.(type) {
		case OnCreateType:
			return value(state, indexes)
		case reflect.Type:
			panic("change to ITypeMapper")
		case ITypeMapper:
			return value
		default:
			panic(unk)
		}
	}

	switch item := node.Node.(type) {
	case *ast.StructType:
		param := ChangeParamNode(node, item)
		typeMapper := compiler.createStructTypeMapper(state, param)
		return initOnCreateType(0, typeMapper, nil)
	case *ast.MapType:
		paramKey := ChangeParamNode[ast.Expr, ast.Expr](node, item.Key)
		rtKeyTypeMapper := compiler.findType(state, paramKey)
		rtKey := rtKeyTypeMapper.MapperKeyType(state)

		paramValue := ChangeParamNode[ast.Expr, ast.Expr](node, item.Value)
		rtValueTypeMapper := compiler.findType(state, paramValue)
		rtValue := rtValueTypeMapper.MapperValueType(state)

		rt := reflect.MapOf(rtKey, rtValue)
		return initOnCreateType(0, &TypeMapperForMap{rtKeyTypeMapper, rtValueTypeMapper, rt}, nil)
	case *ast.IndexExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		indexParam := ChangeParamNode(node, item.Index)
		return initOnCreateType(0, compiler.internalFindType(stackIndex+1, state, param), []Node[ast.Expr]{indexParam})
	case *ast.IndexListExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		var arrIndices []Node[ast.Expr]
		for _, index := range item.Indices {
			indexParam := ChangeParamNode(node, index)
			arrIndices = append(arrIndices, indexParam)
		}
		return initOnCreateType(stackIndex, compiler.internalFindType(stackIndex+1, state, param), arrIndices)
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		unk := compiler.internalFindType(stackIndex+1, state, param)
		switch value := unk.(type) {
		case ast.ImportMapEntry:
			vk := ValueKey{value.Path, item.Sel.Name}
			returnValue, ok := compiler.GlobalTypes[vk]
			if ok {
				return initOnCreateType(stackIndex, returnValue, nil)
			}
			panic("sdfdsfds")
		default:
			panic("sdfdsfds")
		}
	case *ast.Ident:
		if path, ok := node.ImportMap[item.Name]; ok {
			return initOnCreateType(stackIndex, path, nil)
		}
		if onCreateType, ok := compiler.GlobalTypes[ValueKey{"", item.Name}]; ok {
			return initOnCreateType(stackIndex, onCreateType, nil)
		}
		if onCreateType, ok := compiler.GlobalTypes[ValueKey{node.RelPath, item.Name}]; ok {
			return initOnCreateType(stackIndex, onCreateType, nil)
		}
		typeMapper := GetCompilerState[TypeMapper](state)
		if rt, ok := typeMapper[item.Name]; ok {
			return initOnCreateType(stackIndex, rt, nil)
		}
		panic(item.Name)
	default:
		panic(node.Node)
	}
}
