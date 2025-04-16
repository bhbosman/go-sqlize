package internal

import (
	"go/ast"
	"go/token"
)

func (compiler *Compiler) findLhsExpression(state State, node Node[ast.Expr], tok token.Token) AssignStatement {
	return compiler.internalFindLhsExpression(state, node, tok).(AssignStatement)
}

func (compiler *Compiler) internalFindLhsExpression(state State, node Node[ast.Expr], tok token.Token) interface{} {
	switch item := node.Node.(type) {
	case *ast.Ident:
		if item.Name == "_" {
			fn := func() AssignStatement {
				return func(state State, value Node[ast.Node]) {
					// do nothing
				}
			}
			return fn()
		}
		fn := func(currentContext *CurrentContext, key string) AssignStatement {
			return func(state State, value Node[ast.Node]) {
				currentContext.Mm[key] = ValueInformation{value}
			}
		}
		switch tok {
		case token.DEFINE:
			currentContext := GetCompilerState[*CurrentContext](state)
			return fn(currentContext, item.Name)
		case token.ASSIGN:
			//err := syntaxErrorf(state.currentNode, "No assignements allowed")
			//panic(err)
			currentContext := GetCompilerState[*CurrentContext](state)
			currentContext.FindValueByString(item.Name)
			for currentContext != nil {
				if _, ok := currentContext.Mm[item.Name]; ok {
					break
				}
				currentContext = currentContext.Parent
			}
			if currentContext != nil {
				return fn(currentContext, item.Name)
			}
			panic("this should never happen")
		default:
			panic("unhandled default case")
		}
	default:
		panic(item)
	}
}
