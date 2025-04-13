package internal

import (
	"go/ast"
	"go/token"
	"reflect"
)

type (
	EntitySource struct {
		rt reflect.Type
	}
	TrailRecord struct {
		Position token.Pos     // identifier position
		Value    reflect.Value // identifier name
	}
	TrailSource struct {
		Position token.Pos // identifier position
		Alias    string
	}
	EntityField struct {
		Position token.Pos // identifier position
		alias    string
		field    string
	}
	coercion struct {
		Position token.Pos
		to       string
		Node     Node[ast.Node]
		rt       reflect.Type
	}

	CheckForNotNullExpression struct {
		node Node[ast.Node]
	}
	BinaryExpr struct {
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		left  Node[ast.Node]
		right Node[ast.Node]
	}
	builtInNil struct {
	}
	ReflectValueExpression struct {
		Rv reflect.Value
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
	IfThenElseSingleValueCondition struct {
		conditionalStatement []SingleValueCondition
	}
	MultiBinaryExpr struct {
		Op          token.Token // operator
		expressions []Node[ast.Node]
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
		m            reflect.Value
		defaultValue reflect.Value
	}
)

func (c *CheckForNotNullExpression) Pos() token.Pos {
	return token.NoPos
}

func (c *CheckForNotNullExpression) End() token.Pos {
	return token.NoPos
}

//	func (rv *NilValueExpression) Pos() token.Pos {
//		return token.NoPos
//	}
func (multiBinOp *MultiBinaryExpr) Pos() token.Pos {
	return token.NoPos
}
func (ccn *CaseClauseNode) Pos() token.Pos {
	return token.NoPos
}
func (lhsRhsOperator *LhsToMultipleRhsOperator) Pos() token.Pos {
	return token.NoPos
}
func (de *DictionaryExpression) Pos() token.Pos {
	return token.NoPos
}
func (iteSingleCondition *IfThenElseSingleValueCondition) Pos() token.Pos {
	return token.NoPos
}
func (ite *IfThenElseMultiValueCondition) Pos() token.Pos {
	return token.NoPos
}
func (supportedFunction *SupportedFunction) Pos() token.Pos {
	return token.NoPos
}
func (rv *ReflectValueExpression) Pos() token.Pos {
	return token.NoPos
}
func (value *TrailRecord) Pos() token.Pos {
	return value.Position
}
func (value *TrailSource) Pos() token.Pos {
	return value.Position
}
func (entityField *EntityField) Pos() token.Pos {
	return entityField.Position
}
func (coercion *coercion) Pos() token.Pos {
	return coercion.Position
}
func (nilExpression *builtInNil) Pos() token.Pos {
	return token.NoPos
}
func (binop *BinaryExpr) Pos() token.Pos {
	return binop.OpPos
}

func (binop *BinaryExpr) End() token.Pos {
	return binop.OpPos
}
func (nilExpression *builtInNil) End() token.Pos {
	return token.NoPos
}
func (coercion *coercion) End() token.Pos {
	return coercion.Position
}
func (entityField *EntityField) End() token.Pos {
	return entityField.Position
}
func (value *TrailSource) End() token.Pos {
	return value.Position
}
func (value *TrailRecord) End() token.Pos {
	return value.Position
}
func (rv *ReflectValueExpression) End() token.Pos {
	return token.NoPos
}
func (supportedFunction *SupportedFunction) End() token.Pos {
	return token.NoPos
}
func (ite *IfThenElseMultiValueCondition) End() token.Pos {
	return token.NoPos
}
func (iteSingleCondition *IfThenElseSingleValueCondition) End() token.Pos {
	return token.NoPos
}
func (de *DictionaryExpression) End() token.Pos {
	return token.NoPos
}
func (lhsRhsOperator *LhsToMultipleRhsOperator) End() token.Pos {
	return token.NoPos
}
func (ccn *CaseClauseNode) End() token.Pos {
	return token.NoPos
}

//	func (rv *NilValueExpression) End() token.Pos {
//		return token.NoPos
//	}
func (multiBinOp *MultiBinaryExpr) End() token.Pos {
	return token.NoPos
}

type IExpand interface {
	Expand(parentNode Node[ast.Node]) []Node[ast.Node]
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

func (rv *ReflectValueExpression) String() string {
	return rv.Rv.String()
}
