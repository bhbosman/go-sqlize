select
	((((T13.Name + ' ') + T13.Surname) + ' ') + 'ddd') as Name,
	'Bosman' as Surname,
	22 as Data
from
	struct { UserInformationId internal.Node[go/ast.Node]; Name internal.Node[go/ast.Node]; Surname internal.Node[go/ast.Node]; Data internal.Node[go/ast.Node] } T13
