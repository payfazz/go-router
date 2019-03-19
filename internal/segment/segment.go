package segment

import (
	"net/http"

	"github.com/payfazz/go-router/segment/shifter"
)

// CtxType .
type CtxType struct{}

// CtxKey .
var CtxKey CtxType

// SetShifterIndex .
func SetShifterIndex(r *http.Request, index int) {
	s, _ := shifter.With(r, CtxKey, nil)
	num := s.CurrentIndex() + 1 - index
	if num >= 0 {
		for i := 0; i < num; i++ {
			s.Unshift()
		}
	} else if num < 0 {
		num = -num
		for i := 0; i < num; i++ {
			s.Shift()
		}
	}
}
