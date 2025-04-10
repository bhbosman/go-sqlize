package internal

import (
	"go/ast"
)

func (compiler *Compiler) executeBlockStmt(state State, node Node[*ast.BlockStmt]) ([]Node[ast.Node], CallArrayResultType) {
	newContext := &CurrentContext{map[string]Node[ast.Node]{}, GetCompilerState[*CurrentContext](state)}
	state = SetCompilerState(newContext, state)
	for _, item := range node.Node.List {
		param := ChangeParamNode(node, item)
		tempState := state.setCurrentNode(ChangeParamNode[*ast.BlockStmt, ast.Node](node, item))
		statementFn, currentNode := compiler.findStatement(tempState, param)
		tempState = state.setCurrentNode(currentNode)
		arr, rt := compiler.executeAndExpandStatement(tempState, statementFn)
		switch rt {
		case artFCI:
			switch instance := arr[0].Node.(type) {
			case *ast.FolderContextInformation:
				node = Node[*ast.BlockStmt]{node.Key, node.Node, instance.Imports, instance.AbsPath, instance.RelPath, instance.FileName, node.Fs, node.Valid}
			}
		case artValue:
		case artReturn:
			if len(arr) == 0 {
				panic("big error")
			}
			return arr, artReturn
		default:
			continue
		}
	}
	state = SetCompilerState(newContext.Parent, state)
	return nil, artNone
}
