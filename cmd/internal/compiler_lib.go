package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"strconv"
)

const libFolder = "github.com/bhbosman/go-sqlize/lib"

var SomeValueKey = ValueKey{libFolder, "Some"}
var IQueryOptValueKey = ValueKey{libFolder, "IQueryOpt"}
var IRelationshipOptValueKey = ValueKey{libFolder, "IRelationshipOpt"}

func (compiler *Compiler) addLibFunctions() {
	compiler.addLibSomeType()

	compiler.GlobalTypes[ValueKey{libFolder, "Dictionary"}] = compiler.registerLibType()

	compiler.GlobalFunctions[ValueKey{libFolder, "Query"}] = functionInformation{compiler.libQueryImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "QueryTop"}] = functionInformation{compiler.libQueryTopImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "QueryDistinct"}] = functionInformation{compiler.libQueryDistinctImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "Map"}] = functionInformation{compiler.libMapImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "GenerateSql"}] = functionInformation{compiler.libGenerateSqlImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "GenerateSqlTest"}] = functionInformation{compiler.libGenerateSqlTestImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "Atoi"}] = functionInformation{compiler.libAtoiImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "SetSomeValue"}] = functionInformation{compiler.libSetSomeValueImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "SetSomeNone"}] = functionInformation{compiler.libSetSomeNoneImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "IsSomeAssigned"}] = functionInformation{compiler.libIsSomeAssignedImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "SomeData"}] = functionInformation{compiler.libSomeDataImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "SomeData2"}] = functionInformation{compiler.libSomeData2Implementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "GetSomeData"}] = functionInformation{compiler.libGetSomeDataImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "GetSomeData02"}] = functionInformation{compiler.libGetSomeData02Implementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "GetSomeData03"}] = functionInformation{compiler.libGetSomeData03Implementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "GetSomeData04"}] = functionInformation{compiler.libGetSomeData04Implementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "GetSomeData05"}] = functionInformation{compiler.libGetSomeData05Implementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "CreateDictionary"}] = functionInformation{compiler.libCreateDictionaryImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "DictionaryLookup"}] = functionInformation{compiler.libDictionaryLookupImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "DictionaryDefault"}] = functionInformation{compiler.libDictionaryDefaultImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "TypeFor"}] = functionInformation{compiler.libTypeForImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "TestType"}] = functionInformation{compiler.libTestTypeImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "CoreRelationship"}] = functionInformation{compiler.libCoreRelationshipImplementation, Node[*ast.FuncType]{}, true}
	compiler.GlobalFunctions[ValueKey{libFolder, "CombinePredFunctionsWithAnd"}] = functionInformation{compiler.libCombinePredFunctionsWithAndImplementation, Node[*ast.FuncType]{}, true}
}

func (compiler *Compiler) libGenerateSql(state State, argument Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	switch item := argument.Node.(type) {
	case *TrailRecord:
		p01 := ChangeParamNode[ast.Node, *TrailRecord](argument, item)
		s := compiler.trailRecordToSelectStatement(state, p01)
		basicLit := &ast.BasicLit{state.currentNode.Node.Pos(), token.STRING, strconv.Quote(s)}
		nodeValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, basicLit)
		return []Node[ast.Node]{nodeValue}, artValue
	default:
		panic("implementation required")
	}
}

func (compiler *Compiler) libGenerateSqlImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("Lib.GenerateSql implementation requires 1 arguments, got %d", len(arguments)))
		}

		return compiler.libGenerateSql(state, arguments[0])
	}
}

func (compiler *Compiler) libGenerateSqlTestImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("Lib.GenerateSqlTest implementation requires 1 arguments, got %d", len(arguments)))
		}
		args := arguments[0]
		ans, _ := compiler.libGenerateSql(state, args)
		if rv, isLiterate := isLiterateValue(ans[0]); isLiterate && rv.Kind() == reflect.String {
			currentContext := GetCompilerState[*CurrentContext](state)
			if value, b := currentContext.FindValueByString("__stdOut__"); b {
				switch nodeValue := value.Node.(type) {
				case *ReflectValueExpression:
					if wr, ok := nodeValue.Rv.Interface().(io.Writer); ok {
						_, _ = wr.Write([]byte(rv.String()))
					}
					return nil, artNone
				default:
					panic(fmt.Sprintf("libGenerateSqlTestImplementation *ReflectValueExpression not found"))
				}
			}
			panic(fmt.Sprintf("libGenerateSqlTestImplementation __stdOut__ not found"))
		}
		panic(fmt.Sprintf("libGenerateSqlTestImplementation value from GenerateSql not literal"))
	}
}

func (compiler *Compiler) libAtoiImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("Lib.Atoi implementation requires 1 arguments, got %d", len(arguments)))
		}
		nodes, resultType := compiler.strconvAtoiCompiled(state, arguments)
		return nodes[:1], resultType
	}
}

func (compiler *Compiler) libSetSomeValueImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("SetSomeValue implementation requires 1 arguments, got %d", len(arguments)))
		}
		if value, ok := isLiterateValue(arguments[0]); ok {
			tm := typeParams[funcTypeNode.Node.TypeParams.List[0].Names[0].Name]
			mapper := compiler.createSomeType(state, []ITypeMapper{tm})
			typ, vk := mapper.ActualType()
			rv := reflect.New(typ).Elem()
			rv.FieldByName("Assigned").SetBool(true)
			rv.FieldByName("Value").Set(value)
			param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rv, vk})
			return []Node[ast.Node]{param}, artValue
		} else if fvk, ok := arguments[0].Node.(IFindValueKey); ok {
			vk := fvk.GetValueKey()
			rt := reflect.ValueOf(arguments[0].Node).Type()
			typeMapper := &WrapReflectTypeInMapper{rt, vk}
			mapper := compiler.createSomeType(state, []ITypeMapper{typeMapper})
			typ, vk := mapper.ActualType()
			rv := reflect.New(typ).Elem()
			rv.FieldByName("Assigned").SetBool(true)
			rv.FieldByName("Value").Set(reflect.ValueOf(arguments[0].Node))
			param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rv, vk})
			return []Node[ast.Node]{param}, artValue
		} else {
			panic("implemt IFindValueKey ")
		}
	}
}

func (compiler *Compiler) libSetSomeNoneImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		tm := typeParams[funcTypeNode.Node.TypeParams.List[0].Names[0].Name]
		mapper := compiler.createSomeType(state, []ITypeMapper{tm})
		typ, vk := mapper.ActualType()
		rv := reflect.New(typ).Elem()
		rv.FieldByName("Assigned").SetBool(false)
		param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rv, vk})
		return []Node[ast.Node]{param}, artValue
	}
}

func (compiler *Compiler) libIsSomeAssignedImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("IsSomeAssigned implementation requires 1 arguments, got %d", len(arguments)))
		}

		return compiler.IsSomeAssigned(state, arguments[0])
	}
}

func (compiler *Compiler) libSomeDataImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("SomeData implementation requires 1 arguments, got %d", len(arguments)))
		}
		return arguments[0:1], artValue
	}
}

func (compiler *Compiler) libSomeData2Implementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("SomeData2 implementation requires 1 arguments, got %d", len(arguments)))
		}

		v, _ := compiler.IsSomeAssigned(state, arguments[0])
		result := append(arguments, v[0])
		return result, artValue
	}
}

func (compiler *Compiler) libGetSomeDataImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.getGetSomeDataNCompiled(state, funcTypeNode, 1, arguments)
	}
}

func (compiler *Compiler) libGetSomeData02Implementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.getGetSomeDataNCompiled(state, funcTypeNode, 2, arguments)
	}
}

func (compiler *Compiler) libGetSomeData03Implementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.getGetSomeDataNCompiled(state, funcTypeNode, 3, arguments)
	}
}

func (compiler *Compiler) libGetSomeData04Implementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.getGetSomeDataNCompiled(state, funcTypeNode, 4, arguments)
	}
}

func (compiler *Compiler) libGetSomeData05Implementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.getGetSomeDataNCompiled(state, funcTypeNode, 5, arguments)
	}
}

func (compiler *Compiler) libCreateDictionaryImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 2 {
			panic(fmt.Errorf("CreateDictionary implementation requires 2 arguments, got %d", len(arguments)))
		}

		nameAndParams := findAllParamNameAndTypes(ChangeParamNode(funcTypeNode, funcTypeNode.Node.TypeParams))
		key := nameAndParams.arr[0]
		value := nameAndParams.arr[1]
		if rv00, ok00 := isLiterateValue(arguments[0]); ok00 {
			return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](
				state.currentNode,
				&DictionaryExpression{
					rv00,
					arguments[1],
					typeParams[key.name],
					typeParams[value.name],
				})}, artValue
		}
		panic(fmt.Errorf("createDictionary implementation requires literal values"))
	}
}

func (compiler *Compiler) libDictionaryLookupImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return libDictionaryLookupImplementation{compiler, state}.ExecuteStatement()
}

func (compiler *Compiler) IsSomeAssigned(state State, argument Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	switch nodeItem := argument.Node.(type) {
	case *ReflectValueExpression:
		if isSomeData, assigned, _ := compiler.isValueSomeDataType(nodeItem.Rv); isSomeData {
			param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(assigned), nodeItem.Vk})
			return []Node[ast.Node]{param}, artValue
		}
		panic("Not a SomeData")
	default:
		param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &CheckForNotNullExpression{argument, compiler.registerBool()(state, nil)})
		return []Node[ast.Node]{param}, artValue
	}
}

func (compiler *Compiler) libDictionaryDefaultImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("DictionaryLookup implementation requires 1 arguments, got %d", len(arguments)))
		}
		dictionaryExpression := arguments[0].Node.(*DictionaryExpression)
		return []Node[ast.Node]{dictionaryExpression.defaultValue}, artReturn
	}
}

func (compiler *Compiler) libTypeForImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		typeMapper := typeParams[funcTypeNode.Node.TypeParams.List[0].Names[0].Name]
		actualType, key := typeMapper.ActualType()
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(actualType), key})}, artNone
	}
}

func (compiler *Compiler) libTestTypeImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		typeMapper := typeParams[funcTypeNode.Node.TypeParams.List[0].Names[0].Name]
		actualType, key := typeMapper.ActualType()

		println(actualType.String(), key.Key)
		rv00 := arguments[0].Node.(*ReflectValueExpression).Rv
		rt00 := rv00.Type()
		if rt00 != actualType {
			panic("test fails")
		}
		rv01 := arguments[1].Node.(*ReflectValueExpression).Rv
		rt01 := rv01.Interface().(reflect.Type)
		if rt00 != rt01 {
			panic("test fails")
		}
		return nil, artNone
	}
}
