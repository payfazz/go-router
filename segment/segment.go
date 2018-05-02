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
		s, r := shifter.From(r, ctxKey)
		seg, _ := s.Shift()
		next, ok := h[seg]
		if !ok {
			next = def
			s.Unshift()
		}
		next(w, r)
	}
}

func Tag(tag string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, r := shifter.From(r, ctxKey)
		if !s.End() {
			s.Shift()
			s.Tag(tag)
		}
		next(w, r)
	}
}

func Stripper(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, r := shifter.From(r, ctxKey)
		_, rest := s.Split()
		r.URL.Path = "/" + strings.Join(rest, "/")
		r.URL.RawPath = ""
		next(w, r)
	}
}

func Get(r *http.Request, tag string) (string, bool) {
	s, _ := shifter.From(r, ctxKey)
	return s.GetByTag(tag)
}

func Current(r *http.Request) string {
	s, _ := shifter.From(r, ctxKey)
	return s.GetRelative(0)
}

func End(r *http.Request) bool {
	s, _ := shifter.From(r, ctxKey)
	return s.End()
}

func Rest(r *http.Request) []string {
	s, _ := shifter.From(r, ctxKey)
	_, rest := s.Split()
	return rest
}
