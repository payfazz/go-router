// Package host provide host based routing
package host

import (
	"net/http"

	"github.com/payfazz/go-router/defhandler"
)

// H is type for mapping host and its handler
type H map[string]http.HandlerFunc

// Compile into single http.HandlerFunc
func Compile(h H, def http.HandlerFunc) http.HandlerFunc {
	if h == nil {
		h = make(H)
	}
	if def == nil {
		def = defhandler.ResponseCodeWithMessage(http.StatusBadRequest, "Host not found")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		next, ok := h[r.Host]
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
