package work

import "github.com/bhbosman/go-sqlize/lib"

func init() {
	type (
		Master struct {
			Id       int
			MotherId lib.Some[int]
			FatherId lib.Some[int]
			Name     string
			Surname  string
		}
	)
	masterData := lib.Query[Master]()
	lib.GenerateSqlTest[Master](
		lib.Map(masterData,
			func(inputData Master) Master {
				childData01 := lib.Relationship(
					func(target Master) bool {
						if inputDataMotherId, inputDataFatherId, targetMotherId, targetFatherId, ok := lib.GetSomeData04(inputData.MotherId, inputData.FatherId, target.MotherId, target.FatherId); ok {
							return inputDataMotherId == targetMotherId && inputDataFatherId == targetFatherId
						} else if _, ok := lib.GetSomeData(target.MotherId); ok {
							if target.Surname == "Spijkerman" {
								return true
							} else {
								return false
							}
						} else if _, ok := lib.GetSomeData(target.FatherId); ok {
							if target.Surname == "Bosman" {
								return true
							} else {
								return false
							}
						} else if _, ok := lib.GetSomeData(target.FatherId); ok {
							if target.Surname == "Bosman" {
								return true
							} else {
								return false
							}
						} else if _, ok := lib.GetSomeData(target.FatherId); ok {
							if target.Surname == "Bosman" {
								return true
							} else {
								return false
							}
						} else {
							return false
						}
					})
				childData02 := lib.Relationship(
					func(target Master) bool {
						if parentId, ok := lib.GetSomeData(target.MotherId); ok {
							return parentId == childData01.Id
						} else {
							return false
						}
					})
				return Master{
					Id:       inputData.Id,
					MotherId: childData01.MotherId,
					Name:     childData02.Name,
					//Surname:  childData01.Surname,
				}
			}))
}
