package work

import "github.com/bhbosman/go-sqlize/lib"

func init() {
	lib.TestType(int64(0), lib.TypeFor[int64]())
	lib.TestType(int32(0), lib.TypeFor[int32]())
	lib.TestType(int16(0), lib.TypeFor[int16]())
	lib.TestType(int8(0), lib.TypeFor[int8]())
	lib.TestType(int(0), lib.TypeFor[int]())

	lib.TestType(uint64(0), lib.TypeFor[uint64]())
	lib.TestType(uint32(0), lib.TypeFor[uint32]())
	lib.TestType(uint16(0), lib.TypeFor[uint16]())
	lib.TestType(uint8(0), lib.TypeFor[uint8]())
	lib.TestType(uint(0), lib.TypeFor[uint]())

	lib.TestType(0, lib.TypeFor[int]())
	lib.TestType(0.0, lib.TypeFor[float64]())
	lib.TestType(float32(0.0), lib.TypeFor[float32]())
	lib.TestType("", lib.TypeFor[string]())

	lib.TestType(lib.SetSomeValue(int64(0)), lib.TypeFor[lib.Some[int64]]())
	lib.TestType(lib.SetSomeValue(int32(0)), lib.TypeFor[lib.Some[int32]]())
	lib.TestType(lib.SetSomeValue(int16(0)), lib.TypeFor[lib.Some[int16]]())
	lib.TestType(lib.SetSomeValue(int8(0)), lib.TypeFor[lib.Some[int8]]())
	lib.TestType(lib.SetSomeValue(int(0)), lib.TypeFor[lib.Some[int]]())
	//
	lib.TestType(lib.SetSomeValue(uint64(0)), lib.TypeFor[lib.Some[uint64]]())
	lib.TestType(lib.SetSomeValue(uint32(0)), lib.TypeFor[lib.Some[uint32]]())
	lib.TestType(lib.SetSomeValue(uint16(0)), lib.TypeFor[lib.Some[uint16]]())
	lib.TestType(lib.SetSomeValue(uint8(0)), lib.TypeFor[lib.Some[uint8]]())
	lib.TestType(lib.SetSomeValue(uint(0)), lib.TypeFor[lib.Some[uint]]())
	//
	lib.TestType(lib.SetSomeValue(0), lib.TypeFor[lib.Some[int]]())
	lib.TestType(lib.SetSomeValue(0.0), lib.TypeFor[lib.Some[float64]]())
	lib.TestType(lib.SetSomeValue(float32(0.0)), lib.TypeFor[lib.Some[float32]]())
	lib.TestType(lib.SetSomeValue(""), lib.TypeFor[lib.Some[string]]())
}
