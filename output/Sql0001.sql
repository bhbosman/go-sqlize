select
	"Brendan" as Name,
	"Bosman" as Surname,
	CAST(T1.Data as float) as Data
from
	struct { Name internal.Node[go/ast.Node]; Surname internal.Node[go/ast.Node]; Data internal.Node[go/ast.Node] } T1