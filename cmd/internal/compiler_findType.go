package internal

import (
	"go/ast"
	"reflect"
)

type findTypeFlags int

const (
	ArgumentType findTypeFlags = 1 << iota

	ParamType
	TypeParamType
	Default = ArgumentType | ParamType
)

func (compiler *Compiler) findType(state State, node Node[ast.Node], flags findTypeFlags) ITypeMapper {
	if node.Node == nil {
		panic("node cannot be nil")
	}
	return compiler.internalFindType(0, state, node, flags).(ITypeMapper)
}

func (compiler *Compiler) internalFindType(stackIndex int, state State, node Node[ast.Node], flags findTypeFlags) interface{} {
	initOnCreateType := func(stackIndex int, unk interface{}, indexes []Node[ast.Node]) interface{} {
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
	case *ast.FuncLit:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.Type.Results.List[0].Type)
		typeMapper := compiler.findType(state, param, flags)
		return initOnCreateType(0, typeMapper, nil)
	case *ast.StructType:
		param := ChangeParamNode(node, item)
		typeMapper := compiler.createStructTypeMapper(state, param)
		return initOnCreateType(0, typeMapper, nil)
	case *ast.MapType:
		paramKey := ChangeParamNode[ast.Node, ast.Node](node, item.Key)
		rtKeyTypeMapper := compiler.findType(state, paramKey, flags)
		rtKey := rtKeyTypeMapper.MapperKeyType()

		paramValue := ChangeParamNode[ast.Node, ast.Node](node, item.Value)
		rtValueTypeMapper := compiler.findType(state, paramValue, flags)
		rtValue := rtValueTypeMapper.MapperValueType()

		rt := reflect.MapOf(rtKey, rtValue)
		return initOnCreateType(0, &TypeMapperForMap{rtKeyTypeMapper, rtValueTypeMapper, rt}, nil)
	case *ast.IndexExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		indexParam := ChangeParamNode[ast.Node, ast.Node](node, item.Index)
		unk := compiler.internalFindType(stackIndex+1, state, param, flags)
		return initOnCreateType(0, unk, []Node[ast.Node]{indexParam})
	case *ast.IndexListExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		var arrIndices []Node[ast.Node]
		for _, index := range item.Indices {
			indexParam := ChangeParamNode[ast.Node, ast.Node](node, index)
			arrIndices = append(arrIndices, indexParam)
		}
		return initOnCreateType(stackIndex, compiler.internalFindType(stackIndex+1, state, param, flags), arrIndices)
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		unk := compiler.internalFindType(stackIndex+1, state, param, flags)
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
		if flags&TypeParamType == TypeParamType {
			if typeMapper, ok := GetCompilerState[*CurrentContext](state).FindTypeParam(item.Name); ok {
				return initOnCreateType(stackIndex, typeMapper, nil)
			}
		}

		if path, ok := node.ImportMap[item.Name]; ok {
			return initOnCreateType(stackIndex, path, nil)
		}
		typeMapper := GetCompilerState[TypeMapper](state)
		if rt, ok := typeMapper[item.Name]; ok {
			return initOnCreateType(stackIndex, rt, nil)
		}
		currentContext := GetCompilerState[*CurrentContext](state)
		if onCreateType, ok := currentContext.findLocalType(item.Name); ok {
			return initOnCreateType(stackIndex, onCreateType, nil)
		}

		if onCreateType, ok := compiler.GlobalTypes[ValueKey{"", item.Name}]; ok {
			return initOnCreateType(stackIndex, onCreateType, nil)
		}
		if onCreateType, ok := compiler.GlobalTypes[ValueKey{node.RelPath, item.Name}]; ok {
			return initOnCreateType(stackIndex, onCreateType, nil)
		}

		panic(compiler.Fileset.Position(item.Pos()).String())
	case *ReflectValueExpression:
		kind := item.Rv.Kind()
		switch kind {
		case reflect.Map:
			rt := item.Rv.Type()
			return &WrapReflectTypeInMapper{rt}

		default:
			panic(compiler.Fileset.Position(item.Pos()).String())
		}

	default:
		panic(node.Node)
	}
}
