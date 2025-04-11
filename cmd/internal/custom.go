package internal

import (
	"go/ast"
	"go/token"
	"reflect"
)

type TrailRecord struct {
	Position token.Pos     // identifier position
	Value    reflect.Value // identifier name
}

func (value *TrailRecord) Pos() token.Pos {
	return value.Position
}

func (value *TrailRecord) End() token.Pos {
	return value.Position
}

type TrailSource struct {
	Position token.Pos // identifier position
	Alias    string
}

func (value *TrailSource) Pos() token.Pos {
	return value.Position
}

func (value *TrailSource) End() token.Pos {
	return value.Position
}

type EntityField struct {
	Position token.Pos // identifier position
	alias    string
	field    string
}

func (entityField *EntityField) Pos() token.Pos {
	return entityField.Position
}

func (entityField *EntityField) End() token.Pos {
	return entityField.Position
}

type coercion struct {
	Position token.Pos
	to       string
	Node     Node[ast.Node]
	rt       reflect.Type
}

func (coercion *coercion) Pos() token.Pos {
	return coercion.Position
}

func (coercion *coercion) End() token.Pos {
	return coercion.Position
}

type BinaryExpr struct {
	OpPos token.Pos   // position of Op
	Op    token.Token // operator
	left  Node[ast.Node]
	right Node[ast.Node]
}

func (binop *BinaryExpr) Pos() token.Pos {
	return binop.OpPos
}

func (binop *BinaryExpr) End() token.Pos {
	return binop.OpPos
}

type nullValue struct{}

type NilExpression struct {
}

func (nilExpression *NilExpression) Pos() token.Pos {
	return token.NoPos
}

func (nilExpression *NilExpression) End() token.Pos {
	return token.NoPos
}

type ReflectValueExpression struct {
	Rv reflect.Value
}

func (rv *ReflectValueExpression) String() string {
	return rv.Rv.String()
}

func (rv *ReflectValueExpression) Pos() token.Pos {
	return token.NoPos
}

func (rv *ReflectValueExpression) End() token.Pos {
	return token.NoPos
}

type SupportedFunction struct {
	functionName string
	params       []Node[ast.Node]
	rt           reflect.Type
}

func (supportedFunction *SupportedFunction) Pos() token.Pos {
	return token.NoPos
}

func (supportedFunction *SupportedFunction) End() token.Pos {
	return token.NoPos
}

type MultiValueCondition struct {
	condition Node[ast.Node]
	values    []Node[ast.Node]
}

type SingleValueCondition struct {
	condition Node[ast.Node]
	value     Node[ast.Node]
}

type IfThenElseMultiValueCondition struct {
	conditionalStatement []MultiValueCondition
}

func (ite *IfThenElseMultiValueCondition) Expand(parentNode Node[ast.Node]) []Node[ast.Node] {
	var result []Node[ast.Node]
	for range ite.conditionalStatement[0].values {
		result = append(result, ChangeParamNode[ast.Node, ast.Node](parentNode, &IfThenElseSingleValueCondition{}))
	}

	for partialAnswerIdx, partialAnswer := range result {
		if partialAnswerNode, ok := partialAnswer.Node.(*IfThenElseSingleValueCondition); ok {
			for _, stmt := range ite.conditionalStatement {
				idxNode := stmt.values[partialAnswerIdx]
				partialAnswerNode.conditionalStatement = append(partialAnswerNode.conditionalStatement, SingleValueCondition{condition: stmt.condition, value: idxNode})
			}
		}
	}
	return result
}

func (ite *IfThenElseMultiValueCondition) Pos() token.Pos {
	return token.NoPos
}

func (ite *IfThenElseMultiValueCondition) End() token.Pos {
	return token.NoPos
}

type IfThenElseSingleValueCondition struct {
	conditionalStatement []SingleValueCondition
}

func (iteSingleCondition *IfThenElseSingleValueCondition) Pos() token.Pos {
	return token.NoPos
}

func (iteSingleCondition *IfThenElseSingleValueCondition) End() token.Pos {
	return token.NoPos
}

type NilValueExpression struct {
}

func (rv *NilValueExpression) Pos() token.Pos {
	return token.NoPos
}

func (rv *NilValueExpression) End() token.Pos {
	return token.NoPos
}

type MultiBinaryExpr struct {
	Op          token.Token // operator
	expressions []Node[ast.Node]
}

func (multiBinOp *MultiBinaryExpr) Pos() token.Pos {
	return token.NoPos
}

func (multiBinOp *MultiBinaryExpr) End() token.Pos {
	return token.NoPos
}

type IExpand interface {
	Expand(parentNode Node[ast.Node]) []Node[ast.Node]
}
