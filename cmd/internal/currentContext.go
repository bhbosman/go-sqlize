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

func (self *CurrentContext) FindTypeFromNode(node Node[ast.Node]) ([]ITypeMapper, bool) {
	if _, fromNode, b := self.internalFindTypeFromNode(node); b {
		return fromNode, true
	}
	return nil, false
}

func (self *CurrentContext) internalFindTypeFromNode(node Node[ast.Node]) (Node[ast.Node], []ITypeMapper, bool) {
	switch item := node.Node.(type) {
	case *ast.Ident:
		if value, b := self.FindValue(item.Name); b {
			if findTypeMapper, ok := value.Node.(IFindTypeMapper); ok {
				if mapper, ok := findTypeMapper.GetTypeMapper(); ok {
					return value, mapper, true
				}
			}
		}
		return Node[ast.Node]{}, nil, false
	case *ast.SelectorExpr:
		param := ChangeParamNode[ast.Node, ast.Node](node, item.X)
		if _, mapper, b := self.internalFindTypeFromNode(param); b {
			var typeMapperArray []ITypeMapper
			for _, typeMapper := range mapper {
				rt := typeMapper.ActualType(State{})
				sf, _ := rt.FieldByName(item.Sel.Name)
				typeMapperArray = append(typeMapperArray, &WrapReflectTypeInMapper{sf.Type})
			}
			return Node[ast.Node]{}, typeMapperArray, true
		}
	}
	panic("implement me")
}

func (self *CurrentContext) FindValue(value string) (Node[ast.Node], bool) {
	if v, ok := self.Mm[value]; ok {
		return v.Node, true
	}
	if self.Parent == nil {
		return Node[ast.Node]{}, false
	}
	return self.Parent.FindValue(value)
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
