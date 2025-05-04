package internal

import (
	"go/ast"
	"hash"
	"hash/crc32"
)

func (compiler *Compiler) calculateHash(node Node[ast.Node]) uint32 {
	hashCalculator := crc32.NewIEEE()
	compiler.internalCalculateHash(hashCalculator, node)
	return hashCalculator.Sum32()
}

func (compiler *Compiler) internalCalculateHash(hash hash.Hash, node Node[ast.Node]) {
	switch item := node.Node.(type) {
	default:
		panic("internalCalculateHash: unknown node type")
	case EntityField:
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
