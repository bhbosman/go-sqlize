package internal

import (
	"go/ast"
	"go/token"
	"reflect"
)

type TypeMapperForMap struct {
	keyTypeMapper   ITypeMapper
	valueTypeMapper ITypeMapper
	mapRt           reflect.Type
	key             Node[ast.Node]
	value           Node[ast.Node]
}

func (tyfm *TypeMapperForMap) Keys() []Node[ast.Node] {
	return []Node[ast.Node]{tyfm.key, tyfm.value}
}

func (tyfm *TypeMapperForMap) Pos() token.Pos {
	return token.NoPos
}

func (tyfm *TypeMapperForMap) End() token.Pos {
	return token.NoPos
}

func (tyfm *TypeMapperForMap) ActualType() reflect.Type {
	return tyfm.mapRt
}

func (tyfm *TypeMapperForMap) MapperValueType() reflect.Type {
	return tyfm.valueTypeMapper.MapperValueType()
}

func (tyfm *TypeMapperForMap) MapperKeyType() reflect.Type {
	return tyfm.keyTypeMapper.MapperKeyType()
}

func (tyfm *TypeMapperForMap) Kind() reflect.Kind {
	return tyfm.mapRt.Kind()
}

func (tyfm *TypeMapperForMap) Create(option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
	//TODO implement me
	panic("implement me")
}

func (tyfm *TypeMapperForMap) NodeType() reflect.Type {
	return tyfm.mapRt
}
