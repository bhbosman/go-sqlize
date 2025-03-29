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
	Level   string
}

func init() {
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[UserInformationForIf](),
			func(inputData UserInformationForIf) UserInformationForIfView {
				statusFn := func() (string, string) {
					if inputData.Active {
						return "activeX", "activeY"
					} else {
						return "dddd", "eeeeee"
					}
				}
				levelFn := func() string {
					if inputData.Points < 100 {
						return "Value 01"
					} else if inputData.Points > 1000 {
						return "Value 02"
					} else if inputData.Points >= 100 && inputData.Points < 1000 {
						if b := inputData.Active && inputData.Points >= 100 && inputData.Points < 1000; b {
							s := func() string {
								if inputData.Active && inputData.Points >= 100 && inputData.Points < 1000 {
									return "2"
								} else {
									return "1"
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
				x, y := statusFn()
				return UserInformationForIfView{
					Name:    inputData.Name,
					SurName: inputData.SurName,
					Status1: x,
					Status2: y,
					Level:   levelFn(),
				}
			},
		), "./output/if_01.sql")
}
