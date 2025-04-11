package work

import "github.com/bhbosman/go-sqlize/lib"

func init() {
	query := lib.Query[Switch01InputValues]()
	mapFn := func(inputData Switch01InputValues) Switch01InputValuesView {
		fn := func(Points01, Points02, Points03, Points04 lib.Some[int]) (lib.Some[string], string, lib.Some[string]) {
			if p1, p2, p3, p4, ok := lib.GetSomeData04(Points01, Points02, Points03, Points04); ok {
				switch {
				case p1 == 0 || p1 == 323 || p1 == 45345 || p1 == 4534234 || p1 == -34:
					value := p1 + p2 + p3 + p4
					return lib.SetSomeValue(lib.Itoa(value + 4)), lib.Itoa(value + 5), lib.SetSomeValue(lib.Itoa(value + 6))
				case p4 == 1:
					value := p1 + p2 + p3 + p4
					return lib.SetSomeValue(lib.Itoa(value + 7)), lib.Itoa(value + 8), lib.SetSomeValue(lib.Itoa(value + 9))
				case p3 == 2:
					value := p1 + p2 + p3 + p4
					return lib.SetSomeValue(lib.Itoa(value + 10)), lib.Itoa(value + 11), lib.SetSomeValue(lib.Itoa(value + 12))
				default:
					value := p1 + p2 + p3 + p4
					return lib.SetSomeValue(lib.Itoa(value + 1)), lib.Itoa(value + 2), lib.SetSomeValue(lib.Itoa(value + 3))
				}
			} else if p1, p2, p3, ok := lib.GetSomeData03(Points01, Points02, Points03); ok {
				value := p1 + p2 + p3
				return lib.SetSomeValue(lib.Itoa(value)), lib.Itoa(value), lib.SetSomeValue(lib.Itoa(value))
			} else if p1, p2, ok := lib.GetSomeData02(Points01, Points02); ok {
				value := p1 + p2
				return lib.SetSomeValue(lib.Itoa(value)), lib.Itoa(value), lib.SetSomeValue(lib.Itoa(value))
			} else if p1, ok := lib.GetSomeData(Points01); ok {
				value := p1
				return lib.SetSomeValue(lib.Itoa(value)), lib.Itoa(value), lib.SetSomeValue(lib.Itoa(value))
			} else {
				return lib.SetSomeValue("ABC"), "DEF", lib.SetSomeValue("GHI")
			}
		}
		l1, l2, l3 := fn(inputData.Points01, inputData.Points02, inputData.Points03, inputData.Points04)

		return Switch01InputValuesView{
			"",
			"",
			"",
			"",
			l1,
			l2,
			l3,
		}
	}
	lib.GenerateSqlTest(lib.Map(query, mapFn))
}
