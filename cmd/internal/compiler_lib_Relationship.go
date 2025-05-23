package internal

import (
	"fmt"
	"go/ast"
	"go/token"
)

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

func (compiler *Compiler) libCoreRelationshipImplementation(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		if len(arguments) != 2 {
			panic(fmt.Errorf("Lib.CoreRelationship implementation requires 2 arguments, got %d", len(arguments)))
		}
		return compiler.internalLibCoreRelationshipImplementation(state, typeParams, arguments[0], arguments[1])
	}
}

func (compiler *Compiler) internalLibCoreRelationshipImplementation(state State, typeParams map[string]ITypeMapper, arg0, arg1 Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
	switch argItem := arg1.Node.(type) {
	default:
		panic(argItem)
	case *ast.Ident, *ast.FuncLit:
		fn := compiler.findRhsExpression(state, arg1)
		v, _ := compiler.executeAndExpandStatement(state, typeParams, nil, fn)
		return compiler.internalLibCoreRelationshipImplementation(state, typeParams, arg0, v[0])
	case FuncLit:
		param := ChangeParamNode[ast.Node, FuncLit](arg1, argItem)
		callback, _ := compiler.executeFuncLit(state, param, []Node[ast.Node]{arg0}, typeParams)

		//var relationshipOpt []IRelationshipOpt
		//for _, arg := range arguments[2:] {
		//	rv, _ := isLiterateValue(arg)
		//	relationshipOpt = append(relationshipOpt, rv.Interface().(IRelationshipOpt))
		//}

		pp := compiler.coreRelationship(state, arg0.Node.(ITrailMarker), callback[0])
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, pp)}, artValue

	}

	//funcLit, _ := arg1.Node.(FuncLit)
	//p01 := ChangeParamNode(arg1, funcLit)
	//callback, _ := compiler.executeFuncLit(state, p01, []Node[ast.Node]{arg0}, typeParams)

	//var relationshipOpt []IRelationshipOpt
	//for _, arg := range arguments[2:] {
	//	rv, _ := isLiterateValue(arg)
	//	relationshipOpt = append(relationshipOpt, rv.Interface().(IRelationshipOpt))
	//}

	//pp := compiler.coreRelationship(state, arguments[0].Node.(ITrailMarker), callback[0], relationshipOpt...)
	//return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, pp)}, artValue

}

func (compiler *Compiler) coreRelationship(state State, from ITrailMarker, callback Node[ast.Node], relationshipOpt ...IRelationshipOpt) ITrailMarker {
	fn := func() (map[string]ISource, string, bool) {
		hashValue := compiler.calculateHash(callback)
		for key, value := range compiler.JoinInformation {
			itemHashValue := compiler.calculateHash(value.condition)
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
			joinInformation := JoinInformation{item.Alias, ss, callback, jtInner}
			compiler.JoinInformation[item.Alias] = joinInformation
			return from
		}
	}
}
