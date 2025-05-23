package internal

import (
	"fmt"
	"go/ast"
	"go/token"
)

type (
	queryState struct {
		top      int
		distinct bool
	}
	IQueryOptions interface {
		ast.Node
		IIsLiterateValue
		Apply(*queryState)
	}
	queryTop struct {
		count int
	}
	queryDistinct struct {
	}
)

func (qd queryDistinct) Apply(state *queryState) {
	state.distinct = true
}

func (qt queryTop) Apply(state *queryState) {
	state.top = qt.count
}

func (qd queryDistinct) ThisIsALiterateValue() {

}

func (qt queryTop) ThisIsALiterateValue() {
}

func (qd queryDistinct) Pos() token.Pos {
	return token.NoPos
}

func (qd queryDistinct) End() token.Pos {
	return token.NoPos
}

func (qt queryTop) Pos() token.Pos {
	return token.NoPos
}

func (qt queryTop) End() token.Pos {
	return token.NoPos
}

func (compiler *Compiler) libQueryImplementation(state State, funcTypeNode Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(typeParams) != 1 {
			panic(fmt.Errorf("Lib.Query implementation requires 1 type argument, got %d", len(typeParams)))
		}
		typeMapper := typeParams[funcTypeNode.Node.TypeParams.List[0].Names[0].Name]

		var queryOptions []IQueryOptions
		for _, arg := range arguments[0:] {
			rv, _ := isLiterateValue(arg)
			queryOptions = append(queryOptions, rv.Interface().(IQueryOptions))
		}
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, compiler.query(typeMapper, queryOptions...))}, artValue
	}
}

func (compiler *Compiler) libQueryTopImplementation(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		qt := queryTop{int(arguments[0].Node.(*ReflectValueExpression).Rv.Int())}
		return []Node[ast.Node]{ChangeParamNode[*ast.FuncType, ast.Node](node, qt)}, artValue
	}
}

func (compiler *Compiler) libQueryDistinctImplementation(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		qd := queryDistinct{}
		return []Node[ast.Node]{ChangeParamNode[*ast.FuncType, ast.Node](node, qd)}, artValue
	}
}

func (compiler *Compiler) query(typeMapper ITypeMapper, options ...IQueryOptions) ITrailMarker {
	qs := &queryState{}
	for _, option := range options {
		option.Apply(qs)
	}
	alias := compiler.AddEntitySource(typeMapper, *qs)
	trailSource := TrailSource{alias, typeMapper}
	return trailSource
}
