package work

import "github.com/bhbosman/go-sqlize/lib"

type Dict003[TData interface{}] struct {
	A TData
}

func init() {
	//lookup := lib.CreateDictionary(
	//	map[struct {
	//		A struct{ C struct{ Z int } }
	//		B int
	//	}]struct{ A int }{
	//		{struct{ C struct{ Z int } }{struct{ Z int }{1}}, 1}: {11},
	//		{struct{ C struct{ Z int } }{struct{ Z int }{1}}, 2}: {12},
	//		{struct{ C struct{ Z int } }{struct{ Z int }{1}}, 3}: {13},
	//	},
	//	struct{ A int }{0})

	lookup := lib.CreateDictionary(
		map[struct {
			A int
			B int
		}]struct{ A int }{
			{1, 1}: {11},
			{1, 2}: {12},
			{1, 3}: {13},
		},
		struct{ A int }{0})

	query := lib.Query[struct{ A, B int }]()

	mapFn := func(inputData struct{ A, B int }) struct{ A int } {
		dd := lib.DictionaryLookup(lookup, struct {
			A int
			B int
		}{inputData.A, inputData.B})
		return struct{ A int }{dd.A}
	}
	lib.GenerateSqlTest(lib.Map(query, mapFn))
}
