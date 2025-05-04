package internal

import (
	"go/ast"
	"go/token"
	"reflect"
)

type TypeMapperForStruct struct {
	nodeRt             reflect.Type
	actualTypeRt       reflect.Type
	typeMapperInstance reflect.Value
	vk                 ValueKey
}

func (typeMapperForStruct TypeMapperForStruct) TranslateNodeValueToReflectValue(node Node[ast.Node]) reflect.Value {
	rv, _ := isLiterateValue(node)
	newRt := typeMapperForStruct.actualTypeRt
	newRv := reflect.New(newRt).Elem()
	typeMapperForStruct.walk(newRt, newRv, rv)
	return newRv
}

func (typeMapperForStruct TypeMapperForStruct) Keys() []Node[ast.Node] {
	return nil
}

func (typeMapperForStruct TypeMapperForStruct) Pos() token.Pos {
	return token.NoPos
}

func (typeMapperForStruct TypeMapperForStruct) End() token.Pos {
	return token.NoPos
}

func (typeMapperForStruct TypeMapperForStruct) ActualType() (reflect.Type, ValueKey) {
	return typeMapperForStruct.actualTypeRt, typeMapperForStruct.vk
}

func (typeMapperForStruct TypeMapperForStruct) Kind() reflect.Kind {
	return reflect.Struct
}

func (typeMapperForStruct TypeMapperForStruct) walk(newRt reflect.Type, newRv reflect.Value, oldValue reflect.Value) {
	// TODO: remove this function
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

func (typeMapperForStruct TypeMapperForStruct) createDefaultType(parentNode Node[ast.Node]) reflect.Value {
	rv := reflect.New(typeMapperForStruct.nodeRt).Elem()
	for idx := range typeMapperForStruct.nodeRt.NumField() {
		typeMapper := typeMapperForStruct.typeMapperInstance.Field(idx).Interface().(ITypeMapper)
		rt, _ := typeMapper.ActualType()
		rvZero := reflect.Zero(rt)
		_, vk := typeMapper.ActualType()
		node := ChangeParamNode[ast.Node, ast.Node](parentNode, &ReflectValueExpression{rvZero, vk})
		rv.Field(idx).Set(reflect.ValueOf(node))
	}
	return rv
}
