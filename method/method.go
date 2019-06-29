// Package method provide method based routing
package method

import (
	"net/http"

	"github.com/payfazz/go-router/defhandler"
)

var allowedMethod = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

func inArr(v string, xs []string) bool {
	for _, x := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// H is type for mapping method and its handler
type H map[string]http.HandlerFunc

// Compile into single http.HandlerFunc. if def is nil, default handler is defhandler.StatusMethodNotAllowed
func (h H) Compile(def http.HandlerFunc) http.HandlerFunc {
	if def == nil {
		def = defhandler.StatusMethodNotAllowed
	}
	for k, v := range h {
		if !inArr(k, allowedMethod) {
			panic("method: method '" + k + "' is not allowed.")
		}
		if v == nil {
			h[k] = defhandler.StatusNotImplemented
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		next, ok := h[r.Method]
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
	return func(next http.HandlerFunc) http.HandlerFunc {
		h := H{}
		for _, m := range methods {
			h[m] = next
		}
		return h.C()
	}
}
