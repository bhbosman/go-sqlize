select
	123 as Value01,
	CAST(T2.Value01 as varchar) as Value02
from
	struct { Value01 internal.Node[go/ast.Node]; Value02 internal.Node[go/ast.Node] } T2
