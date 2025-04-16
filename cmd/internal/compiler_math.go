package internal

import (
	"go/ast"
	"math"
	"reflect"
)

func (compiler *Compiler) addMathFunctions() {
	compiler.GlobalFunctions[ValueKey{"math", "Sin"}] = functionInformation{compiler.mathSinImplementation, Node[*ast.FuncType]{}, false}
	compiler.GlobalFunctions[ValueKey{"math", "Cos"}] = functionInformation{compiler.mathCosImplementation, Node[*ast.FuncType]{}, false}
}

func (compiler *Compiler) mathSinImplementation(state State) ExecuteStatement {

	return func(state State, typeParams ITypeMapperArray, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(math.Sin)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		supportedFunction := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &SupportedFunction{"sin", arguments, reflect.TypeFor[float64]()})
		return []Node[ast.Node]{supportedFunction}, artValue
	}
}

func (compiler *Compiler) mathCosImplementation(state State) ExecuteStatement {

	return func(state State, typeParams ITypeMapperArray, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(math.Cos)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		supportedFunction := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &SupportedFunction{"cos", arguments, reflect.TypeFor[float64]()})
		return []Node[ast.Node]{supportedFunction}, artValue
	}
}
