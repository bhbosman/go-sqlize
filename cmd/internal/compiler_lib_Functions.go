package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const libFolder = "github.com/bhbosman/go-sqlize/lib"

func (compiler *Compiler) addLibFunctions() {
	compiler.GlobalFunctions[ValueKey{libFolder, "Query"}] = compiler.libQueryImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "Map"}] = compiler.libMapImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "GenerateSql"}] = compiler.libGenerateSqlImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "GenerateSqlTest"}] = compiler.libGenerateSqlTestImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "Atoi"}] = compiler.libAtoiImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "SetSomeValue"}] = compiler.libSetSomeValueImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "SetSomeNone"}] = compiler.libSetSomeNoneImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "IsSomeAssigned"}] = compiler.libIsSomeAssignedImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "SomeData"}] = compiler.libSomeDataImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "SomeData2"}] = compiler.libSomeData2Implementation
	compiler.GlobalFunctions[ValueKey{libFolder, "GetSomeData"}] = compiler.libGetSomeDataImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "GetSomeData02"}] = compiler.libGetSomeData02Implementation
	compiler.GlobalFunctions[ValueKey{libFolder, "GetSomeData03"}] = compiler.libGetSomeData03Implementation
	compiler.GlobalFunctions[ValueKey{libFolder, "GetSomeData04"}] = compiler.libGetSomeData04Implementation
	compiler.GlobalFunctions[ValueKey{libFolder, "GetSomeData05"}] = compiler.libGetSomeData05Implementation
	compiler.GlobalFunctions[ValueKey{libFolder, "CreateDictionary"}] = compiler.libCreateDictionaryImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "DictionaryLookup"}] = compiler.libDictionaryLookupImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "DictionaryDefault"}] = compiler.libDictionaryDefaultImplementation
}

func (compiler *Compiler) libQueryImplementation(_ State, typeParams []Node[ast.Expr], _ []Node[ast.Node]) ExecuteStatement {
	if len(typeParams) != 1 {
		panic(fmt.Errorf("Lib.Query implementation requires 1 type argument, got %d", len(typeParams)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var arr []reflect.Type
		for _, expr := range typeParams {
			typeMapper := compiler.findType(state, expr)
			rt := typeMapper.NodeType(state)
			arr = append(arr, rt)
		}
		rt := arr[0]
		alias := compiler.AddEntitySource(rt)
		qt := &TrailSource{state.currentNode.Node.Pos(), alias}
		returnValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, qt)
		return []Node[ast.Node]{returnValue}, artValue
	}
}

func (compiler *Compiler) libMapImplementation(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 2 {
		panic(fmt.Errorf("Lib.Map implementation requires 2 arguments, got %d", len(arguments)))
	}
	if _, ok := arguments[1].Node.(*ast.FuncLit); !ok {
		panic("map implementation requires function literal")
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		funcLit, _ := arguments[1].Node.(*ast.FuncLit)
		var names []*ast.Ident
		if funcLit.Type.Params != nil {
			for _, field := range funcLit.Type.Params.List {
				names = append(names, field.Names...)
			}
		}
		m := map[string]Node[ast.Node]{}
		m[names[0].Name] = arguments[0]
		newContext := &CurrentContext{m, GetCompilerState[*CurrentContext](state)}
		state = SetCompilerState(newContext, state)
		param := ChangeParamNode[ast.Node, *ast.BlockStmt](state.currentNode, funcLit.Body)
		values, _ := compiler.executeBlockStmt(state, param)
		state = SetCompilerState(newContext.Parent, state)
		return values, artValue
	}
}

func (compiler *Compiler) libGenerateSqlImplementation(state State, i []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("Lib.GenerateSql implementation requires 1 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		switch item := arguments[0].Node.(type) {
		case *TrailRecord:
			s := compiler.trailRecordToSql(item)
			basicLit := &ast.BasicLit{state.currentNode.Node.Pos(), token.STRING, strconv.Quote(s)}
			nodeValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, basicLit)
			return []Node[ast.Node]{nodeValue}, artValue
		default:
			panic("implementation required")
		}
	}
}

func (compiler *Compiler) libGenerateSqlTestImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("Lib.GenerateSqlTest implementation requires 1 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		args := arguments[0:1]
		ans, _ := compiler.libGenerateSqlImplementation(state, params, args)(state)
		if rv, isLiterate := isLiterateValue(ans[0]); isLiterate && rv.Kind() == reflect.String {
			currentContext := GetCompilerState[*CurrentContext](state)
			if value, b := currentContext.FindValue("__stdOut__"); b {
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

func (compiler *Compiler) libAtoiImplementation(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("Lib.Atoi implementation requires 1 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		nodes, resultType := compiler.strconvAtoiImplementation(state, nil, arguments)(state)
		return nodes[:1], resultType
	}
}

func (compiler *Compiler) libSetSomeValueImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("SetSomeValue implementation requires 1 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		sd := SomeDataWithNode{arguments[0], true}
		param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(sd)})
		return []Node[ast.Node]{param}, artValue
	}
}

func (compiler *Compiler) libSetSomeNoneImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		sd := SomeDataWithNode{assigned: false}
		param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(sd)})
		return []Node[ast.Node]{param}, artValue
	}
}

func (compiler *Compiler) libIsSomeAssignedImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("IsSomeAssigned implementation requires 1 arguments, got %d", len(arguments)))
	}
	// Todo: do some optimize when arguments[0] is a literal
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		switch nodeItem := arguments[0].Node.(type) {
		case *ReflectValueExpression:
			unk := nodeItem.Rv.Interface()
			switch rvInstance := unk.(type) {
			case SomeDataWithNode:
				param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(rvInstance.assigned)})
				return []Node[ast.Node]{param}, artValue
			default:
				panic("ddddd")
			}
		default:
			param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &CheckForNotNullExpression{arguments[0]})
			return []Node[ast.Node]{param}, artValue
		}
	}
}

func (compiler *Compiler) libSomeDataImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("SomeData implementation requires 1 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		return arguments[0:1], artValue
	}
}

func (compiler *Compiler) libSomeData2Implementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("SomeData2 implementation requires 1 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		v, _ := compiler.libIsSomeAssignedImplementation(state, params, arguments)(state)
		result := append(arguments, v[0])
		return result, artValue
	}
}

func (compiler *Compiler) libGetSomeDataImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("GetSomeData implementation requires 1 arguments, got %d", len(arguments)))
	}
	return compiler.getGetSomeDataN(state, params, arguments)
}

func (compiler *Compiler) getGetSomeDataN(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		fn := func() Node[ast.Node] {
			var binaryOperations []Node[ast.Node]
			for _, arg := range arguments {
				arr, _ := compiler.libIsSomeAssignedImplementation(state, params, []Node[ast.Node]{arg})(state)
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
		return append(arguments, fn()), artValue
	}
}

func (compiler *Compiler) libGetSomeData02Implementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 2 {
		panic(fmt.Errorf("GetSomeData02 implementation requires 2 arguments, got %d", len(arguments)))
	}
	return compiler.getGetSomeDataN(state, params, arguments)
}

func (compiler *Compiler) libGetSomeData03Implementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 3 {
		panic(fmt.Errorf("GetSomeData03 implementation requires 3 arguments, got %d", len(arguments)))
	}
	return compiler.getGetSomeDataN(state, params, arguments)
}

func (compiler *Compiler) libGetSomeData04Implementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 4 {
		panic(fmt.Errorf("GetSomeData04 implementation requires 4 arguments, got %d", len(arguments)))
	}
	return compiler.getGetSomeDataN(state, params, arguments)
}

func (compiler *Compiler) libGetSomeData05Implementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 5 {
		panic(fmt.Errorf("GetSomeData05 implementation requires 5 arguments, got %d", len(arguments)))
	}
	return compiler.getGetSomeDataN(state, params, arguments)
}

func (compiler *Compiler) libCreateDictionaryImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 2 {
		panic(fmt.Errorf("CreateDictionary implementation requires 2 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		if rv00, ok00 := isLiterateValue(arguments[0]); ok00 {
			if rv01, ok01 := isLiterateValue(arguments[1]); ok01 {
				return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &DictionaryExpression{rv00, rv01})}, artValue
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

func (compiler *Compiler) libDictionaryLookupImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 2 {
		panic(fmt.Errorf("DictionaryLookup implementation requires 2 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var conditionalStatement []SingleValueCondition
		dictionaryExpression := arguments[0].Node.(*DictionaryExpression)
		{
			inputData := arguments[1]
			rvMap := dictionaryExpression.m
			keyArr := rvMap.MapKeys()
			sorter := &rvArraySorter{keyArr}
			sort.Sort(sorter)

			for _, rvKey := range keyArr {
				rvValue := rvMap.MapIndex(rvKey)
				fn := func() Node[ast.Node] {
					switch {
					case rvKey.CanFloat() || rvKey.CanInt() || rvKey.Kind() == reflect.String:
						left := inputData
						right := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rvKey})
						be := &BinaryExpr{token.NoPos, token.EQL, left, right}
						return ChangeParamNode[ast.Node, ast.Node](state.currentNode, be)
					case rvKey.Kind() == reflect.Struct:
						switch leftItem := inputData.Node.(type) {
						case *TrailRecord:
							if leftItem.Value.NumField() == rvKey.NumField() {
								var expressions []Node[ast.Node]
								for idx := 0; idx < rvKey.NumField(); idx++ {
									left := leftItem.Value.Field(idx).Interface().(Node[ast.Node])
									right := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rvKey.Field(idx)})
									be := &BinaryExpr{token.NoPos, token.EQL, left, right}
									expressions = append(expressions, ChangeParamNode[ast.Node, ast.Node](state.currentNode, be))
								}
								mbe := &MultiBinaryExpr{token.LAND, expressions}
								return ChangeParamNode[ast.Node, ast.Node](state.currentNode, mbe)
							}
						}
						panic("sdsfdsfd")
					default:
						panic("find out")
					}
				}
				condition := fn()
				singleValueCondition := SingleValueCondition{condition: condition, value: ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rvValue})}
				conditionalStatement = append(conditionalStatement, singleValueCondition)
			}
		}
		{
			rvDefault := dictionaryExpression.defaultValue
			condition := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(true)})
			singleValueCondition := SingleValueCondition{condition: condition, value: ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rvDefault})}
			conditionalStatement = append(conditionalStatement, singleValueCondition)
		}
		ite := &IfThenElseSingleValueCondition{conditionalStatement}
		resultValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, ite)
		return []Node[ast.Node]{resultValue}, artReturn
	}
}

func (compiler *Compiler) libDictionaryDefaultImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("DictionaryLookup implementation requires 2 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		dictionaryExpression := arguments[0].Node.(*DictionaryExpression)
		rve := &ReflectValueExpression{dictionaryExpression.defaultValue}
		resultValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, rve)
		return []Node[ast.Node]{resultValue}, artReturn
	}
}
