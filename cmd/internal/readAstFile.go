package internal

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadAstFile(incomingState *ReadAstFileState, rawAst *RawAstRead, astFile *ast.File, relPath string, absPath string, file string, fs *token.FileSet) {
	excluded := map[string]bool{
		"time":                   true,
		"io":                     true,
		"os":                     true,
		"reflect":                true,
		"math":                   true,
		"strconv":                true,
		"github.com/google/uuid": true,
		"path/filepath":          true,
	}

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
						funcDecl := Node[*ast.FuncDecl]{key, node, imports, absPath, relPath, file, fs, true}
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
						rawAst.TypeSpecMap[key] = Node[*ast.TypeSpec]{key, node, imports, absPath, relPath, file, fs, true}
						return true
					}
					return true
				default:
					return true
				}
			case *ast.SwitchStmt:
				return true

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
		},
	)
}
