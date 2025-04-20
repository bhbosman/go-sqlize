package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
)

var intValueKey = ValueKey{"", "int"}
var float64ValueKey = ValueKey{"", "float64"}
var stringValueKey = ValueKey{"", "string"}
var boolValueKey = ValueKey{"", "bool"}

func (compiler *Compiler) addBuiltInFunctions() {
	compiler.GlobalTypes[intValueKey] = compiler.registerInt()
	compiler.GlobalTypes[stringValueKey] = compiler.registerString()
	compiler.GlobalTypes[float64ValueKey] = compiler.registerFloat64()
	compiler.GlobalFunctions[float64ValueKey] = functionInformation{compiler.coercionFloat64UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "float32"}] = functionInformation{compiler.coercionFloat32UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[intValueKey] = functionInformation{compiler.coercionIntUnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[stringValueKey] = functionInformation{compiler.coercionStringUnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "panic"}] = functionInformation{compiler.builtInPanic, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "nil"}] = functionInformation{compiler.builtInNil, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "true"}] = functionInformation{compiler.builtInTrue, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "false"}] = functionInformation{compiler.builtInFalse, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "println"}] = functionInformation{compiler.builtInPrintln, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "print"}] = functionInformation{compiler.builtInPrint, Node[*ast.FuncType]{}, false}
}

func (compiler *Compiler) registerLibType() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &ReflectTypeHolder{
			nil,
			func() reflect.Type {
				panic("dfgdfgf")
			},
			func() (reflect.Type, ValueKey) {
				panic("dfgdfgf")
			},
			func() reflect.Kind {
				return reflect.Invalid
			},
			func() reflect.Type {
				panic("dfgdfgf")
			},
		}
	}
}

func (compiler *Compiler) registerFloat64() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[float64](), float64ValueKey}
	}
}

func (compiler *Compiler) registerReflectType() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		fn := func() reflect.Type {
			return reflect.TypeFor[reflect.Type]()
		}
		return &ReflectTypeHolder{
			nil,
			fn,
			func() (reflect.Type, ValueKey) {
				return fn(), ValueKey{"reflect", "Type"}
			},
			func() reflect.Kind {
				return reflect.TypeFor[reflect.Type]().Kind()
			},
			fn,
		}
	}
}

func (compiler *Compiler) registerString() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[string](), stringValueKey}

	}
}

type ITranslateNodeValueToReflectValue interface {
	TranslateNodeValueToReflectValue(node Node[ast.Node]) reflect.Value
}

type ITypeMapper interface {
	ast.Node
	NodeType() reflect.Type
	ActualType() (reflect.Type, ValueKey)
	MapperKeyType() reflect.Type
	Kind() reflect.Kind
	Keys() []Node[ast.Node]
}

type ITypeMapperArray []ITypeMapper

type WrapReflectTypeInMapper struct {
	rt reflect.Type
	vk ValueKey
}

func (typeWrapper *WrapReflectTypeInMapper) Keys() []Node[ast.Node] {
	return nil
}

func (typeWrapper *WrapReflectTypeInMapper) GetTypeMapper() (ITypeMapperArray, bool) {
	return ITypeMapperArray{typeWrapper}, true
}

func (typeWrapper *WrapReflectTypeInMapper) Pos() token.Pos {
	return token.NoPos
}

func (typeWrapper *WrapReflectTypeInMapper) End() token.Pos {
	return token.NoPos
}

func (typeWrapper *WrapReflectTypeInMapper) NodeType() reflect.Type {
	return typeWrapper.rt
}

func (typeWrapper *WrapReflectTypeInMapper) ActualType() (reflect.Type, ValueKey) {
	return typeWrapper.rt, typeWrapper.vk
}

func (typeWrapper *WrapReflectTypeInMapper) MapperKeyType() reflect.Type {
	return typeWrapper.rt
}

func (typeWrapper *WrapReflectTypeInMapper) Kind() reflect.Kind {
	return typeWrapper.rt.Kind()
}

type ReflectTypeHolder struct {
	fnTranslateNodeValueToReflectValue func(node Node[ast.Node]) reflect.Value
	fnNodeType                         func() reflect.Type
	fnActualType                       func() (reflect.Type, ValueKey)
	fnKind                             func() reflect.Kind
	fnMapperKeyType                    func() reflect.Type
}

func (rth *ReflectTypeHolder) TranslateNodeValueToReflectValue(node Node[ast.Node]) reflect.Value {
	return rth.fnTranslateNodeValueToReflectValue(node)
}

func (rth *ReflectTypeHolder) Keys() []Node[ast.Node] {
	return nil
}

func (rth *ReflectTypeHolder) Pos() token.Pos {
	return token.NoPos
}

func (rth *ReflectTypeHolder) End() token.Pos {
	return token.NoPos
}

func (rth *ReflectTypeHolder) ActualType() (reflect.Type, ValueKey) {
	return rth.fnActualType()
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

func (compiler *Compiler) registerInt() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[int](), intValueKey}
	}
}

func (compiler *Compiler) registerBool() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[bool](), boolValueKey}
	}
}

func (compiler *Compiler) builtInNil(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &builtInNil{})}, artValue
	}
}

func (compiler *Compiler) builtInTrue(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(true), boolValueKey})}, artValue
	}
}

func (compiler *Compiler) builtInFalse(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(false), boolValueKey})}, artValue
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
		returnValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rv.Convert(rt), ValueKey{}})
		return []Node[ast.Node]{returnValue}, artValue
	}
	returnValue := ChangeParamNode[ast.Node, ast.Node](
		state.currentNode,
		coercion{state.currentNode.Node.Pos(), rt.String(), compiledArguments[0], rt},
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
