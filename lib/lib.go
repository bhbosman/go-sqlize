package lib

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func Query[TInputData interface{}]() []TInputData { panic("implement Query") }

func Map[TInputData interface{}, TOutputData interface{}](inputData []TInputData, cb func(inputData TInputData) TOutputData) []TOutputData {
	panic("implement Map")
}

func Itoa(i int) string {
	return strconv.Itoa(i)
}

func Atoi(s string) int {
	panic("implement Atoi")
}

func GenerateSqlFile[TInputData interface{}](data []TInputData, outputFile string) {
	sql := GenerateSql(data)
	wd, _ := os.Getwd()
	fileName := filepath.Join(wd, outputFile)
	dir := filepath.Dir(fileName)
	_ = os.MkdirAll(dir, os.ModePerm)
	writer, _ := os.Create(fileName)
	_, _ = io.WriteString(writer, sql)
	_ = writer.Close()
}

func GenerateSqlStdOut[TInputData interface{}](data []TInputData) {
	sql := GenerateSql(data)
	writer := os.Stdout
	_, _ = io.WriteString(writer, sql)
	_ = writer.Close()
}

func GenerateSqlTest[TInputData interface{}](data []TInputData) {
	panic("implement GenerateSqlTest")
}

func GenerateSql[TInputData interface{}](data []TInputData) string {
	panic("implement GenerateSql")
}

type Some[TData interface{}] struct {
	data       TData
	isAssigned bool
}

//func (self Some[TData]) Set(data TData) Some[TData] {
//	return Some[TData]{data, true}
//}
//
//func (self Some[TData]) ToNone() Some[TData] {
//	return Some[TData]{reflect.Zero(reflect.TypeFor[TData]()).Interface().(TData), false}
//}
//
//func (self Some[TData]) Get() TData {
//	if self.isAssigned {
//		return self.data
//	}
//	return self.ToNone().data
//}
//
//func (self Some[TData]) IsAssigned() bool {
//	return self.isAssigned
//}

func SetSomeValue[TData interface{}](data TData) Some[TData] {
	panic("implement SetSomeValue")

}

func SetSomeNone[TData interface{}]() Some[TData] {
	panic("implement SetSomeNone")

}

func IsSomeAssigned[TData interface{}](Some[TData]) bool {
	panic("implement IsSomeAssigned")

}

func SomeData[TData interface{}](Some[TData]) TData {
	panic("implement SomeData")

}

//func GetValue[TData interface{}](data Some[TData]) TData {
//	return data.Get()
//}
//func IsAssigned[TData interface{}](data Some[TData]) bool {
//	return data.isAssigned
//}
