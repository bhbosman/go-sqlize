select
	CAST(T4.Value02 as int) as Value01,
	CAST(T4.Value01 as varchar) as Value02
from
	struct { Value01 internal.Node[go/ast.Node]; Value02 internal.Node[go/ast.Node] } T4
