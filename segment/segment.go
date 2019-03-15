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

// Compile into single http.Handler. if def is nil, it will use defhandler.StatusNotFound
func Compile(h H, def http.HandlerFunc) http.HandlerFunc {
	if h == nil {
		h = make(H)
	}
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
		seg, _ := s.Shift()
		next, ok := h[seg]
		if !ok {
			next = def
			s.Unshift()
		}
		next(w, r)
	}
}

// C same as Compile with def equal to nil
func C(h H) http.HandlerFunc {
	return Compile(h, nil)
}

// Compile into single http.HandlerFunc
func (h H) Compile(def http.HandlerFunc) http.HandlerFunc {
	return Compile(h, def)
}

// C same as Compile with def equal to nil
func (h H) C() http.HandlerFunc {
	return C(h)
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
// if otherwise is nil, defhandler.StatusNotFound is used
func End(h http.HandlerFunc, otherwise http.HandlerFunc) http.HandlerFunc {
	return EndOr(otherwise)(h)
}

// EndOr same as End, but return middleware
func EndOr(otherwise http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return Compile(H{
			"": h,
		}, otherwise)
	}
}

// E is *DEPREDECATED*, use MustEnd
func E(h http.HandlerFunc) http.HandlerFunc {
	return End(h, nil)
}

// MustEnd same as End with otherwise equal to nil.
func MustEnd(h http.HandlerFunc) http.HandlerFunc {
	return End(h, nil)
}

// Stripper return helper for stripping processed segment from r.URL.Path
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

		s, r2 = shifter.Reset(r2, ctxKey, nil)
		next(w, r2)
	}
}

// Get return tagged segment
func Get(r *http.Request, tag string) (string, bool) {
	s, _ := shifter.With(r, ctxKey, nil)
	return s.GetByTag(tag)
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
