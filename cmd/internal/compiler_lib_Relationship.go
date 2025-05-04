package internal

import (
	"go/ast"
	"go/token"
	"hash"
	"hash/crc32"
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
		funcLit, _ := arguments[1].Node.(FuncLit)
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

func (compiler *Compiler) calculateHash(node Node[ast.Node]) uint32 {
	hashCalculator := crc32.NewIEEE()
	compiler.internalCalculateHash(hashCalculator, node)
	return hashCalculator.Sum32()
}

func (compiler *Compiler) internalCalculateHash(hash hash.Hash, node Node[ast.Node]) {
	switch item := node.Node.(type) {
	default:
		panic(item)
		panic("internalCalculateHash: unknown node type")
	case EntityField:
		hash.Write([]byte(item.alias))
		hash.Write([]byte(item.field))
		compiler.internalCalculateHash(hash, ChangeParamNode[ast.Node, ast.Node](node, item.typeMapper))
		hash.Write([]byte(item.field))
	case *ReflectValueExpression:
		hash.Write([]byte(item.Rv.String()))
		hash.Write([]byte(item.Vk.String()))
	case ITypeMapper:
		actualType, key := item.ActualType()
		hash.Write([]byte(actualType.String()))
		hash.Write([]byte(key.String()))
		hash.Write([]byte{byte(item.Kind())})
	case BinaryExpr:
		hash.Write([]byte{byte(item.Op)})
		compiler.internalCalculateHash(hash, item.left)
		compiler.internalCalculateHash(hash, item.right)
		compiler.internalCalculateHash(hash, ChangeParamNode[ast.Node, ast.Node](node, item.typeMapper))
	case *CheckForNotNullExpression:
	case MultiBinaryExpr:
		hash.Write([]byte{byte(item.Op)})
		for _, expression := range item.expressions {
			compiler.internalCalculateHash(hash, expression)
		}
		compiler.internalCalculateHash(hash, ChangeParamNode[ast.Node, ast.Node](node, item.typeMapper))
	case IfThenElseSingleValueCondition:
		for _, condition := range item.conditionalStatement {
			compiler.internalCalculateHash(hash, ChangeParamNode[ast.Node, ast.Node](node, condition))
		}
	case SingleValueCondition:
		compiler.internalCalculateHash(hash, item.condition)
		compiler.internalCalculateHash(hash, item.value)
	}
}
