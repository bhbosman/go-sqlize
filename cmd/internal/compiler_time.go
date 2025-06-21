package internal

import (
	"fmt"
	"go/ast"
	"reflect"
	"time"
)

var timeTimeValueKey = ValueKey{"time", "Time"}

func (compiler *Compiler) addTimeFunctions() {
	compiler.GlobalTypes[timeTimeValueKey] = compiler.registerTimeTime()
	compiler.TypesToValueKeys[reflect.TypeFor[time.Time]()] = timeTimeValueKey
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
