package internal

import "go/ast"

type TypeMapper map[string]ITypeMapper

type LocalTypesMap map[string]OnCreateType

type CurrentContext struct {
	Mm         map[string]Node[ast.Node]
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

func (self *CurrentContext) FindValue(value string) (Node[ast.Node], bool) {
	if v, ok := self.Mm[value]; ok {
		return v, true
	}
	if self.Parent == nil {
		return Node[ast.Node]{}, false
	}
	return self.Parent.FindValue(value)
}
