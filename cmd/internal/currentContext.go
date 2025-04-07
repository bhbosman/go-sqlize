package internal

import "go/ast"

//type ICurrentContext interface {
//	currentContext() ICurrentContext
//	Flatten() *CurrentContext
//}

type CurrentContext struct {
	Mm     map[string]Node[ast.Node]
	Parent *CurrentContext
}

//func (self *CurrentContext) currentContext() ICurrentContext {
//	return self
//}

func (self *CurrentContext) FindValue(value string) (Node[ast.Node], bool) {
	if v, ok := self.Mm[value]; ok {
		return v, true
	}
	if self.Parent == nil {
		return Node[ast.Node]{}, false
	}
	return self.Parent.FindValue(value)
}

func (self *CurrentContext) Flatten() *CurrentContext {
	result := func() *CurrentContext {
		if self.Parent != nil {
			return self.Parent.Flatten()
		}
		return &CurrentContext{map[string]Node[ast.Node]{}, nil}
	}()
	for k, v := range self.Mm {
		result.Mm[k] = v
	}
	return result
}
