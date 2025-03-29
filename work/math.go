package work

import (
	"github.com/bhbosman/go-sqlize/lib"
	"math"
)

type MathData struct {
	Value01 float64
	Value02 float64
}

type MathView struct {
	Value01 float64
	Value02 float64
}

func init() {
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[MathData](),
			func(inputData MathData) MathView {
				return MathView{
					Value01: math.Sin(123 + 123),
					Value02: math.Cos(123),
				}
			},
		), "./output/mathFunctionsData_01.sql")
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[MathData](),
			func(inputData MathData) MathView {
				return MathView{
					Value01: math.Sin(123),
					Value02: math.Cos(inputData.Value01),
				}
			},
		), "./output/mathFunctionsData_02.sql")
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[MathData](),
			func(inputData MathData) MathView {
				return MathView{
					Value01: math.Sin(inputData.Value02),
					Value02: math.Cos(123),
				}
			},
		), "./output/mathFunctionsData_03.sql")
	lib.GenerateSqlFile(
		lib.Map(
			lib.Query[MathData](),
			func(inputData MathData) MathView {
				return MathView{
					Value01: math.Sin(inputData.Value02),
					Value02: math.Cos(inputData.Value01),
				}
			},
		), "./output/mathFunctionsData_04.sql")

}
