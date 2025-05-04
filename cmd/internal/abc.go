package internal

import (
	"encoding/json"
	"fmt"
	"github.com/ugurcsen/gods-generic/queues"
	"go/ast"
	"go/token"
	"strconv"
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
) *ast.AssignStmt {
	fci := &FolderContextInformation{nodePos, fileName, absPath, relPath, importMap}
	marshal, err := json.Marshal(fci)
	if err != nil {
		return nil
	}
	m := []struct {
		key   string
		value string
	}{
		{"json", string(marshal)},
	}
	var lhs []ast.Expr
	var rhs []ast.Expr
	for i := 0; i < len(m); i++ {
		lhs = append(lhs, &ast.Ident{nodePos, m[i].key, nil})
		quoteValue := strconv.Quote(m[i].value)
		rhs = append(rhs, &ast.BasicLit{nodePos, token.STRING, quoteValue})
	}
	return &ast.AssignStmt{lhs, nodePos, 0xFFFF, rhs}
}
