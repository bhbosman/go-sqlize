package cmd

import (
	"fmt"
	"github.com/bhbosman/go-sqlize/cmd/internal"
	"github.com/spf13/cobra"
	"github.com/ugurcsen/gods-generic/queues/linkedlistqueue"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
)

var RootCmd = &cobra.Command{
	Use:   "sqlize",
	Short: "sqlize",
	Long:  "sqlize",
	Run: func(cmd *cobra.Command, args []string) {
		for _, argItem := range args {
			absPath := func(argItem string) string {
				if !filepath.IsAbs(argItem) {
					wd, _ := os.Getwd()
					return filepath.Join(wd, argItem)
				}
				return argItem
			}(argItem)

			inputFolder, inputFile := func(absPath string) (string, string) {
				fileInfo, err := os.Stat(absPath)
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err.Error())
					os.Exit(1)
				}
				if !fileInfo.IsDir() {
					return filepath.Split(absPath)
				}
				return absPath, ""
			}(absPath)

			fs := token.NewFileSet()
			rawAst := &internal.RawAstRead{internal.FunctionMap{}, internal.StructMethodMap{}, internal.TypeSpecMap{}, []internal.Node[*ast.FuncDecl]{}}
			state := &internal.ReadAstFileState{
				linkedlistqueue.New[string](),
				map[string]bool{},
				build.Default.GOPATH,
				inputFolder,
			}

			state.PathsToBeReadQueue.Enqueue(inputFolder)
			for !state.PathsToBeReadQueue.Empty() {
				if absPath, ok := state.PathsToBeReadQueue.Dequeue(); ok {
					___relPath, err := filepath.Rel(filepath.Join(build.Default.GOPATH, "src"), absPath)
					if err != nil {
						_, _ = fmt.Fprintln(os.Stderr, err.Error())
						os.Exit(1)
					}
					pkg, err := build.Default.Import(___relPath, ".", 0)
					if err != nil {
						_, _ = fmt.Fprintln(os.Stderr, err.Error())
						os.Exit(1)
					}
					for _, file := range pkg.GoFiles {
						astFile, err := parser.ParseFile(fs, filepath.Join(pkg.Dir, file), nil, parser.SkipObjectResolution)
						if err != nil {
							_, _ = fmt.Fprintln(os.Stderr, err.Error())
							os.Exit(1)
						}
						internal.ReadAstFile(state, rawAst, astFile, pkg.ImportPath, absPath, file, fs)
					}
				}
			}
			compiler := &internal.Compiler{}
			compiler.Init(
				rawAst.FunctionMap,
				rawAst.StructMethodMap,
				rawAst.TypeSpecMap,
				rawAst.InitFunctions,
			)
			var ss []string
			if inputFile != "" {
				ss = append(ss, inputFile)
			}

			currentContext := &internal.CurrentContext{
				map[string]internal.Node[ast.Node]{
					"__stdOut__": {Node: &internal.ReflectValueExpression{reflect.ValueOf(cmd.OutOrStdout())}},
				},
				nil,
			}

			compiler.Compile(currentContext, ss...)
		}
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().StringP("inputFolder", "i", "./work", "input folder")
}
