package internal

import "reflect"

type TypeMapperForMap struct {
	keyTypeMapper   ITypeMapper
	valueTypeMapper ITypeMapper
	mapRt           reflect.Type
}

func (tyfm *TypeMapperForMap) ActualType(state State) reflect.Type {
	return tyfm.mapRt
}

func (tyfm *TypeMapperForMap) MapperValueType(state State) reflect.Type {
	return tyfm.valueTypeMapper.MapperValueType(state)
}

func (tyfm *TypeMapperForMap) MapperKeyType(state State) reflect.Type {
	return tyfm.keyTypeMapper.MapperKeyType(state)
}

func (tyfm *TypeMapperForMap) Kind() reflect.Kind {
	return tyfm.mapRt.Kind()
}

func (tyfm *TypeMapperForMap) Create(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
	//TODO implement me
	panic("implement me")
}

func (tyfm *TypeMapperForMap) NodeType(state State) reflect.Type {
	return tyfm.mapRt
}
