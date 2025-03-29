package cmd

import (
	"fmt"
	"github.com/bhbosman/go-sqlize/cmd/internal"
	"go/ast"

	"github.com/spf13/cobra"
	"github.com/ugurcsen/gods-generic/queues/linkedlistqueue"

	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use:   "sqlize",
	Short: "sqlize",
	Long:  "sqlize",
	Run: func(cmd *cobra.Command, args []string) {

		for _, __inputFolder := range args {
			inputFolder := __inputFolder
			if !filepath.IsAbs(inputFolder) {
				wd, _ := os.Getwd()
				inputFolder = filepath.Join(wd, inputFolder)
			}
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
					println(0)

				}
			}
			compiler := &internal.Compiler{}
			compiler.Init(
				rawAst.FunctionMap,
				rawAst.StructMethodMap,
				rawAst.TypeSpecMap,
				rawAst.InitFunctions,
			)
			compiler.Compile()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("inputFolder", "i", "./work", "input folder")
}
