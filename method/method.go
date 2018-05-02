// Package method provide method based routing
package method

import (
	"net/http"

	"github.com/payfazz/go-router/defhandler"
)

// H is type for mapping method and its handler
type H map[string]http.HandlerFunc

// Compile into single http.Handler. if def is nil, it will use defhandler.StatusMethodNotAllowed
func Compile(h H, def http.HandlerFunc) http.HandlerFunc {
	if h == nil {
		h = make(H)
	}
	if def == nil {
		def = defhandler.StatusMethodNotAllowed
	}
	return func(w http.ResponseWriter, r *http.Request) {
		next, ok := h[r.Method]
		if !ok {
			next = def
		}
		next(w, r)
	}
}
