select
	T11.UserInformationId as UserInformationId,
	'Brendan' as Name,
	'Bosman' as Surname,
	((((((CAST(T11.Data as float) + 12) + 5292) + 2) + 2) + 2) + 2) as Data
from
	struct { UserInformationId internal.Node[go/ast.Node]; Name internal.Node[go/ast.Node]; Surname internal.Node[go/ast.Node]; Data internal.Node[go/ast.Node] } T11
