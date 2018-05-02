// Package segment provide segment based routing
package segment

import (
	"context"
	"net/http"
	"strings"

	"github.com/payfazz/go-router/defhandler"
)

// H is type for mapping segment and its handler
type H map[string]http.HandlerFunc

// Compile into single http.Handler. if def is nil, it will use defhandler.StatusNotFound
func Compile(h H, def http.HandlerFunc) http.HandlerFunc {
	if h == nil {
		h = make(H)
	}
	if def == nil {
		def = defhandler.StatusNotFound
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var next http.HandlerFunc
		s, r := getState(r)
		if seg, ok := s.shift(); ok {
			if next, ok = h[seg]; !ok {
				next = def
			}
		} else {
			next = def
		}
		next(w, r)
	}
}

// Len return number of segment in path
func Len(r *http.Request) int {
	s := r.Context().Value(key).(*state)
	return len(s.list)
}

// Cur return current index of segment in path
func Cur(r *http.Request) int {
	s := r.Context().Value(key).(*state)
	return s.next - 1
}

// Get return current segment in path
func Get(r *http.Request) string {
	return GetN(r, 0)
}

// GetN return segment in path, relative to current index
func GetN(r *http.Request, d int) string {
	s := r.Context().Value(key).(*state)
	return s.get(s.next - 1 + d)
}

// End indicate end segment of path
func End(r *http.Request) bool {
	s := r.Context().Value(key).(*state)
	return s.next == len(s.list)
}

// Rest segment in the path
func Rest(r *http.Request) []string {
	s := r.Context().Value(key).(*state)
	ret := make([]string, len(s.list)-s.next)
	copy(ret, s.list[s.next:])
	return ret
}

func getState(r *http.Request) (*state, *http.Request) {
	if tmp := r.Context().Value(key); tmp != nil {
		return tmp.(*state), r
	}

	s := &state{
		strings.Split(
			strings.TrimPrefix(r.URL.EscapedPath(), "/"),
			"/",
		),
		0,
	}

	r = r.WithContext(context.WithValue(
		r.Context(), key, s,
	))

	return s, r
}
