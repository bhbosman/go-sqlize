package work

import "github.com/bhbosman/go-sqlize/lib"

func init() {
	type InputData struct {
	}
	_ = lib.Query[InputData](lib.QueryTop(10), lib.QueryDistinct())
}
