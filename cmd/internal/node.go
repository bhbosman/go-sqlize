package internal

import (
	"go/ast"
	"go/token"
)

type Node[TType ast.Node] struct {
	Key       ValueKey
	Node      TType
	ImportMap ast.FileImports
	AbsPath   string
	RelPath   string
	FileName  string
	Fs        *token.FileSet
	Valid     bool
}

type ICheckNode interface {
	IsValidNode() bool
}

func isValidNodes(nodes []Node[ast.Node]) bool {
	for _, node := range nodes {
		b := isValidNode(node)
		if !b {
			return false
		}
	}
	return true
}
func isValidNode(node Node[ast.Node]) bool {
	switch value := node.Node.(type) {
	case ICheckNode:
		return value.IsValidNode()
	}
	return true
}

func ChangeParamNode[TExisting ast.Node, TTarget ast.Node](
	old Node[TExisting],
	withNode TTarget,
) Node[TTarget] {
	return Node[TTarget]{
		old.Key,
		withNode,
		old.ImportMap,
		old.AbsPath,
		old.RelPath,
		old.FileName,
		old.Fs,
		old.Valid,
	}
}
