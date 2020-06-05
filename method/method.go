// Package method provide method based routing
package method

import (
	"net/http"
	"strings"

	"github.com/payfazz/go-router/defhandler"
)

// H is type for mapping method and its handler
type H map[string]http.HandlerFunc

// Compile into single http.HandlerFunc. if def is nil, default handler is defhandler.StatusMethodNotAllowed
func (h H) Compile(def http.HandlerFunc) http.HandlerFunc {
	if def == nil {
		def = defhandler.StatusMethodNotAllowed
	}
	realH := make(H)
	for k, v := range h {
		if v == nil {
			v = defhandler.StatusNotImplemented
		}
		if _, ok := realH[strings.ToUpper(k)]; ok {
			panic("method: duplicate handler for " + strings.ToUpper(k))
		}
		realH[strings.ToUpper(k)] = v
	}
	return func(w http.ResponseWriter, r *http.Request) {
		next, ok := realH[strings.ToUpper(r.Method)]
		if !ok {
			next = def
		}
		next(w, r)
	}
}

// C same as Compile with def equal to nil
func (h H) C() http.HandlerFunc {
	return h.Compile(nil)
}

// Must return middleware that only allowed method that specified in methods
func Must(methods ...string) func(http.HandlerFunc) http.HandlerFunc {
	// little optimization for single value
	if len(methods) == 1 {
		expected := strings.ToUpper(methods[0])
		return func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				got := strings.ToUpper(r.Method)
				if got != expected {
					defhandler.StatusMethodNotAllowed(w, r)
					return
				}

				next(w, r)
			}
		}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		h := H{}
		for _, m := range methods {
			h[m] = next
		}
		return h.C()
	}
}
