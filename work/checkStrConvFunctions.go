package work

import "github.com/bhbosman/go-sqlize/lib"

type StrConvFunctionsData struct {
	Value01 int
	Value02 string
}

type StrConvFunctionsView struct {
	Value01 int
	Value02 string
}

func init() {
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[StrConvFunctionsData](),
			func(inputData StrConvFunctionsData) StrConvFunctionsView {
				return StrConvFunctionsView{
					Value01: lib.Atoi("123"),
					Value02: lib.Itoa(123),
				}
			},
		), "./output/sqlStrconvFunctionsData_01.sql")
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[StrConvFunctionsData](),
			func(inputData StrConvFunctionsData) StrConvFunctionsView {
				return StrConvFunctionsView{
					Value01: lib.Atoi("123"),
					Value02: lib.Itoa(inputData.Value01),
				}
			},
		), "./output/sqlStrconvFunctionsData_02.sql")
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[StrConvFunctionsData](),
			func(inputData StrConvFunctionsData) StrConvFunctionsView {
				return StrConvFunctionsView{
					Value01: lib.Atoi(inputData.Value02),
					Value02: lib.Itoa(123),
				}
			},
		), "./output/sqlStrconvFunctionsData_03.sql")
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[StrConvFunctionsData](),
			func(inputData StrConvFunctionsData) StrConvFunctionsView {
				return StrConvFunctionsView{
					Value01: lib.Atoi(inputData.Value02),
					Value02: lib.Itoa(inputData.Value01),
				}
			},
		), "./output/sqlStrconvFunctionsData_04.sql")

}
