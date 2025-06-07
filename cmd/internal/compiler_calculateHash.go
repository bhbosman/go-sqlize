package internal

import (
	"go/ast"
	"hash"
	"hash/crc32"
	"reflect"
	"strconv"
)

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
	case BooleanCondition:
		hash.Write([]byte{byte(item.op)})
		for _, condition := range item.conditions {
			compiler.internalCalculateHash(hash, condition)
		}
	case EntityField:
		hash.Write([]byte(item.alias))
		hash.Write([]byte(item.field))
		p01 := ChangeParamNode[ast.Node, ast.Node](node, item.typeMapper)
		compiler.internalCalculateHash(hash, p01)
		hash.Write([]byte(item.field))
	case *ReflectValueExpression:
		kind := item.Rv.Kind()
		switch kind {
		default:
			panic(item.Rv.Kind())
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			hash.Write([]byte(strconv.FormatInt(item.Rv.Int(), 10)))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			hash.Write([]byte(strconv.FormatUint(item.Rv.Uint(), 10)))
		case reflect.Float32, reflect.Float64:
			hash.Write([]byte(strconv.FormatFloat(item.Rv.Float(), 'g', -1, 64)))
		case reflect.String:
			hash.Write([]byte(item.Rv.String()))
		case reflect.Bool:
			hash.Write([]byte(strconv.FormatBool(item.Rv.Bool())))
		}
		hash.Write([]byte(item.Rv.Type().String()))
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
		p01 := ChangeParamNode[ast.Node, ast.Node](node, item.typeMapper)
		compiler.internalCalculateHash(hash, p01)
	case *CheckForNotNullExpression:
		compiler.internalCalculateHash(hash, item.node)
		p01 := ChangeParamNode[ast.Node, ast.Node](node, item.typeMapper)
		compiler.internalCalculateHash(hash, p01)
	case MultiBinaryExpr:
		hash.Write([]byte{byte(item.Op)})
		for _, expression := range item.expressions {
			compiler.internalCalculateHash(hash, expression)
		}
		p01 := ChangeParamNode[ast.Node, ast.Node](node, item.typeMapper)
		compiler.internalCalculateHash(hash, p01)
	case IfThenElseSingleValueCondition:
		for _, condition := range item.conditionalStatement {
			p01 := ChangeParamNode[ast.Node, ast.Node](node, condition)
			compiler.internalCalculateHash(hash, p01)
		}
	case SingleValueCondition:
		compiler.internalCalculateHash(hash, item.condition)
		compiler.internalCalculateHash(hash, item.value)
	}
}
