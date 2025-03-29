select
	"Brendan" as Name,
	"Bosman" as Surname,
	((((((CAST(T2.Data as float) + 12) + (12 * 441)) + 2) + 2) + 2) + 2) as Data
from
	struct { Name internal.Node[go/ast.Node]; Surname internal.Node[go/ast.Node]; Data internal.Node[go/ast.Node] } T2