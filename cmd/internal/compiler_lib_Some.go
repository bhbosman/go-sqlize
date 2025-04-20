package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
)

func (compiler *Compiler) addLibSomeType() {
	compiler.GlobalTypes[SomeValueKey] = compiler.registerSomeType()
}

func (compiler *Compiler) registerSomeType() OnCreateType {
	return compiler.createSomeType
}

func (compiler *Compiler) createSomeType(state State, typeParams []ITypeMapper) ITypeMapper {
	typ, vk := typeParams[0].ActualType()
	structFieldTag := reflect.StructTag(fmt.Sprintf(`Type:"%v" TData:"%v"`, someTypeType, vk.String()))
	var sfArr []reflect.StructField
	sfAssigned := reflect.StructField{Name: "Assigned", Type: reflect.TypeOf(true), Tag: structFieldTag}
	sfValue := reflect.StructField{Name: "Value", Type: typ, Tag: structFieldTag}
	sfArr = append(sfArr, sfAssigned, sfValue)
	rt := reflect.StructOf(sfArr)
	fn := func() reflect.Type {
		return rt
	}
	return &ReflectTypeHolder{
		nil,
		fn,
		func() (reflect.Type, ValueKey) {
			return rt, SomeValueKey
		},
		func() reflect.Kind {
			return reflect.Struct
		},
		fn,
	}
}

const someTypeType = "__built_in_Some_Type__"

func (compiler *Compiler) isValueSomeDataType(rv reflect.Value) (bool, bool, reflect.Value) {
	rt := rv.Type()
	if rt.Kind() == reflect.Struct && rt.NumField() >= 2 && rt.Name() == "" {
		if sv := rt.Field(0).Tag.Get("Type"); sv == someTypeType {
			return true, rv.FieldByName("Assigned").Bool(), rv.FieldByName("Value")
		}
	}
	return false, false, reflect.Value{}
}

func (compiler *Compiler) extractSomeDataType(rv reflect.Value) (reflect.Type, bool) {
	if ok, _, _ := compiler.isValueSomeDataType(rv); ok {
		return rv.FieldByName("Value").Type(), true
	}
	return nil, false
}

func (compiler *Compiler) extractSomeDataTag(rv reflect.Value, tagName string) (string, bool) {
	if ok, _, _ := compiler.isValueSomeDataType(rv); ok {
		tag := rv.Type().Field(0).Tag
		return tag.Get(tagName), true
	}
	return "", false
}

func (compiler *Compiler) getGetSomeDataNCompiled(state State, funcTypeNode Node[*ast.FuncType], n int, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	fn := func() Node[ast.Node] {
		var binaryOperations []Node[ast.Node]
		for _, arg := range compiledArguments {
			arr, _ := compiler.IsSomeAssigned(state, arg)
			for _, item := range arr {
				switch nodeItem := item.Node.(type) {
				case *ReflectValueExpression:
					if nodeItem.Rv.Kind() == reflect.Bool {
						if nodeItem.Rv.Bool() {
							continue
						} else {
							return item
						}
					}
				}
				binaryOperations = append(binaryOperations, item)
			}
		}
		return ChangeParamNode[ast.Node, ast.Node](state.currentNode, MultiBinaryExpr{token.LAND, binaryOperations})
	}
	return append(compiledArguments, fn()), artValue
}
