// Package segment provide segment based routing
package segment

import (
	"net/http"
	"strings"

	"github.com/payfazz/go-router/defhandler"
	"github.com/payfazz/go-router/segment/shifter"
)

type ctxType struct{}

var ctxKey ctxType

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
		s, r := shifter.With(r, ctxKey, nil)
		seg, _ := s.Shift()
		next, ok := h[seg]
		if !ok {
			next = def
			s.Unshift()
		}
		next(w, r)
	}
}

// Tag return helper that will tag current segment and process to next segment
func Tag(tag string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, r := shifter.With(r, ctxKey, nil)
		if !s.End() {
			s.Shift()
			s.Tag(tag)
		}
		next(w, r)
	}
}

// End return helper that only will execute h when its position is last segment of the path
func End(h http.HandlerFunc, def http.HandlerFunc) http.HandlerFunc {
	return Compile(H{
		"": h,
	}, def)
}

// Stripper return helper for stripping processed segment from r.URL.Path
func Stripper(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, r := shifter.With(r, ctxKey, nil)
		_, rest := s.Split()
		r.URL.Path = "/" + strings.Join(rest, "/")
		r.URL.RawPath = ""
		next(w, r)
	}
}

// Get return tagged segment
func Get(r *http.Request, tag string) (string, bool) {
	s, _ := shifter.With(r, ctxKey, nil)
	return s.GetByTag(tag)
}

// Rest return rest of the segment
func Rest(r *http.Request) []string {
	s, _ := shifter.With(r, ctxKey, nil)
	_, rest := s.Split()
	return rest
}
