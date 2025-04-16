package work

import (
	"github.com/bhbosman/go-sqlize/lib"
)

func init() {
	a := 12
	lib.TestType(a, lib.TypeFor[int]())

	b := lib.SetSomeValue(12)
	lib.TestType(b, lib.TypeFor[lib.Some[int]]())
}
