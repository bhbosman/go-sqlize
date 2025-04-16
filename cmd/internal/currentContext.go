package internal

import (
	"go/ast"
)

type TypeMapper map[string]ITypeMapper

type LocalTypesMap map[string]OnCreateType

type ValueInformation struct {
	Node Node[ast.Node]
}
type ValueInformationMap map[string]ValueInformation
type CurrentContext struct {
	Mm         ValueInformationMap
	TypeParams map[string]ITypeMapper
	LocalTypes LocalTypesMap
	Parent     *CurrentContext
}

func (self *CurrentContext) flattenTypeParams() map[string]ITypeMapper {
	result := func() map[string]ITypeMapper {
		if self.Parent == nil {
			return map[string]ITypeMapper{}
		} else {
			return self.Parent.flattenTypeParams()
		}
	}()

	for key, value := range self.TypeParams {
		result[key] = value
	}
	return result
}

func (self *CurrentContext) addLocalTypes(key string, value OnCreateType) {
	self.LocalTypes[key] = value
}

func (self *CurrentContext) findLocalType(key string) (OnCreateType, bool) {
	if v, ok := self.LocalTypes[key]; ok {
		return v, true
	}
	if self.Parent == nil {
		return nil, false
	}
	return self.Parent.findLocalType(key)
}

func (self *CurrentContext) FindTypeFromNode(node Node[ast.Node]) (ITypeMapperArray, bool) {
	if _, fromNode, b := self.internalFindTypeFromNode(node); b {
		return fromNode, true
	}
	return nil, false
}

func (self *CurrentContext) internalFindTypeFromNode(node Node[ast.Node]) (Node[ast.Node], ITypeMapperArray, bool) {
	switch item := node.Node.(type) {
	case *ast.Ident:
		if value, b := self.FindValueByString(item.Name); b {
			if findTypeMapper, ok := value.Node.(IFindTypeMapper); ok {
				if mapper, ok := findTypeMapper.GetTypeMapper(""); ok {
					return value, mapper, true
				}
			}
		}
		return Node[ast.Node]{}, nil, false
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		if _, mapper, b := self.internalFindTypeFromNode(param); b {
			var typeMapperArray ITypeMapperArray
			for _, typeMapper := range mapper {
				rt := typeMapper.ActualType()
				sf, _ := rt.FieldByName(item.Sel.Name)
				typeMapperArray = append(typeMapperArray, &WrapReflectTypeInMapper{sf.Type})
			}
			return Node[ast.Node]{}, typeMapperArray, true
		}
	}
	panic("implement me")
}

func (self *CurrentContext) FindValueByString(value string) (Node[ast.Node], bool) {
	if v, ok := self.Mm[value]; ok {
		return v.Node, true
	}
	if self.Parent == nil {
		return Node[ast.Node]{}, false
	}
	return self.Parent.FindValueByString(value)
}

func (self *CurrentContext) FindValueByNode(node Node[ast.Node]) (Node[ast.Node], bool) {
	if byNode, b := self.Parent.internalFindValueByNode(node); b {
		return byNode.(Node[ast.Node]), true
	}
	return Node[ast.Node]{}, false
}

func (self *CurrentContext) internalFindValueByNode(node Node[ast.Node]) (interface{}, bool) {
	switch item := node.Node.(type) {
	default:
		return nil, false
	case *ast.Ident:
		return self.FindValueByString(item.Name)
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		if unk, ok := self.internalFindValueByNode(param); ok {
			switch unkItem := unk.(type) {
			default:
				panic("implement me")
			case Node[ast.Node]:
				switch nodeItem := unkItem.Node.(type) {
				default:
					panic("implement me")
				case *TrailRecord:
					return nodeItem.Value.FieldByName(item.Sel.Name).Interface().(Node[ast.Node]), true
				}
			}
		}
		return nil, false
	}
}

func (self *CurrentContext) FindTypeParam(value string) (ITypeMapper, bool) {
	if v, ok := self.TypeParams[value]; ok {
		return v, true
	}
	if self.Parent == nil {
		return nil, false
	}
	return self.Parent.FindTypeParam(value)
}
