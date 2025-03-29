select
	T1.Name as Name,
	T1.SurName as SurName,
	(
		PartialExpressions(2)
		T1.Active
			'activeX'
		true
			'dddd'
	) as Status1,
	(
		PartialExpressions(2)
		T1.Active
			'activeY'
		true
			'eeeeee'
	) as Status2,
	(
		PartialExpressions(10)
		(T1.Points < 100)
			'Value 01'
		(T1.Points > 1000)
			'Value 02'
		((T1.Points >= 100) && (T1.Points < 1000))
			(
			PartialExpressions(2)
			((T1.Active && (T1.Points >= 100)) && (T1.Points < 1000))
				(
				PartialExpressions(2)
				((T1.Active && (T1.Points >= 100)) && (T1.Points < 1000))
					'2'
				true
					'1'
			)
			true
				'Value 03'
		)
		((T1.Points >= 1000) && (T1.Points < 1000))
			'Value 04'
		((T1.Points >= 100) && (T1.Points < 1000))
			'Value 05'
		((T1.Points >= 100) && (T1.Points < 1000))
			'Value 06'
		((T1.Points >= 100) && (T1.Points < 1000))
			'Value 07'
		((T1.Points >= 100) && (T1.Points < 1000))
			'Value 08'
		((T1.Points >= 100) && (T1.Points < 1000))
			'Value 09'
		true
			'Value 10'
	) as Level
from
	struct { Name internal.Node[go/ast.Node]; SurName internal.Node[go/ast.Node]; Active internal.Node[go/ast.Node]; Points internal.Node[go/ast.Node] } T1
