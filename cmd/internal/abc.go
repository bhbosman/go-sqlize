package internal

import (
	"fmt"
	"github.com/ugurcsen/gods-generic/queues"
	"go/ast"
	"go/token"
)

type ValueKey struct {
	Folder string
	Key    string
}

func (v *ValueKey) String() string {
	if v.Folder == "" {
		return fmt.Sprintf("%v", v.Key)
	}
	return fmt.Sprintf("%v.%v", v.Folder, v.Key)
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

func CreateFciStmtNode(importMap FileImports,
	nodePos token.Pos,
	absPath string,
	relPath string,
	fileName string,
) *FolderContextInformation {
	fci := &FolderContextInformation{nodePos, fileName, absPath, relPath, importMap}
	return fci
}
