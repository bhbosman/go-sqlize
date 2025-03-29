package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

const libFolder = "github.com/bhbosman/go-sqlize/lib"

var someTypeValueKey = ValueKey{libFolder, "Some"}

type GlobalMethodHandlerKey struct {
	rt         reflect.Type
	methodName string
}

type CompilerState int

const (
	CompilerState_InitCalled CompilerState = 1 << iota
)

type TypeMapper map[string]reflect.Type

type IABC interface{}

type State struct {
	arr         []IABC
	currentNode Node[ast.Node]
}

func (s State) setCurrentNode(node Node[ast.Node]) State {
	return State{s.arr, node}
}

type ExecuteStatement func(state State) ([]Node[ast.Node], CallArrayResultType)

type AssignStatement func(state State, value Node[ast.Node])

type OnCreateExecuteStatement func(State, []Node[ast.Expr], []Node[ast.Node]) ExecuteStatement

type OnCreateType func(State, []Node[ast.Expr]) reflect.Type

type Compiler struct {
	CompilerState   CompilerState
	InitFunctions   []Node[*ast.FuncDecl]
	GlobalFunctions map[ValueKey]OnCreateExecuteStatement
	GlobalTypes     map[ValueKey]OnCreateType
	Sources         map[string]interface{}
	NextAlias       int
}

func (compiler *Compiler) Init(
	FunctionMap FunctionMap,
	StructMethodMap StructMethodMap,
	TypeSpecMap TypeSpecMap,
	InitFunctions []Node[*ast.FuncDecl],
) {
	compiler.CompilerState |= CompilerState_InitCalled
	compiler.Sources = map[string]interface{}{}
	compiler.GlobalTypes = map[ValueKey]OnCreateType{
		ValueKey{"", "int"}:     compiler.registerInt,
		ValueKey{"", "string"}:  compiler.registerString,
		ValueKey{"", "float64"}: compiler.registerFloat64,
	}

	compiler.GlobalFunctions = map[ValueKey]OnCreateExecuteStatement{
		ValueKey{"", "float64"}:            compiler.coercionFloat64,
		ValueKey{"", "panic"}:              compiler.builtInPanic,
		ValueKey{libFolder, "Query"}:       compiler.libQueryImplementation,
		ValueKey{libFolder, "Map"}:         compiler.libMapImplementation,
		ValueKey{libFolder, "GenerateSql"}: compiler.libGenerateSqlImplementation,
		ValueKey{"os", "Getwd"}:            compiler.osGetWdImplementation,
	}
	for key, value := range FunctionMap {
		if _, ok := compiler.GlobalFunctions[key]; !ok {
			fn := func(vk ValueKey, node Node[*ast.FuncDecl]) OnCreateExecuteStatement {
				return func(states State, typeParams []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
					return func(state State) ([]Node[ast.Node], CallArrayResultType) {
						funcLit := &ast.FuncLit{node.Node.Type, node.Node.Body}
						param := ChangeParamNode[*ast.FuncDecl, *ast.FuncLit](node, funcLit)
						onCreateExecuteStatement := compiler.onFuncLitExecutionStatement(param)
						executeStatement := onCreateExecuteStatement(state, typeParams, arguments)
						return executeStatement(state)
					}
				}
			}
			compiler.GlobalFunctions[key] = fn(key, value)
		} else {
			println("d")
		}

	}

	for key, value := range TypeSpecMap {
		fn := func(vk ValueKey, node Node[*ast.TypeSpec]) OnCreateType {
			return func(states State, exprs []Node[ast.Expr]) reflect.Type {
				if node.Node.TypeParams == nil || node.Node.TypeParams.NumFields() == len(exprs) {
					var dd []*ast.Ident
					if node.Node.TypeParams != nil {
						for _, field := range node.Node.TypeParams.List {
							dd = append(dd, field.Names...)
						}
					}
					switch v := value.Node.Type.(type) {
					case *ast.StructType:
						if v.Fields != nil {
							var structFields []reflect.StructField
							for _, field := range v.Fields.List {
								for _, fieldName := range field.Names {
									structField := reflect.StructField{
										Name: fieldName.Name,
										PkgPath: func() string {
											switch token.IsExported(fieldName.Name) {
											case true:
												return ""
											default:
												return "PkgPath" // required for unexported items
											}
										}(),
										Type:      reflect.TypeFor[Node[ast.Node]](),
										Tag:       reflect.StructTag(""),
										Offset:    0,
										Index:     nil,
										Anonymous: false,
									}
									structFields = append(structFields, structField)
								}
							}
							states = RemoveCompilerState[TypeMapper](states)
							rt := reflect.StructOf(structFields)
							return rt
						}
					}
				}
				panic("sdfdsfs")
			}
		}
		compiler.GlobalTypes[key] = fn(key, value)
	}

	compiler.InitFunctions = InitFunctions
}

func (compiler *Compiler) libGenerateSqlImplementation(state State, i []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var output string
		switch item := arguments[1].Node.(type) {
		case *ast.BasicLit:
			output, _ = strconv.Unquote(item.Value)
		}
		wd, _ := os.Getwd()
		fileName := filepath.Join(wd, output)
		dir := filepath.Dir(fileName)
		_ = os.MkdirAll(dir, os.ModePerm)

		create, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer create.Close()

		switch item := arguments[0].Node.(type) {
		case *TrailRecord:
			sources := compiler.findSources(item)
			sources = compiler.calculateSourceDependency(sources)
			_, _ = fmt.Fprintf(create, "select\n")
			compiler.projectTrailRecord(create, 1, item)
			_, _ = fmt.Fprintf(create, "from\n")
			compiler.projectSources(create, 1, sources)
			return nil, artNone
		default:
			panic("implementation required")
		}
	}
}

func (compiler *Compiler) libMapImplementation(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 2 {
		panic(fmt.Errorf("map implementation requires 2 arguments, got %d", len(arguments)))
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
		state = SetCompilerState(newContext.parent, state)
		return values, artValue
	}
}

func (compiler *Compiler) registerFloat64(State, []Node[ast.Expr]) reflect.Type {
	return reflect.TypeFor[float64]()
}

func (compiler *Compiler) registerString(State, []Node[ast.Expr]) reflect.Type {
	return reflect.TypeFor[string]()
}

func (compiler *Compiler) registerInt(State, []Node[ast.Expr]) reflect.Type {
	return reflect.TypeFor[int]()
}

func (compiler *Compiler) builtInPanic(state State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("built-in panic requires 1 argument, got %d", len(arguments)))
	}

	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		switch arg := arguments[0].Node.(type) {
		case *ast.BasicLit:
			panic(fmt.Errorf(arg.Value))
		default:
			panic("implementation required")
		}
	}
}

func (compiler *Compiler) coercionFloat64(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
	}

	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		returnValue := ChangeParamNode[ast.Node, ast.Node](
			state.currentNode,
			&coercion{state.currentNode.Node.Pos(), "float64", arguments[0]},
		)
		return []Node[ast.Node]{returnValue}, artValue
	}
}

func (compiler *Compiler) libQueryImplementation(_ State, typeParams []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var arr []reflect.Type
		for _, expr := range typeParams {
			rt := compiler.findType(state, expr)
			arr = append(arr, rt)
		}
		rt := arr[0]
		alias := compiler.AddEntitySource(rt)
		qt := &TrailSource{state.currentNode.Node.Pos(), alias}
		returnValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, qt)
		return []Node[ast.Node]{returnValue}, artValue
	}
}

func (compiler *Compiler) Compile() {
	for _, initFunction := range compiler.InitFunctions {
		currentNode := ChangeParamNode[*ast.FuncDecl, ast.Node](initFunction, initFunction.Node)
		compiler.CompileFunc(
			State{
				[]IABC{&CurrentContext{map[string]Node[ast.Node]{}, nil}},
				currentNode}, initFunction)
	}
}

type CallArrayResultType int

const (
	artNone CallArrayResultType = iota
	artReturn
	artValue
	artFCI
)

type sss struct {
}

func (compiler *Compiler) executeBlockStmt(state State, node Node[*ast.BlockStmt]) ([]Node[ast.Node], CallArrayResultType) {
	newContext := &CurrentContext{map[string]Node[ast.Node]{}, GetCompilerState[*CurrentContext](state)}
	state = SetCompilerState(newContext, state)

	for _, item := range node.Node.List {
		param := ChangeParamNode(node, item)
		tempState := state.setCurrentNode(ChangeParamNode[*ast.BlockStmt, ast.Node](node, item))
		statementFn, currentNode := compiler.findStatement(tempState, param)
		tempState = state.setCurrentNode(currentNode)
		arr, rt := statementFn(tempState)
		switch rt {
		case artFCI:
			switch instance := arr[0].Node.(type) {
			case *ast.FolderContextInformation:
				node = Node[*ast.BlockStmt]{
					node.Key,
					node.Node,
					instance.Imports,
					instance.AbsPath,
					instance.RelPath,
					instance.FileName,
					node.Fs,
					node.Valid,
				}
			}
		case artValue:
		case artReturn:
			return arr, rt
		default:
			continue
		}
	}
	state = SetCompilerState(newContext.parent, state)
	return nil, artNone
}

func (compiler *Compiler) CompileFunc(state State, fn Node[*ast.FuncDecl]) ([]Node[ast.Node], CallArrayResultType) {
	param := ChangeParamNode(fn, fn.Node.Body)
	return compiler.executeBlockStmt(state, param)
}

func (compiler *Compiler) findType(state State, node Node[ast.Expr]) reflect.Type {
	return compiler.internalFindType(0, state, node).(reflect.Type)
}

func (compiler *Compiler) internalFindType(stackIndex int, state State, node Node[ast.Expr]) interface{} {
	initOnCreateType := func(stackIndex int, unk interface{}, indexes []Node[ast.Expr]) interface{} {
		if stackIndex != 0 {
			return unk
		}
		switch value := unk.(type) {
		case OnCreateType:
			return value(state, indexes)
		case reflect.Type:
			return value
		default:
			panic(unk)
		}
	}

	switch item := node.Node.(type) {
	case *ast.IndexExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		indexParam := ChangeParamNode(node, item.Index)
		return initOnCreateType(0, compiler.internalFindType(stackIndex+1, state, param), []Node[ast.Expr]{indexParam})
	case *ast.IndexListExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		var arrIndices []Node[ast.Expr]
		for _, index := range item.Indices {
			indexParam := ChangeParamNode(node, index)
			arrIndices = append(arrIndices, indexParam)
		}
		return initOnCreateType(stackIndex, compiler.internalFindType(stackIndex+1, state, param), arrIndices)
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Expr, ast.Expr](node, item.X)
		unk := compiler.internalFindType(stackIndex+1, state, param)
		switch value := unk.(type) {
		case ast.ImportMapEntry:
			vk := ValueKey{value.Path, item.Sel.Name}
			returnValue, ok := compiler.GlobalTypes[vk]
			if ok {
				return initOnCreateType(stackIndex, returnValue, nil)
			}
			panic("sdfdsfds")
		default:
			panic("sdfdsfds")
		}
	case *ast.Ident:
		if path, ok := node.ImportMap[item.Name]; ok {
			return initOnCreateType(stackIndex, path, nil)
		}
		if onCreateType, ok := compiler.GlobalTypes[ValueKey{"", item.Name}]; ok {
			return initOnCreateType(stackIndex, onCreateType, nil)
		}
		typeMapper := GetCompilerState[TypeMapper](state)
		if rt, ok := typeMapper[item.Name]; ok {
			return initOnCreateType(stackIndex, rt, nil)
		}
		panic(item.Name)
	default:
		panic(node.Node)
	}

}

type EntitySource struct {
	rt reflect.Type
}

func (compiler *Compiler) AddEntitySource(rt reflect.Type) string {
	compiler.NextAlias++
	reference := fmt.Sprintf("T%v", compiler.NextAlias)
	compiler.Sources[reference] = &EntitySource{rt}
	return reference

}

func (compiler *Compiler) findSources(item *TrailRecord) []string {
	m := map[string]bool{}
	for idx := 0; idx < item.Value.NumField(); idx++ {
		if node, ok := item.Value.Field(idx).Interface().(Node[ast.Node]); ok {
			compiler.internalFindSources(node, m)
		}
	}
	ss := make([]string, 0, len(m))
	for key, _ := range m {
		ss = append(ss, key)
	}
	return ss
}

func (compiler *Compiler) internalFindSources(node Node[ast.Node], m map[string]bool) {
	if !node.Valid {
		return
	}
	switch nodeItem := node.Node.(type) {
	case *EntityField:
		m[nodeItem.alias] = true
		break
	case *ast.BasicLit:
		break
	case *coercion:
		compiler.internalFindSources(nodeItem.Node, m)
		break
	case *BinaryExpr:
		for _, v := range nodeItem.Values {
			compiler.internalFindSources(v, m)

		}

	default:
		panic(node)
	}
}

func (compiler *Compiler) calculateSourceDependency(sources []string) []string {
	return sources
}

func (compiler *Compiler) projectTrailRecord(w io.Writer, tabCount int, item *TrailRecord) {
	for idx := 0; idx < item.Value.NumField(); idx++ {
		if node, ok := item.Value.Field(idx).Interface().(Node[ast.Node]); ok {
			compiler.internalProjectTrailRecord(w, tabCount, idx == item.Value.NumField()-1, 0, item.Value.Type().Field(idx).Name, node)
		}
	}
}

func (compiler *Compiler) internalProjectTrailRecord(w io.Writer, count int, last bool, stackCount int, name string, node Node[ast.Node]) {
	if !node.Valid {
		return
	}

	if stackCount == 0 {
		_, _ = io.WriteString(w, strings.Repeat("\t", count))
	}
	switch nodeItem := node.Node.(type) {
	case *EntityField:
		_, _ = io.WriteString(w, fmt.Sprintf("%v.%v", nodeItem.alias, nodeItem.field))
	case *coercion:
		_, _ = io.WriteString(w, "CAST(")
		param := ChangeParamNode[ast.Node, ast.Node](node, nodeItem.Node.Node)
		compiler.internalProjectTrailRecord(w, count, last, stackCount+1, name, param)
		_, _ = io.WriteString(w, " as ")
		switch nodeItem.to {
		case "float64":
			_, _ = io.WriteString(w, "float")
		}
		_, _ = io.WriteString(w, ")")
	case *ast.BasicLit:
		switch nodeItem.Kind {
		case token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING:
			_, _ = io.WriteString(w, nodeItem.Value)
		default:
			panic("unhandled default case")
		}
	case *BinaryExpr:
		_, _ = io.WriteString(w, "(")
		for idx, v := range nodeItem.Values {
			if idx != 0 {
				switch nodeItem.Op {
				case token.ADD: // +
					_, _ = io.WriteString(w, " + ")
				case token.SUB: // -
					_, _ = io.WriteString(w, " - ")
				case token.MUL: // *
					_, _ = io.WriteString(w, " * ")
				case token.QUO: // /
					_, _ = io.WriteString(w, " / ")
				default:
					panic("unhandled default case")
				}
			}
			compiler.internalProjectTrailRecord(w, count, last, stackCount+1, name, v)
		}
		_, _ = io.WriteString(w, ")")
	default:
		panic("implement me")
	}
	if stackCount == 0 {
		_, _ = io.WriteString(w, fmt.Sprintf(" as %v", name))
		if !last {
			_, _ = io.WriteString(w, ",")
		}
		_, _ = io.WriteString(w, "\n")
	}
}

func (compiler *Compiler) projectSources(w io.Writer, tabCount int, sources []string) {
	_, _ = io.WriteString(w, strings.Repeat("\t", tabCount))
	for _, source := range sources {
		query, _ := compiler.Sources[source]
		switch item := query.(type) {
		case *EntitySource:
			_, _ = io.WriteString(w, fmt.Sprintf("%v %v", item.rt.String(), source))
		default:
			panic("unhandled default case")
		}
	}
}

func (compiler *Compiler) nodesToValues(state State, node []Node[ast.Node]) []reflect.Value {
	return nil
}

func (compiler *Compiler) valuesToNodes(state State, values []reflect.Value) []Node[ast.Node] {
	var arr []Node[ast.Node]
	for _, node := range values {
		arr = append(arr, compiler.valueToNode(state, node))
	}
	return arr
}

func (compiler *Compiler) valueToNode(state State, value reflect.Value) Node[ast.Node] {
	kind := value.Kind()
	switch kind {
	case reflect.String:
		s := strconv.Quote(value.String())
		basicLit := &ast.BasicLit{
			ValuePos: state.currentNode.Node.Pos(),
			Kind:     token.STRING,
			Value:    s,
		}
		return ChangeParamNode[ast.Node, ast.Node](state.currentNode, basicLit)
	case reflect.Interface:
		panic("ffff")
	default:
		panic("unhandled default case")
	}

	panic("implement me")
}

func (compiler *Compiler) osGetWdImplementation(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		input := compiler.nodesToValues(state, arguments)
		output := reflect.ValueOf(os.Getwd).Call(input)
		return compiler.valuesToNodes(state, output), artValue
	}
}

func RemoveCompilerState[a IABC](state State) State {
	var result []IABC
	for _, compileState := range state.arr {
		if reflect.TypeFor[a]() == reflect.TypeOf(compileState) {
			continue
		}
		result = append(result, compileState)
	}

	return State{result, state.currentNode}
}

func SetCompilerState[a IABC](data a, state State) State {
	result := State{[]IABC{data}, state.currentNode}
	for _, compileState := range state.arr {
		if reflect.TypeFor[a]() == reflect.TypeOf(compileState) {
			continue
		}
		result.arr = append(result.arr, compileState)
	}
	return result
}

func GetCompilerState[a IABC](state State) a {
	for _, compileState := range state.arr {
		if reflect.TypeFor[a]() == reflect.TypeOf(compileState) {
			return compileState.(a)
		}
	}
	var unk IABC = nil
	vv, _ := unk.(a)
	return vv
}

type CurrentAstNode struct {
	node ast.Node
}
