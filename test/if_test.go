package test

import (
	"github.com/bhbosman/go-sqlize/cmd"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestName(t *testing.T) {

	ss := []struct {
		testName string
		setout   bool
		fileName string
		result   string
	}{
		{"if test 01", true, "./work/ifStatement001.go", "select\n\t[T1].[Name] as Name,\n\t[T1].[SurName] as SurName,\n\tcase\n\t\twhen ([T1].[Active] <> true) then\n\t\t\t'activeX'\n\t\telse\n\t\t\t'dddd'\n\tend as Status1,\n\tcase\n\t\twhen ([T1].[Active] <> true) then\n\t\t\t'activeY'\n\t\telse\n\t\t\t'eeeeee'\n\tend as Status2,\n\tcase\n\t\twhen ([T1].[Points] < 100) then\n\t\t\t'Value 01'\n\t\twhen ([T1].[Points] > 1000) then\n\t\t\t'Value 02'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\tcase\n\t\t\t\twhen (([T1].[Active] AND ([T1].[Points] >= 100)) AND ([T1].[Points] < 1000)) then\n\t\t\t\t\tcase\n\t\t\t\t\t\twhen ([T1].[Active] OR (([T1].[Points] >= 100) AND ([T1].[Points] < 1000))) then\n\t\t\t\t\t\t\t'122'\n\t\t\t\t\t\telse\n\t\t\t\t\t\t\t'122233'\n\t\t\t\t\tend\n\t\t\t\telse\n\t\t\t\t\t'Value 03'\n\t\t\tend\n\t\twhen (([T1].[Points] >= 1000) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 04'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 05'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 06'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 07'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 08'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 09'\n\t\telse\n\t\t\t'Value 10'\n\tend as Level1,\n\tcase\n\t\twhen ([T1].[Points] < 100) then\n\t\t\t'Value 01'\n\t\twhen ([T1].[Points] > 1000) then\n\t\t\t'Value 02'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\tcase\n\t\t\t\twhen (([T1].[Active] AND ([T1].[Points] >= 100)) AND ([T1].[Points] < 1000)) then\n\t\t\t\t\tcase\n\t\t\t\t\t\twhen ([T1].[Active] OR (([T1].[Points] >= 100) AND ([T1].[Points] < 1000))) then\n\t\t\t\t\t\t\t'122'\n\t\t\t\t\t\telse\n\t\t\t\t\t\t\t'122233'\n\t\t\t\t\tend\n\t\t\t\telse\n\t\t\t\t\t'Value 03'\n\t\t\tend\n\t\twhen (([T1].[Points] >= 1000) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 04'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 05'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 06'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 07'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 08'\n\t\twhen (([T1].[Points] >= 100) AND ([T1].[Points] < 1000)) then\n\t\t\t'Value 09'\n\t\telse\n\t\t\t'Value 10'\n\tend as Level2\nfrom\n\tstruct { Name internal.Node[go/ast.Node]; SurName internal.Node[go/ast.Node]; Active internal.Node[go/ast.Node]; Points internal.Node[go/ast.Node] } [T1]\n"},
		{"if test 02", true, "./work/ifStatement002.go", "select\n\t[T1].[Name] as Name,\n\t[T1].[SurName] as SurName,\n\tcase\n\t\twhen ([T1].[Active] <> true) then\n\t\t\t'activeX'\n\t\telse\n\t\t\t'dddd'\n\tend as Status1,\n\tcase\n\t\twhen ([T1].[Active] <> true) then\n\t\t\t'activeY'\n\t\telse\n\t\t\t'eeeeee'\n\tend as Status2,\n\tcase\n\t\twhen [T1].[Active] then\n\t\t\t /*SomeDataWithNode(assigned:false) */ nil\n\t\telse\n\t\t\t /*SomeDataWithNode(assigned:true) */ case\n\t\t\t\twhen (((([T1].[Points01] is not null) AND ([T1].[Points02] is not null)) AND ([T1].[Points01] < 100)) AND ([T1].[Points02] < 200)) then\n\t\t\t\t\t'Value 01'\n\t\t\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points01] > 1000)) then\n\t\t\t\t\t'Value 02'\n\t\t\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t\t\tcase\n\t\t\t\t\t\twhen ((([T1].[Active] AND ([T1].[Points01] is not null)) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t\t\t\t\tcase\n\t\t\t\t\t\t\t\twhen ([T1].[Active] OR ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000))) then\n\t\t\t\t\t\t\t\t\t'122'\n\t\t\t\t\t\t\t\telse\n\t\t\t\t\t\t\t\t\t'122233'\n\t\t\t\t\t\t\tend\n\t\t\t\t\t\telse\n\t\t\t\t\t\t\t'Value 03'\n\t\t\t\t\tend\n\t\t\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 1000)) AND ([T1].[Points01] < 1000)) then\n\t\t\t\t\t'Value 04'\n\t\t\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t\t\t'Value 05'\n\t\t\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t\t\t'Value 06'\n\t\t\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t\t\t'Value 07'\n\t\t\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t\t\t'Value 08'\n\t\t\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t\t\t'Value 09'\n\t\t\t\telse\n\t\t\t\t\t'Value 10'\n\t\t\tend\n\tend as Level1,\n\tcase\n\t\twhen (((([T1].[Points01] is not null) AND ([T1].[Points02] is not null)) AND ([T1].[Points01] < 100)) AND ([T1].[Points02] < 200)) then\n\t\t\t'Value 01'\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points01] > 1000)) then\n\t\t\t'Value 02'\n\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\tcase\n\t\t\t\twhen ((([T1].[Active] AND ([T1].[Points01] is not null)) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t\t\tcase\n\t\t\t\t\t\twhen ([T1].[Active] OR ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000))) then\n\t\t\t\t\t\t\t'122'\n\t\t\t\t\t\telse\n\t\t\t\t\t\t\t'122233'\n\t\t\t\t\tend\n\t\t\t\telse\n\t\t\t\t\t'Value 03'\n\t\t\tend\n\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 1000)) AND ([T1].[Points01] < 1000)) then\n\t\t\t'Value 04'\n\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t'Value 05'\n\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t'Value 06'\n\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t'Value 07'\n\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t'Value 08'\n\t\twhen ((([T1].[Points01] is not null) AND ([T1].[Points01] >= 100)) AND ([T1].[Points01] < 1000)) then\n\t\t\t'Value 09'\n\t\telse\n\t\t\t'Value 10'\n\tend as Level2,\n\t /*SomeDataWithNode(assigned:false) */ nil as Level3\nfrom\n\tstruct { Name internal.Node[go/ast.Node]; SurName internal.Node[go/ast.Node]; Active internal.Node[go/ast.Node]; Points01 internal.Node[go/ast.Node]; Points02 internal.Node[go/ast.Node]; Points03 internal.Node[go/ast.Node]; Points04 internal.Node[go/ast.Node] } [T1]\n"},
		{"if test 03", true, "./work/ifStatement003.go", "select\n\tcase\n\t\twhen ([T1].[Active] <> true) then\n\t\t\t'activeY'\n\t\telse\n\t\t\t'eeeeee'\n\tend as Status2,\nfrom\n\tstruct { Name internal.Node[go/ast.Node]; SurName internal.Node[go/ast.Node]; Active internal.Node[go/ast.Node]; Points01 internal.Node[go/ast.Node]; Points02 internal.Node[go/ast.Node]; Points03 internal.Node[go/ast.Node]; Points04 internal.Node[go/ast.Node] } [T1]\n"},
		{"switch test 01", true, "./work/switchStmt001.go", "select\n\t'' as Name,\n\t'' as SurName,\n\t'' as Status1,\n\t'' as Status2,\n\tcase\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null) AND ([T1].[Points04] is not null)) then\n\t\t\tcase\n\t\t\t\twhen (([T1].[Points01] = 0) OR ([T1].[Points01] = 323) OR ([T1].[Points01] = 45345) OR ([T1].[Points01] = 4534234) OR ([T1].[Points01] = (34 * -1))) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 4) as varchar)\n\t\t\t\twhen (([T1].[Points01] = 1)) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 7) as varchar)\n\t\t\t\twhen (([T1].[Points01] = 2)) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 10) as varchar)\n\t\t\t\telse\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 1) as varchar)\n\t\t\tend\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) as varchar)\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST(([T1].[Points01] + [T1].[Points02]) as varchar)\n\t\twhen (([T1].[Points01] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST([T1].[Points01] as varchar)\n\t\telse\n\t\t\t /*SomeDataWithRv(assigned:true)*/ 'ABC'\n\tend as Level1,\n\tcase\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null) AND ([T1].[Points04] is not null)) then\n\t\t\tcase\n\t\t\t\twhen (([T1].[Points01] = 0) OR ([T1].[Points01] = 323) OR ([T1].[Points01] = 45345) OR ([T1].[Points01] = 4534234) OR ([T1].[Points01] = (34 * -1))) then\n\t\t\t\t\tCAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 5) as varchar)\n\t\t\t\twhen (([T1].[Points01] = 1)) then\n\t\t\t\t\tCAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 8) as varchar)\n\t\t\t\twhen (([T1].[Points01] = 2)) then\n\t\t\t\t\tCAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 11) as varchar)\n\t\t\t\telse\n\t\t\t\t\tCAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 2) as varchar)\n\t\t\tend\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null)) then\n\t\t\tCAST((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) as varchar)\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null)) then\n\t\t\tCAST(([T1].[Points01] + [T1].[Points02]) as varchar)\n\t\twhen (([T1].[Points01] is not null)) then\n\t\t\tCAST([T1].[Points01] as varchar)\n\t\telse\n\t\t\t'DEF'\n\tend as Level2,\n\tcase\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null) AND ([T1].[Points04] is not null)) then\n\t\t\tcase\n\t\t\t\twhen (([T1].[Points01] = 0) OR ([T1].[Points01] = 323) OR ([T1].[Points01] = 45345) OR ([T1].[Points01] = 4534234) OR ([T1].[Points01] = (34 * -1))) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 6) as varchar)\n\t\t\t\twhen (([T1].[Points01] = 1)) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 9) as varchar)\n\t\t\t\twhen (([T1].[Points01] = 2)) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 12) as varchar)\n\t\t\t\telse\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 3) as varchar)\n\t\t\tend\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) as varchar)\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST(([T1].[Points01] + [T1].[Points02]) as varchar)\n\t\twhen (([T1].[Points01] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST([T1].[Points01] as varchar)\n\t\telse\n\t\t\t /*SomeDataWithRv(assigned:true)*/ 'GHI'\n\tend as Level3\nfrom\n\tstruct { Name internal.Node[go/ast.Node]; SurName internal.Node[go/ast.Node]; Active internal.Node[go/ast.Node]; Points01 internal.Node[go/ast.Node]; Points02 internal.Node[go/ast.Node]; Points03 internal.Node[go/ast.Node]; Points04 internal.Node[go/ast.Node] } [T1]\n"},
		{"switch test 02", true, "./work/switchStmt002.go", "select\n\t'' as Name,\n\t'' as SurName,\n\t'' as Status1,\n\t'' as Status2,\n\tcase\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null) AND ([T1].[Points04] is not null)) then\n\t\t\tcase\n\t\t\t\twhen ((((([T1].[Points01] = 0) OR ([T1].[Points01] = 323)) OR ([T1].[Points01] = 45345)) OR ([T1].[Points01] = 4534234)) OR ([T1].[Points01] = (34 * -1))) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 4) as varchar)\n\t\t\t\twhen ([T1].[Points04] = 1) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 7) as varchar)\n\t\t\t\twhen ([T1].[Points03] = 2) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 10) as varchar)\n\t\t\t\telse\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 1) as varchar)\n\t\t\tend\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) as varchar)\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST(([T1].[Points01] + [T1].[Points02]) as varchar)\n\t\twhen (([T1].[Points01] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST([T1].[Points01] as varchar)\n\t\telse\n\t\t\t /*SomeDataWithRv(assigned:true)*/ 'ABC'\n\tend as Level1,\n\tcase\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null) AND ([T1].[Points04] is not null)) then\n\t\t\tcase\n\t\t\t\twhen ((((([T1].[Points01] = 0) OR ([T1].[Points01] = 323)) OR ([T1].[Points01] = 45345)) OR ([T1].[Points01] = 4534234)) OR ([T1].[Points01] = (34 * -1))) then\n\t\t\t\t\tCAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 5) as varchar)\n\t\t\t\twhen ([T1].[Points04] = 1) then\n\t\t\t\t\tCAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 8) as varchar)\n\t\t\t\twhen ([T1].[Points03] = 2) then\n\t\t\t\t\tCAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 11) as varchar)\n\t\t\t\telse\n\t\t\t\t\tCAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 2) as varchar)\n\t\t\tend\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null)) then\n\t\t\tCAST((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) as varchar)\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null)) then\n\t\t\tCAST(([T1].[Points01] + [T1].[Points02]) as varchar)\n\t\twhen (([T1].[Points01] is not null)) then\n\t\t\tCAST([T1].[Points01] as varchar)\n\t\telse\n\t\t\t'DEF'\n\tend as Level2,\n\tcase\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null) AND ([T1].[Points04] is not null)) then\n\t\t\tcase\n\t\t\t\twhen ((((([T1].[Points01] = 0) OR ([T1].[Points01] = 323)) OR ([T1].[Points01] = 45345)) OR ([T1].[Points01] = 4534234)) OR ([T1].[Points01] = (34 * -1))) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 6) as varchar)\n\t\t\t\twhen ([T1].[Points04] = 1) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 9) as varchar)\n\t\t\t\twhen ([T1].[Points03] = 2) then\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 12) as varchar)\n\t\t\t\telse\n\t\t\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) + [T1].[Points04]) + 3) as varchar)\n\t\t\tend\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null) AND ([T1].[Points03] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST((([T1].[Points01] + [T1].[Points02]) + [T1].[Points03]) as varchar)\n\t\twhen (([T1].[Points01] is not null) AND ([T1].[Points02] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST(([T1].[Points01] + [T1].[Points02]) as varchar)\n\t\twhen (([T1].[Points01] is not null)) then\n\t\t\t /*SomeDataWithNode(assigned:true) */ CAST([T1].[Points01] as varchar)\n\t\telse\n\t\t\t /*SomeDataWithRv(assigned:true)*/ 'GHI'\n\tend as Level3\nfrom\n\tstruct { Name internal.Node[go/ast.Node]; SurName internal.Node[go/ast.Node]; Active internal.Node[go/ast.Node]; Points01 internal.Node[go/ast.Node]; Points02 internal.Node[go/ast.Node]; Points03 internal.Node[go/ast.Node]; Points04 internal.Node[go/ast.Node] } [T1]\n"},
		{"Dictionary Test 01", false, "./work/dict0001.go", ""},
		{"Dictionary Test 02", false, "./work/dict0002.go", ""},
		{"Dictionary Test 03", false, "./work/dict0003.go", ""},
		{"Dictionary Test 04", false, "./work/dict0004.go", ""},
		{"Some 001", false, "./work/Some00001.go", ""},
		//{"Relation 001", false, "./work/relationship0001.go", ""},

	}
	for _, s := range ss {
		t.Run(
			s.testName,
			func(t *testing.T) {
				cmd.RootCmd.SetArgs([]string{s.fileName})
				sb := new(strings.Builder)
				cmd.RootCmd.SetOut(sb)

				err := cmd.RootCmd.Execute()
				if err != nil {
					os.Exit(1)
				}
				if s.setout {
					require.Equal(t, s.result, sb.String())

				} else {
					println(sb.String())
				}
			},
		)
	}
}
