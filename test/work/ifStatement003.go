package work

import "github.com/bhbosman/go-sqlize/lib"

func init() {

	statusFn := func(inputData UserInformationForIf02) (string, string) {
		if !inputData.Active {
			return "activeX", "activeY"
		} else {
			return "dddd", "eeeeee"
		}
	}
	query := lib.Query[UserInformationForIf02]()
	mapFn := func(inputData UserInformationForIf02) UserInformationForIfView02 {
		_, y := statusFn(inputData)
		return UserInformationForIfView02{
			//Status1: x,
			Status2: y,
		}
	}
	lib.GenerateSqlTest(lib.Map(query, mapFn))
}
