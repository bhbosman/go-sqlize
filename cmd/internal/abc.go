package internal

import (
	"github.com/ugurcsen/gods-generic/queues"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"path/filepath"
	"strconv"
	"strings"
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

type ValueKey struct {
	Folder string
	Key    string
}

type FunctionMap map[ValueKey]Node[*ast.FuncDecl]
type SqlizeSelectMap map[ValueKey]Node[*ast.FuncDecl]
type StructMethodMap map[ValueKey]map[string]*ast.FuncDecl
type TypeSpecMap map[ValueKey]Node[*ast.TypeSpec]

type RawAstRead struct {
	FunctionMap     FunctionMap
	StructMethodMap StructMethodMap
	TypeSpecMap     TypeSpecMap
	InitFunctions   []Node[*ast.FuncDecl]
}

func ReadAstFile(incomingState *ReadAstFileState, rawAst *RawAstRead, astFile *ast.File, relPath string, absPath string, file string, fs *token.FileSet) {
	imports := ast.FileImports{}
	for _, spec := range astFile.Imports {
		pathValue, _ := strconv.Unquote(spec.Path.Value)
		ime := ast.ImportMapEntry{
			func(spec *ast.ImportSpec) string {
				if spec.Name == nil {
					pathValue, _ := strconv.Unquote(spec.Path.Value)
					pathSplits := strings.Split(pathValue, "/")
					return pathSplits[len(pathSplits)-1]
				}
				return spec.Name.Name
			}(spec),
			pathValue}
		imports[ime.Key] = ime
		if _, ok := excluded[ime.Path]; !ok {
			if _, ok := incomingState.PathsRead[ime.Path]; !ok {
				incomingState.PathsRead[ime.Path] = true

				pathToAdd := filepath.Join(incomingState.GoPath, "src", ime.Path)
				incomingState.PathsToBeReadQueue.Enqueue(pathToAdd)
			}
		}
	}
	astutil.Apply(
		astFile,
		func(cursor *astutil.Cursor) bool { return true },
		func(cursor *astutil.Cursor) bool {
			switch nodeType := cursor.Node().(type) {
			case ast.Expr:
				switch /*node := */ nodeType.(type) {
				default:
					return true
				}
			case ast.Decl:
				switch node := nodeType.(type) {
				case *ast.FuncDecl:
					if node.Recv == nil {
						key := ValueKey{relPath, node.Name.Name}
						funcDecl := Node[*ast.FuncDecl]{
							key,
							node,
							imports,
							absPath,
							relPath,
							file,
							fs,
							true,
						}
						switch {
						case node.Name != nil && node.Name.Name == "init":
							if incomingState.InputFolder == absPath {
								rawAst.InitFunctions = append(rawAst.InitFunctions, funcDecl)
							}
						default:
							rawAst.FunctionMap[key] = funcDecl
						}
					} else {
						structName := node.Recv.List[0].Type.(*ast.Ident).Name
						key := ValueKey{relPath, structName}
						if _, ok := rawAst.StructMethodMap[key]; !ok {
							rawAst.StructMethodMap[key] = map[string]*ast.FuncDecl{}
						}
						rawAst.StructMethodMap[key][node.Name.Name] = node
					}
					return true
				default:
					return true
				}
			case ast.Spec:
				switch node := nodeType.(type) {
				case *ast.TypeSpec:
					key := ValueKey{relPath, node.Name.Name}
					switch node.Type.(type) {
					case *ast.StructType:
						rawAst.TypeSpecMap[key] = Node[*ast.TypeSpec]{
							key,
							node,
							imports,
							absPath,
							relPath,
							file,
							fs,
							true,
						}
						return true
					}
					return true
				default:
					return true
				}
			case ast.Stmt:
				// check if part of a list
				if cursor.Index() < 0 {
					return true
				}

				// part of list
				switch cursor.Parent().(type) {
				case *ast.BlockStmt:
					if cursor.Index() == 0 {
						fciNode := CreateFciStmtNode(imports, cursor.Parent().Pos(), absPath, relPath, file)
						cursor.InsertBefore(fciNode)
					}
					return true
				default:
					return true
				}
			default:
				return true
			}
		})
}

type ReadAstFileState struct {
	PathsToBeReadQueue queues.Queue[string]
	PathsRead          map[string]bool
	GoPath             string
	InputFolder        string
}

var excluded = map[string]bool{
	"time":                   true,
	"os":                     true,
	"reflect":                true,
	"strconv":                true,
	"github.com/google/uuid": true,
	//"github.com/bhbosman/go-sqlize/lib": true,
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
