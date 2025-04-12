package work

import "github.com/bhbosman/go-sqlize/lib"

func init() {
	query := lib.Query[Switch01InputValues]()
	mapFn := func(inputData Switch01InputValues) Switch01InputValuesView {
		l1 := func(Points lib.Some[int]) lib.Some[string] {
			dict := lib.CreateDictionary[int, lib.Some[string]](
				map[int]lib.Some[string]{
					1: lib.SetSomeValue("1"),
					2: lib.SetSomeNone[string](),
					3: lib.SetSomeValue("3"),
					4: lib.SetSomeValue("4"),
					5: lib.SetSomeValue("5"),
				},
				lib.SetSomeNone[string]())

			if p, ok := lib.GetSomeData(Points); ok {
				return lib.DictionaryLookup(dict, p)
			} else {
				return lib.DictionaryDefault(dict)
			}
		}(inputData.Points01)
		l2 := func(Points lib.Some[int]) string {
			dict := lib.CreateDictionary[int, string](
				map[int]string{
					1: "11",
					2: "22",
					3: "33",
					4: "44",
					5: "55",
				},
				"99")
			if p, ok := lib.GetSomeData(Points); ok {
				return lib.DictionaryLookup(dict, p)
			} else {
				return lib.DictionaryDefault(dict)
			}
		}(inputData.Points02)
		l3 := func(Points lib.Some[int]) string {
			dict := lib.CreateDictionary[int, string](
				map[int]string{
					1: "111",
					2: "222",
					3: "333",
					4: "444",
					5: "555",
				},
				"99")
			if p, ok := lib.GetSomeData(Points); ok {
				return lib.DictionaryLookup(dict, p)
			} else {
				return lib.DictionaryDefault(dict)
			}
		}(inputData.Points03)

		return Switch01InputValuesView{
			Name:    "",
			SurName: "",
			Status1: "",
			Status2: "",
			Level1:  l1,
			Level2:  l2,
			Level3:  lib.SetSomeValue(l3),
		}
	}
	lib.GenerateSqlTest(lib.Map(query, mapFn))
}
