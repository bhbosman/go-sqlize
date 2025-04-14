package work

import "github.com/bhbosman/go-sqlize/lib"

func init() {
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
		dd := lib.DictionaryLookup(
			lookup,
			struct {
				A int
				B int
			}{inputData.A, inputData.B},
		)
		return struct{ A int }{dd.A}
	}
	lib.GenerateSqlTest(lib.Map(query, mapFn))
}
