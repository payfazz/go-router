// Package method provide method based routing
package method

import (
	"net/http"

	"github.com/payfazz/go-router/defhandler"
)

// AllowedMethod in Compile.
var AllowedMethod = []string{
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

// H is type for mapping method and its handler
type H map[string]http.HandlerFunc

// Compile into single http.HandlerFunc. If def is nil, it will use defhandler.StatusMethodNotAllowed
func Compile(h H, def http.HandlerFunc) http.HandlerFunc {
	if h == nil {
		h = make(H)
	}
	if def == nil {
		def = defhandler.StatusMethodNotAllowed
	}
	for k, v := range h {
		if !inArr(k, AllowedMethod) {
			panic("method: method \"" + k + "\" is not allowed.")
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

func inArr(v string, xs []string) bool {
	for _, x := range xs {
		if v == x {
			return true
		}
	}
	return false
}
