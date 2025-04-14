package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"sort"
)

type libDictionaryLookupImplementation struct {
	state     State
	params    []Node[ast.Expr]
	arguments []Node[ast.Node]
}

func (impl libDictionaryLookupImplementation) ExecuteStatement() ExecuteStatement {
	if len(impl.arguments) != 2 {
		panic(fmt.Errorf("DictionaryLookup implementation requires 2 arguments, got %d", len(impl.arguments)))
	}
	return impl.Run
}

func (impl libDictionaryLookupImplementation) Run(state State) ([]Node[ast.Node], CallArrayResultType) {
	var conditionalStatement []SingleValueCondition
	dictionaryExpression := impl.arguments[0].Node.(*DictionaryExpression)
	{
		rvMap := dictionaryExpression.m
		keyArr := rvMap.MapKeys()
		sorter := &rvArraySorter{keyArr}
		sort.Sort(sorter)

		for _, rvKey := range keyArr {
			rvValue := rvMap.MapIndex(rvKey)
			inputData := impl.arguments[1]
			expressions := impl.walk(inputData, rvKey)

			mbe := &MultiBinaryExpr{token.LAND, expressions}
			mbeNode := ChangeParamNode[ast.Node, ast.Node](impl.state.currentNode, mbe)

			singleValueCondition := SingleValueCondition{mbeNode, ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rvValue})}
			conditionalStatement = append(conditionalStatement, singleValueCondition)
		}
	}
	{
		rvDefault := dictionaryExpression.defaultValue
		condition := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(true)})
		singleValueCondition := SingleValueCondition{condition: condition, value: ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{rvDefault})}
		conditionalStatement = append(conditionalStatement, singleValueCondition)
	}
	ite := &IfThenElseSingleValueCondition{conditionalStatement}
	resultValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, ite)
	return []Node[ast.Node]{resultValue}, artReturn
}

func (impl libDictionaryLookupImplementation) walk(inputData Node[ast.Node], rvKey reflect.Value) []Node[ast.Node] {
	switch {
	case rvKey.CanFloat() || rvKey.CanInt() || rvKey.Kind() == reflect.String:
		left := inputData
		right := ChangeParamNode[ast.Node, ast.Node](impl.state.currentNode, &ReflectValueExpression{rvKey})
		be := &BinaryExpr{token.NoPos, token.EQL, left, right}
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](impl.state.currentNode, be)}
	case rvKey.Kind() == reflect.Struct:
		switch leftItem := inputData.Node.(type) {
		case *TrailRecord:
			if leftItem.Value.NumField() == rvKey.NumField() {
				var expressions []Node[ast.Node]
				for idx := 0; idx < rvKey.NumField(); idx++ {
					left := leftItem.Value.Field(idx).Interface().(Node[ast.Node])
					nodes := impl.walk(left, rvKey.Field(idx))
					expressions = append(expressions, nodes...)
				}
				return expressions
			}
		}
		panic("sdsfdsfd")
	default:
		panic("find out")
	}
}
