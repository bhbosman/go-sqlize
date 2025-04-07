package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"strconv"
)

func (compiler *Compiler) addLibFunctions() {
	compiler.GlobalFunctions[ValueKey{libFolder, "Query"}] = compiler.libQueryImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "Map"}] = compiler.libMapImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "GenerateSql"}] = compiler.libGenerateSqlImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "GenerateSqlTest"}] = compiler.libGenerateSqlTestImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "Atoi"}] = compiler.libAtoiImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "SetSomeValue"}] = compiler.libSetSomeValueImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "SetSomeNone"}] = compiler.libSetSomeNoneImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "IsSomeAssigned"}] = compiler.libIsSomeAssignedImplementation
	compiler.GlobalFunctions[ValueKey{libFolder, "SomeData"}] = compiler.libSomeDataImplementation
}

func (compiler *Compiler) libQueryImplementation(_ State, typeParams []Node[ast.Expr], _ []Node[ast.Node]) ExecuteStatement {
	if len(typeParams) != 1 {
		panic(fmt.Errorf("Lib.Query implementation requires 1 type argument, got %d", len(typeParams)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		var arr []reflect.Type
		for _, expr := range typeParams {
			rt := compiler.findType(state, expr)
			arr = append(arr, rt)
		}
		rt := arr[0]
		alias := compiler.AddEntitySource(rt)
		qt := &TrailSource{state.currentNode.Node.Pos(), alias}
		returnValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, qt)
		return []Node[ast.Node]{returnValue}, artValue
	}
}

func (compiler *Compiler) libMapImplementation(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 2 {
		panic(fmt.Errorf("Lib.Map implementation requires 2 arguments, got %d", len(arguments)))
	}
	if _, ok := arguments[1].Node.(*ast.FuncLit); !ok {
		panic("map implementation requires function literal")
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		funcLit, _ := arguments[1].Node.(*ast.FuncLit)
		var names []*ast.Ident
		if funcLit.Type.Params != nil {
			for _, field := range funcLit.Type.Params.List {
				names = append(names, field.Names...)
			}
		}
		m := map[string]Node[ast.Node]{}
		m[names[0].Name] = arguments[0]
		newContext := &CurrentContext{m, GetCompilerState[*CurrentContext](state)}
		state = SetCompilerState(newContext, state)
		param := ChangeParamNode[ast.Node, *ast.BlockStmt](state.currentNode, funcLit.Body)
		values, _ := compiler.executeBlockStmt(state, param)
		state = SetCompilerState(newContext.Parent, state)
		return values, artValue
	}
}

func (compiler *Compiler) libGenerateSqlImplementation(state State, i []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("Lib.GenerateSql implementation requires 1 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		switch item := arguments[0].Node.(type) {
		case *TrailRecord:
			s := compiler.trailRecordToSql(item)
			basicLit := &ast.BasicLit{state.currentNode.Node.Pos(), token.STRING, strconv.Quote(s)}
			nodeValue := ChangeParamNode[ast.Node, ast.Node](state.currentNode, basicLit)
			return []Node[ast.Node]{nodeValue}, artValue
		default:
			panic("implementation required")
		}
	}
}

func (compiler *Compiler) libGenerateSqlTestImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("Lib.GenerateSqlTest implementation requires 1 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		args := arguments[0:1]
		ans, _ := compiler.libGenerateSqlImplementation(state, params, args)(state)
		if rv, isLiterate := isLiterateValue(ans[0]); isLiterate && rv.Kind() == reflect.String {
			currentContext := GetCompilerState[*CurrentContext](state)
			if value, b := currentContext.FindValue("__stdOut__"); b {
				switch nodeValue := value.Node.(type) {
				case *ReflectValueExpression:
					if wr, ok := nodeValue.Rv.Interface().(io.Writer); ok {
						_, _ = wr.Write([]byte(rv.String()))
					}
					return nil, artNone
				default:
					panic(fmt.Sprintf("libGenerateSqlTestImplementation *ReflectValueExpression not found"))
				}
			}
			panic(fmt.Sprintf("libGenerateSqlTestImplementation __stdOut__ not found"))
		}
		panic(fmt.Sprintf("libGenerateSqlTestImplementation value from GenerateSql not literal"))

	}
}

func (compiler *Compiler) libAtoiImplementation(_ State, _ []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("Lib.Atoi implementation requires 1 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		nodes, resultType := compiler.strconvAtoiImplementation(state, nil, arguments)(state)
		return nodes[:1], resultType
	}
}

func (compiler *Compiler) libSetSomeValueImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	if len(arguments) != 1 {
		panic(fmt.Errorf("SetSomeValue implementation requires 1 arguments, got %d", len(arguments)))
	}
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		return arguments[0:1], artValue
	}
}

func (compiler *Compiler) libSetSomeNoneImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, &NilValueExpression{})}, artValue
	}
}

func (compiler *Compiler) libIsSomeAssignedImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		binExpr := &BinaryExpr{
			OpPos: 0,
			Op:    token.NEQ,
			left:  arguments[0],
			right: ChangeParamNode[ast.Node, ast.Node](state.currentNode, &NilValueExpression{}),
		}
		return []Node[ast.Node]{ChangeParamNode[ast.Node, ast.Node](state.currentNode, binExpr)}, artValue
	}
}

func (compiler *Compiler) libSomeDataImplementation(state State, params []Node[ast.Expr], arguments []Node[ast.Node]) ExecuteStatement {
	return func(state State) ([]Node[ast.Node], CallArrayResultType) {
		return arguments[0:1], artValue
	}
}
