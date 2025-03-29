select
	T5.Name as Name,
	T5.SurName as SurName,
	(
		PartialExpressions(2)
		T5.Active
			'activeX'
		true
			'dddd'
	) as Status1,
	(
		PartialExpressions(2)
		T5.Active
			'activeY'
		true
			'eeeeee'
	) as Status2,
	(
		PartialExpressions(10)
		(T5.Points < 100)
			'Value 01'
		(T5.Points > 1000)
			'Value 02'
		((T5.Points >= 100) && (T5.Points < 1000))
			(
			PartialExpressions(2)
			((T5.Active && (T5.Points >= 100)) && (T5.Points < 1000))
				(
				PartialExpressions(2)
				((T5.Active && (T5.Points >= 100)) && (T5.Points < 1000))
					'2'
				true
					'1'
			)
			true
				'Value 03'
		)
		((T5.Points >= 1000) && (T5.Points < 1000))
			'Value 04'
		((T5.Points >= 100) && (T5.Points < 1000))
			'Value 05'
		((T5.Points >= 100) && (T5.Points < 1000))
			'Value 06'
		((T5.Points >= 100) && (T5.Points < 1000))
			'Value 07'
		((T5.Points >= 100) && (T5.Points < 1000))
			'Value 08'
		((T5.Points >= 100) && (T5.Points < 1000))
			'Value 09'
		true
			'Value 10'
	) as Level
from
	struct { Name internal.Node[go/ast.Node]; SurName internal.Node[go/ast.Node]; Active internal.Node[go/ast.Node]; Points internal.Node[go/ast.Node] } T5
