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
	Temp       bool
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

func (self *CurrentContext) flattenVariables() map[string]ValueInformation {
	result := func() map[string]ValueInformation {
		if self.Parent == nil {
			return map[string]ValueInformation{}
		} else {
			return self.Parent.flattenVariables()
		}
	}()

	for key, value := range self.Mm {
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

func (self *CurrentContext) ReplaceRoot(root *CurrentContext) {
	if self.Parent != nil {
		if !self.Parent.Temp {
			self.Parent.ReplaceRoot(root)
		} else {
			root.Parent = self.Parent
			self.Parent = root
		}
	} else {
		self.Parent = root
	}
}

func (self *CurrentContext) RemoveRoot(root *CurrentContext) {
	if self.Parent != nil {
		if self.Parent != root {
			self.Parent.RemoveRoot(root)
		} else {
			self.Parent = root.Parent
		}
	}
}
