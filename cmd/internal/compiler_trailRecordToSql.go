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
			compiler.internalProjectTrailRecord(w, tabCount, idx == item.Value.NumField()-1, 0, item.Value.Type().Field(idx).Name, node)
		}
	}
}

func (compiler *Compiler) internalProjectTrailRecord(w io.Writer, count int, last bool, stackCount int, name string, node Node[ast.Node]) {
	if !node.Valid {
		return
	}
	if stackCount == 0 {
		_, _ = io.WriteString(w, strings.Repeat("\t", count))
	}
	switch nodeItem := node.Node.(type) {
	case *EntityField:
		_, _ = io.WriteString(w, fmt.Sprintf("%v.%v", nodeItem.alias, nodeItem.field))
	case *coercion:
		_, _ = io.WriteString(w, "CAST(")
		param := ChangeParamNode[ast.Node, ast.Node](node, nodeItem.Node.Node)
		compiler.internalProjectTrailRecord(w, count, last, stackCount+1, name, param)
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
		compiler.internalProjectTrailRecord(w, count, last, stackCount+1, name, nodeItem.left)
		switch nodeItem.Op {
		case token.ADD: // +
			_, _ = io.WriteString(w, " + ")
		case token.SUB: // -
			_, _ = io.WriteString(w, " - ")
		case token.MUL: // *
			_, _ = io.WriteString(w, " * ")
		case token.QUO: // /
			_, _ = io.WriteString(w, " / ")
		case token.LSS: // <
			_, _ = io.WriteString(w, " < ")
		case token.GTR: // >
			_, _ = io.WriteString(w, " > ")
		case token.LAND:
			_, _ = io.WriteString(w, " && ")
		case token.GEQ:
			_, _ = io.WriteString(w, " >= ")
		case token.NEQ:
			_, _ = io.WriteString(w, " != ")
		default:
			panic("unhandled default case")
		}
		compiler.internalProjectTrailRecord(w, count, last, stackCount+1, name, nodeItem.right)
		_, _ = io.WriteString(w, ")")
	case *ReflectValueExpression:
		kind := nodeItem.rv.Kind()
		switch {
		case kind == reflect.String:
			_, _ = io.WriteString(w, fmt.Sprintf("'%v'", nodeItem.rv.String()))
		case nodeItem.rv.CanInt():
			_, _ = io.WriteString(w, fmt.Sprintf("%v", nodeItem.rv.Int()))
		case nodeItem.rv.CanUint():
			_, _ = io.WriteString(w, fmt.Sprintf("%v", nodeItem.rv.Int()))
		case kind == reflect.Float64:
			_, _ = io.WriteString(w, fmt.Sprintf("%v", nodeItem.rv.Float()))
		case kind == reflect.Bool:
			_, _ = io.WriteString(w, fmt.Sprintf("%v", nodeItem.rv.Bool()))
		default:
			panic("unhandled default case")
		}
	case *SupportedFunction:
		_, _ = io.WriteString(w, fmt.Sprintf("%v(", nodeItem.functionName))
		for idx, param := range nodeItem.params {
			compiler.internalProjectTrailRecord(w, count, last, stackCount+1, name, param)
			if idx != len(nodeItem.params)-1 {
				_, _ = io.WriteString(w, ", ")
			}
		}
		_, _ = io.WriteString(w, fmt.Sprintf(")"))
	case *PartialExpression:
		_, _ = io.WriteString(w, "(\n")
		_, _ = io.WriteString(w, fmt.Sprintf("%v%v(%v)\n", strings.Repeat("\t", count+1), "PartialExpressions", len(nodeItem.conditionalStatement)))
		for _, expr := range nodeItem.conditionalStatement {
			_, _ = io.WriteString(w, fmt.Sprintf("%v", strings.Repeat("\t", count+1)))
			compiler.internalProjectTrailRecord(w, count+1, last, stackCount+1, name, expr.condition)
			_, _ = io.WriteString(w, "\n")
			_, _ = io.WriteString(w, fmt.Sprintf("%v", strings.Repeat("\t", count+1+1)))
			compiler.internalProjectTrailRecord(w, count+1, last, stackCount+1, name, expr.value)
			_, _ = io.WriteString(w, "\n")
		}
		_, _ = io.WriteString(w, fmt.Sprintf("%v)", strings.Repeat("\t", count)))
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
