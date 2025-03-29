package internal

import "go/ast"

//type ICurrentContext interface {
//	currentContext() ICurrentContext
//	Flatten() *CurrentContext
//}

type CurrentContext struct {
	mm     map[string]Node[ast.Node]
	parent *CurrentContext
}

//func (self *CurrentContext) currentContext() ICurrentContext {
//	return self
//}

func (self *CurrentContext) FindValue(value string) (Node[ast.Node], bool) {
	if v, ok := self.mm[value]; ok {
		return v, true
	}
	if self.parent == nil {
		return Node[ast.Node]{}, false
	}
	return self.parent.FindValue(value)
}

func (self *CurrentContext) Flatten() *CurrentContext {
	result := func() *CurrentContext {
		if self.parent != nil {
			return self.parent.Flatten()
		}
		return &CurrentContext{map[string]Node[ast.Node]{}, nil}
	}()
	for k, v := range self.mm {
		result.mm[k] = v
	}
	return result
}
