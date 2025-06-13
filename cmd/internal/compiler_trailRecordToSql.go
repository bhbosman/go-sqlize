package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"strings"
)

func (compiler *Compiler) trailRecordToSelectStatement(state State, node Node[*TrailRecord]) string {
	sb := &strings.Builder{}
	sources := compiler.findSourcesFromNode(ChangeParamNode[*TrailRecord, ast.Node](node, node.Node))
	sources = compiler.findAdditionalSourcesFromJoins(sources)
	sources = compiler.findAdditionalSourcesFromAssociations(sources)
	orderedSources := compiler.calculateSourcesOrder(sources)

	_, _ = fmt.Fprintf(sb, "select\n")
	compiler.projectTrailNode(sb, 1, node)
	compiler.projectSources(state, sb, 1, orderedSources)
	s := sb.String()
	return s

}

func (compiler *Compiler) projectTrailNode(w io.Writer, tabCount int, item Node[*TrailRecord]) {
	for idx := 0; idx < item.Node.Value.NumField(); idx++ {
		if node, ok := item.Node.Value.Field(idx).Interface().(Node[ast.Node]); ok && node.Valid {
			last := idx == item.Node.Value.NumField()-1
			_, _ = io.WriteString(w, strings.Repeat("\t", tabCount))
			compiler.internalProjectNode(w, tabCount, node)
			_, _ = io.WriteString(w, fmt.Sprintf(" as %v", item.Node.Value.Type().Field(idx).Name))
			if !last {
				_, _ = io.WriteString(w, ",")
			}
			_, _ = io.WriteString(w, "\n")
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

func (compiler *Compiler) internalProjectRv(w io.Writer, tabCount int, node Node[*ReflectValueExpression]) {
	kind := node.Node.Rv.Kind()
	switch {
	case kind == reflect.Invalid:
		_, _ = io.WriteString(w, "nil")
	case kind == reflect.String:
		_, _ = io.WriteString(w, fmt.Sprintf("'%v'", node.Node.Rv.String()))
	case node.Node.Rv.CanInt():
		_, _ = io.WriteString(w, fmt.Sprintf("%v", node.Node.Rv.Int()))
	case node.Node.Rv.CanUint():
		_, _ = io.WriteString(w, fmt.Sprintf("%v", node.Node.Rv.Int()))
	case node.Node.Rv.CanFloat():
		_, _ = io.WriteString(w, fmt.Sprintf("%v", node.Node.Rv.Float()))
	case kind == reflect.Interface:
		break
	case kind == reflect.Map:
		break
	case kind == reflect.Pointer:
		p01 := ChangeParamNode(node, &ReflectValueExpression{node.Node.Rv.Elem(), node.Node.Vk})
		compiler.internalProjectRv(w, tabCount, p01)
		break
	case kind == reflect.Bool:
		_, _ = io.WriteString(w, fmt.Sprintf("%v", node.Node.Rv.Bool()))
	case kind == reflect.Struct:
		if node.Node.Rv.CanInterface() {
			unk := node.Node.Rv.Interface()
			switch expr := unk.(type) {
			case Node[ast.Node]:
				compiler.internalProjectNode(w, tabCount, expr)
			default:
				if dataType, assigned, rvSomeType := compiler.isValueSomeDataType(node.Node.Rv); dataType {
					sss, _ := compiler.extractSomeDataTag(node.Node.Rv, "TData")
					_, _ = io.WriteString(w, fmt.Sprintf("/* Some[%v](assigned: %v) */\n", sss, assigned))
					_, _ = io.WriteString(w, strings.Repeat("\t", tabCount))

					if assigned {
						p01 := ChangeParamNode(node, &ReflectValueExpression{rvSomeType, ValueKey{"ccccc", "ddddd"}})
						compiler.internalProjectRv(w, tabCount, p01)
					} else {
						_, _ = io.WriteString(w, "nil")
					}
				} else {
					if node00, ok := node.Node.Rv.Interface().(ast.Node); ok {
						p01 := ChangeParamNode[*ReflectValueExpression, ast.Node](node, node00)
						compiler.internalProjectUnk(w, tabCount, p01)
					}

				}
			}
		} else {
			_, _ = io.WriteString(w, fmt.Sprintf("rv.CanInterface() == false"))
		}
	default:
		panic("unhandled default case")
	}
}

func (compiler *Compiler) internalProjectUnk(w io.Writer, tabCount int, node Node[ast.Node]) {
	switch nodeItem := node.Node.(type) {
	default:
		panic("implement me")
	case TrailSource:
		_, _ = io.WriteString(w, fmt.Sprintf("XXXXXXXXXXXXX %v XXXXXXXXXXXXX", nodeItem.Alias))

	case BooleanCondition:
		_, _ = io.WriteString(w, fmt.Sprintf("-- BooleanCondition%v\n", compiler.nodeOperator(nodeItem.op)))
		_, _ = io.WriteString(w, fmt.Sprintf("%v", strings.Repeat("\t", tabCount)))
		_, _ = io.WriteString(w, "(\n")

		for idx, condition := range nodeItem.conditions {
			_, _ = io.WriteString(w, fmt.Sprintf("%v", strings.Repeat("\t", tabCount+1)))
			_, _ = io.WriteString(w, fmt.Sprintf("-- %v\n", reflect.ValueOf(condition.Node).Type().String()))
			_, _ = io.WriteString(w, fmt.Sprintf("%v", strings.Repeat("\t", tabCount+1)))
			compiler.internalProjectNode(w, tabCount+1, condition)
			if idx != len(nodeItem.conditions)-1 {
				_, _ = io.WriteString(w, "\n")
				_, _ = io.WriteString(w, fmt.Sprintf("%v", strings.Repeat("\t", tabCount)))
				_, _ = io.WriteString(w, compiler.nodeOperator(nodeItem.op))
				_, _ = io.WriteString(w, "\n")
			}
		}
		_, _ = io.WriteString(w, "\n")
		_, _ = io.WriteString(w, fmt.Sprintf("%v", strings.Repeat("\t", tabCount)))
		_, _ = io.WriteString(w, ")")

	case *SupportedFunction:
		_, _ = io.WriteString(w, fmt.Sprintf("%v(", nodeItem.functionName))
		for idx, param := range nodeItem.params {
			compiler.internalProjectNode(w, tabCount, param)
			if idx != len(nodeItem.params)-1 {
				_, _ = io.WriteString(w, ", ")
			}
		}
		_, _ = io.WriteString(w, fmt.Sprintf(")"))
	case *LhsToMultipleRhsOperator:
		_, _ = io.WriteString(w, "(")
		for idx, rhs := range nodeItem.Rhs {
			_, _ = io.WriteString(w, "(")
			compiler.internalProjectNode(w, tabCount, nodeItem.Lhs)
			_, _ = io.WriteString(w, compiler.nodeOperator(nodeItem.LhsToRhsOp))
			compiler.internalProjectNode(w, tabCount, rhs)
			_, _ = io.WriteString(w, ")")
			if idx != len(nodeItem.Rhs)-1 {
				_, _ = io.WriteString(w, compiler.nodeOperator(nodeItem.betweenTerminalsOp))
			}
		}
		_, _ = io.WriteString(w, ")")
	case IfThenElseSingleValueCondition:
		_, _ = io.WriteString(w, "case\n")

		for _, expr := range nodeItem.conditionalStatement {
			tabCount++
			_, isLiteral := isLiterateValue(expr.condition)
			if !isLiteral {
				_, _ = io.WriteString(w, fmt.Sprintf("%vwhen ", strings.Repeat("\t", tabCount)))
				compiler.internalProjectNode(w, tabCount, expr.condition)
				_, _ = io.WriteString(w, " then\n")
			} else {
				_, _ = io.WriteString(w, fmt.Sprintf("%velse\n", strings.Repeat("\t", tabCount)))
			}
			tabCount++
			_, _ = io.WriteString(w, fmt.Sprintf("%v", strings.Repeat("\t", tabCount)))
			compiler.internalProjectNode(w, tabCount, expr.value)
			tabCount--
			_, _ = io.WriteString(w, "\n")
			tabCount--
		}
		_, _ = io.WriteString(w, fmt.Sprintf("%vend", strings.Repeat("\t", tabCount)))
	//case *MultiBinaryExpr:
	//	compiler.internalProjectUnk(w, tabCount, stackCount+1, *nodeItem)
	case MultiBinaryExpr:
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
			compiler.internalProjectNode(w, tabCount, expr)
		}
		_, _ = io.WriteString(w, ")")
	case BinaryExpr:
		_, _ = io.WriteString(w, "(")
		compiler.internalProjectNode(w, tabCount, nodeItem.left)
		_, _ = io.WriteString(w, compiler.nodeOperator(nodeItem.Op))
		compiler.internalProjectNode(w, tabCount, nodeItem.right)
		_, _ = io.WriteString(w, ")")
	case EntityField:
		_, _ = io.WriteString(w, fmt.Sprintf("[%v].[%v]", nodeItem.alias, nodeItem.field))
	case *CheckForNotNullExpression:
		p01 := ChangeParamNode[ast.Node, ast.Node](node, *nodeItem)
		compiler.internalProjectUnk(w, tabCount, p01)
	case CheckForNotNullExpression:
		_, _ = io.WriteString(w, "(")
		compiler.internalProjectNode(w, tabCount, nodeItem.node)
		_, _ = io.WriteString(w, " is not null)")
	case *coercion:
		p01 := ChangeParamNode[ast.Node, ast.Node](node, *nodeItem)
		compiler.internalProjectUnk(w, tabCount, p01)
	case coercion:
		_, _ = io.WriteString(w, "CAST(")
		compiler.internalProjectNode(w, tabCount, nodeItem.Node)
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
	}
}

func (compiler *Compiler) internalProjectNode(w io.Writer, tabCount int, node Node[ast.Node]) {
	if !node.Valid {
		return
	}
	switch nodeItem := node.Node.(type) {
	default:
		compiler.internalProjectUnk(w, tabCount, node)
	case *ReflectValueExpression:
		kind := nodeItem.Rv.Kind()
		switch kind {
		default:
			p01 := ChangeParamNode(node, nodeItem)
			compiler.internalProjectRv(w, tabCount, p01)
		}
	}
}

func (compiler *Compiler) findAdditionalSourcesFromJoins(sources map[string]ISource) map[string]ISource {
	for _, valueOuter := range compiler.JoinInformation {
		for keyInner, valueInner := range valueOuter.rhs {
			sources[keyInner] = valueInner
		}
	}
	for key, value := range compiler.JoinInformation {
		if _, ok := sources[key]; ok {
			sources[key] = value
		}
	}
	return sources
}

func (compiler *Compiler) projectSources(state State, w io.Writer, tabCount int, sources []ISource) {
	if len(sources) > 0 {
		for _, source := range sources {
			switch item := source.(type) {
			default:
				panic("unhandled default case")
			case PrimarySource:
				_, _ = fmt.Fprintf(w, "from\n")
				_, _ = io.WriteString(w, strings.Repeat("\t", tabCount))
				query, _ := compiler.Sources[item.sourceName]
				switch item := query.(type) {
				default:
					panic("unhandled default case")
				case *EntitySource:
					rt, _ := item.typeMapper.ActualType()
					_, _ = io.WriteString(w, fmt.Sprintf("%v [%v]", rt.String(), source.SourceName()))
				}
				_, _ = fmt.Fprintf(w, "\n")
			case JoinInformation:
				switch item.joinType {
				default:
					panic("dd")
				case jtLeftInner:
					_, _ = fmt.Fprintf(w, "inner join\n")
					_, _ = io.WriteString(w, strings.Repeat("\t", tabCount))
				case jtLeftOuter:
					_, _ = fmt.Fprintf(w, "left outer join\n")
					_, _ = io.WriteString(w, strings.Repeat("\t", tabCount))
				}
				query, _ := compiler.Sources[item.lhs]
				switch item := query.(type) {
				default:
					panic("unhandled default case")
				case *EntitySource:
					rt, _ := item.typeMapper.ActualType()
					_, _ = io.WriteString(w, fmt.Sprintf("%v [%v]", rt.String(), source.SourceName()))
				}
				_, _ = fmt.Fprintf(w, "\n")
				_, _ = io.WriteString(w, strings.Repeat("\t", tabCount))
				_, _ = fmt.Fprintf(w, "on\n")
				_, _ = io.WriteString(w, strings.Repeat("\t", tabCount+1))

				p2 := ChangeParamNode[BooleanCondition, ast.Node](item.condition, item.condition.Node)
				compiler.internalProjectNode(w, tabCount+1, p2)
				_, _ = fmt.Fprintf(w, "\n")
			}
		}
	}
}
