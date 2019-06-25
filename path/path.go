// Package path provide path based routing
package path

import (
	"net/http"
)

// H is type for mapping path and its handler
//
// parameterized segment of path can be prefixed with ":", example:
//	/info/:userid/name
//
// NOTE: handler registered here is compared by prefix path match,
// so request like /a/b/c/d/e will be still handled by /a/b.
type H map[string]http.HandlerFunc

func compile(h H, notfoundHandler http.HandlerFunc) http.HandlerFunc {
	b := &builderT{make(tree)}
	for k, v := range h {
		b.add(k, v)
	}
	return b.compile(notfoundHandler)
}

// Compile into single http.HandlerFunc
func (h H) Compile(def http.HandlerFunc) http.HandlerFunc {
	return compile(h, def)
}

// C same as Compile with def equal to nil
func (h H) C() http.HandlerFunc {
	return compile(h, nil)
}
