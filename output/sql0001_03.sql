select
	T12.UserInformationId as UserInformationId,
	((((T12.Name + ' ') + T12.Surname) + ' ') + 'ddd') as Name,
	'Bosman' as Surname,
	((((((CAST(T12.Data as float) + 12) + 5292) + 2) + 2) + 2) + 2) as Data
from
	struct { UserInformationId internal.Node[go/ast.Node]; Name internal.Node[go/ast.Node]; Surname internal.Node[go/ast.Node]; Data internal.Node[go/ast.Node] } T12
