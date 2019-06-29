package segment

import (
	"net/http"

	"github.com/payfazz/go-router/segment/shifter"
)

type ctxKeyType struct{}

// CtxKey .
var CtxKey ctxKeyType

// HandlerFunc .
//
// we need this for improve efficiency of passing Shifter instance
type HandlerFunc func(s *shifter.Shifter, w http.ResponseWriter, r *http.Request)

// IntoStd .
func IntoStd(handler HandlerFunc) http.HandlerFunc {
	if handler == nil {
		return nil
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s, r2 := shifter.From(r, CtxKey)
		handler(s, w, r2)
	}
}

// FromStd .
func FromStd(handler http.HandlerFunc) HandlerFunc {
	if handler == nil {
		return nil
	}
	return func(s *shifter.Shifter, w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}
