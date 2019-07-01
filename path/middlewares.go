package path

import (
	"net/http"
	"net/url"
	"strings"
)

// WithTrailingSlash is middleware for redirect request to url that with trailing slash
func WithTrailingSlash(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.String(), "/") {
			// rfc2616, only "HEAD" and "GET"
			switch r.Method {
			case http.MethodHead, http.MethodGet:
				newURL := new(url.URL)
				*newURL = *r.URL

				hint := newURL.EscapedPath() == newURL.RawPath
				newURL.Path += "/"
				if hint {
					newURL.RawPath += "/"
				}

				http.Redirect(w, r, newURL.String(), http.StatusMovedPermanently)
				return
			}
		}

		next(w, r)
	}
}

// WithoutTrailingSlash is middleware for redirect request to url that without trailing slash
func WithoutTrailingSlash(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.String()
		if path != "/" && strings.HasSuffix(path, "/") {
			// rfc2616, only "HEAD" and "GET"
			switch r.Method {
			case http.MethodHead, http.MethodGet:
				newURL := new(url.URL)
				*newURL = *r.URL

				hint := newURL.EscapedPath() == newURL.RawPath
				newURL.Path = newURL.Path[:len(newURL.Path)-1]
				if hint {
					newURL.RawPath = newURL.RawPath[:len(newURL.RawPath)-1]
				}

				http.Redirect(w, r, newURL.String(), http.StatusMovedPermanently)
				return
			}
		}

		next(w, r)
	}
}
