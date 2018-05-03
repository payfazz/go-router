// Package path provide path based routing
package path

import (
	"net/http"
)

// H is type for mapping path and its handler
//
// parameterized segment of path can be prefixed with ":", example:
//	/info/:userid/name
type H map[string]http.HandlerFunc

// Compile into single http.Handler. if notfoundHandler is nil, it will use defhandler.StatusNotFound
func Compile(h H, notfoundHandler http.HandlerFunc) http.HandlerFunc {
	b := &builderT{make(tree)}
	for k, v := range h {
		b.add(k, v)
	}
	return b.compile(notfoundHandler)
}
