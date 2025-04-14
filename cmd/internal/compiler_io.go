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

func (compiler *Compiler) ioWriteStringImplementation(state State) ExecuteStatement {

	return func(state State, typeParams []ITypeMapper, unprocessedArgs []Node[ast.Expr]) ([]Node[ast.Node], CallArrayResultType) {
		arguments := compiler.compileArguments(state, unprocessedArgs, typeParams)
		rv := reflect.ValueOf(io.WriteString)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic(fmt.Errorf("io.WriteString only accept literal values"))
	}
}
