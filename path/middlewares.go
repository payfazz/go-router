package path

import (
	"net/http"
	"strings"
)

// WithTrailingSlash is middleware for redirect request to url that with trailing slash
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

// WithoutTrailingSlash is middleware for redirect request to url that without trailing slash
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
