package internal

import (
	"go/ast"
	"go/token"
	"reflect"
)

type (
	IFindValueKey interface {
		GetValueKey() ValueKey
	}

	IFindTypeMapper interface {
		ast.Node
		GetTypeMapper(State) (ITypeMapperArray, bool)
	}

	EntitySource struct {
		typeMapper ITypeMapper
		queryState queryState
	}
	TrailRecord struct {
		Position   token.Pos     // identifier position
		Value      reflect.Value // identifier name
		typeMapper ITypeMapper
	}

	TrailSource struct {
		Alias      string
		typeMapper ITypeMapper
	}
	EntityField struct {
		alias      string
		typeMapper ITypeMapper
		field      string
	}
	coercion struct {
		Position token.Pos
		to       string
		Node     Node[ast.Node]
		rt       reflect.Type
		vk       ValueKey
	}

	CheckForNotNullExpression struct {
		node       Node[ast.Node]
		typeMapper ITypeMapper
	}
	BinaryExpr struct {
		Op         token.Token // operator
		left       Node[ast.Node]
		right      Node[ast.Node]
		typeMapper ITypeMapper
	}
	builtInNil struct {
	}
	ReflectValueExpression struct {
		Rv reflect.Value
		// todo: add ITypeMapper here
		Vk ValueKey
	}
	MultiValueCondition struct {
		condition Node[ast.Node]
		values    []Node[ast.Node]
	}
	SingleValueCondition struct {
		condition Node[ast.Node]
		value     Node[ast.Node]
	}
	SupportedFunction struct {
		functionName string
		params       []Node[ast.Node]
		rt           reflect.Type
	}
	IfThenElseMultiValueCondition struct {
		conditionalStatement []MultiValueCondition
	}
	BooleanCondition struct {
		conditions []Node[ast.Node]
	}
	IfThenElseSingleValueCondition struct {
		conditionalStatement []SingleValueCondition
	}
	MultiBinaryExpr struct {
		Op          token.Token // operator
		expressions []Node[ast.Node]
		typeMapper  ITypeMapper
	}
	CaseClauseNode struct {
		arr   []Node[ast.Node]
		nodes []Node[ast.Node]
	}
	LhsToMultipleRhsOperator struct {
		LhsToRhsOp         token.Token
		betweenTerminalsOp token.Token // operator
		Lhs                Node[ast.Node]
		Rhs                []Node[ast.Node]
	}
	DictionaryExpression struct {
		m               reflect.Value
		defaultValue    Node[ast.Node]
		keyTypeMapper   ITypeMapper
		valueTypeMapper ITypeMapper
	}
	VaradicArgument struct {
		data []Node[ast.Node]
	}
	functionInformation struct {
		fn               OnCreateExecuteStatement
		funcType         Node[*ast.FuncType]
		funcTypeRequired bool
	}
	callExpressionInformation struct {
		fn OnCreateExecuteStatement
		ce Node[*ast.CallExpr]
	}

	TrailArray struct {
		arr        []Node[ast.Node]
		typeMapper ITypeMapper
	}
	FuncLit struct {
		Type       *ast.FuncType
		Body       *ast.BlockStmt
		values     map[string]ValueInformation
		typeMapper ITypeMapper
	}
	CallExpression struct {
		CallExpression Node[*ast.CallExpr]
		values         map[string]ValueInformation
	}
)

func (functList FuncLit) GetTypeMapper(state State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{functList.typeMapper}, true
}

func (c callExpressionInformation) Pos() token.Pos {
	return token.NoPos
}

func (c callExpressionInformation) End() token.Pos {
	return token.NoPos
}

func (c CallExpression) Pos() token.Pos {
	return token.NoPos
}

func (c CallExpression) End() token.Pos {
	return token.NoPos
}

func (entityField EntityField) GetValueKey() ValueKey {
	_, vk := entityField.typeMapper.ActualType()
	return vk
}

func (functList FuncLit) Pos() token.Pos {
	return functList.Body.Pos()
}

func (functList FuncLit) End() token.Pos {
	return functList.Body.End()
}

func (boolcond BooleanCondition) Pos() token.Pos {
	return token.NoPos
}

func (boolcond BooleanCondition) End() token.Pos {
	return token.NoPos
}

func (ta TrailArray) Pos() token.Pos {
	return token.NoPos
}

func (ta TrailArray) End() token.Pos {
	return token.NoPos
}

func (value TrailSource) trailMarker() {}

func (e EntitySource) sourceType() {}

func (s SingleValueCondition) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (s SingleValueCondition) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (c CheckForNotNullExpression) GetTypeMapper(state State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{c.typeMapper}, true
}

func (multiBinOp MultiBinaryExpr) GetTypeMapper(state State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{multiBinOp.typeMapper}, true
}

func (v VaradicArgument) Pos() token.Pos {
	return token.NoPos
}

func (v VaradicArgument) End() token.Pos {
	return token.NoPos
}

func (binOp BinaryExpr) GetValueKey() ValueKey {
	return ValueKey{"builtin", "BinaryExpr"}
}

func (value *TrailRecord) GetValueKey() ValueKey {
	return ValueKey{"builtin", "TrailRecord"}
}

func (e EntitySource) Pos() token.Pos {
	return token.NoPos
}

func (e EntitySource) End() token.Pos {
	return token.NoPos
}

func (e EntitySource) GetTypeMapper(State) (ITypeMapperArray, bool) {
	//TODO implement me
	panic("implement me")
}

func (iteSingleCondition IfThenElseSingleValueCondition) GetValueKey() ValueKey {
	return ValueKey{"builtin", "IfThenElseSingleValueCondition"}
}

func (coercion coercion) GetValueKey() ValueKey {
	return ValueKey{"builtin", "coercion"}
}

func (rv *ReflectValueExpression) GetTypeMapper(State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{&WrapReflectTypeInMapper{rv.Rv.Type(), rv.Vk}}, true
}

func (coercion coercion) GetTypeMapper(State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{&WrapReflectTypeInMapper{coercion.rt, coercion.vk}}, true
}

func (iteSingleCondition IfThenElseSingleValueCondition) GetTypeMapper(state State) (ITypeMapperArray, bool) {
	switch v := iteSingleCondition.conditionalStatement[0].value.Node.(type) {
	case BinaryExpr:
		return v.GetTypeMapper(state)
	case IfThenElseSingleValueCondition:
		return v.GetTypeMapper(state)
	case *ReflectValueExpression:
		return ITypeMapperArray{&WrapReflectTypeInMapper{v.Rv.Type(), v.Vk}}, true
	case IFindTypeMapper:
		panic("ffff")
	default:
		panic("dfdsfds")
	}
	panic("ddd")
}

func (de *DictionaryExpression) GetTypeMapper(State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{de.keyTypeMapper, de.valueTypeMapper}, true
}

func (entityField EntityField) GetTypeMapper(State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{entityField.typeMapper}, true
}

func (value *TrailRecord) GetTypeMapper(State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{value.typeMapper}, true
}

func (value TrailSource) GetTypeMapper(State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{value.typeMapper}, true
}

func (c CheckForNotNullExpression) Pos() token.Pos {
	return token.NoPos
}

func (c CheckForNotNullExpression) End() token.Pos {
	return token.NoPos
}

//	func (rv *NilValueExpression) Pos() token.Pos {
//		return token.NoPos
//	}
func (multiBinOp MultiBinaryExpr) Pos() token.Pos {
	return token.NoPos
}
func (ccn *CaseClauseNode) Pos() token.Pos {
	return token.NoPos
}
func (lhsRhsOperator LhsToMultipleRhsOperator) Pos() token.Pos {
	return token.NoPos
}
func (de *DictionaryExpression) Pos() token.Pos {
	return token.NoPos
}
func (iteSingleCondition IfThenElseSingleValueCondition) Pos() token.Pos {
	return token.NoPos
}
func (ite IfThenElseMultiValueCondition) Pos() token.Pos {
	return token.NoPos
}
func (supportedFunction SupportedFunction) Pos() token.Pos {
	return token.NoPos
}
func (rv *ReflectValueExpression) Pos() token.Pos {
	return token.NoPos
}
func (value *TrailRecord) Pos() token.Pos {
	return value.Position
}
func (value TrailSource) Pos() token.Pos {
	return token.NoPos
}
func (entityField EntityField) Pos() token.Pos {
	return token.NoPos
}
func (coercion coercion) Pos() token.Pos {
	return coercion.Position
}
func (nilExpression *builtInNil) Pos() token.Pos {
	return token.NoPos
}

func (binOp BinaryExpr) Pos() token.Pos {
	return token.NoPos
}

func (binOp BinaryExpr) End() token.Pos {
	return token.NoPos
}

func (binOp BinaryExpr) GetTypeMapper(State) (ITypeMapperArray, bool) {
	return ITypeMapperArray{binOp.typeMapper}, true
}

func (nilExpression *builtInNil) End() token.Pos {
	return token.NoPos
}
func (coercion coercion) End() token.Pos {
	return coercion.Position
}
func (entityField EntityField) End() token.Pos {
	return token.NoPos
}
func (value TrailSource) End() token.Pos {
	return token.NoPos
}
func (value *TrailRecord) End() token.Pos {
	return value.Position
}
func (rv *ReflectValueExpression) End() token.Pos {
	return token.NoPos
}
func (supportedFunction SupportedFunction) End() token.Pos {
	return token.NoPos
}
func (ite IfThenElseMultiValueCondition) End() token.Pos {
	return token.NoPos
}
func (iteSingleCondition IfThenElseSingleValueCondition) End() token.Pos {
	return token.NoPos
}
func (de *DictionaryExpression) End() token.Pos {
	return token.NoPos
}
func (lhsRhsOperator LhsToMultipleRhsOperator) End() token.Pos {
	return token.NoPos
}
func (ccn *CaseClauseNode) End() token.Pos {
	return token.NoPos
}

//	func (rv *NilValueExpression) End() token.Pos {
//		return token.NoPos
//	}
func (multiBinOp MultiBinaryExpr) End() token.Pos {
	return token.NoPos
}

type IExpand interface {
	Expand(parentNode Node[ast.Node]) []Node[ast.Node]
}

func (ite IfThenElseMultiValueCondition) Expand(parentNode Node[ast.Node]) []Node[ast.Node] {
	var arr [][]SingleValueCondition
	for range ite.conditionalStatement[0].values {
		arr = append(arr, []SingleValueCondition{})
	}

	for partialAnswerIdx, _ := range arr {
		for _, stmt := range ite.conditionalStatement {
			idxNode := stmt.values[partialAnswerIdx]
			arr[partialAnswerIdx] = append(arr[partialAnswerIdx], SingleValueCondition{condition: stmt.condition, value: idxNode})
		}
	}

	var result []Node[ast.Node]
	for idx, _ := range ite.conditionalStatement[0].values {
		p01 := ChangeParamNode[ast.Node, ast.Node](parentNode, IfThenElseSingleValueCondition{arr[idx]})
		result = append(result, p01)
	}
	return result
}

func (rv *ReflectValueExpression) String() string {
	return rv.Rv.String()
}

func (f functionInformation) Pos() token.Pos {
	return token.NoPos
}

func (f functionInformation) End() token.Pos {
	return token.NoPos
}
