package work

import (
	"github.com/bhbosman/go-sqlize/lib"
	"time"
)

func init() {
	type (
		Master struct {
			Id        int
			MotherId  lib.Some[int]
			FatherId  lib.Some[int]
			Name      string
			Surname   string
			BirthDate lib.Some[time.Time]
		}
		View struct {
			Id            int
			Name          string
			Surname       string
			MotherId      lib.Some[int]
			MotherName    lib.Some[string]
			MotherSurname lib.Some[string]
			FatherId      lib.Some[int]
			FatherName    lib.Some[string]
			FatherSurname lib.Some[string]
			OldestParent  lib.Some[time.Time]
		}
	)
	masterData := lib.Query[Master]()
	lib.GenerateSqlTest(
		lib.Map(masterData,
			func(inputData Master) View {
				motherId := lib.OptionalRelationship(
					func(target Master) bool {
						if inputDataMotherId, targetMotherId, ok := lib.GetSomeData02(inputData.MotherId, target.MotherId); ok {
							return inputDataMotherId == targetMotherId
						} else {
							return false
						}
					},
				)

				fatherId := lib.OptionalRelationship(
					func(target Master) bool {
						if inputDataFatherId, targetFatherId, ok := lib.GetSomeData02(inputData.FatherId, target.FatherId); ok {
							return inputDataFatherId == targetFatherId
						} else {
							return false
						}
					},
				)

				return View{
					Id:      inputData.Id,
					Name:    inputData.Name,
					Surname: inputData.Surname,
					MotherId: func() lib.Some[int] {
						if value, ok := lib.GetSomeData(motherId); ok {
							return lib.SetSomeValue[int](value.Id)
						} else {
							return lib.SetSomeNone[int]()
						}
					}(),
					MotherName: func() lib.Some[string] {
						if value, ok := lib.GetSomeData(motherId); ok {
							return lib.SetSomeValue[string](value.Name)
						} else {
							return lib.SetSomeNone[string]()
						}
					}(),
					MotherSurname: func() lib.Some[string] {
						if value, ok := lib.GetSomeData(motherId); ok {
							return lib.SetSomeValue[string](value.Surname)
						} else {
							return lib.SetSomeNone[string]()
						}
					}(),
					FatherId: func() lib.Some[int] {
						if value, ok := lib.GetSomeData(fatherId); ok {
							return lib.SetSomeValue[int](value.Id)
						} else {
							return lib.SetSomeNone[int]()
						}
					}(),
					FatherName: func() lib.Some[string] {
						if value, ok := lib.GetSomeData(fatherId); ok {
							return lib.SetSomeValue[string](value.Name)
						} else {
							return lib.SetSomeNone[string]()
						}
					}(),
					FatherSurname: func() lib.Some[string] {
						if value, ok := lib.GetSomeData(fatherId); ok {
							return lib.SetSomeValue[string](value.Surname)
						} else {
							return lib.SetSomeNone[string]()
						}
					}(),
					OldestParent: func() lib.Some[time.Time] {
						if fatherValue, motherValue, ok := lib.GetSomeData02(fatherId, motherId); ok {
							if fatherBirthDate, motherBirthDate, ok := lib.GetSomeData02(fatherValue.BirthDate, motherValue.BirthDate); ok {
								if fatherBirthDate.Before(motherBirthDate) {
									return lib.SetSomeValue(fatherBirthDate)
								} else {
									return lib.SetSomeValue(motherBirthDate)
								}
							} else {
								return lib.SetSomeNone[time.Time]()
							}
						} else {
							return lib.SetSomeNone[time.Time]()
						}
					}(),
				}
			},
		),
	)
}
