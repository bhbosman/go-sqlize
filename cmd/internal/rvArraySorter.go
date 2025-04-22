package internal

import (
	"reflect"
	"strings"
)

type rvArraySorter struct {
	rvArray []reflect.Value
}

func (rvArray *rvArraySorter) Len() int {
	return len(rvArray.rvArray)
}

func (rvArray *rvArraySorter) walkStruct(rvI, rvJ reflect.Value) int {
	for i := 0; i < rvI.NumField(); i++ {
		rvIField := rvI.Field(i)
		rvJField := rvJ.Field(i)
		ans := rvArray.walkField(rvIField, rvJField)
		if ans != 0 {
			return ans
		}
	}
	return 0
}

func (rvArray *rvArraySorter) walkField(ith, jth reflect.Value) int {
	switch {
	case ith.CanInt():
		return int(ith.Int() - jth.Int())
	case ith.CanFloat():
		return int(ith.Float() - jth.Float())
	case ith.Kind() == reflect.String:
		return strings.Compare(ith.String(), jth.String())
	case ith.Kind() == reflect.Struct:
		return rvArray.walkStruct(ith, jth)
	default:
		return -1
	}
}

func (rvArray *rvArraySorter) Less(i, j int) bool {
	if rvArray.rvArray[i].Kind() == rvArray.rvArray[j].Kind() {
		ith := rvArray.rvArray[i]
		jth := rvArray.rvArray[j]
		return rvArray.walkField(ith, jth) < 0
	}
	return false
}

func (rvArray *rvArraySorter) Swap(i, j int) {
	rvArray.rvArray[i], rvArray.rvArray[j] = rvArray.rvArray[j], rvArray.rvArray[i]
}
