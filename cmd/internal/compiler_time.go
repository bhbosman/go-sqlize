package internal

import (
	"fmt"
	"go/ast"
	"reflect"
	"time"
)

func (compiler *Compiler) addTimeFunctions() {
	compiler.GlobalFunctions[ValueKey{"time", "Now"}] = functionInformation{compiler.timeNowImplementation, Node[*ast.FuncType]{}, false}
}

func (compiler *Compiler) timeNowImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(time.Now)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic(fmt.Errorf("time.Now only accept literal values"))
	}
}
