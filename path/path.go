// Package path provide path based routing
package path

import (
	"net/http"
	"strings"
)

// H is type for mapping path and its handler
//
// parameterized segment of path can be prefixed with ":", example:
//	/info/:userid/name
//
// NOTE: handler registered here is compared by prefix path match,
// so request like /a/b/c/d/e will be still handled by /a/b.
type H map[string]http.HandlerFunc

// Compile into single http.Handler. if notfoundHandler is nil, it will use defhandler.StatusNotFound
func Compile(h H, notfoundHandler http.HandlerFunc) http.HandlerFunc {
	if h == nil {
		h = make(H)
	}
	b := &builderT{make(tree)}
	for k, v := range h {
		b.add(k, v)
	}
	return b.compile(notfoundHandler)
}

// C same as Compile with notfoundHandler equal to nil
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

// WithTrailingSlash return helper for redirect request to url that with trailing slash
func WithTrailingSlash(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()
		if !strings.HasSuffix(path, "/") {
			// rfc2616, only "HEAD" and "GET"
			switch r.Method {
			case http.MethodHead, http.MethodGet:
				http.Redirect(w, r, path+"/", http.StatusMovedPermanently)
				return
			}
		}

		next(w, r)
	}
}

// WithoutTrailingSlash return helper for redirect request to url that without trailing slash
func WithoutTrailingSlash(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()
		if path != "/" && strings.HasSuffix(path, "/") {
			// rfc2616, only "HEAD" and "GET"
			switch r.Method {
			case http.MethodHead, http.MethodGet:
				http.Redirect(w, r, path[:len(path)-1], http.StatusMovedPermanently)
				return
			}
		}

		next(w, r)
	}
}
