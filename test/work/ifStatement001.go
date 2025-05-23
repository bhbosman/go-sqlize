package work

import "github.com/bhbosman/go-sqlize/lib"

type UserInformationForIf struct {
	Name    string
	SurName string
	Active  bool
	Points  int
}
type UserInformationForIfView struct {
	Name    string
	SurName string
	Status1 string
	Status2 string
	Level1  string
	Level2  string
}

func init() {
	levelFn := func(inputData UserInformationForIf) string {
		if inputData.Points < 100 {
			return "Value 01"
		} else if inputData.Points > 1000 {
			return "Value 02"
		} else if inputData.Points >= 100 && inputData.Points < 1000 {
			if b := inputData.Active && inputData.Points >= 100 && inputData.Points < 1000; b {
				s := func() string {
					d := "12"
					if inputData.Active || inputData.Points >= 100 && inputData.Points < 1000 {
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
		} else if inputData.Points >= 1000 && inputData.Points < 1000 {
			return "Value 04"
		} else if inputData.Points >= 100 && inputData.Points < 1000 {
			return "Value 05"
		} else if inputData.Points >= 100 && inputData.Points < 1000 {
			return "Value 06"
		} else if inputData.Points >= 100 && inputData.Points < 1000 {
			return "Value 07"
		} else if inputData.Points >= 100 && inputData.Points < 1000 {
			return "Value 08"
		} else if inputData.Points >= 100 && inputData.Points < 1000 {
			return "Value 09"
		} else {
			return "Value 10"
		}
	}

	statusFn := func(inputData UserInformationForIf) (string, string) {
		if !inputData.Active {
			return "activeX", "activeY"
		} else {
			return "dddd", "eeeeee"
		}
	}
	query := lib.Query[UserInformationForIf]()
	mapFn2 := func(inputData UserInformationForIf) UserInformationForIfView {
		x, y := statusFn(inputData)
		return UserInformationForIfView{
			Name:    inputData.Name,
			SurName: inputData.SurName,
			Status1: x,
			Status2: y,
			Level1:  levelFn(inputData),
			Level2:  levelFn(inputData),
		}
	}
	mapValue := lib.Map(query, mapFn2)
	lib.GenerateSqlTest(mapValue)
}
