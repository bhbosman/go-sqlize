package work

import "github.com/bhbosman/go-sqlize/lib"

type LevelData struct {
	Level1 lib.Some[string]
	Level2 string
	Level3 lib.Some[string]
}

func init() {
	//v := lib.SetSomeNone[string]()
	v := lib.SetSomeValue("ddd")
	query := lib.Query[Switch01InputValues]()
	mapFn := func(inputData Switch01InputValues) Switch01InputValuesView {
		d := lib.CreateDictionary(map[int]LevelData{
			1: {lib.SetSomeValue("1"), "2", lib.SetSomeValue("3")},
			2: {lib.SetSomeNone[string](), "32", lib.SetSomeNone[string]()},
		}, LevelData{})

		l1 := func(Points01 lib.Some[int]) lib.Some[string] {
			if p1, _, ok := lib.GetSomeData02(Points01, v); ok {
				return lib.DictionaryLookup(d, p1).Level1
			} else {
				return lib.DictionaryDefault(d).Level1
			}
		}(inputData.Points01)

		return Switch01InputValuesView{
			Level1: l1,
		}
	}
	lib.GenerateSqlTest(lib.Map(query, mapFn))
}
