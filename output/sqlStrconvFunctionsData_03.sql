select
	CAST(T3.Value02 as int) as Value01,
	'123' as Value02
from
	struct { Value01 internal.Node[go/ast.Node]; Value02 internal.Node[go/ast.Node] } T3
