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

	masterData := lib.Query[Master]()

	lib.GenerateSqlTest(lib.Map(masterData, func(inputData Master) Master {
		childData := lib.Relationship(func(target Master) bool {
			if parentId, ok := lib.GetSomeData(target.ParentId); ok {
				return parentId == inputData.Id
			} else {
				return false
			}
		})
		return Master{
			Id:       childData.Id,
			ParentId: childData.ParentId,
			Name:     childData.Name,
			Surname:  childData.Surname,
		}
	}))
}
