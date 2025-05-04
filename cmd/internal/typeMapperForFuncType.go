package internal

import (
	"go/token"
	"reflect"
)

type ITypeMapperForFuncType interface {
	ITypeMapper
	InCount() int
	OutCount() int
	In(int) ITypeMapper
	Out(int) ITypeMapper
}

type TypeMapperForFuncType struct {
	rt     reflect.Type
	vk     ValueKey
	inData []struct {
		rt reflect.Type
		vk ValueKey
	}
	outData []struct {
		rt reflect.Type
		vk ValueKey
	}
}

func (typeMapperForFuncType TypeMapperForFuncType) InCount() int {
	return len(typeMapperForFuncType.inData)
}

func (typeMapperForFuncType TypeMapperForFuncType) OutCount() int {
	return len(typeMapperForFuncType.outData)
}

func (typeMapperForFuncType TypeMapperForFuncType) In(i int) ITypeMapper {
	data := typeMapperForFuncType.inData[i]
	return &WrapReflectTypeInMapper{data.rt, data.vk}

}

func (typeMapperForFuncType TypeMapperForFuncType) Out(i int) ITypeMapper {
	data := typeMapperForFuncType.outData[i]
	return &WrapReflectTypeInMapper{data.rt, data.vk}
}

func (typeMapperForFuncType TypeMapperForFuncType) Pos() token.Pos {
	return token.NoPos
}

func (typeMapperForFuncType TypeMapperForFuncType) End() token.Pos {
	return token.NoPos
}

func (typeMapperForFuncType TypeMapperForFuncType) ActualType() (reflect.Type, ValueKey) {
	return typeMapperForFuncType.rt, typeMapperForFuncType.vk
}

func (typeMapperForFuncType TypeMapperForFuncType) Kind() reflect.Kind {

	return typeMapperForFuncType.rt.Kind()
}
