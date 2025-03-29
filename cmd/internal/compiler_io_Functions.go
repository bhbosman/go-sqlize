package internal

import (
	"fmt"
	"go/ast"
	"io"
	"reflect"
)

func (compiler *Compiler) addIoFunctions() {
	compiler.GlobalFunctions[ValueKey{"io", "WriteString"}] = compiler.ioWriteStringImplementation
}

func (compiler *Compiler) ioWriteStringImplementation(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		rv := reflect.ValueOf(io.WriteString)
		if outputNodes, art, b := compiler.genericCall(state, rv, arguments); b {
			return outputNodes, art
		}
		panic(fmt.Errorf("io.WriteString only accept literal values"))
	}
}
