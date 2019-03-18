// Package segment provide segment based routing
package segment

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/payfazz/go-router/defhandler"
	"github.com/payfazz/go-router/segment/shifter"
)

type ctxType struct{}

var ctxKey ctxType

// H is type for mapping segment and its handler
type H map[string]http.HandlerFunc

func compile(h H, def http.HandlerFunc) http.HandlerFunc {
	if def == nil {
		def = defhandler.StatusNotFound
	}
	for k, v := range h {
		if v == nil {
			h[k] = defhandler.StatusNotImplemented
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var next http.HandlerFunc
		s, r := shifter.With(r, ctxKey, nil)
		end := s.End()
		seg, _ := s.Shift()
		next, ok := h[seg]
		if !ok {
			next = def
			if !end {
				s.Unshift()
			}
		}
		next(w, r)
	}
}

// Compile into single http.HandlerFunc. if def is nil, default handler is defhandler.StatusNotFound
func (h H) Compile(def http.HandlerFunc) http.HandlerFunc {
	return compile(h, def)
}

// C same as Compile with def equal to nil
func (h H) C() http.HandlerFunc {
	return compile(h, nil)
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

// EndOr return middleware that only will execute h when its position is last segment of the path
// if otherwise is nil, defhandler.StatusNotFound is used.
func EndOr(otherwise http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return H{
			"": next,
		}.Compile(otherwise)
	}
}

// MustEnd same as EndOr with otherwise equal to nil.
func MustEnd(h http.HandlerFunc) http.HandlerFunc {
	return EndOr(nil)(h)
}

// Stripper is middleware for stripping processed segment from r.URL.Path
func Stripper(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, r := shifter.With(r, ctxKey, nil)
		_, rest := s.Split()

		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = "/" + strings.Join(rest, "/")
		r2.URL.RawPath = r2.URL.Path
		if r2.URL.EscapedPath() != r2.URL.RawPath {
			r2.URL.RawPath = ""
		}

		s, r2 = shifter.Reset(r2, ctxKey, nil)
		next(w, r2)
	}
}

// Get return tagged segment
func Get(r *http.Request, tag string) (string, bool) {
	s, _ := shifter.With(r, ctxKey, nil)
	return s.GetByTag(tag)
}

// Param do the same thing as Get, but panic when tag is not found in the segment
func Param(r *http.Request, tag string) string {
	s, ok := Get(r, tag)
	if !ok {
		panic("segment: param " + tag + " not found in the segment")
	}
	return s
}

// Rest return rest of the segment
func Rest(r *http.Request) []string {
	_, rest := Split(r)
	return rest
}

// Split return processed and the rest of the segments
func Split(r *http.Request) ([]string, []string) {
	s, _ := shifter.With(r, ctxKey, nil)
	return s.Split()
}

// UnshiftInternalShifter by num of segment
func UnshiftInternalShifter(r *http.Request, num int) {
	s, _ := shifter.With(r, ctxKey, nil)
	for i := 0; i < num; i++ {
		s.Unshift()
	}
}
