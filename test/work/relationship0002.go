package work

import "github.com/bhbosman/go-sqlize/lib"

func init() {
	type (
		Master struct {
			Id       int
			ParentId lib.Some[int]
			Name     string
			Surname  string
		}
	)
	motherRelation := func(source Master) func(Master) bool {
		n := 12
		f := 34
		return func(target Master) bool {
			if parentId, ok := lib.GetSomeData(source.ParentId); ok {
				return parentId == target.Id+n+f
			} else {
				return false
			}
		}
	}

	masterData := lib.Query[Master]()

	lib.GenerateSqlTest[Master](
		lib.Map(
			masterData,
			func(inputData Master) Master {
				motherData := lib.Relationship(motherRelation(inputData))
				v := lib.SetSomeValue(motherData.Id)
				vv, _ := lib.GetSomeData(v)
				return Master{
					Id:       inputData.Id,
					ParentId: lib.SetSomeValue(vv),
					Name:     motherData.Name,
					Surname:  motherData.Surname,
				}
			}))
}
