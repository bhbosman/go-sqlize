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
		childData01 := lib.Relationship(func(target Master) bool {
			if parentId, ok := lib.GetSomeData(target.ParentId); ok {
				return parentId == inputData.Id
			} else {
				return false
			}
		})
		childData02 := lib.Relationship(func(target Master) bool {
			if parentId, ok := lib.GetSomeData(target.ParentId); ok {
				return parentId == inputData.Id
			} else {
				return false
			}
		})
		childData03 := lib.Relationship(func(target Master) bool {
			if parentId, ok := lib.GetSomeData(target.ParentId); ok {
				return parentId == inputData.Id
			} else {
				return false
			}
		})
		//lib.CombinePredFunctionsWithAnd([]func(target Master) bool{
		//	func(target Master) bool { return false }},
		//)

		//ss := lib.CombinePredFunctionsWithAnd(
		//	func(target Master) bool {
		//		if parentId, ok := lib.GetSomeData(target.ParentId); ok {
		//			return parentId == inputData.Id
		//		} else {
		//			return false
		//		}
		//	},
		//	func(target Master) bool {
		//		if parentId, ok := lib.GetSomeData(target.ParentId); ok {
		//			return parentId == inputData.Id
		//		} else {
		//			return false
		//		}
		//	},
		//	func(target Master) bool {
		//		if parentId, ok := lib.GetSomeData(target.ParentId); ok {
		//			return parentId == inputData.Id
		//		} else {
		//			return false
		//		}
		//	},
		//)
		//childData04 := lib.Relationship(ss)
		return Master{
			Id:       childData01.Id,
			ParentId: childData02.ParentId,
			Name:     childData03.Name,
			//Surname:  childData04.Surname,
		}
	}))
}
