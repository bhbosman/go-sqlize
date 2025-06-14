package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
)

func (compiler *Compiler) addLibSomeType() {
	compiler.GlobalTypes[SomeValueKey] = compiler.registerSomeType()
	compiler.GlobalTypes[IQueryOptValueKey] = func(state State, mappers []ITypeMapper) ITypeMapper {
		return &ReflectTypeHolder{
			func() reflect.Type {
				return reflect.TypeFor[IQueryOptions]()
			},
			func() (reflect.Type, ValueKey) {
				return reflect.TypeFor[IQueryOptions](), IQueryOptValueKey
			},
			func() reflect.Kind {
				return reflect.Interface
			},
			func() reflect.Type {
				return reflect.TypeFor[IQueryOptions]()
			},
		}
	}
	compiler.GlobalTypes[IRelationshipOptValueKey] = func(state State, mappers []ITypeMapper) ITypeMapper {

		return &ReflectTypeHolder{
			func() reflect.Type {
				return reflect.TypeFor[IRelationshipOpt]()
			},
			func() (reflect.Type, ValueKey) {
				return reflect.TypeFor[IRelationshipOpt](), IRelationshipOptValueKey
			},
			func() reflect.Kind {
				return reflect.Interface
			},
			func() reflect.Type {
				return reflect.TypeFor[IRelationshipOpt]()
			},
		}
	}
}

func (compiler *Compiler) registerSomeType() OnCreateType {
	return compiler.createSomeType
}

type TypeMapperForSomeType struct {
	someRt reflect.Type
	dataRt reflect.Type
	vk     ValueKey
}

func (tm TypeMapperForSomeType) Pos() token.Pos {
	return token.NoPos
}

func (tm TypeMapperForSomeType) End() token.Pos {
	return token.NoPos
}

func (tm TypeMapperForSomeType) ActualType() (reflect.Type, ValueKey) {
	return tm.someRt, tm.vk
}

func (tm TypeMapperForSomeType) Kind() reflect.Kind {
	return reflect.Struct
}

func (compiler *Compiler) createSomeType(state State, typeParams []ITypeMapper) ITypeMapper {
	dataType, dataValueKey := typeParams[0].ActualType()
	if dataValueKey.Key == "" {
		println(dataValueKey.Key)
	}
	structFieldTag := reflect.StructTag(fmt.Sprintf(`Type:"%v" TData:"%v"`, someTypeType, dataValueKey.String()))
	var sfArr []reflect.StructField
	sfAssigned := reflect.StructField{Name: "Assigned", Type: reflect.TypeOf(true), Tag: structFieldTag}
	sfValue := reflect.StructField{Name: "Value", Type: dataType, Tag: structFieldTag}
	sfArr = append(sfArr, sfAssigned, sfValue)
	rt := reflect.StructOf(sfArr)

	return &TypeMapperForSomeType{rt, dataType, ValueKey{SomeValueKey.Folder, fmt.Sprintf("%v_%v.%v", SomeValueKey.Key, dataValueKey.Folder, dataValueKey.Key)}}
}

const someTypeType = "__built_in_Some_Type__"

func (compiler *Compiler) isTypeSomeDataType(rt reflect.Type) bool {
	if rt.Kind() == reflect.Struct && rt.NumField() >= 2 && rt.Name() == "" {
		if sv := rt.Field(0).Tag.Get("Type"); sv == someTypeType {
			return true
		}
	}
	return false
}

func (compiler *Compiler) isValueSomeDataType(rv reflect.Value) (bool, bool, reflect.Value) {
	rt := rv.Type()
	if compiler.isTypeSomeDataType(rt) {
		return true, rv.FieldByName("Assigned").Bool(), rv.FieldByName("Value")
	}
	return false, false, reflect.Value{}
}

func (compiler *Compiler) extractSomeDataType(rv reflect.Value) (reflect.Type, bool) {
	if ok, _, _ := compiler.isValueSomeDataType(rv); ok {
		return rv.FieldByName("Value").Type(), true
	}
	return nil, false
}

func (compiler *Compiler) extractSomeDataTypeMapper(rt reflect.Type) (ITypeMapper, bool) {
	if compiler.isTypeSomeDataType(rt) {
		if sf, ok := rt.FieldByName("Value"); ok {
			return &WrapReflectTypeInMapper{sf.Type, ValueKey{}}, true
		}
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
		return ChangeParamNode[ast.Node, ast.Node](state.currentNode, MultiBinaryExpr{token.LAND, binaryOperations, compiler.registerBool()(state, nil)})
	}

	var extractedArguments []Node[ast.Node]
	for _, argument := range compiledArguments {
		if _, ok := isLiterateValue(argument); ok {
			switch argItem := argument.Node.(type) {
			default:
				extractedArguments = append(extractedArguments, argument)
			case *ReflectValueExpression:
				if ok, assigned, rv := compiler.isValueSomeDataType(argItem.Rv); ok && assigned {
					unk := rv.Interface()
					if astNode, ok := unk.(ast.Node); ok {
						p01 := ChangeParamNode[ast.Node, ast.Node](argument, astNode)
						extractedArguments = append(extractedArguments, p01)
					} else {
						p01 := ChangeParamNode[ast.Node, ast.Node](argument, &ReflectValueExpression{reflect.ValueOf(unk), ValueKey{"dddddd", "CCCCCCC"}})
						extractedArguments = append(extractedArguments, p01)
					}
				} else {
					panic(argItem.Rv)
				}
			}
		} else {
			extractedArguments = append(extractedArguments, argument)
		}
	}
	return append(extractedArguments, fn()), artValue
}
