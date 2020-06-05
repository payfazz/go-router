package segment

import (
	"context"
	"net/http"
	"strings"

	"github.com/payfazz/go-router/segment/shifter"
)

type ctxKeyType struct{}

var ctxKey ctxKeyType

// TryShifterFrom .
func TryShifterFrom(r *http.Request) (*shifter.Shifter, *http.Request) {
	tmp := r.Context().Value(ctxKey)
	if tmp != nil {
		return tmp.(*shifter.Shifter), r
	}
	return NewShifterFor(r, strings.Split(
		strings.TrimPrefix(r.URL.EscapedPath(), "/"), "/",
	))
}

// NewShifterFor .
//
// will clone r
func NewShifterFor(r *http.Request, list []string) (*shifter.Shifter, *http.Request) {
	s := shifter.New(list)

	// WithContext will clone r
	r2 := r.WithContext(context.WithValue(
		r.Context(), ctxKey, s,
	))

	return s, r2
}
