package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"strconv"
	"strings"
)

const libFolder = "github.com/bhbosman/go-sqlize/lib"

func (compiler *Compiler) addLibFunctions() {
	compiler.GlobalFunctions[ValueKey{libFolder, "Query"}] = functionInformation{compiler.libQueryImplementation, Node[*ast.FuncType]{}, true}
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
	//compiler.GlobalFunctions[ValueKey{libFolder, "TypeFor"}] = functionInformation{compiler.libTypeForImplementation, Node[*ast.FuncType]{}, true}
	//compiler.GlobalFunctions[ValueKey{libFolder, "TestType"}] = functionInformation{compiler.libTestTypeImplementation, Node[*ast.FuncType]{}, true}
	//compiler.GlobalFunctions[ValueKey{libFolder, "CoreRelationship"}] = compiler.libCoreRelationshipImplementation
	//compiler.GlobalFunctions[ValueKey{libFolder, "Relationship"}] = compiler.libRelationshipImplementation
}

func (compiler *Compiler) libQueryImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(typeParams) != 1 {
			panic(fmt.Errorf("Lib.Query implementation requires 1 type argument, got %d", len(typeParams)))
		}

		typeMapper := typeParams[funcTypeNode.Node.TypeParams.List[0].Names[0].Name]
		alias := compiler.AddEntitySource(typeMapper)
		qt := &TrailSource{state.currentNode.Node.Pos(), alias, typeMapper}
		returnValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, qt)
		return []Node[ast.Node]{returnValue}, artValue
	}
}

func (compiler *Compiler) executeFuncLit(state State, funcLit Node[*ast.FuncLit], arguments []Node[ast.Node], typeParams map[string]ITypeMapper) ([]Node[ast.Node], CallArrayResultType) {
	nameAndParams := findAllParamNameAndTypes(ChangeParamNode(funcLit, funcLit.Node.Type.Params))
	mm := ValueInformationMap{}
	for i, param := range nameAndParams {
		mm[param.name] = ValueInformation{arguments[i]}
	}

	newContext := &CurrentContext{
		mm,
		map[string]ITypeMapper{},
		LocalTypesMap{},
		GetCompilerState[*CurrentContext](state),
	}
	state = SetCompilerState(newContext, state)
	param := ChangeParamNode[*ast.FuncLit, *ast.BlockStmt](funcLit, funcLit.Node.Body)
	values, _ := compiler.executeBlockStmt(state, param, typeParams, arguments)
	state = SetCompilerState(newContext.Parent, state)
	return values, artValue
}

func (compiler *Compiler) libMapImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 2 {
			panic(fmt.Errorf("Lib.Map implementation requires 2 arguments, got %d", len(arguments)))
		}
		if _, ok := arguments[1].Node.(*ast.FuncLit); !ok {
			panic("map implementation requires function literal")
		}
		if funcLit, ok := arguments[1].Node.(*ast.FuncLit); ok {
			return compiler.executeFuncLit(state, ChangeParamNode(arguments[1], funcLit), arguments, typeParams)
		}
		panic("map implementation argument 1 is not a function literal")

	}
}

func (compiler *Compiler) libGenerateSql(state State, argument Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	switch item := argument.Node.(type) {
	case *TrailRecord:
		s := compiler.trailRecordToSql(item)
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
			rt := typeParams["TData"].ActualType()
			value = value.Convert(rt)
			sd := SomeDataWithRv{rt, value.Interface(), true}
			param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(sd)})
			return []Node[ast.Node]{param}, artValue
		} else {
			sd := SomeDataWithNode{arguments[0], true}
			param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(sd)})
			return []Node[ast.Node]{param}, artValue
		}
	}
}

func (compiler *Compiler) libSetSomeNoneImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		sd := SomeDataWithNode{assigned: false}
		param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(sd)})
		return []Node[ast.Node]{param}, artValue
	}
}

func (compiler *Compiler) IsSomeAssigned(state State, argument Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	switch nodeItem := argument.Node.(type) {
	case *ReflectValueExpression:
		unk := nodeItem.Rv.Interface()
		switch rvInstance := unk.(type) {
		case SomeDataWithNode:
			param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(rvInstance.assigned)})
			return []Node[ast.Node]{param}, artValue
		case SomeDataWithRv:
			param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(rvInstance.assigned)})
			return []Node[ast.Node]{param}, artValue
		default:
			panic("ddddd")
		}
	default:
		param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &CheckForNotNullExpression{argument})
		return []Node[ast.Node]{param}, artValue
	}
}

func (compiler *Compiler) libIsSomeAssignedImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	// Todo: do some optimize when arguments[0] is a literal
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
		return compiler.getGetSomeDataNCompiled(state, arguments)
	}
}

func (compiler *Compiler) getGetSomeDataNCompiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
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
		return ChangeParamNode[ast.Node, ast.Node](state.currentNode, &MultiBinaryExpr{token.LAND, binaryOperations})
	}
	return append(compiledArguments, fn()), artValue
}

func (compiler *Compiler) libGetSomeData02Implementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.getGetSomeDataNCompiled(state, arguments)
	}
}

func (compiler *Compiler) libGetSomeData03Implementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.getGetSomeDataNCompiled(state, arguments)
	}
}

func (compiler *Compiler) libGetSomeData04Implementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.getGetSomeDataNCompiled(state, arguments)
	}
}

func (compiler *Compiler) libGetSomeData05Implementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.getGetSomeDataNCompiled(state, arguments)
	}
}

func (compiler *Compiler) libCreateDictionaryImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 2 {
			panic(fmt.Errorf("CreateDictionary implementation requires 2 arguments, got %d", len(arguments)))
		}

		nameAndParams := findAllParamNameAndTypes(ChangeParamNode(funcTypeNode, funcTypeNode.Node.TypeParams))
		key := nameAndParams[0]
		value := nameAndParams[1]
		if rv00, ok00 := isLiterateValue(arguments[0]); ok00 {
			if rv01, ok01 := isLiterateValue(arguments[1]); ok01 {
				return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](
					state.currentNode,
					&DictionaryExpression{
						rv00,
						rv01,
						key.name,
						value.name,
						&WrapReflectTypeInMapper{rv00.Type()},
						typeParams[key.name],
						typeParams[value.name],
					})}, artValue
			}
		}
		panic(fmt.Errorf("createDictionary implementation requires literal values"))
	}
}

type rvArraySorter struct {
	rvArray []reflect.Value
}

func (rvArray *rvArraySorter) Len() int {
	return len(rvArray.rvArray)
}

func (rvArray *rvArraySorter) Less(i, j int) bool {
	if rvArray.rvArray[i].Kind() == rvArray.rvArray[j].Kind() {
		ith := rvArray.rvArray[i]
		jth := rvArray.rvArray[j]
		switch {
		case ith.CanInt():
			return ith.Int() < jth.Int()
		case ith.CanFloat():
			return ith.Float() < jth.Float()
		case ith.Kind() == reflect.String:
			return strings.Compare(ith.String(), jth.String()) < 0
		default:
			return false
		}
	}
	return false
}

func (rvArray *rvArraySorter) Swap(i, j int) {
	rvArray.rvArray[i], rvArray.rvArray[j] = rvArray.rvArray[j], rvArray.rvArray[i]
}

func (compiler *Compiler) libDictionaryLookupImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return libDictionaryLookupImplementation{compiler, state}.ExecuteStatement()
}

func (compiler *Compiler) libDictionaryDefaultImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("DictionaryLookup implementation requires 1 arguments, got %d", len(arguments)))
		}
		dictionaryExpression := arguments[0].Node.(*DictionaryExpression)
		rve := &ReflectValueExpression{dictionaryExpression.defaultValue}
		resultValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, rve)
		return []Node[ast.Node]{resultValue}, artReturn
	}
}

func findAllParamNameAndTypes(node Node[*ast.FieldList]) []struct {
	name string
	node Node[ast.Node]
} {
	result := make([]struct {
		name string
		node Node[ast.Node]
	}, 0, node.Node.NumFields())

	if node.Node != nil {
		for _, g := range node.Node.List {
			if len(g.Names) == 0 {
				result = append(result, struct {
					name string
					node Node[ast.Node]
				}{name: "_", node: ChangeParamNode[*ast.FieldList, ast.Node](node, g.Type)})

			} else {
				for _, n := range g.Names {
					result = append(result, struct {
						name string
						node Node[ast.Node]
					}{name: n.Name, node: ChangeParamNode[*ast.FieldList, ast.Node](node, g.Type)})
				}
			}
		}
	}
	return result
}

//func (compiler *Compiler) libCoreRelationshipImplementation(state State, parentNode Node[ast.Node], arguments []ast.Expr, typeParams []Node[ast.Node]) ExecuteStatement {
//	arguments := compiler.compileArguments(state, parentNode, arguments)
//	typeParams = func() []Node[ast.Node] {
//		if len(typeParams) == 0 && len(arguments) == 2 {
//			if ft, ok := arguments[1].Node.(*ast.FuncLit); ok {
//				allParams := findAllParamNameAndTypes(ft.Type.Params)
//				param := ChangeParamNode(arguments[1], allParams[0].node)
//				return []Node[ast.Node]{param}
//			} else {
//				panic(fmt.Errorf("calculation of relationship typeParams fails"))
//			}
//		}
//		return typeParams
//	}()
//
//	if len(typeParams) != 1 {
//		panic(fmt.Errorf("relationship implementation requires 1 typeParams, got %d", len(typeParams)))
//	}
//
//	if len(arguments) != 2 {
//		panic(fmt.Errorf("relationship implementation requires 2 arguments, got %d", len(arguments)))
//	}
//
//	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
//		// todo: do some more here, register the callack argument[1] someshere
//		return arguments[0:1], artValue
//	}
//}

//func (compiler *Compiler) libRelationshipImplementation(state State, parentNode Node[ast.Node], arguments []ast.Expr, typeParams []Node[ast.Node]) ExecuteStatement {
//	arguments := compiler.compileArguments(state, parentNode, arguments)
//	typeParams = func() []Node[ast.Node] {
//		if len(typeParams) == 0 && len(arguments) == 1 {
//			if ft, ok := arguments[0].Node.(*ast.FuncLit); ok {
//				allParams := findAllParamNameAndTypes(ft.Type.Params)
//				param := ChangeParamNode(arguments[0], allParams[0].node)
//				return []Node[ast.Node]{param}
//			} else {
//				panic(fmt.Errorf("calculation of relationship typeParams fails"))
//			}
//		}
//		return typeParams
//	}()
//
//	if len(typeParams) != 1 {
//		panic(fmt.Errorf("relationship implementation requires 1 typeParams, got %d", len(typeParams)))
//	}
//
//	if len(arguments) != 1 {
//		panic(fmt.Errorf("relationship implementation requires 1 arguments, got %d", len(arguments)))
//	}
//
//	es := compiler.libQueryImplementation(state, parentNode, arguments, typeParams)
//	arr, _ := compiler.executeAndExpandStatement(state, es)
//	return compiler.libCoreRelationshipImplementation(state, typeParams, append(arr, arguments...))
//}

//func (compiler *Compiler) libTypeForImplementation(state State) ExecuteStatement {
//	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
//		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(typeParams[0].ActualType())})}, artNone
//	}
//}

//func (compiler *Compiler) libTestTypeImplementation(state State) ExecuteStatement {
//	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
//		nodes := append(typeParams.toNodeArray(), arguments...)
//		rv := reflect.ValueOf(TestType)
//		if outputNodes, art, b := compiler.genericCall(state, rv, nodes); b {
//			return outputNodes, art
//		}
//		panic("fsdds")
//	}
//}
