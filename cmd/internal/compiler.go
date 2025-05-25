package internal

import (
	"errors"
	"fmt"
	"github.com/dominikbraun/graph"
	"go/ast"
	"go/token"
	"os"
	"reflect"
	"sort"
)

type IIsLiterateValue interface {
	ThisIsALiterateValue()
}
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

type ExecuteStatement func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType)

type AssignStatement func(state State, value Node[ast.Node])

type OnCreateExecuteStatement func(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement

type OnCreateType func(State, []ITypeMapper) ITypeMapper

type ISourceType interface {
	sourceType()
}

type ITrailMarker interface {
	ast.Node
	trailMarker()
}

type joinType int

const (
	jtInner joinType = iota
)

type JoinInformation struct {
	lhs       string
	rhs       map[string]ISource
	condition Node[BooleanCondition]
	joinType  joinType
}

func (j JoinInformation) Dependencies() map[string]ISource {
	return j.rhs
}

func (j JoinInformation) SourceName() string {
	return j.lhs
}

type Compiler struct {
	CompilerState   CompilerState
	InitFunctions   []Node[*ast.FuncDecl]
	GlobalFunctions map[ValueKey]functionInformation
	GlobalTypes     map[ValueKey]OnCreateType
	Sources         map[string]ISourceType
	JoinInformation map[string]JoinInformation

	NextAlias int
	Fileset   *token.FileSet
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
	compiler.Sources = map[string]ISourceType{}
	compiler.GlobalTypes = map[ValueKey]OnCreateType{}
	compiler.JoinInformation = map[string]JoinInformation{}

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
						currentContext := GetCompilerState[*CurrentContext](state)
						flattenValues := currentContext.flattenVariables()

						p01 := ChangeParamNode(node, node.Node.Type)
						typeMapper := compiler.createTypeMapperForFuncType(state, p01)
						funcLit := FuncLit{node.Node.Type, node.Node.Body, flattenValues, typeMapper}
						param := ChangeParamNode[*ast.FuncDecl, ast.Node](node, funcLit)
						return []Node[ast.Node]{param}, artValue
						//onCreateExecuteStatement := compiler.onFuncLitExecutionStatement(param)
						//executeStatement := onCreateExecuteStatement(state, funcTypeNode)
						//return executeStatement(state, typeParams, unprocessedArgs)
					}
				}
				p01 := ChangeParamNode(value, value.Node.Type)
				return functionInformation{fn, p01, true}
			}
			compiler.GlobalFunctions[key] = fn(key, value)
		} else {
			p01 := ChangeParamNode(value, value.Node.Type)
			compiler.GlobalFunctions[key] = functionInformation{current.fn, p01, current.funcTypeRequired}
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
					[]IABC{
						&CurrentContext{ValueInformationMap{}, map[string]ITypeMapper{}, LocalTypesMap{}, false, currentContext},
					},
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
	return compiler.executeBlockStmt(state, param)
}

func (compiler *Compiler) AddEntitySource(rt ITypeMapper, qs queryState) string {
	compiler.NextAlias++
	reference := fmt.Sprintf("T%v", compiler.NextAlias)
	compiler.Sources[reference] = &EntitySource{rt, qs}
	return reference
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

	case IfThenElseSingleValueCondition:
		var singleValueConditions []SingleValueCondition
		for _, conditionalStatement := range nodeItem.conditionalStatement {

			if rve, ok00 := conditionalStatement.value.Node.(*TrailRecord); ok00 && rve.Value.Kind() == reflect.Struct {
				value := rve.Value.FieldByName(sel.Name).Interface().(Node[ast.Node])
				singleValueCondition := SingleValueCondition{conditionalStatement.condition, value}
				singleValueConditions = append(singleValueConditions, singleValueCondition)
			} else if _, ok := conditionalStatement.value.Node.(IfThenElseSingleValueCondition); ok {
				return node, false
			} else {
				panic("unhandled default case")
			}
		}
		ifThenElseSingleValueCondition := IfThenElseSingleValueCondition{singleValueConditions}
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

	fn := func(node Node[*ast.StructType], arr []ParamNameAndTypes) []FieldInformation {
		var result []FieldInformation
		for _, ss := range arr {
			field := FieldInformation{ss.name, ss.node}
			result = append(result, field)
		}
		return result
	}
	p0 := ChangeParamNode(node, node.Node.Fields)
	fieldList := fn(node, findAllParamNameAndTypes(p0).arr)
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

type ParamNameAndTypes struct {
	name string
	node Node[ast.Node]
}

type findAllParamNameAndTypesResult struct {
	arr        []ParamNameAndTypes
	isVariadic bool
}

func findAllParamNameAndTypes(node Node[*ast.FieldList]) findAllParamNameAndTypesResult {
	result := make([]ParamNameAndTypes, 0, node.Node.NumFields())
	isVariadic := false
	if node.Node != nil {
		for idx, g := range node.Node.List {
			if idx == node.Node.NumFields()-1 {
				_, isVariadic = g.Type.(*ast.Ellipsis)
			}
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
	return findAllParamNameAndTypesResult{result, isVariadic}
}

func (compiler *Compiler) executeFuncLit(state State, funcLit Node[FuncLit], arguments []Node[ast.Node], typeParams map[string]ITypeMapper) ([]Node[ast.Node], CallArrayResultType) {
	mm := ValueInformationMap{}
	if funcLit.Node.Type != nil {
		p0 := ChangeParamNode(funcLit, funcLit.Node.Type.Params)
		nameAndParamsResult := findAllParamNameAndTypes(p0)
		for i, param := range nameAndParamsResult.arr {
			mm[param.name] = ValueInformation{arguments[i]}
		}
	}
	newContext := &CurrentContext{mm, map[string]ITypeMapper{}, LocalTypesMap{}, false, GetCompilerState[*CurrentContext](state)}

	newRoot := &CurrentContext{funcLit.Node.values, map[string]ITypeMapper{}, LocalTypesMap{}, true, nil}
	var values []Node[ast.Node]
	newContext.ReplaceRoot(newRoot)
	{
		state = SetCompilerState(newContext, state)
		{
			param := ChangeParamNode[FuncLit, *ast.BlockStmt](funcLit, funcLit.Node.Body)
			values, _ = compiler.executeBlockStmt(state, param)
		}
		state = SetCompilerState(newContext.Parent, state)
	}
	newContext.RemoveRoot(newRoot)

	return values, artValue
}

func (compiler *Compiler) findAdditionalSourcesFromAssociations(sources map[string]ISource) map[string]ISource {
	return sources
}

func (compiler *Compiler) calculateSourcesOrder(sources map[string]ISource) []ISource {
	sourceGraph := graph.New[string, ISource](
		func(source ISource) string {
			return source.SourceName()
		},
		graph.Directed(),
		graph.Acyclic(),
		graph.PreventCycles())

	for sourceKey, sourceValue := range sources {
		err := sourceGraph.AddVertex(sourceValue)
		if err != nil {
			switch {
			default:
				panic(err)
			case errors.Is(err, graph.ErrVertexAlreadyExists):
				// do nothing
				break
			}
		}
		for dependencyKey, dependencyValue := range sourceValue.Dependencies() {
			err := sourceGraph.AddVertex(dependencyValue)
			if err != nil {
				switch {
				default:
					panic(err)
				case errors.Is(err, graph.ErrVertexAlreadyExists):
					// do nothing
					break
				}
			}
			err = sourceGraph.AddEdge(sourceKey, dependencyKey)
			if err != nil {
				switch {
				default:
					panic(err)
				case errors.Is(err, graph.ErrEdgeAlreadyExists):
					// do nothing
					break
				}
			}
		}
	}
	predecessorMap, err := sourceGraph.AdjacencyMap()
	if err != nil {
		panic(err)
	}
	ss := compiler.processGraphMap(nil, predecessorMap)
	var result []ISource
	for _, s := range ss {
		result = append(result, sources[s])
	}

	return result
}

func (compiler *Compiler) processGraphMap(ss []string, graphMap map[string]map[string]graph.Edge[string]) []string {
	findSourcesWithNoPredecessor := func(predecessorMap map[string]map[string]graph.Edge[string]) []string {
		var result []string
		for key, value := range predecessorMap {
			if len(value) == 0 {
				result = append(result, key)
			}
		}
		return result
	}

	removePredecessorFromMap := func(
		key string,
		predecessorMap map[string]map[string]graph.Edge[string]) map[string]map[string]graph.Edge[string] {
		delete(predecessorMap, key)
		for _, value := range predecessorMap {
			if _, ok := value[key]; ok {
				delete(value, key)
			}
		}
		return predecessorMap
	}

	noPredecessors := findSourcesWithNoPredecessor(graphMap)
	for _, predecessor := range noPredecessors {
		graphMap = removePredecessorFromMap(predecessor, graphMap)
	}
	sort.Strings(noPredecessors)

	if len(graphMap) > 0 {
		child := compiler.processGraphMap(noPredecessors, graphMap)
		ss = append(ss, child...)
	} else {
		ss = append(ss, noPredecessors...)
	}
	return ss
}
