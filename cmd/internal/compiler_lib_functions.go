package internal

import "reflect"

func TestType(rt reflect.Type, a interface{}, typeFor reflect.Type) {
	rv := reflect.ValueOf(a)
	rv = rv.Convert(rt)
	rv = rv.Convert(typeFor)
}
