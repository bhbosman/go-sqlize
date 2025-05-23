package internal

import (
	"fmt"
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

func (compiler *Compiler) createTypeMapperForFuncType(state State, node Node[*ast.FuncType]) TypeMapperForFuncType {
	itemParamsNamesAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(node, node.Node.Params))

	var inData []struct {
		rt reflect.Type
		vk ValueKey
	}
	for idx, arrItem := range itemParamsNamesAndTypeExpressions.arr {
		typeMapper := compiler.findType(state, arrItem.node, Default|TypeParamType)
		rt, vk := typeMapper.ActualType()
		if itemParamsNamesAndTypeExpressions.isVariadic && (idx == len(itemParamsNamesAndTypeExpressions.arr)-1) {
			rt = reflect.SliceOf(rt)
		}

		inData = append(inData, struct {
			rt reflect.Type
			vk ValueKey
		}{rt, vk})
	}

	resultNamesAndTypeExpressions := findAllParamNameAndTypes(ChangeParamNode(node, node.Node.Results))
	var outData []struct {
		rt reflect.Type
		vk ValueKey
	}
	for _, arrItem := range resultNamesAndTypeExpressions.arr {
		typeMapper := compiler.findType(state, arrItem.node, Default|TypeParamType)
		rt, vk := typeMapper.ActualType()
		outData = append(
			outData,
			struct {
				rt reflect.Type
				vk ValueKey
			}{rt, vk},
		)
	}

	var inArr []reflect.Type
	for _, arrItem := range inData {
		inArr = append(inArr, arrItem.rt)
	}
	var outArr []reflect.Type
	for _, arrItem := range outArr {
		outArr = append(outArr, arrItem)
	}
	funcTypeRt := reflect.FuncOf(inArr, outArr, itemParamsNamesAndTypeExpressions.isVariadic)
	return TypeMapperForFuncType{funcTypeRt, node.Key, inData, outData}
}

func (compiler *Compiler) internalFindType(stackIndex int, state State, node Node[ast.Node], flags findTypeFlags) interface{} {

	initOnCreateType := func(stackIndex int, unk interface{}, indexes []ITypeMapper) interface{} {
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
	default:
		panic(node.Node)
	case *ast.CallExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.Fun)
		return compiler.internalFindType(stackIndex+1, state, param, flags)
	case *ast.FuncLit:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.Type)
		return compiler.internalFindType(stackIndex+1, state, param, flags)
	case *ast.FuncType:
		newContext := &CurrentContext{ValueInformationMap{}, map[string]ITypeMapper{}, LocalTypesMap{}, false, GetCompilerState[*CurrentContext](state)}
		state = SetCompilerState(newContext, state)

		knownTypeParams := newContext.flattenTypeParams()
		fmt.Printf("\t knownTypeParams:\n")
		for key, value := range knownTypeParams {
			typ, _ := value.ActualType()
			fmt.Printf("\t\t %s -> %s\n", key, typ.String())
		}

		nameAndTypeParams := findAllParamNameAndTypes(ChangeParamNode(node, item.TypeParams))
		requiredTypeParams := map[string]bool{}
		for _, ss := range nameAndTypeParams.arr {
			if _, ok := knownTypeParams[ss.name]; !ok {
				requiredTypeParams[ss.name] = true
			}
		}
		if len(requiredTypeParams) > 0 {
			panic("implement this")
		}

		p01 := ChangeParamNode(node, item)
		typeMapper := compiler.createTypeMapperForFuncType(state, p01)
		return initOnCreateType(0, typeMapper, nil)
	case FuncLit:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.Type)
		typeMapper := compiler.findType(state, param, flags)
		return initOnCreateType(0, typeMapper, nil)
	case *ast.ArrayType:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.Elt)
		typeMapper := compiler.findType(state, param, flags)
		return initOnCreateType(0, typeMapper, nil)
	case *ast.Ellipsis:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.Elt)
		typeMapper := compiler.findType(state, param, flags)
		return initOnCreateType(0, typeMapper, nil)
	case *ast.StructType:
		param := ChangeParamNode(node, item)
		typeMapper := compiler.createStructTypeMapper(state, param)
		return initOnCreateType(0, typeMapper, nil)
	case *ast.MapType:
		paramKey := ChangeParamNode[ast.Node, ast.Node](node, item.Key)
		rtKeyTypeMapper := compiler.findType(state, paramKey, flags)
		actualType, _ := rtKeyTypeMapper.ActualType()
		rtKey := actualType

		paramValue := ChangeParamNode[ast.Node, ast.Node](node, item.Value)
		rtValueTypeMapper := compiler.findType(state, paramValue, flags)

		rt := reflect.MapOf(rtKey, reflect.TypeFor[Node[ast.Node]]())
		keyParam := ChangeParamNode[ast.Node, ast.Node](node, item.Key)
		valueParam := ChangeParamNode[ast.Node, ast.Node](node, item.Value)

		return initOnCreateType(0, &TypeMapperForMap{rtKeyTypeMapper, rtValueTypeMapper, rt, keyParam, valueParam}, nil)
	case *ast.IndexExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		indexParam := ChangeParamNode[ast.Node, ast.Node](node, item.Index)
		unk := compiler.internalFindType(stackIndex+1, state, param, flags)

		return initOnCreateType(0, unk, []ITypeMapper{compiler.findType(state, indexParam, Default|TypeParamType)})
	case *ast.IndexListExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		var arrIndices []ITypeMapper
		for _, index := range item.Indices {
			indexParam := ChangeParamNode[ast.Node, ast.Node](node, index)
			arrIndices = append(arrIndices, compiler.findType(state, indexParam, Default|TypeParamType))
		}
		return initOnCreateType(stackIndex, compiler.internalFindType(stackIndex+1, state, param, flags), arrIndices)
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		unk := compiler.internalFindType(stackIndex+1, state, param, flags)
		switch value := unk.(type) {
		case ImportMapEntry:
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

		if v, ok := GetCompilerState[*CurrentContext](state).FindValueByString(item.Name); ok {
			if findTypeMapper, ok := v.Node.(IFindTypeMapper); ok {
				if typeMapperArray, ok := findTypeMapper.GetTypeMapper(state); ok {
					return initOnCreateType(stackIndex, typeMapperArray[0], nil)
				}
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
		if d, ok := compiler.GlobalFunctions[ValueKey{"", item.Name}]; ok {
			p := ChangeParamNode[*ast.FuncType, ast.Node](d.funcType, d.funcType.Node)
			return compiler.internalFindType(stackIndex+1, state, p, flags)
		}

		if d, ok := compiler.GlobalFunctions[ValueKey{node.RelPath, item.Name}]; ok {
			p := ChangeParamNode[*ast.FuncType, ast.Node](d.funcType, d.funcType.Node)
			return compiler.internalFindType(stackIndex+1, state, p, flags)
		}

		panic(compiler.Fileset.Position(item.Pos()).String())
	case *ReflectValueExpression:
		kind := item.Rv.Kind()
		switch kind {
		case reflect.Map:
			rt := item.Rv.Type()
			return &WrapReflectTypeInMapper{rt, item.Vk}

		default:
			panic(compiler.Fileset.Position(item.Pos()).String())
		}

	}
}
