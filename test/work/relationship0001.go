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

	lib.GenerateSqlTest[Master](lib.Map(masterData, func(inputData Master) Master {
		fn := func(target Master) bool {
			//return target.Name == inputData.Name
			if parentId, ok := lib.GetSomeData(target.ParentId); ok {
				return parentId == inputData.Id

			} else if target.Name == inputData.Name {
				return true
			} else {
				return false
			}
		}
		childData01 := lib.Relationship(fn)
		childData02 := lib.Relationship(fn)
		childData03 := lib.Relationship(fn)
		return Master{
			Id:       childData01.Id,
			ParentId: childData02.ParentId,
			Name:     childData03.Name,
			Surname:  childData03.Surname,
		}
	}))
}
