package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"strings"
)

func (compiler *Compiler) trailRecordToSql(item *TrailRecord) string {
	sb := &strings.Builder{}
	sources := compiler.findSources(item)
	sources = compiler.calculateSourceDependency(sources)
	_, _ = fmt.Fprintf(sb, "select\n")
	compiler.projectTrailRecord(sb, 1, item)
	if len(sources) > 0 {
		_, _ = fmt.Fprintf(sb, "from\n")
		compiler.projectSources(sb, 1, sources)
		_, _ = fmt.Fprintf(sb, "\n")
	}
	s := sb.String()
	return s
}

func (compiler *Compiler) projectTrailRecord(w io.Writer, tabCount int, item *TrailRecord) {
	for idx := 0; idx < item.Value.NumField(); idx++ {
		if node, ok := item.Value.Field(idx).Interface().(Node[ast.Node]); ok {
			compiler.internalProjectNode(w, tabCount, idx == item.Value.NumField()-1, 0, item.Value.Type().Field(idx).Name, node)
		}
	}
}

func (compiler *Compiler) nodeOperator(op token.Token) string {
	switch op {
	case token.ADD: // +
		return " + "
	case token.SUB: // -
		return " - "
	case token.MUL: // *
		return " * "
	case token.QUO: // /
		return " / "
	case token.LSS: // <
		return " < "
	case token.GTR: // >
		return " > "
	case token.LAND:
		return " AND "
	case token.LOR:
		return " OR "
	case token.GEQ:
		return " >= "
	case token.NEQ:
		return " <> "
	case token.EQL:
		return " = "
	default:
		panic("unhandled default case")
	}
}

func (compiler *Compiler) internalProjectRv(w io.Writer, tabCount int, last bool, stackCount int, name string, rv reflect.Value) {
	kind := rv.Kind()
	switch {
	case kind == reflect.Invalid:
		_, _ = io.WriteString(w, "nil")
	case kind == reflect.String:
		_, _ = io.WriteString(w, fmt.Sprintf("'%v'", rv.String()))
	case rv.CanInt():
		_, _ = io.WriteString(w, fmt.Sprintf("%v", rv.Int()))
	case rv.CanUint():
		_, _ = io.WriteString(w, fmt.Sprintf("%v", rv.Int()))
	case rv.CanFloat():
		_, _ = io.WriteString(w, fmt.Sprintf("%v", rv.Float()))
	case kind == reflect.Interface:
		break
	case kind == reflect.Map:
		break
	case kind == reflect.Pointer:
		break
	case kind == reflect.Bool:
		_, _ = io.WriteString(w, fmt.Sprintf("%v", rv.Bool()))
	case kind == reflect.Struct:
		if rv.CanInterface() {
			unk := rv.Interface()
			switch expr := unk.(type) {
			case SomeDataWithNode:
				_, _ = io.WriteString(w, fmt.Sprintf(" /*SomeDataWithNode(assigned:%v) */ ", expr.assigned))
				if expr.assigned {
					compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, expr.node)
				} else {
					_, _ = io.WriteString(w, fmt.Sprintf("nil"))
				}
			case SomeDataWithRv:
				_, _ = io.WriteString(w, fmt.Sprintf(" /*SomeDataWithRv(assigned:%v)*/ ", expr.assigned))
				if expr.assigned {
					compiler.internalProjectRv(w, tabCount, last, stackCount+1, name, expr.rv)
				} else {
					_, _ = io.WriteString(w, fmt.Sprintf("nil"))
				}
			case Node[ast.Node]:
				compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, expr)
			default:
				_, _ = io.WriteString(w, fmt.Sprintf("("))
				for idx := range rv.NumField() {
					fieldRv := rv.Field(idx)
					compiler.internalProjectRv(w, tabCount, last, stackCount+1, name, fieldRv)
					if idx != rv.NumField()-1 {
						_, _ = io.WriteString(w, fmt.Sprintf(","))
					}
				}
				_, _ = io.WriteString(w, fmt.Sprintf(")"))
			}
		} else {
			_, _ = io.WriteString(w, fmt.Sprintf("rv.CanInterface() == false"))
		}
	default:
		panic("unhandled default case")
	}
}

func (compiler *Compiler) internalProjectNode(w io.Writer, tabCount int, last bool, stackCount int, name string, node Node[ast.Node]) {
	if !node.Valid {
		return
	}
	if stackCount == 0 {
		_, _ = io.WriteString(w, strings.Repeat("\t", tabCount))
	}
	switch nodeItem := node.Node.(type) {
	case *CheckForNotNullExpression:
		_, _ = io.WriteString(w, "(")
		compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, nodeItem.node)
		_, _ = io.WriteString(w, " is not null)")
	case *EntityField:
		_, _ = io.WriteString(w, fmt.Sprintf("[%v].[%v]", nodeItem.alias, nodeItem.field))
	case *coercion:
		_, _ = io.WriteString(w, "CAST(")
		param := ChangeParamNode[ast.Node, ast.Node](node, nodeItem.Node.Node)
		compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, param)
		_, _ = io.WriteString(w, " as ")
		switch nodeItem.to {
		case "float64":
			_, _ = io.WriteString(w, "float")
		case "int":
			_, _ = io.WriteString(w, "int")
		case "string":
			_, _ = io.WriteString(w, "varchar")
		default:
			panic(node)
		}
		_, _ = io.WriteString(w, ")")
	case *BinaryExpr:
		_, _ = io.WriteString(w, "(")
		compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, nodeItem.left)
		_, _ = io.WriteString(w, compiler.nodeOperator(nodeItem.Op))
		compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, nodeItem.right)
		_, _ = io.WriteString(w, ")")
	case *MultiBinaryExpr:
		_, _ = io.WriteString(w, "(")
		for idx, expr := range nodeItem.expressions {
			if idx != 0 {
				switch nodeItem.Op {
				case token.LAND:
					_, _ = io.WriteString(w, " AND ")
				case token.LOR:
					_, _ = io.WriteString(w, " OR ")
				default:
					panic("unhandled default case")
				}
			}
			compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, expr)
		}
		_, _ = io.WriteString(w, ")")

	case *ReflectValueExpression:
		kind := nodeItem.Rv.Kind()
		switch kind {
		//case reflect.Invalid:
		//	_, _ = io.WriteString(w, fmt.Sprintf("nil"))
		default:
			compiler.internalProjectRv(w, tabCount, last, stackCount+1, name, nodeItem.Rv)
		}

	case *SupportedFunction:
		_, _ = io.WriteString(w, fmt.Sprintf("%v(", nodeItem.functionName))
		for idx, param := range nodeItem.params {
			compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, param)
			if idx != len(nodeItem.params)-1 {
				_, _ = io.WriteString(w, ", ")
			}
		}
		_, _ = io.WriteString(w, fmt.Sprintf(")"))

	case *IfThenElseSingleValueCondition:
		_, _ = io.WriteString(w, "case\n")

		for _, expr := range nodeItem.conditionalStatement {
			tabCount++
			_, isLiteral := isLiterateValue(expr.condition)
			if !isLiteral {
				_, _ = io.WriteString(w, fmt.Sprintf("%vwhen ", strings.Repeat("\t", tabCount)))
				compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, expr.condition)
				_, _ = io.WriteString(w, " then\n")
			} else {
				_, _ = io.WriteString(w, fmt.Sprintf("%velse\n", strings.Repeat("\t", tabCount)))
			}
			tabCount++
			_, _ = io.WriteString(w, fmt.Sprintf("%v", strings.Repeat("\t", tabCount)))
			compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, expr.value)
			tabCount--
			_, _ = io.WriteString(w, "\n")
			tabCount--
		}
		_, _ = io.WriteString(w, fmt.Sprintf("%vend", strings.Repeat("\t", tabCount)))

	case *LhsToMultipleRhsOperator:
		_, _ = io.WriteString(w, "(")
		for idx, rhs := range nodeItem.Rhs {
			_, _ = io.WriteString(w, "(")
			compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, nodeItem.Lhs)
			_, _ = io.WriteString(w, compiler.nodeOperator(nodeItem.LhsToRhsOp))
			compiler.internalProjectNode(w, tabCount, last, stackCount+1, name, rhs)
			_, _ = io.WriteString(w, ")")
			if idx != len(nodeItem.Rhs)-1 {
				_, _ = io.WriteString(w, compiler.nodeOperator(nodeItem.betweenTerminalsOp))
			}
		}
		_, _ = io.WriteString(w, ")")
	default:
		panic("implement me")
	}
	if stackCount == 0 {
		_, _ = io.WriteString(w, fmt.Sprintf(" as %v", name))
		if !last {
			_, _ = io.WriteString(w, ",")
		}
		_, _ = io.WriteString(w, "\n")
	}
}
