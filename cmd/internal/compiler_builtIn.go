package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
)

var intValueKey = ValueKey{"", "int"}
var int64ValueKey = ValueKey{"", "int64"}
var int32ValueKey = ValueKey{"", "int32"}
var int16ValueKey = ValueKey{"", "int16"}
var int08ValueKey = ValueKey{"", "int8"}

var uintValueKey = ValueKey{"", "uint"}
var uint64ValueKey = ValueKey{"", "uint64"}
var uint32ValueKey = ValueKey{"", "uint32"}
var uint16ValueKey = ValueKey{"", "uint16"}
var uint08ValueKey = ValueKey{"", "uint8"}
var float32ValueKey = ValueKey{"", "float32"}
var float64ValueKey = ValueKey{"", "float64"}
var stringValueKey = ValueKey{"", "string"}
var boolValueKey = ValueKey{"", "bool"}

func (compiler *Compiler) addBuiltInFunctions() {
	compiler.GlobalTypes[intValueKey] = compiler.registerInt()
	compiler.GlobalTypes[int64ValueKey] = compiler.registerInt64()
	compiler.GlobalTypes[int32ValueKey] = compiler.registerInt32()
	compiler.GlobalTypes[int16ValueKey] = compiler.registerInt16()
	compiler.GlobalTypes[int08ValueKey] = compiler.registerInt08()

	compiler.GlobalTypes[uintValueKey] = compiler.registerUint()
	compiler.GlobalTypes[uint64ValueKey] = compiler.registerUint64()
	compiler.GlobalTypes[uint32ValueKey] = compiler.registerUint32()
	compiler.GlobalTypes[uint16ValueKey] = compiler.registerUint16()
	compiler.GlobalTypes[uint08ValueKey] = compiler.registerUint08()

	compiler.GlobalTypes[boolValueKey] = compiler.registerBool()
	compiler.GlobalTypes[stringValueKey] = compiler.registerString()
	compiler.GlobalTypes[float64ValueKey] = compiler.registerFloat64()
	compiler.GlobalTypes[float32ValueKey] = compiler.registerFloat32()

	compiler.GlobalFunctions[float64ValueKey] = functionInformation{compiler.coercionFloat64UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[float32ValueKey] = functionInformation{compiler.coercionFloat32UnCompiled, Node[*ast.FuncType]{}, false}

	compiler.GlobalFunctions[intValueKey] = functionInformation{compiler.coercionIntUnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[int64ValueKey] = functionInformation{compiler.coercionInt64UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[int32ValueKey] = functionInformation{compiler.coercionInt32UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[int16ValueKey] = functionInformation{compiler.coercionInt16UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[int08ValueKey] = functionInformation{compiler.coercionInt08UnCompiled, Node[*ast.FuncType]{}, false}

	compiler.GlobalFunctions[uintValueKey] = functionInformation{compiler.coercionUintUnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[uint64ValueKey] = functionInformation{compiler.coercionUint64UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[uint32ValueKey] = functionInformation{compiler.coercionUint32UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[uint16ValueKey] = functionInformation{compiler.coercionUint16UnCompiled, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[uint08ValueKey] = functionInformation{compiler.coercionUint08UnCompiled, Node[*ast.FuncType]{}, false}

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

func (compiler *Compiler) registerFloat32() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[float32](), float32ValueKey}
	}
}

func (compiler *Compiler) registerReflectType() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		fn := func() reflect.Type {
			return reflect.TypeFor[reflect.Type]()
		}
		return &ReflectTypeHolder{
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
	ActualType() (reflect.Type, ValueKey)
	Kind() reflect.Kind
}

type ITypeMapperArray []ITypeMapper

type WrapReflectTypeInMapper struct {
	rt reflect.Type
	vk ValueKey
}

func (typeWrapper *WrapReflectTypeInMapper) Keys() []Node[ast.Node] {
	return nil
}

func (typeWrapper *WrapReflectTypeInMapper) GetTypeMapper(state State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{typeWrapper}, true
}

func (typeWrapper *WrapReflectTypeInMapper) Pos() token.Pos {
	return token.NoPos
}

func (typeWrapper *WrapReflectTypeInMapper) End() token.Pos {
	return token.NoPos
}

func (typeWrapper *WrapReflectTypeInMapper) ActualType() (reflect.Type, ValueKey) {
	return typeWrapper.rt, typeWrapper.vk
}

func (typeWrapper *WrapReflectTypeInMapper) Kind() reflect.Kind {
	return typeWrapper.rt.Kind()
}

type ReflectTypeHolder struct {
	fnNodeType      func() reflect.Type
	fnActualType    func() (reflect.Type, ValueKey)
	fnKind          func() reflect.Kind
	fnMapperKeyType func() reflect.Type
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

func (rth *ReflectTypeHolder) Kind() reflect.Kind {
	return rth.fnKind()
}

func (compiler *Compiler) registerInt() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[int](), intValueKey}
	}
}

func (compiler *Compiler) registerInt64() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[int64](), int64ValueKey}
	}
}

func (compiler *Compiler) registerInt32() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[int32](), int32ValueKey}
	}
}

func (compiler *Compiler) registerInt16() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[int16](), int16ValueKey}
	}
}

func (compiler *Compiler) registerInt08() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[int8](), int08ValueKey}
	}
}

func (compiler *Compiler) registerUint64() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[uint64](), uint64ValueKey}
	}
}

func (compiler *Compiler) registerUint() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[uint](), uintValueKey}
	}
}

func (compiler *Compiler) registerUint32() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[uint32](), uint32ValueKey}
	}
}

func (compiler *Compiler) registerUint16() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[uint16](), uint16ValueKey}
	}
}

func (compiler *Compiler) registerUint08() OnCreateType {
	return func(state State, i []ITypeMapper) ITypeMapper {
		return &WrapReflectTypeInMapper{reflect.TypeFor[uint8](), uint08ValueKey}
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

func (compiler *Compiler) coercionInt64UnCompiled(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionInt64Compiled(state, arguments)
	}
}

func (compiler *Compiler) coercionInt32UnCompiled(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionInt32Compiled(state, arguments)
	}
}

func (compiler *Compiler) coercionInt16UnCompiled(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionInt16Compiled(state, arguments)
	}
}

func (compiler *Compiler) coercionInt08UnCompiled(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionInt08Compiled(state, arguments)
	}
}

func (compiler *Compiler) coercionInt64Compiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[int64](), int64ValueKey, compiledArguments)
}

func (compiler *Compiler) coercionInt32Compiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[int32](), int32ValueKey, compiledArguments)
}

func (compiler *Compiler) coercionInt16Compiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[int16](), int16ValueKey, compiledArguments)
}

func (compiler *Compiler) coercionInt08Compiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[int8](), int08ValueKey, compiledArguments)
}

func (compiler *Compiler) coercionIntUnCompiled(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionIntCompiled(state, arguments)
	}
}

func (compiler *Compiler) coercionUintUnCompiled(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionUintCompiled(state, arguments)
	}
}

func (compiler *Compiler) coercionUint64UnCompiled(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionUint64Compiled(state, arguments)
	}
}

func (compiler *Compiler) coercionUint32UnCompiled(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionUint32Compiled(state, arguments)
	}
}
func (compiler *Compiler) coercionUint16UnCompiled(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionUint16Compiled(state, arguments)
	}
}

func (compiler *Compiler) coercionUint08UnCompiled(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}
		return compiler.coercionUint08Compiled(state, arguments)
	}
}

func (compiler *Compiler) coercionIntCompiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[int](), intValueKey, compiledArguments)
}

func (compiler *Compiler) coercionUintCompiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[uint](), uintValueKey, compiledArguments)
}

func (compiler *Compiler) coercionUint64Compiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[uint64](), uint64ValueKey, compiledArguments)
}

func (compiler *Compiler) coercionUint32Compiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[uint32](), uint32ValueKey, compiledArguments)
}

func (compiler *Compiler) coercionUint16Compiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[uint16](), uint16ValueKey, compiledArguments)
}

func (compiler *Compiler) coercionUint08Compiled(state State, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[uint8](), uint08ValueKey, compiledArguments)
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
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[string](), stringValueKey, compiledArguments)
}

func (compiler *Compiler) genericCoercionCompiled(state State, rt reflect.Type, vk ValueKey, compiledArguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	rv, isLiterate := isLiterateValue(compiledArguments[0])
	if isLiterate {
		returnValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rv.Convert(rt), vk})
		return []Node[ast.Node]{returnValue}, artValue
	}
	returnValue := ChangeParamNode[ast.Node, ast.Node](
		state.currentNode,
		coercion{state.currentNode.Node.Pos(), rt.String(), compiledArguments[0], rt, vk},
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
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[float32](), float32ValueKey, compiledArguments)
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
	return compiler.genericCoercionCompiled(state, reflect.TypeFor[float64](), float64ValueKey, compiledArguments)
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
