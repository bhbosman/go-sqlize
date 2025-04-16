package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
)

type SomeDataWithRv struct {
	rv       reflect.Value
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

	compiler.GlobalFunctions[ValueKey{"", "float64"}] = functionInformation{compiler.coercionFloat64, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "float32"}] = functionInformation{compiler.coercionFloat32, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "int"}] = functionInformation{compiler.coercionInt, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"", "string"}] = functionInformation{compiler.coercionString, Node[*ast.FuncType]{}, false}
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
			func(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
				panic("dfgdfgf")
			},
			func(state State) reflect.Type {
				panic("dfgdfgf")
			},
			func(state State) reflect.Type {
				panic("dfgdfgf")
			},
			func() reflect.Kind {
				return reflect.Invalid
			},
			func(state State) reflect.Type {
				panic("dfgdfgf")
			},
			func(state State) reflect.Type {
				panic("dfgdfgf")
			},
		}
	}
}

func (compiler *Compiler) registerSomeType() OnCreateType {
	return func(state State, typeParams []Node[ast.Node]) ITypeMapper {
		fn := func(state State) reflect.Type {
			return reflect.TypeFor[SomeDataWithRv]()
		}
		return &ReflectTypeHolder{
			func(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
				switch unk := rv.Interface().(type) {
				case SomeDataWithNode:
					if !unk.assigned {
						rt := reflect.TypeFor[SomeDataWithRv]()
						return reflect.Zero(rt)
					}
					switch nodeItem := unk.node.Node.(type) {
					case *ReflectValueExpression:
						v := SomeDataWithRv{nodeItem.Rv, true}
						return reflect.ValueOf(v)
					}
				default:
				}
				v := SomeDataWithRv{rv, true}
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
		fn := func(state State) reflect.Type {
			return reflect.TypeFor[float64]()
		}
		return &ReflectTypeHolder{
			func(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
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

func (compiler *Compiler) registerString() OnCreateType {
	return func(state State, i []Node[ast.Node]) ITypeMapper {
		fn := func(state State) reflect.Type {
			return reflect.TypeFor[string]()
		}
		return &ReflectTypeHolder{
			func(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
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
	Create(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value
	NodeType(state State) reflect.Type
	ActualType(state State) reflect.Type
	MapperKeyType(state State) reflect.Type
	MapperValueType(state State) reflect.Type
	Kind() reflect.Kind
}

type WrapReflectTypeInMapper struct {
	rt reflect.Type
}

func (typeWrapper *WrapReflectTypeInMapper) Pos() token.Pos {
	return token.NoPos
}

func (typeWrapper *WrapReflectTypeInMapper) End() token.Pos {
	return token.NoPos
}

func (typeWrapper *WrapReflectTypeInMapper) Create(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
	return rv
}

func (typeWrapper *WrapReflectTypeInMapper) NodeType(state State) reflect.Type {
	return typeWrapper.rt
}

func (typeWrapper *WrapReflectTypeInMapper) ActualType(state State) reflect.Type {
	return typeWrapper.rt
}

func (typeWrapper *WrapReflectTypeInMapper) MapperKeyType(state State) reflect.Type {
	//TODO implement me
	panic("implement me")
}

func (typeWrapper *WrapReflectTypeInMapper) MapperValueType(state State) reflect.Type {
	//TODO implement me
	panic("implement me")
}

func (typeWrapper *WrapReflectTypeInMapper) Kind() reflect.Kind {
	return typeWrapper.rt.Kind()
}

type ReflectTypeHolder struct {
	fnCreate          func(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value
	fnNodeType        func(state State) reflect.Type
	fnActualType      func(state State) reflect.Type
	fnKind            func() reflect.Kind
	fnMapperKeyType   func(state State) reflect.Type
	fnMapperValueType func(state State) reflect.Type
}

func (rth *ReflectTypeHolder) Pos() token.Pos {
	return token.NoPos
}

func (rth *ReflectTypeHolder) End() token.Pos {
	return token.NoPos
}

func (rth *ReflectTypeHolder) ActualType(state State) reflect.Type {
	return rth.fnActualType(state)
}

func (rth *ReflectTypeHolder) MapperValueType(state State) reflect.Type {
	return rth.fnMapperValueType(state)
}

func (rth *ReflectTypeHolder) MapperKeyType(state State) reflect.Type {
	return rth.fnMapperKeyType(state)
}

func (rth *ReflectTypeHolder) Kind() reflect.Kind {
	return rth.fnKind()
}

func (rth *ReflectTypeHolder) NodeType(state State) reflect.Type {
	return rth.fnNodeType(state)
}

func (rth *ReflectTypeHolder) Create(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
	return rth.fnCreate(state, option, rv)
}

func (compiler *Compiler) registerInt() OnCreateType {
	return func(state State, i []Node[ast.Node]) ITypeMapper {
		fn := func(state State) reflect.Type {
			return reflect.TypeFor[int]()
		}
		return &ReflectTypeHolder{
			func(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
				return rv
			},
			fn,
			fn,
			func() reflect.Kind {
				return reflect.Int
			},
			fn,
			fn,
		}
	}
}

func (compiler *Compiler) registerBool() OnCreateType {
	return func(state State, i []Node[ast.Node]) ITypeMapper {
		fn := func(state State) reflect.Type {
			return reflect.TypeFor[bool]()
		}
		return &ReflectTypeHolder{
			func(state State, option TypeMapperCreateOption, rv reflect.Value) reflect.Value {
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

func (compiler *Compiler) builtInNil(state State) ExecuteStatement {
	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &builtInNil{})}, artValue
	}
}

func (compiler *Compiler) builtInTrue(state State) ExecuteStatement {
	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(true)})}, artValue
	}
}

func (compiler *Compiler) builtInFalse(state State) ExecuteStatement {
	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(false)})}, artValue
	}
}

func (compiler *Compiler) builtInPrintln(state State) ExecuteStatement {
	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		// todo: implement at some point
		return nil, artNone
	}
}

func (compiler *Compiler) builtInPrint(state State) ExecuteStatement {
	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		// todo: implement at some point
		return nil, artNone
	}
}

func (compiler *Compiler) coercionInt(state State) ExecuteStatement {

	return compiler.genericCoercion(reflect.TypeFor[int]())
}

func (compiler *Compiler) coercionString(state State) ExecuteStatement {
	return compiler.genericCoercion(reflect.TypeFor[string]())
}

func (compiler *Compiler) genericCoercion(rt reflect.Type) ExecuteStatement {
	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		arguments := compiler.compileArguments(state, unprocessedArgs, typeParams)
		if len(arguments) != 1 {
			panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
		}

		rv, isLiterate := isLiterateValue(arguments[0])
		if isLiterate {
			returnValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rv.Convert(rt)})
			return []Node[ast.Node]{returnValue}, artValue
		}
		returnValue := ChangeParamNode[ast.Node, ast.Node](
			state.currentNode,
			&coercion{state.currentNode.Node.Pos(), rt.String(), arguments[0], rt},
		)
		return []Node[ast.Node]{returnValue}, artValue
	}
}

func (compiler *Compiler) coercionFloat32(state State) ExecuteStatement {
	return compiler.genericCoercion(reflect.TypeFor[float32]())
}

func (compiler *Compiler) coercionFloat64(state State) ExecuteStatement {
	return compiler.genericCoercion(reflect.TypeFor[float64]())
}

func (compiler *Compiler) builtInPanic(state State) ExecuteStatement {
	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		arguments := compiler.compileArguments(state, unprocessedArgs, typeParams)
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
