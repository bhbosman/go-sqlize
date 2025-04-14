package internal

import (
	"go/ast"
	"reflect"
)

type TypeMapperForStruct struct {
	nodeRt             reflect.Type
	actualTypeRt       reflect.Type
	typeMapperInstance reflect.Value
}

func (typeMapperForStruct *TypeMapperForStruct) ActualType(state State) reflect.Type {
	return typeMapperForStruct.actualTypeRt
}

func (typeMapperForStruct *TypeMapperForStruct) MapperValueType(state State) reflect.Type {
	return typeMapperForStruct.nodeRt
}

func (typeMapperForStruct *TypeMapperForStruct) MapperKeyType(state State) reflect.Type {
	return typeMapperForStruct.actualTypeRt
}

func (typeMapperForStruct *TypeMapperForStruct) Kind() reflect.Kind {
	return reflect.Struct
}

func (typeMapperForStruct *TypeMapperForStruct) walk(newRt reflect.Type, newRv reflect.Value, oldValue reflect.Value) {
	switch newRt.Kind() {
	case reflect.Struct:
		for fieldIdx := 0; fieldIdx < newRt.NumField(); fieldIdx++ {
			fieldIdxRt := newRt.Field(fieldIdx).Type
			fieldIdxRv := newRv.Field(fieldIdx)
			fieldIdxOldValue := oldValue.Field(fieldIdx)
			node := fieldIdxOldValue.Interface().(Node[ast.Node])
			if fieldRv, ok := isLiterateValue(node); ok {
				typeMapperForStruct.walk(fieldIdxRt, fieldIdxRv, fieldRv)
			} else {
				panic("to map to a map key, the value must be a literate value")
			}
		}
	default:
		newRv.Set(oldValue.Convert(newRt))
	}
}

func (typeMapperForStruct *TypeMapperForStruct) Create(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
	switch option {
	case tmcoMapKey:
		if rv.Type() == typeMapperForStruct.nodeRt {
			newRt := typeMapperForStruct.actualTypeRt
			newRv := reflect.New(newRt).Elem()
			typeMapperForStruct.walk(newRt, newRv, rv)
			return newRv
		}
		panic("must be of type typeMapperForStruct.nodeRt")
	case tmcoMapValue:
		if rv.Type() == typeMapperForStruct.nodeRt {
			return rv
		}
		panic("must be of type typeMapperForStruct.nodeRt")
	default:
		return rv
	}
}

func (typeMapperForStruct *TypeMapperForStruct) NodeType(state State) reflect.Type {
	return typeMapperForStruct.nodeRt
}

func (typeMapperForStruct *TypeMapperForStruct) createDefaultType(state State, parentNode Node[ast.Node]) reflect.Value {
	rv := reflect.New(typeMapperForStruct.nodeRt).Elem()
	for idx := range typeMapperForStruct.nodeRt.NumField() {
		typeMapper := typeMapperForStruct.typeMapperInstance.Field(idx).Interface().(ITypeMapper)
		rvZero := reflect.Zero(typeMapper.NodeType(state))
		node := ChangeParamNode[ast.Node, ast.Node](parentNode, &ReflectValueExpression{rvZero})
		rv.Field(idx).Set(reflect.ValueOf(node))
	}
	return rv
}
