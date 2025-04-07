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

type ConditionalStatement struct {
	condition Node[ast.Node]
	values    []Node[ast.Node]
}

type IfThenElseCondition struct {
	conditionalStatement []ConditionalStatement
}

func (ite *IfThenElseCondition) Pos() token.Pos {
	return token.NoPos
}

func (ite *IfThenElseCondition) End() token.Pos {
	return token.NoPos
}

type PartialExpression struct {
	conditionalStatement []struct {
		condition Node[ast.Node]
		value     Node[ast.Node]
	}
}

func (partialExpression *PartialExpression) IsValidNode() bool {
	for _, conditionalStatement := range partialExpression.conditionalStatement {
		if rv, b := isLiterateValue(conditionalStatement.condition); b && rv.Kind() == reflect.Bool && rv.Bool() {
			return true
		}
	}
	return false
}

func (partialExpression *PartialExpression) Pos() token.Pos {
	return token.NoPos
}

func (partialExpression *PartialExpression) End() token.Pos {
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
