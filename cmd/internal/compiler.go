package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"strings"
)

type GlobalMethodHandlerKey struct {
	rt         reflect.Type
	methodName string
}

type CompilerState int

const (
	CompilerState_InitCalled CompilerState = 1 << iota
)

type TypeMapper map[string]reflect.Type

type ExecuteStatement func(state State) ([]Node[ast.Node], CallArrayResultType)

type AssignStatement func(state State, value Node[ast.Node])

type OnCreateExecuteStatement func(state State, typeParams []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement

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
		ValueKey{"", "int"}:         compiler.registerInt,
		ValueKey{"", "string"}:      compiler.registerString,
		ValueKey{"", "float64"}:     compiler.registerFloat64,
		ValueKey{libFolder, "Some"}: compiler.registerLibType,
	}

	compiler.GlobalFunctions = map[ValueKey]OnCreateExecuteStatement{
		ValueKey{"path/filepath", "Join"}: compiler.pathFilepathJoinImplementation,
		ValueKey{"path/filepath", "Dir"}:  compiler.pathFilepathDirImplementation,
	}
	compiler.addBuiltInFunctions()
	compiler.addStrconvFunctions()
	compiler.addLibFunctions()
	compiler.addOsFunctions()
	compiler.addIoFunctions()
	compiler.addMathFunctions()

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
		}
	}

	for key, value := range TypeSpecMap {
		if key.Folder == libFolder && key.Key == "Some" {
			continue
		}

		fn := func(vk ValueKey, node Node[*ast.TypeSpec]) OnCreateType {
			return func(state State, exprs []Node[ast.Expr]) reflect.Type {
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
							state = RemoveCompilerState[TypeMapper](state)
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

func (compiler *Compiler) genericCall(state State, rv reflect.Value, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType, bool) {
	input, allLiterals := compiler.nodesToValues(state, rv, arguments)
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
					[]IABC{&CurrentContext{map[string]Node[ast.Node]{}, currentContext}},
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

type sss struct {
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
		if onCreateType, ok := compiler.GlobalTypes[ValueKey{node.RelPath, item.Name}]; ok {
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

func (compiler *Compiler) calculateSourceDependency(sources []string) []string {
	return sources
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

func (compiler *Compiler) nodesToValues(state State, rvFunc reflect.Value, nodes []Node[ast.Node]) ([]reflect.Value, bool) {
	funcRt := rvFunc.Type()
	var arr []reflect.Value
	for idx, node := range nodes {
		rv, b := compiler.nodeToValue(state, node)
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

func (compiler *Compiler) nodeToValue(_ State, node Node[ast.Node]) (reflect.Value, bool) {
	if value, isLiterateValue := isLiterateValue(node); isLiterateValue {
		return value, true
	}
	return reflect.Value{}, false
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
	case reflect.Interface:
		return ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{value.Elem()})
	case reflect.String, reflect.Pointer, reflect.Int, reflect.Float32, reflect.Float64:
		return ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{value})
	default:
		panic("unhandled default case")
	}
}

func (compiler *Compiler) executeAndExpandStatement(state State, executeStatement ExecuteStatement) ([]Node[ast.Node], CallArrayResultType) {
	var result []Node[ast.Node]
	arr, v := executeStatement(state)
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
