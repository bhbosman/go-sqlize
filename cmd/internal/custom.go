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

//func (value *Value) exprNode() {
//}

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
}

func (coercion *coercion) Pos() token.Pos {
	return coercion.Position
}

func (coercion *coercion) End() token.Pos {
	return coercion.Position
}

type BinaryExpr struct {
	OpPos  token.Pos        // position of Op
	Op     token.Token      // operator
	Values []Node[ast.Node] // right operand
}

func (binop *BinaryExpr) Pos() token.Pos {
	return binop.OpPos
}

func (binop *BinaryExpr) End() token.Pos {
	return binop.OpPos
}
