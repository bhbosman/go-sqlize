package work

import "github.com/bhbosman/go-sqlize/lib"

type (
	Master struct {
		Id       int
		ParentId lib.Some[int]
		Name     string
		Surname  string
	}
)

func motherRelation(source Master, ass int) func(Master) bool {
	n := 12
	f := 34
	return func(target Master) bool {
		if parentId, ok := lib.GetSomeData(source.ParentId); ok {
			return parentId == target.Id+n+f+ass
		} else if source.Name == "dddd" {
			return true
		} else {
			return false
		}
	}
}

func init() {
	masterData := lib.Query[Master]()

	lib.GenerateSqlTest[Master](
		lib.Map(
			masterData,
			func(inputData Master) Master {
				preds := lib.CombinePredFunctionsWithAnd(
					[]func(master Master) bool{
						motherRelation(inputData, 22222),
						motherRelation(inputData, 23333),
						motherRelation(inputData, 24444)},
				)
				//preds := motherRelation(inputData)
				motherData := lib.Relationship(preds)
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
