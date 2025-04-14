package internal

import (
	"github.com/ugurcsen/gods-generic/queues"
	"go/ast"
	"go/token"
)

type ValueKey struct {
	Folder string
	Key    string
}

type FunctionMap map[ValueKey]Node[*ast.FuncDecl]

type StructMethodMap map[ValueKey]map[string]*ast.FuncDecl
type TypeSpecMap map[ValueKey]Node[*ast.TypeSpec]

type RawAstRead struct {
	FunctionMap     FunctionMap
	StructMethodMap StructMethodMap
	TypeSpecMap     TypeSpecMap
	InitFunctions   []Node[*ast.FuncDecl]
}

type ReadAstFileState struct {
	PathsToBeReadQueue queues.Queue[string]
	PathsRead          map[string]bool
	GoPath             string
	InputFolder        string
}

func CreateFciStmtNode(importMap ast.FileImports,
	nodePos token.Pos,
	absPath string,
	relPath string,
	fileName string,
) *ast.FolderContextInformation {
	fci := &ast.FolderContextInformation{nodePos, fileName, absPath, relPath, importMap}
	return fci
}
