package work

import "github.com/bhbosman/go-sqlize/lib"

type Switch01InputValues struct {
	Name     string
	SurName  string
	Active   bool
	Points01 lib.Some[int]
	Points02 lib.Some[int]
	Points03 lib.Some[int]
	Points04 lib.Some[int]
}
type Switch01InputValuesView struct {
	Name    string
	SurName string
	Status1 string
	Status2 string
	Level1  lib.Some[string]
	Level2  string
	Level3  lib.Some[string]
}

func init() {
	query := lib.Query[Switch01InputValues]()
	mapFn := func(inputData Switch01InputValues) Switch01InputValuesView {

		fn := func(Points01, Points02, Points03, Points04 lib.Some[int]) (lib.Some[string], string, lib.Some[string]) {
			if p1, p2, p3, p4, ok := lib.GetSomeData04(Points01, Points02, Points03, Points04); ok {
				switch p1 {
				case 0, 323, 45345, 4534234, -34:
					value := p1 + p2 + p3 + p4
					return lib.SetSomeValue(lib.Itoa(value)), lib.Itoa(value), lib.SetSomeValue(lib.Itoa(value))
				case 1:
					value := p1 + p2 + p3 + p4
					return lib.SetSomeValue(lib.Itoa(value)), lib.Itoa(value), lib.SetSomeValue(lib.Itoa(value))
				case 2:
					value := p1 + p2 + p3 + p4
					return lib.SetSomeValue(lib.Itoa(value)), lib.Itoa(value), lib.SetSomeValue(lib.Itoa(value))
				default:
					value := p1 + p2 + p3 + p4
					return lib.SetSomeValue(lib.Itoa(value)), lib.Itoa(value), lib.SetSomeValue(lib.Itoa(value))
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
