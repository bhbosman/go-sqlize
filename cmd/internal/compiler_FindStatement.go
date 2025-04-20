package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

func syntaxErrorf(node Node[ast.Node], format string, a ...any) error {
	p := node.Fs.Position(node.Node.Pos())
	position := fmt.Sprintf("%v:%v:%v", p.Filename, p.Line, p.Column)
	ss := fmt.Sprintf(format, a...)
	return fmt.Errorf("syntax error at (%v): %v", position, ss)
}

func notFound(typeName, methodName string) error {
	return fmt.Errorf("handler not found for %v in %v", typeName, methodName)
}

func createError(methodName, message string) error {
	return fmt.Errorf("%v in %v", message, methodName)
}

func (compiler *Compiler) findStatement(state State, node Node[ast.Stmt]) (ExecuteStatement, Node[ast.Node]) {
	switch item := node.Node.(type) {
	case *FolderContextInformation:
		return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
			value := ChangeParamNode[ast.Stmt, ast.Node](node, item)
			return []Node[ast.Node]{value}, artFCI
		}, ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)
	case *ast.IfStmt:
		value := ChangeParamNode(node, item)
		return compiler.createIfStmtExecution(value), ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)
	case *ast.SwitchStmt:
		value := ChangeParamNode(node, item)
		return compiler.createSwitchStmtExecution(value), ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)

	case *ast.CaseClause:
		value := ChangeParamNode(node, item)
		return compiler.createCaseClauseExecution(value), ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)

	case *ast.AssignStmt:
		value := ChangeParamNode(node, item)
		return compiler.createAssignStatementExecution(value), ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)
	case *ast.ExprStmt:
		param := ChangeParamNode[ast.Stmt, ast.Node](node, item.X)
		tempState := state.setCurrentNode(ChangeParamNode[ast.Stmt, ast.Node](node, item.X))
		return compiler.findRhsExpression(tempState, param), ChangeParamNode[ast.Stmt, ast.Node](node, item.X)
	case *ast.ReturnStmt:
		value := ChangeParamNode(node, item)
		return compiler.createReturnStmtExecution(value), ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)
	case *ast.BlockStmt:
		value := ChangeParamNode(node, item)
		return compiler.createBlockStmtExecution(value), ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)
	case *ast.DeclStmt:
		value := ChangeParamNode(node, item)
		return compiler.createDeclStmtExecution(value), ChangeParamNode[ast.Stmt, ast.Node](node, node.Node)
	default:
		panic(notFound(reflect.TypeOf(item).String(), "findStatement"))
	}
}

func (compiler *Compiler) handleSpec(state State, node Node[ast.Spec]) {
	switch expr := node.Node.(type) {
	case *ast.TypeSpec:
		param := ChangeParamNode(node, expr)
		currentContext := GetCompilerState[*CurrentContext](state)
		onCreateType := compiler.readTypeSpec(param)
		currentContext.addLocalTypes(expr.Name.Name, onCreateType)
	default:
		panic(notFound(reflect.TypeOf(expr).String(), "spec"))
	}
}

func (compiler *Compiler) createDeclStmtExecution(node Node[*ast.DeclStmt]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		switch nodeItem := node.Node.Decl.(type) {
		case *ast.GenDecl:
			for _, spec := range nodeItem.Specs {
				param := ChangeParamNode(node, spec)
				compiler.handleSpec(state, param)
			}
		}
		return nil, artNone
	}
}

func (compiler *Compiler) createBlockStmtExecution(node Node[*ast.BlockStmt]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		return compiler.executeBlockStmt(state, node, typeParams, unprocessedArgs)
	}
}

func isLiterateValue(node Node[ast.Node]) (reflect.Value, bool) {
	switch item := node.Node.(type) {
	case *TrailRecord:
		rv := item.Value
		for idx := range rv.NumField() {
			switch rvIdxField := rv.Field(idx).Interface().(type) {
			case Node[ast.Node]:
				if _, b := isLiterateValue(rvIdxField); !b {
					return reflect.Value{}, false
				}
			}
		}
		return item.Value, true
	case *IfThenElseSingleValueCondition:
		return reflect.Value{}, false
	case *CheckForNotNullExpression:
		return isLiterateValue(item.node)
	case *coercion:
		rv, isLiterate := isLiterateValue(item.Node)
		return rv, isLiterate
	case *EntityField:
		return reflect.Value{}, false
	case *ReflectValueExpression:
		return item.Rv, true
	case *ast.BasicLit:
		switch item.Kind {
		case token.INT:
			intValue, _ := strconv.ParseInt(item.Value, 10, 64)
			return reflect.ValueOf(intValue), true
		case token.FLOAT:
			floatValue, _ := strconv.ParseFloat(item.Value, 64)
			return reflect.ValueOf(floatValue), true
		case token.IMAG:
			panic("ssfds")
		case token.CHAR:
			panic("ssfds")
		case token.STRING:
			stringValue, _ := strconv.Unquote(item.Value)
			return reflect.ValueOf(stringValue), true
		default:
			panic(notFound(item.Kind.String(), "isLiterateValue"))
		}
	case *builtInNil:
		return reflect.Value{}, true
		// TODO: *BinaryExpr: should always be false, this needs to be fixed where *BinaryExpr: is created
	case *BinaryExpr:
		return reflect.Value{}, false
	case *MultiBinaryExpr:
		return reflect.Value{}, false
	case *LhsToMultipleRhsOperator:
		return reflect.Value{}, false
	default:
		panic(notFound(reflect.TypeOf(item).String(), "isLiterateValue"))
	}
}

func (compiler *Compiler) createReturnStmtExecution(node Node[*ast.ReturnStmt]) ExecuteStatement {
	return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
		var result []Node[ast.Node]
		for _, expr := range node.Node.Results {
			param := ChangeParamNode[*ast.ReturnStmt, ast.Node](node, expr)
			tempState := state.setCurrentNode(ChangeParamNode[*ast.ReturnStmt, ast.Node](node, expr))
			fn := compiler.findRhsExpression(tempState, param)

			param01 := ChangeParamNode[*ast.ReturnStmt, ast.Node](node, expr)
			state = state.setCurrentNode(param01)
			v, _ := compiler.executeAndExpandStatement(state, typeParams, unprocessedArgs, fn)
			result = append(result, v...)
		}
		return result, artReturn
	}
}

func (compiler *Compiler) createAssignStatementExecution(node Node[*ast.AssignStmt]) ExecuteStatement {
	switch node.Node.Tok {
	case token.DEFINE, token.ASSIGN:
		return func(state State, typeParams map[string]ITypeMapper, unprocessedArgs []Node[ast.Node]) ([]Node[ast.Node], CallArrayResultType) {
			var rhsArray []Node[ast.Node]

			for _, rhsExpression := range node.Node.Rhs {
				param := ChangeParamNode[*ast.AssignStmt, ast.Node](node, rhsExpression)
				tempState := state.setCurrentNode(ChangeParamNode[*ast.AssignStmt, ast.Node](node, rhsExpression))
				fn := compiler.findRhsExpression(tempState, param)
				tempState = tempState.setCurrentNode(ChangeParamNode[*ast.AssignStmt, ast.Node](node, rhsExpression))
				arr, _ := compiler.executeAndExpandStatement(tempState, typeParams, unprocessedArgs, fn)
				rhsArray = append(rhsArray, arr...)
			}

			for idx, lhsExpression := range node.Node.Lhs {
				param := ChangeParamNode(node, lhsExpression)
				state = state.setCurrentNode(ChangeParamNode[*ast.AssignStmt, ast.Node](node, lhsExpression))
				assignStatement := compiler.findLhsExpression(state, param, node.Node.Tok)
				assignStatement(state, rhsArray[idx])
			}
			return nil, artNone
		}
	default:
		panic("dddd")
	}
}
