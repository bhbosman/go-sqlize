package work

import "github.com/bhbosman/go-sqlize/lib"

type UserInformationForIf02 struct {
	Name     string
	SurName  string
	Active   bool
	Points01 lib.Some[int]
	Points02 lib.Some[int]
	Points03 lib.Some[int]
	Points04 lib.Some[int]
}
type UserInformationForIfView02 struct {
	Name    string
	SurName string
	Status1 string
	Status2 string
	Level1  lib.Some[string]
	Level2  string
	Level3  lib.Some[string]
}

func init() {
	levelFn := func(inputData UserInformationForIf02) string {
		if value01, value02, ok := lib.GetSomeData02(inputData.Points01, inputData.Points02); ok && value01 < 100 && value02 < 200 {
			return "Value 01"
		} else if b := lib.IsSomeAssigned(inputData.Points01); b && lib.SomeData(inputData.Points01) > 1000 {
			return "Value 02"
		} else if lib.IsSomeAssigned(inputData.Points01) && lib.SomeData(inputData.Points01) >= 100 && lib.SomeData(inputData.Points01) < 1000 {
			if b := inputData.Active && lib.IsSomeAssigned(inputData.Points01) && lib.SomeData(inputData.Points01) >= 100 && lib.SomeData(inputData.Points01) < 1000; b {
				s := func() string {
					d := "12"
					if inputData.Active || lib.IsSomeAssigned(inputData.Points01) && lib.SomeData(inputData.Points01) >= 100 && lib.SomeData(inputData.Points01) < 1000 {
						d = d + "2"
						return d
					} else {
						d = d + "233"
						return d
					}
				}
				return s()
			} else {
				return "Value 03"
			}
		} else if lib.IsSomeAssigned(inputData.Points01) && lib.SomeData(inputData.Points01) >= 1000 && lib.SomeData(inputData.Points01) < 1000 {
			return "Value 04"
		} else if lib.IsSomeAssigned(inputData.Points01) && lib.SomeData(inputData.Points01) >= 100 && lib.SomeData(inputData.Points01) < 1000 {
			return "Value 05"
		} else if lib.IsSomeAssigned(inputData.Points01) && lib.SomeData(inputData.Points01) >= 100 && lib.SomeData(inputData.Points01) < 1000 {
			return "Value 06"
		} else if lib.IsSomeAssigned(inputData.Points01) && lib.SomeData(inputData.Points01) >= 100 && lib.SomeData(inputData.Points01) < 1000 {
			return "Value 07"
		} else if lib.IsSomeAssigned(inputData.Points01) && lib.SomeData(inputData.Points01) >= 100 && lib.SomeData(inputData.Points01) < 1000 {
			return "Value 08"
		} else if lib.IsSomeAssigned(inputData.Points01) && lib.SomeData(inputData.Points01) >= 100 && lib.SomeData(inputData.Points01) < 1000 {
			return "Value 09"
		} else {
			return "Value 10"
		}
	}

	statusFn := func(inputData UserInformationForIf02) (string, string) {
		if !inputData.Active {
			return "activeX", "activeY"
		} else {
			return "dddd", "eeeeee"
		}
	}
	query := lib.Query[UserInformationForIf02]()
	mapFn := func(inputData UserInformationForIf02) UserInformationForIfView02 {
		x, y := statusFn(inputData)
		return UserInformationForIfView02{
			Name:    inputData.Name,
			SurName: inputData.SurName,
			Status1: x,
			Status2: y,
			Level1: func() lib.Some[string] {
				if inputData.Active {
					return lib.SetSomeNone[string]()

				} else {
					return lib.SetSomeValue(levelFn(inputData))
				}
			}(),
			Level2: levelFn(inputData),
			Level3: lib.SetSomeNone[string](),
		}
	}
	lib.GenerateSqlTest(lib.Map(query, mapFn))
}
