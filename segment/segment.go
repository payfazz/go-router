// Package segment provide segment based routing
package segment

import (
	"net/http"
	"net/url"
	"strings"

	internalsegment "github.com/payfazz/go-router/internal/segment"
)

// H is type for mapping segment and its handler
type H map[string]http.HandlerFunc

// Compile into single http.HandlerFunc. if def is nil, default handler is defhandler.StatusNotFound
func (h H) Compile(def http.HandlerFunc) http.HandlerFunc {
	hInternal := make(internalsegment.H)
	for k, v := range h {
		hInternal[k] = internalsegment.FromStd(v)
	}
	return internalsegment.IntoStd(
		hInternal.Compile(
			internalsegment.FromStd(def),
		),
	)
}

// C same as Compile with def equal to nil
func (h H) C() http.HandlerFunc {
	return h.Compile(nil)
}

// Tag return helper that will tag current segment and process to next segment
//
// The tagged segment can be retrieved later via Get function.
func Tag(tag string, next http.HandlerFunc) http.HandlerFunc {
	return internalsegment.IntoStd(
		internalsegment.Tag(
			tag,
			internalsegment.FromStd(next),
		),
	)
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
		s, r2 := internalsegment.TryShifterFrom(r)
		_, rest := s.Split()

		newURL := new(url.URL)
		*newURL = *r2.URL
		newURL.Path = "/" + strings.Join(rest, "/")
		newURL.RawPath = newURL.Path
		if newURL.EscapedPath() != newURL.RawPath {
			newURL.RawPath = ""
		}

		// temporary set r2.URL to nil, so it will not be cloned
		oldURL := r2.URL
		r2.URL = nil

		// create shifter based on newURL on r2, this will clone r2
		_, r3 := internalsegment.NewShifterFor(r2, strings.Split(
			strings.TrimPrefix(newURL.EscapedPath(), "/"), "/",
		))
		r3.URL = newURL

		// change r2.URL back new
		r2.URL = oldURL

		next(w, r3)
	}
}

// Get return tagged segment
func Get(r *http.Request, tag string) (string, bool) {
	s, _ := internalsegment.TryShifterFrom(r)
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
	s, _ := internalsegment.TryShifterFrom(r)
	return s.Split()
}
