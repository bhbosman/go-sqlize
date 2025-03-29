select
	sin(T9.Value02) as Value01,
	cos(T9.Value01) as Value02
from
	struct { Value01 internal.Node[go/ast.Node]; Value02 internal.Node[go/ast.Node] } T9
