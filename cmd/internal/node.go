package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
)

type Node[TType ast.Node] struct {
	Key       ValueKey
	Node      TType
	ImportMap FileImports
	AbsPath   string
	RelPath   string
	FileName  string
	Fs        *token.FileSet
	Valid     bool
}

type SortNodes struct {
	nodes  []Node[ast.Node]
	lessFn func(i, j int) bool
}

func (sortNode *SortNodes) Len() int {
	return len(sortNode.nodes)
}

func (sortNode *SortNodes) Less(i, j int) bool {
	return sortNode.lessFn(i, j)
}

func (sortNode *SortNodes) Swap(i, j int) {
	sortNode.nodes[i], sortNode.nodes[j] = sortNode.nodes[j], sortNode.nodes[i]
}

func NodeString(node Node[ast.Node]) (string, bool) {
	if unk, ok := node.Node.(fmt.Stringer); ok {
		return unk.String(), true
	}
	return "", false
}

func NodesString(nodes []Node[ast.Node]) string {
	var s []string
	for _, node := range nodes {
		if nodeString, b := NodeString(node); b {
			s = append(s, nodeString)
		} else {
			rt := reflect.TypeOf(node.Node)
			s = append(s, rt.String())
		}
	}
	return strings.Join(s, ",")
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

func ChangeParamNodeWithValueKey[TExisting ast.Node, TTarget ast.Node](
	old Node[TExisting],
	withNode TTarget,
	vk ValueKey,
) Node[TTarget] {
	return Node[TTarget]{
		vk,
		withNode,
		old.ImportMap,
		old.AbsPath,
		old.RelPath,
		old.FileName,
		old.Fs,
		old.Valid,
	}
}
