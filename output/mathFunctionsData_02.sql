select
	-0.45990349068959124 as Value01,
	cos(T7.Value01) as Value02
from
	struct { Value01 internal.Node[go/ast.Node]; Value02 internal.Node[go/ast.Node] } T7
