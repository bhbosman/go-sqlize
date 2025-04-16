package internal

import (
	"fmt"
	"go/ast"
	"io"
	"reflect"
)

func (compiler *Compiler) addIoFunctions() {
	compiler.GlobalFunctions[ValueKey{"io", "WriteString"}] = functionInformation{compiler.ioWriteStringImplementation, Node[*ast.FuncType]{}, false}
}

func (compiler *Compiler) ioWriteStringImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(io.WriteString)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic(fmt.Errorf("io.WriteString only accept literal values"))
	}
}
