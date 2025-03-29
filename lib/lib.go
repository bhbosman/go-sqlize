package lib

import "os"

func Query[TInputData interface{}]() []TInputData { panic("implement me") }

//type MapCallback[TInputData interface{}, TOutputData interface{}] func(inputData TInputData) TOutputData

func Map[TInputData interface{}, TOutputData interface{}](inputData []TInputData, cb func(inputData TInputData) TOutputData) []TOutputData {
	panic("implement me")
}

func Save() {
	getwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	println(getwd)
}

func GenerateSql[TInputData interface{}](data []TInputData, outputFile string) {
	panic("implement me")
}

//type Some[TData interface{}] struct {
//	data       TData
//	isAssigned bool
//}
//
//func (self Some[TData]) Set(data TData) Some[TData] {
//	return Some[TData]{data, true}
//}
//func (self Some[TData]) ToNone() Some[TData] {
//	return Some[TData]{reflect.Zero(reflect.TypeFor[TData]()).Interface().(TData), false}
//}
//func (self Some[TData]) Get() TData {
//	if self.isAssigned {
//		return self.data
//	}
//	return self.ToNone().data
//}
//func (self Some[TData]) IsAssigned() bool {
//	return self.isAssigned
//}
//func SetValue[TData interface{}](data TData) Some[TData] {
//	return Some[TData]{}.Set(data)
//}
//func SetValueV02[TData interface{}](data TData) Some[TData] {
//	return Some[TData]{}.Set(data)
//}
//
//func SetNone[TData interface{}]() Some[TData] {
//	return Some[TData]{}.ToNone()
//}
//func GetValue[TData interface{}](data Some[TData]) TData {
//	return data.Get()
//}
//func IsAssigned[TData interface{}](data Some[TData]) bool {
//	return data.isAssigned
//}
