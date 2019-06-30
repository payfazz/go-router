package segment

import (
	"context"
	"net/http"
	"strings"

	"github.com/payfazz/go-router/segment/shifter"
)

type ctxKeyType struct{}

var ctxKey ctxKeyType

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
		s, r2 := TryShifterFrom(r)
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

// TryShifterFrom .
func TryShifterFrom(r *http.Request) (*shifter.Shifter, *http.Request) {
	tmp := r.Context().Value(ctxKey)
	if tmp != nil {
		return tmp.(*shifter.Shifter), r
	}
	return NewShifterFor(r)
}

// NewShifterFor .
func NewShifterFor(r *http.Request) (*shifter.Shifter, *http.Request) {
	list := strings.Split(
		strings.TrimPrefix(r.URL.EscapedPath(), "/"), "/",
	)

	s := shifter.New(list)
	r = r.WithContext(context.WithValue(
		r.Context(), ctxKey, s,
	))
	return s, r
}
