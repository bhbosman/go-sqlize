package work

import (
	"github.com/bhbosman/go-sqlize/lib"
)

type UserInformation struct {
	UserInformationId int
	Name              string
	Surname           string
	Data              int
}

type View struct {
	UserInformationId int
	Name              string
	Surname           string
	Data              float64
}

func init() {
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[UserInformation](),
			func(inputData UserInformation) View {
				return View{
					inputData.UserInformationId,
					"Brendan",
					"Bosman",
					float64(inputData.Data),
				}
			},
		), "./output/sql0001_01.sql")
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[UserInformation](),
			func(inputData UserInformation) View {
				return View{
					inputData.UserInformationId,
					"Brendan",
					"Bosman",
					float64(inputData.Data) + 12 + 12*441 + 2 + 2 + 2 + 2,
				}
			},
		), "./output/sql0001_02.sql")
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[UserInformation](),
			func(inputData UserInformation) View {
				f := func(s UserInformation, sss string) string {
					return s.Name + " " + s.Surname + " " + sss
				}
				return View{
					inputData.UserInformationId,
					f(inputData, "ddd"),
					"Bosman",
					float64(inputData.Data) + 12 + 12*441 + 2 + 2 + 2 + 2,
				}
			},
		), "./output/sql0001_03.sql")
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[UserInformation](),
			func(inputData UserInformation) View {
				f := func(s UserInformation, sss string) string {
					return s.Name + " " + s.Surname + " " + sss
				}
				return View{
					Name:    f(inputData, "ddd"),
					Surname: "Bosman",
					Data:    float64(22),
				}
			},
		), "./output/sql0001_04.sql")
}
