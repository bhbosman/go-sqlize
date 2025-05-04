package internal

import (
	"go/ast"
	"go/token"
)

type (
	relationshipState struct {
		top      int
		distinct bool
	}
	iRelationshipOpt interface {
		ast.Node
		iIsLiterateValue
		apply(*relationshipState)
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

func (relationJoinType relationshipJoinType) thisIsALiterateValue() {
}

func (relationJoinType relationshipJoinType) apply(state *relationshipState) {
}

func (compiler *Compiler) libCoreRelationshipImplementation(state State, node Node[*ast.FuncType]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, arguments []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		funcLit, _ := arguments[1].Node.(*ast.FuncLit)
		callback, _ := compiler.executeFuncLit(state, ChangeParamNode(arguments[1], funcLit), arguments, typeParams)

		var relationshipOpt []iRelationshipOpt
		for _, arg := range arguments[2:] {
			rv, _ := isLiterateValue(arg)
			relationshipOpt = append(relationshipOpt, rv.Interface().(iRelationshipOpt))
		}
		return []Node[ast.Node]{
			ChangeParamNode[ast.Node, ast.Node](
				state.currentNode,
				compiler.coreRelationship(state, arguments[0].Node.(ITrailMarker), callback[0], relationshipOpt...),
			),
		}, artValue
	}
}

func (compiler *Compiler) coreRelationship(state State, from ITrailMarker, callback Node[ast.Node], relationshipOpt ...iRelationshipOpt) ITrailMarker {
	ss, key, found := func() (map[string]ISource, string, bool) {
		hashValue := compiler.calculateHash(callback)
		for key, value := range compiler.JoinInformation {
			itemHashValue := compiler.calculateHash(value.condition)
			if hashValue == itemHashValue {
				return value.rhs, key, true
			}
		}
		return nil, "", false
	}()
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
