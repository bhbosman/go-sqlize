select
	T10.UserInformationId as UserInformationId,
	'Brendan' as Name,
	'Bosman' as Surname,
	CAST(T10.Data as float) as Data
from
	struct { UserInformationId internal.Node[go/ast.Node]; Name internal.Node[go/ast.Node]; Surname internal.Node[go/ast.Node]; Data internal.Node[go/ast.Node] } T10
