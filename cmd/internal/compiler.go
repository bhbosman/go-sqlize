package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"reflect"
	"strings"
)

type GlobalMethodHandlerKey struct {
	rt         reflect.Type
	methodName string
}

type CompilerState int

type CurrentCompositeCreateType struct {
	typeMapper ITypeMapper
}

const (
	CompilerState_InitCalled CompilerState = 1 << iota
)

type ExecuteStatement func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType)

type AssignStatement func(state State, value Node[ast.Node])

type OnCreateExecuteStatement func(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement

type OnCreateType func(State, []ITypeMapper) ITypeMapper

type functionInformation struct {
	fn               OnCreateExecuteStatement
	funcType         Node[*ast.FuncType]
	funcTypeRequired bool
}
type Compiler struct {
	CompilerState   CompilerState
	InitFunctions   []Node[*ast.FuncDecl]
	GlobalFunctions map[ValueKey]functionInformation
	GlobalTypes     map[ValueKey]OnCreateType
	Sources         map[string]interface{}
	NextAlias       int
	Fileset         *token.FileSet
}

func (compiler *Compiler) Init(
	FunctionMap FunctionMap,
	StructMethodMap StructMethodMap,
	TypeSpecMap TypeSpecMap,
	InitFunctions []Node[*ast.FuncDecl],
	Fileset *token.FileSet,
) {

	compiler.CompilerState |= CompilerState_InitCalled
	compiler.Fileset = Fileset
	compiler.Sources = map[string]interface{}{}
	compiler.GlobalTypes = map[ValueKey]OnCreateType{
		ValueKey{"", "bool"}:        compiler.registerBool(),
		ValueKey{"", "int"}:         compiler.registerInt(),
		ValueKey{"", "string"}:      compiler.registerString(),
		ValueKey{"", "float64"}:     compiler.registerFloat64(),
		ValueKey{"reflect", "Type"}: compiler.registerReflectType(),
	}

	compiler.GlobalFunctions = map[ValueKey]functionInformation{
		ValueKey{"path/filepath", "Join"}: {compiler.pathFilepathJoinImplementation, Node[*ast.FuncType]{}, false},
		ValueKey{"path/filepath", "Dir"}:  {compiler.pathFilepathDirImplementation, Node[*ast.FuncType]{}, false},
	}
	compiler.addBuiltInFunctions()
	compiler.addOsFunctions()
	compiler.addIoFunctions()
	compiler.addStrconvFunctions()
	compiler.addMathFunctions()
	compiler.addReflectFunctions()
	compiler.addLibFunctions()

	for key, value := range FunctionMap {
		if current, ok := compiler.GlobalFunctions[key]; !ok {
			fn := func(vk ValueKey, node Node[*ast.FuncDecl]) functionInformation {
				fn := func(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
					return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
						funcLit := &ast.FuncLit{node.Node.Type, node.Node.Body}
						param := ChangeParamNode[*ast.FuncDecl, *ast.FuncLit](node, funcLit)
						onCreateExecuteStatement := compiler.onFuncLitExecutionStatement(param)
						executeStatement := onCreateExecuteStatement(state, funcTypeNode)
						return executeStatement(state, typeParams, unprocessedArgs)
					}
				}
				return functionInformation{fn, ChangeParamNode(value, value.Node.Type), true}
			}
			compiler.GlobalFunctions[key] = fn(key, value)
		} else {
			compiler.GlobalFunctions[key] = functionInformation{current.fn, ChangeParamNode(value, value.Node.Type), current.funcTypeRequired}
		}
	}
	for key, information := range compiler.GlobalFunctions {
		if information.funcTypeRequired {
			if information.funcType.Valid && information.fn != nil {
				continue
			}
			_, _ = fmt.Fprintf(os.Stderr, "not fully completed %v", key)
		}
	}

	for key, typeSpecNode := range TypeSpecMap {
		if key.Folder == libFolder && key.Key == "Some" {
			continue
		}
		compiler.GlobalTypes[key] = compiler.readTypeSpec(typeSpecNode)
	}
	compiler.InitFunctions = InitFunctions
}

func (compiler *Compiler) readTypeSpec(node Node[*ast.TypeSpec]) OnCreateType {
	return func(state State, expressions []ITypeMapper) ITypeMapper {
		if node.Node.TypeParams == nil || node.Node.TypeParams.NumFields() == len(expressions) {
			var dd []*ast.Ident
			if node.Node.TypeParams != nil {
				for _, field := range node.Node.TypeParams.List {
					dd = append(dd, field.Names...)
				}
			}
			typeMapper := TypeMapper{}
			for i := 0; i < len(dd); i++ {
				typeMapper[dd[i].Name] = expressions[i]
			}
			state = SetCompilerState[TypeMapper](typeMapper, state)

			switch v := node.Node.Type.(type) {
			case *ast.StructType:
				param := ChangeParamNode(node, v)
				return compiler.createStructTypeMapper(state, param)
			}
		}
		panic("sdfdsfs")
	}
}

func (compiler *Compiler) genericCall(state State, rv reflect.Value, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType, bool) {
	input, allLiterals := compiler.nodesToValues(rv, arguments)
	if allLiterals {
		output := rv.Call(input)
		return compiler.valuesToNodes(state, output), artValue, true
	}
	return nil, artNone, false
}

func (compiler *Compiler) Compile(currentContext *CurrentContext, fileNames ...string) {
	m := map[string]bool{}
	for _, fileName := range fileNames {
		m[fileName] = true
	}
	for _, initFunction := range compiler.InitFunctions {
		if _, ok := m[initFunction.FileName]; len(fileNames) == 0 || ok {
			currentNode := ChangeParamNode[*ast.FuncDecl, ast.Node](initFunction, initFunction.Node)
			compiler.CompileFunc(
				State{
					[]IABC{&CurrentContext{ValueInformationMap{}, map[string]ITypeMapper{}, LocalTypesMap{}, currentContext}},
					currentNode}, initFunction)
		}
	}
}

type CallArrayResultType int

const (
	artNone CallArrayResultType = 1 << iota
	//artPartialReturn
	artReturn
	artValue
	artFCI
	artReturnAndContinue
)

func (compiler *Compiler) CompileFunc(state State, fn Node[*ast.FuncDecl]) ([]Node[ast.Node], CallArrayResultType) {
	param := ChangeParamNode(fn, fn.Node.Body)
	return compiler.executeBlockStmt(state, param, nil, nil)
}

func (compiler *Compiler) AddEntitySource(rt ITypeMapper) string {
	compiler.NextAlias++
	reference := fmt.Sprintf("T%v", compiler.NextAlias)
	compiler.Sources[reference] = &EntitySource{rt}
	return reference
}

func (compiler *Compiler) calculateSourceDependency(sources []string) []string {
	return sources
}

func (compiler *Compiler) projectSources(w io.Writer, tabCount int, sources []string) {
	_, _ = io.WriteString(w, strings.Repeat("\t", tabCount))
	for _, source := range sources {
		query, _ := compiler.Sources[source]
		switch item := query.(type) {
		case *EntitySource:
			rt := item.rt.NodeType()
			_, _ = io.WriteString(w, fmt.Sprintf("%v [%v]", rt.String(), source))
		default:
			panic("unhandled default case")
		}
	}
}

func (compiler *Compiler) nodesToValues(rvFunc reflect.Value, nodes []Node[ast.Node]) ([]reflect.Value, bool) {
	funcRt := rvFunc.Type()
	var arr []reflect.Value
	for idx, node := range nodes {
		rv, b := compiler.nodeToValue(node)
		if !b {
			return nil, false
		}
		rt := func(funcRt reflect.Type, idx int) reflect.Type {
			if funcRt.IsVariadic() && idx >= funcRt.NumIn()-1 {
				idx = min(idx, funcRt.NumIn()-1)
				return funcRt.In(idx).Elem()
			}
			return funcRt.In(idx)
		}(funcRt, idx)
		arr = append(arr, rv.Convert(rt))
	}
	return arr, true
}

func (compiler *Compiler) nodeToValue(node Node[ast.Node]) (reflect.Value, bool) {
	if value, isLiterateValue := isLiterateValue(node); isLiterateValue {
		return value, true
	}
	return reflect.Value{}, false
}

func (compiler *Compiler) valuesToNodes(state State, values []reflect.Value) []Node[ast.Node] {
	var arr []Node[ast.Node]
	for _, node := range values {
		arr = append(arr, compiler.valueToNode(node))
	}
	return arr
}

func (compiler *Compiler) valueToNode(value reflect.Value) Node[ast.Node] {
	kind := value.Kind()
	switch kind {
	case reflect.String:
		return Node[ast.Node]{Valid: true, Node: &ReflectValueExpression{value, stringValueKey}}
	case reflect.Int:
		return Node[ast.Node]{Valid: true, Node: &ReflectValueExpression{value, intValueKey}}
	case reflect.Float64:
		return Node[ast.Node]{Valid: true, Node: &ReflectValueExpression{value, float64ValueKey}}
	default:
		panic("unhandled default case")
	}
}

func (compiler *Compiler) executeAndExpandStatement(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node], executeStatement ExecuteStatement) ([]Node[ast.Node], CallArrayResultType) {
	var result []Node[ast.Node]
	arr, v := executeStatement(state, typeParams, unprocessedArgs)
	for _, instance := range arr {
		if expand, ok := instance.Node.(IExpand); ok {
			ex := expand.Expand(state.currentNode)
			result = append(result, ex...)
		} else {
			result = append(result, instance)
		}
	}
	return result, v
}

func (compiler *Compiler) expandNodeWithSelector(node Node[ast.Node], sel *ast.Ident) (Node[ast.Node], bool) {
	switch nodeItem := node.Node.(type) {
	case *TrailRecord:
		return nodeItem.Value.FieldByName(sel.Name).Interface().(Node[ast.Node]), true
	case *ReflectValueExpression:
		if nodeItem.Rv.Kind() == reflect.Struct {
			return nodeItem.Rv.FieldByName(sel.Name).Interface().(Node[ast.Node]), true
		}

	case *IfThenElseSingleValueCondition:
		var singleValueConditions []SingleValueCondition
		for _, conditionalStatement := range nodeItem.conditionalStatement {

			if rve, ok00 := conditionalStatement.value.Node.(*TrailRecord); ok00 && rve.Value.Kind() == reflect.Struct {
				value := rve.Value.FieldByName(sel.Name).Interface().(Node[ast.Node])
				singleValueCondition := SingleValueCondition{conditionalStatement.condition, value}
				singleValueConditions = append(singleValueConditions, singleValueCondition)
			} else if _, ok := conditionalStatement.value.Node.(*IfThenElseSingleValueCondition); ok {
				return node, false
			}
		}
		ifThenElseSingleValueCondition := &IfThenElseSingleValueCondition{singleValueConditions}
		return ChangeParamNode[ast.Node, ast.Node](node, ifThenElseSingleValueCondition), true
	}
	return node, false
}

type FieldInformation struct {
	Name string
	Type Node[ast.Node]
}

func (compiler *Compiler) createStructTypeMapper(state State, node Node[*ast.StructType]) ITypeMapper {
	type StructTypeToTypeUsage int
	const (
		StructTypeWithNodeType StructTypeToTypeUsage = iota
		StructTypeWithTypeMapper
		StructTypeWithActualTypes
	)
	structTypeToType := func(List []FieldInformation, useActual StructTypeToTypeUsage) reflect.Type {
		fieldTypeFn := func(Type Node[ast.Node], useActual StructTypeToTypeUsage) reflect.Type {
			switch useActual {
			case StructTypeWithNodeType:
				return reflect.TypeFor[Node[ast.Node]]()
			case StructTypeWithTypeMapper:
				return reflect.TypeFor[ITypeMapper]()
			case StructTypeWithActualTypes:
				//param := ChangeParamNode[ast.Node, ast.Node](state.currentNode, Type)
				typeMapper := compiler.findType(state, Type, Default)
				typ, _ := typeMapper.ActualType()
				return typ
			default:
				panic("fsdfds")
			}
		}
		var structFields []reflect.StructField
		for _, field := range List {
			fieldType := fieldTypeFn(field.Type, useActual)

			structField := reflect.StructField{
				Name: field.Name,
				PkgPath: func() string {
					switch token.IsExported(field.Name) {
					case true:
						return ""
					default:
						return "PkgPath" // required for unexported items
					}
				}(),
				Type:      fieldType,
				Tag:       reflect.StructTag(""),
				Offset:    0,
				Index:     nil,
				Anonymous: false,
			}
			structFields = append(structFields, structField)

		}
		return reflect.StructOf(structFields)
	}

	fn := func(node Node[*ast.StructType], arr []struct {
		name string
		node Node[ast.Node]
	}) []FieldInformation {
		var result []FieldInformation
		for _, ss := range arr {
			field := FieldInformation{ss.name, ss.node}
			result = append(result, field)
		}
		return result
	}
	fieldList := fn(node, findAllParamNameAndTypes(ChangeParamNode(node, node.Node.Fields)))
	nodeRt := structTypeToType(fieldList, StructTypeWithNodeType)
	rtWithITypeMapper := structTypeToType(fieldList, StructTypeWithTypeMapper)
	actualTypeRt := structTypeToType(fieldList, StructTypeWithActualTypes)
	typeMapperInstance := reflect.New(rtWithITypeMapper).Elem()
	for _, field := range fieldList {
		param := field.Type
		fieldType := compiler.findType(state, param, Default)
		typeMapperInstance.FieldByName(field.Name).Set(reflect.ValueOf(fieldType))
	}
	return &TypeMapperForStruct{nodeRt, actualTypeRt, typeMapperInstance, node.Key}
}

func (compiler *Compiler) builtInStructMethods(rv reflect.Value) OnCreateExecuteStatement {
	return func(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
		return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
			if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
				return outputNodes, art
			}
			panic(fmt.Errorf("builtInStructMethods only accept literal values"))
		}
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
