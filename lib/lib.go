package lib

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// ToDo: define this type and refactor all the functions that use the syntax
//type BooleanExpression[TTarget interface{}] func(TTarget) bool

type IQueryOpt interface {
}

func QueryTop(int) IQueryOpt {
	panic("implement QueryTop")
}

func QueryDistinct() IQueryOpt {
	panic("implement QueryDistinct")
}

func Query[TInputData interface{}](...IQueryOpt) []TInputData { panic("implement Query") }

func Map[TInputData interface{}, TOutputData interface{}]([]TInputData, func(TInputData) TOutputData) []TOutputData {
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

type Some[TData interface{}] struct{}

func SetSomeValue[TData interface{}](TData) Some[TData] {
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

func SomeData2[TData interface{}](Some[TData]) (TData, bool) {
	panic("implement SomeData2")
}

func GetSomeData[TData interface{}](Some[TData]) (TData, bool) {
	panic("implement GetSomeData")
}

func GetSomeData02[TData01, TData02 interface{}](Some[TData01], Some[TData02]) (TData01, TData02, bool) {
	panic("implement GetSomeData02")
}

func GetSomeData03[TData01, TData02, TData03 interface{}](Some[TData01], Some[TData02], Some[TData03]) (TData01, TData02, TData03, bool) {
	panic("implement GetSomeData03")
}

func GetSomeData04[TData01, TData02, TData03, TData04 interface{}](Some[TData01], Some[TData02], Some[TData03], Some[TData04]) (TData01, TData02, TData03, TData04, bool) {
	panic("implement GetSomeData04")
}

func GetSomeData05[TData01, TData02, TData03, TData04, TData05 interface{}](Some[TData01], Some[TData02], Some[TData03], Some[TData04], Some[TData05]) (TData01, TData02, TData03, TData04, TData05, bool) {
	panic("implement GetSomeData05")
}

type Dictionary[TKey comparable, TValue interface{}] struct {
}

func CreateDictionary[TKey comparable, TValue interface{}](m map[TKey]TValue, defaultValue TValue) Dictionary[TKey, TValue] {
	panic("implement CreateDictionary")
}

func DictionaryLookup[TKey comparable, TValue interface{}](Dictionary[TKey, TValue], TKey) TValue {
	panic("implement DictionaryLookup")
}
func DictionaryDefault[TKey comparable, TValue interface{}](Dictionary[TKey, TValue]) TValue {
	panic("implement DictionaryDefault")
}

type IRelationshipOpt interface{}

func CombinePredFunctionsWithAnd[TTarget interface{}](...func(TTarget) bool) func(TTarget) bool {
	panic("implement CombinePredFunctionsWithAnd")
}

func CoreRelationship[TTarget interface{}]([]TTarget, func(TTarget) bool, ...IRelationshipOpt) TTarget {
	panic("implement CoreRelationship")
}

func Relationship[TTarget interface{}](pred func(TTarget) bool, opts ...IRelationshipOpt) TTarget {
	return CoreRelationship[TTarget](Query[TTarget](), pred, opts...)
}

func CoreOptionalRelationship[TTarget interface{}]([]TTarget, func(TTarget) bool, ...IRelationshipOpt) Some[TTarget] {
	panic("implement CoreOptionalRelationship")
}

func OptionalRelationship[TTarget interface{}](pred func(TTarget) bool, opts ...IRelationshipOpt) Some[TTarget] {
	return CoreOptionalRelationship[TTarget](Query[TTarget](), pred, opts...)
}
