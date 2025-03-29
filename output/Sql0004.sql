select
	((((T4.Name + " ") + T4.Surname) + " ") + "ddd") as Name,
	"Bosman" as Surname,
from
	struct { Name internal.Node[go/ast.Node]; Surname internal.Node[go/ast.Node]; Data internal.Node[go/ast.Node] } T4