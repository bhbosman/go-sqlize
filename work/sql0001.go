package work

import (
	"github.com/bhbosman/go-sqlize/data"
	"github.com/bhbosman/go-sqlize/lib"
)

func init() {
	lib.Save()
	lib.GenerateSql(
		lib.Map(
			lib.Query[data.Table](),
			func(inputData data.Table) data.View {
				return data.View{
					"Brendan",
					"Bosman",
					float64(inputData.Data),
				}
			},
		), "./output/Sql0001.sql")
	lib.GenerateSql(
		lib.Map(
			lib.Query[data.Table](),
			func(inputData data.Table) data.View {
				return data.View{
					"Brendan",
					"Bosman",
					float64(inputData.Data) + 12 + 12*441 + 2 + 2 + 2 + 2,
				}
			},
		), "./output/Sql0002.sql")
	lib.GenerateSql(
		lib.Map(
			lib.Query[data.Table](),
			func(inputData data.Table) data.View {
				f := func(s data.Table, sss string) string {
					return s.Name + " " + s.Surname + " " + sss
				}
				return data.View{
					f(inputData, "ddd"),
					"Bosman",
					float64(inputData.Data) + 12 + 12*441 + 2 + 2 + 2 + 2,
				}
			},
		), "./output/Sql0003.sql")
	lib.GenerateSql(
		lib.Map(
			lib.Query[data.Table](),
			func(inputData data.Table) data.View {
				f := func(s data.Table, sss string) string {
					return s.Name + " " + s.Surname + " " + sss
				}
				return data.View{
					Name:    f(inputData, "ddd"),
					Surname: "Bosman",
				}
			},
		), "./output/Sql0004.sql")
}
