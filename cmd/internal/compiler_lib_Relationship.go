package internal

import (
	"fmt"
	"go/ast"
	"go/token"
)

func (compiler *Compiler) libCoreRelationshipImpl(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 2 {
			panic(fmt.Errorf("Lib.CoreRelationship implementation requires 2 arguments, got %d", len(arguments)))
		}
		callback, _ := compiler.internalLibCoreRelationshipImplementation(state, typeParams, arguments[0], arguments[1])
		pp := compiler.internalRelation(state, arguments[0].Node.(ITrailMarker), callback[0], jtLeftInner)
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, pp)}, artValue
	}
}

func (compiler *Compiler) libCoreOptRelationshipImpl(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 2 {
			panic(fmt.Errorf("Lib.CoreRelationship implementation requires 2 arguments, got %d", len(arguments)))
		}
		callback, _ := compiler.internalLibCoreRelationshipImplementation(state, typeParams, arguments[0], arguments[1])
		pp := compiler.internalRelation(state, arguments[0].Node.(ITrailMarker), callback[0], jtLeftOuter)
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, pp)}, artValue
	}
}

func (compiler *Compiler) internalRelation(state State, from ITrailMarker, callback Node[ast.Node], joinType joinType, relationshipOpt ...IRelationshipOpt) ITrailMarker {
	fn := func() (map[string]ISource, string, bool) {
		hashValue := compiler.calculateHash(callback)
		for key, value := range compiler.JoinInformation {
			p0 := ChangeParamNode[BooleanCondition, ast.Node](value.condition, value.condition.Node)
			itemHashValue := compiler.calculateHash(p0)
			if hashValue == itemHashValue {
				return value.rhs, key, true
			}
		}
		return nil, "", false
	}
	ss, key, found := fn()

	if !found {
		ss = compiler.findSourcesFromNode(callback)
		switch fromItem := from.(type) {
		default:
			panic(fromItem)
		case TrailSource:
			delete(ss, fromItem.Alias)
		}
	}

	if found {
		switch item := compiler.Sources[key].(type) {
		default:
			panic("dddd")
		case *EntitySource:
			return TrailSource{key, item.typeMapper}
		}
	} else {
		switch item := from.(type) {
		default:
			panic("dddd")
		case TrailSource:
			booleanExpression := compiler.transformToBooleanExpression(state, token.LOR, callback)
			joinInformation := JoinInformation{item.Alias, ss, booleanExpression, joinType}
			compiler.JoinInformation[item.Alias] = joinInformation
			return from
		}
	}
}

type (
	relationshipState struct {
		top      int
		distinct bool
	}
	IRelationshipOpt interface {
		ast.Node
		IIsLiterateValue
		Apply(*relationshipState)
	}
	relationshipJoinType struct {
	}
)

func (relationJoinType relationshipJoinType) Pos() token.Pos {
	return token.NoPos
}

func (relationJoinType relationshipJoinType) End() token.Pos {
	return token.NoPos
}

func (relationJoinType relationshipJoinType) ThisIsALiterateValue() {
}

func (relationJoinType relationshipJoinType) Apply(state *relationshipState) {
}

func (compiler *Compiler) internalLibCoreRelationshipImplementation(state State, typeParams map[string]ITypeMapper, arg0, arg1 Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	switch argItem := arg1.Node.(type) {
	default:
		panic(argItem)
	case MultiBinaryExpr:
		m := map[uint32]bool{}
		var arr []Node[ast.Node]
		for _, expression := range argItem.expressions {
			v, _ := compiler.internalLibCoreRelationshipImplementation(state, typeParams, arg0, expression)
			boolExpression := compiler.transformToBooleanExpression(state, token.LOR, v[0])
			p0 := ChangeParamNode[BooleanCondition, ast.Node](boolExpression, boolExpression.Node)
			hashValue := compiler.calculateHash(p0)
			if _, ok := m[hashValue]; !ok {
				m[hashValue] = true
				arr = append(arr, p0)
			}
		}

		bcNode := ChangeParamNode[ast.Node, ast.Node](arg1, BooleanCondition{arr, argItem.Op})
		return []Node[ast.Node]{bcNode}, artValue
	case *ast.Ident, *ast.FuncLit, *ast.CallExpr:
		fn := compiler.findRhsExpression(state, arg1)
		v, _ := compiler.executeAndExpandStatement(state, typeParams, nil, fn)
		return compiler.internalLibCoreRelationshipImplementation(state, typeParams, arg0, v[0])
	case FuncLit:
		param := ChangeParamNode[ast.Node, FuncLit](arg1, argItem)
		return compiler.executeFuncLit(state, param, []Node[ast.Node]{arg0}, typeParams)
	}
}
