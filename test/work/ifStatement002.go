package work

import "github.com/bhbosman/go-sqlize/lib"

type UserInformationForIf02 struct {
	Name    string
	SurName string
	Active  bool
	Points  lib.Some[int]
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
		if lib.IsSomeAssigned(inputData.Points) && lib.SomeData(inputData.Points) < 100 {
			return "Value 01"
		} else if b := lib.IsSomeAssigned(inputData.Points); b && lib.SomeData(inputData.Points) > 1000 {
			return "Value 02"
		} else if lib.IsSomeAssigned(inputData.Points) && lib.SomeData(inputData.Points) >= 100 && lib.SomeData(inputData.Points) < 1000 {
			if b := inputData.Active && lib.IsSomeAssigned(inputData.Points) && lib.SomeData(inputData.Points) >= 100 && lib.SomeData(inputData.Points) < 1000; b {
				s := func() string {
					d := "12"
					if inputData.Active || lib.IsSomeAssigned(inputData.Points) && lib.SomeData(inputData.Points) >= 100 && lib.SomeData(inputData.Points) < 1000 {
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
		} else if lib.IsSomeAssigned(inputData.Points) && lib.SomeData(inputData.Points) >= 1000 && lib.SomeData(inputData.Points) < 1000 {
			return "Value 04"
		} else if lib.IsSomeAssigned(inputData.Points) && lib.SomeData(inputData.Points) >= 100 && lib.SomeData(inputData.Points) < 1000 {
			return "Value 05"
		} else if lib.IsSomeAssigned(inputData.Points) && lib.SomeData(inputData.Points) >= 100 && lib.SomeData(inputData.Points) < 1000 {
			return "Value 06"
		} else if lib.IsSomeAssigned(inputData.Points) && lib.SomeData(inputData.Points) >= 100 && lib.SomeData(inputData.Points) < 1000 {
			return "Value 07"
		} else if lib.IsSomeAssigned(inputData.Points) && lib.SomeData(inputData.Points) >= 100 && lib.SomeData(inputData.Points) < 1000 {
			return "Value 08"
		} else if lib.IsSomeAssigned(inputData.Points) && lib.SomeData(inputData.Points) >= 100 && lib.SomeData(inputData.Points) < 1000 {
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

	lib.GenerateSqlTest(
		lib.Map(
			lib.Query[UserInformationForIf02](),
			func(inputData UserInformationForIf02) UserInformationForIfView02 {
				x, y := statusFn(inputData)
				return UserInformationForIfView02{
					Name:    inputData.Name,
					SurName: inputData.SurName,
					Status1: x,
					Status2: y,
					Level1:  lib.SetSomeValue("ddddd"),
					Level2:  levelFn(inputData),
					Level3:  lib.SetSomeNone[string](),
				}
			},
		))
}
