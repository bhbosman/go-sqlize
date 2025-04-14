package work

import "github.com/bhbosman/go-sqlize/lib"

type AA[TData interface{}] struct{ A TData }

func init() {
	type (
		TKey[TTData interface{}] struct {
			A struct{ C struct{ Z int } }
			B TTData
		}
		TValue struct {
			A int
		}
		TData = struct {
			A int
			B int
		}
		TView = struct {
			A int
		}
	)

	lookup := lib.CreateDictionary(
		map[TKey[int]]TValue{
			{struct{ C struct{ Z int } }{struct{ Z int }{1}}, 1}: {11},
			{struct{ C struct{ Z int } }{struct{ Z int }{1}}, 2}: {12},
			{struct{ C struct{ Z int } }{struct{ Z int }{1}}, 3}: {13},
		},
		TValue{0},
	)

	query := lib.Query[TData]()
	//
	mapFn := func(inputData TData) TValue {
		de := TKey[int]{struct{ C struct{ Z int } }{C: struct{ Z int }{Z: inputData.A}}, inputData.B}
		aa := lib.DictionaryLookup(lookup, de)
		return TValue{aa.A}
	}
	lib.GenerateSqlTest(lib.Map(query, mapFn))
}
