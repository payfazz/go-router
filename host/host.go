// Package host provide host based routing
package host

import (
	"net/http"

	"github.com/payfazz/go-router/defhandler"
)

// TODO(win): support wildcard host (example: *.payfazz.com)

// H is type for mapping host and its handler
type H map[string]http.HandlerFunc

// Compile into single http.HandlerFunc. if def is nil, default handler is defhandler.StatusBadRequest
func (h H) Compile(def http.HandlerFunc) http.HandlerFunc {
	if def == nil {
		def = defhandler.StatusBadRequest
	}
	for k, v := range h {
		if v == nil {
			h[k] = defhandler.StatusNotImplemented
		}
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
func (h H) C() http.HandlerFunc {
	return h.Compile(nil)
}

// Must return middleware that only allowed host that specified in hosts
func Must(hosts ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		h := H{}
		for _, m := range hosts {
			h[m] = next
		}
		return h.C()
	}
}
