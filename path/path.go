// Package path provide path based routing
package path

import (
	"net/http"
)

// H is type for mapping path and its handler
//
// parameterized segment of path can be prefixed with ":", example:
//
//	/info/:userid/name
//
// handler registered here is compared by prefix path match,
//
//	h := H{
//		"a/b": myhandler,
//	}.C()
//
// request to "/a/b/c/d/e" will be still handled by myhandler, see segment.MustEnd
type H map[string]http.HandlerFunc

// Compile into single http.HandlerFunc
func (h H) Compile(notfoundHandler http.HandlerFunc) http.HandlerFunc {
	b := make(tree)
	for k, v := range h {
		b.add(k, v)
	}
	return b.compile(notfoundHandler)
}

// C same as Compile with def equal to nil
func (h H) C() http.HandlerFunc {
	return h.Compile(nil)
}
