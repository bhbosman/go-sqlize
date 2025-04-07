package internal

import (
	"fmt"
	"go/ast"
	"reflect"
)

func (compiler *Compiler) addBuiltInFunctions() {
	compiler.GlobalTypes[ValueKey{"", "int"}] = compiler.registerInt
	compiler.GlobalTypes[ValueKey{"", "string"}] = compiler.registerString
	compiler.GlobalTypes[ValueKey{"", "float64"}] = compiler.registerFloat64

	compiler.GlobalFunctions[ValueKey{"", "float64"}] = compiler.coercionFloat64
	compiler.GlobalFunctions[ValueKey{"", "float32"}] = compiler.coercionFloat32
	compiler.GlobalFunctions[ValueKey{"", "int"}] = compiler.coercionInt
	compiler.GlobalFunctions[ValueKey{"", "string"}] = compiler.coercionString
	compiler.GlobalFunctions[ValueKey{"", "panic"}] = compiler.builtInPanic
	compiler.GlobalFunctions[ValueKey{"", "nil"}] = compiler.builtInNil
	compiler.GlobalFunctions[ValueKey{"", "true"}] = compiler.builtInTrue
	compiler.GlobalFunctions[ValueKey{"", "false"}] = compiler.builtInFalse
	compiler.GlobalFunctions[ValueKey{"", "println"}] = compiler.builtInPrintln
	compiler.GlobalFunctions[ValueKey{"", "print"}] = compiler.builtInPrint
}

func (compiler *Compiler) registerLibType(state State, i []Node[ast.Expr]) reflect.Type {
	panic("dfgdfgfd")
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

func (compiler *Compiler) builtInNil(State, []Node[ast.Expr], []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &NilExpression{})}, artValue
	}
}

func (compiler *Compiler) builtInTrue(State, []Node[ast.Expr], []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(true)})}, artValue
	}
}

func (compiler *Compiler) builtInFalse(State, []Node[ast.Expr], []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(false)})}, artValue
	}
}

func (compiler *Compiler) builtInPrintln(State, []Node[ast.Expr], []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		// todo: implement at some point
		return nil, artNone
	}
}
func (compiler *Compiler) builtInPrint(State, []Node[ast.Expr], []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		// todo: implement at some point
		return nil, artNone
	}
}

func (compiler *Compiler) coercionInt(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return compiler.genericCoercion(reflect.TypeFor[int](), arguments)
}

func (compiler *Compiler) coercionString(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return compiler.genericCoercion(reflect.TypeFor[string](), arguments)
}

func (compiler *Compiler) genericCoercion(rt reflect.Type, arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("coercion panic requires 1 argument, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		rv, isLiterate := isLiterateValue(arguments[0])
		if isLiterate {
			returnValue := ChangeParamNode[ast.Node, ast.Node](
				state.currentNode,
				&ReflectValueExpression{rv.Convert(rt)},
			)
			return []Node[ast.Node]{returnValue}, artValue
		}
		returnValue := ChangeParamNode[ast.Node, ast.Node](
			state.currentNode,
			&coercion{state.currentNode.Node.Pos(), rt.String(), arguments[0], rt},
		)
		return []Node[ast.Node]{returnValue}, artValue
	}
}

func (compiler *Compiler) coercionFloat32(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return compiler.genericCoercion(reflect.TypeFor[float32](), arguments)
}

func (compiler *Compiler) coercionFloat64(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return compiler.genericCoercion(reflect.TypeFor[float64](), arguments)
}

func (compiler *Compiler) builtInPanic(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("built-in panic requires 1 argument, got %d", len(arguments)))
	}

	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
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
