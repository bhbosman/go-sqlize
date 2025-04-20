package internal

import (
	"go/ast"
	"go/token"
	"reflect"
	"sort"
)

type libDictionaryLookupImplementation struct {
	compiler *Compiler
	state    State
}

func (impl libDictionaryLookupImplementation) ExecuteStatement() ExecuteStatement {

	return impl.Run
}

func (impl libDictionaryLookupImplementation) Run(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	var conditionalStatement []SingleValueCondition
	dictionaryExpression := arguments[0].Node.(*DictionaryExpression)
	{
		rvMap := dictionaryExpression.m
		keyArr := rvMap.MapKeys()
		sorter := &rvArraySorter{keyArr}
		sort.Sort(sorter)

		for _, rvKey := range keyArr {
			rvValue := rvMap.MapIndex(rvKey)
			inputData := arguments[1]
			expressions := impl.walk(inputData, rvKey)

			mbe := &MultiBinaryExpr{token.LAND, expressions}
			mbeNode := ChangeParamNode[ast.Node, ast.Node](impl.state.currentNode, mbe)

			singleValueCondition := SingleValueCondition{mbeNode, rvValue.Interface().(Node[ast.Node])}
			conditionalStatement = append(conditionalStatement, singleValueCondition)
		}
	}
	{
		rvDefault := dictionaryExpression.defaultValue
		condition := ChangeParamNode[ast.Node, ast.Node](state.currentNode, &ReflectValueExpression{reflect.ValueOf(true)})
		singleValueCondition := SingleValueCondition{condition: condition, value: rvDefault}
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
