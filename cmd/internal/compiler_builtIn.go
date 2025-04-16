package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
)

type SomeDataWithRv struct {
	rt       reflect.Type
	rv       interface{}
	assigned bool
}

type SomeDataWithNode struct {
	node     Node[ast.Node]
	assigned bool
}

func (compiler *Compiler) addBuiltInFunctions() {
	compiler.GlobalTypes[ValueKey{"", "int"}] = compiler.registerInt()
	compiler.GlobalTypes[ValueKey{"", "string"}] = compiler.registerString()
	compiler.GlobalTypes[ValueKey{"", "float64"}] = compiler.registerFloat64()

	compiler.GlobalFunctions[ValueKey{"", "float64"}] = functionInformation{compiler.coercionFloat64UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "float32"}] = functionInformation{compiler.coercionFloat32UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "int"}] = functionInformation{compiler.coercionIntUnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "string"}] = functionInformation{compiler.coercionStringUnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "panic"}] = functionInformation{compiler.builtInPanic, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "nil"}] = functionInformation{compiler.builtInNil, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "true"}] = functionInformation{compiler.builtInTrue, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "false"}] = functionInformation{compiler.builtInFalse, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "println"}] = functionInformation{compiler.builtInPrintln, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "print"}] = functionInformation{compiler.builtInPrint, Node[*ast.FuncType]{}, false}
}

func (compiler *Compiler) registerLibType() OnCreateType {
	return func(state State, i []Node[ast.Node]) ITypeMapper {
		return &ReflectTypeHolder{
			func(option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
				panic("dfgdfgf")
			},
			func() reflect.Type {
				panic("dfgdfgf")
			},
			func() reflect.Type {
				panic("dfgdfgf")
			},
			func() reflect.Kind {
				return reflect.Invalid
			},
			func() reflect.Type {
				panic("dfgdfgf")
			},
			func() reflect.Type {
				panic("dfgdfgf")
			},
		}
	}
}

func (compiler *Compiler) registerSomeType() OnCreateType {
	return func(state State, typeParams []Node[ast.Node]) ITypeMapper {
		fn := func() reflect.Type {
			return reflect.TypeFor[SomeDataWithRv]()
		}
		return &ReflectTypeHolder{
			func(option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
				switch unk := rv.Interface().(type) {
				case SomeDataWithNode:
					if !unk.assigned {
						rt := reflect.TypeFor[SomeDataWithRv]()
						return reflect.Zero(rt)
					}
					switch nodeItem := unk.node.Node.(type) {
					case *ReflectValueExpression:
						v := SomeDataWithRv{reflect.TypeOf(nodeItem.Rv), nodeItem.Rv.Interface(), true}
						return reflect.ValueOf(v)
					}
				default:
				}
				v := SomeDataWithRv{reflect.TypeOf(rv), rv.Interface(), true}
				return reflect.ValueOf(v)
			},
			fn,
			fn,
			func() reflect.Kind {
				return reflect.Struct
			},
			fn,
			fn,
		}
	}
}

func (compiler *Compiler) registerFloat64() OnCreateType {
	return func(state State, i []Node[ast.Node]) ITypeMapper {
		fn := func() reflect.Type {
			return reflect.TypeFor[float64]()
		}
		return &ReflectTypeHolder{
			func(option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
				return rv
			},
			fn,
			fn,
			func() reflect.Kind {
				return reflect.Float64
			},
			fn,
			fn,
		}
	}
}

func (compiler *Compiler) registerReflectType() OnCreateType {
	return func(state State, i []Node[ast.Node]) ITypeMapper {
		fn := func() reflect.Type {
			return reflect.TypeFor[reflect.Type]()
		}
		return &ReflectTypeHolder{
			func(option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
				return rv
			},
			fn,
			fn,
			func() reflect.Kind {
				return reflect.TypeFor[reflect.Type]().Kind()
			},
			fn,
			fn,
		}
	}
}

func (compiler *Compiler) registerString() OnCreateType {
	return func(state State, i []Node[ast.Node]) ITypeMapper {
		fn := func() reflect.Type {
			return reflect.TypeFor[string]()
		}
		return &ReflectTypeHolder{
			func(option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
				return rv
			},
			fn,
			fn,
			func() reflect.Kind {
				return reflect.String
			},
			fn,
			fn,
		}
	}
}

type TypeMapperCreateOption int

const (
	tmcoDefault TypeMapperCreateOption = iota
	tmcoMapKey
	tmcoMapValue
)

type ITypeMapper interface {
	ast.Node
	Create(option TypeMapperCreateOption, rv reflect.Value) reflect.Value
	NodeType() reflect.Type
	ActualType() reflect.Type
	MapperKeyType() reflect.Type
	MapperValueType() reflect.Type
	Kind() reflect.Kind
}

type ITypeMapperArray []ITypeMapper

func (receiver ITypeMapperArray) toNodeArray() []Node[ast.Node] {
	var result []Node[ast.Node]
	for _, element := range receiver {
		result = append(result, Node[ast.Node]{Node: &ReflectValueExpression{reflect.ValueOf(element.ActualType())}, Valid: true})
	}
	return result
}

type WrapReflectTypeInMapper struct {
	rt reflect.Type
}

func (typeWrapper *WrapReflectTypeInMapper) GetTypeMapper(string) (ITypeMapperArray, bool) {
	return ITypeMapperArray{typeWrapper}, true
}

func (typeWrapper *WrapReflectTypeInMapper) Pos() token.Pos {
	return token.NoPos
}

func (typeWrapper *WrapReflectTypeInMapper) End() token.Pos {
	return token.NoPos
}

func (typeWrapper *WrapReflectTypeInMapper) Create(option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
	return rv
}

func (typeWrapper *WrapReflectTypeInMapper) NodeType() reflect.Type {
	return typeWrapper.rt
}

func (typeWrapper *WrapReflectTypeInMapper) ActualType() reflect.Type {
	return typeWrapper.rt
}

func (typeWrapper *WrapReflectTypeInMapper) MapperKeyType() reflect.Type {
	return typeWrapper.rt
}

func (typeWrapper *WrapReflectTypeInMapper) MapperValueType() reflect.Type {
	return typeWrapper.rt
}

func (typeWrapper *WrapReflectTypeInMapper) Kind() reflect.Kind {
	return typeWrapper.rt.Kind()
}

type ReflectTypeHolder struct {
	fnCreate          func(option TypeMapperCreateOption, rv reflect.Value) reflect.Value
	fnNodeType        func() reflect.Type
	fnActualType      func() reflect.Type
	fnKind            func() reflect.Kind
	fnMapperKeyType   func() reflect.Type
	fnMapperValueType func() reflect.Type
}

func (rth *ReflectTypeHolder) Pos() token.Pos {
	return token.NoPos
}

func (rth *ReflectTypeHolder) End() token.Pos {
	return token.NoPos
}

func (rth *ReflectTypeHolder) ActualType() reflect.Type {
	return rth.fnActualType()
}

func (rth *ReflectTypeHolder) MapperValueType() reflect.Type {
	return rth.fnMapperValueType()
}

func (rth *ReflectTypeHolder) MapperKeyType() reflect.Type {
	return rth.fnMapperKeyType()
}

func (rth *ReflectTypeHolder) Kind() reflect.Kind {
	return rth.fnKind()
}

func (rth *ReflectTypeHolder) NodeType() reflect.Type {
	return rth.fnNodeType()
}

func (rth *ReflectTypeHolder) Create(option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
	return rth.fnCreate(option, rv)
}

func (compiler *Compiler) registerInt() OnCreateType {
	return func(state State, i []Node[ast.Node]) ITypeMapper {

		return &WrapReflectTypeInMapper{reflect.TypeFor[int]()}

	}
}

func (compiler *Compiler) registerBool() OnCreateType {
	return func(state State, i []Node[ast.Node]) ITypeMapper {
		fn := func() reflect.Type {
			return reflect.TypeFor[bool]()
		}
		return &ReflectTypeHolder{
			func(option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
				return rv
			},
			fn,
			fn,
			func() reflect.Kind {
				return reflect.Bool
			},
			fn,
			fn,
		}
	}
}

func (compiler *Compiler) builtInNil(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &builtInNil{})}, artValue
	}
}

func (compiler *Compiler) builtInTrue(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(true)})}, artValue
	}
}

func (compiler *Compiler) builtInFalse(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(false)})}, artValue
	}
}

func (compiler *Compiler) builtInPrintln(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		// todo: implement at some point
		return nil, artNone
	}
}

func (compiler *Compiler) builtInPrint(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		// todo: implement at some point
		return nil, artNone
	}
}

func (compiler *Compiler) coercionIntUnCompiled(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionIntCompiled(state, arguments)
	}
}

func (compiler *Compiler) coercionIntCompiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[int](), compiledArguments)
}

func (compiler *Compiler) coercionStringUnCompiled(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionStringCompiled(state, arguments)
	}
}

func (compiler *Compiler) coercionStringCompiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[string](), compiledArguments)
}

func (compiler *Compiler) genericCoercionCompiled(state State, rt reflect.Type, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	rv, isLiterate := isLiterateValue(compiledArguments[0])
	if isLiterate {
		returnValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rv.Convert(rt)})
		return []Node[ast.Node]{returnValue}, artValue
	}
	returnValue := ChangeParamNode[ast.Node, ast.Node](
		state.currentNode,
		&coercion{state.currentNode.Node.Pos(), rt.String(), compiledArguments[0], rt},
	)
	return []Node[ast.Node]{returnValue}, artValue
}

func (compiler *Compiler) coercionFloat32UnCompiled(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionFloat32Compiled(state, arguments)
	}
}

func (compiler *Compiler) coercionFloat32Compiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[float32](), compiledArguments)
}

func (compiler *Compiler) coercionFloat64UnCompiled(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionFloat64Compiled(state, arguments)
	}
}

func (compiler *Compiler) coercionFloat64Compiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[float64](), compiledArguments)
}

func (compiler *Compiler) builtInPanic(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("built-in panic requires 1 argument, got %d", len(arguments)))
		}
		switch arg := arguments[0].Node.(type) {
		case *ast.BasicLit:
			panic(fmt.Errorf(arg.Value))
		case *ReflectValueExpression:
			panic(fmt.Errorf(arg.Rv.String()))
		default:
			panic("implementation required")
		}
	}
}
