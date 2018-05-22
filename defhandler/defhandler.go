// Package defhandler provide default handler construction
package defhandler

import (
	"fmt"
	"net/http"
)

var (
	StatusBadRequest           = genDefHandler(http.StatusBadRequest)
	StatusUnauthorized         = genDefHandler(http.StatusUnauthorized)
	StatusForbidden            = genDefHandler(http.StatusForbidden)
	StatusNotFound             = genDefHandler(http.StatusNotFound)
	StatusMethodNotAllowed     = genDefHandler(http.StatusMethodNotAllowed)
	StatusUnsupportedMediaType = genDefHandler(http.StatusUnsupportedMediaType)
	StatusUnprocessableEntity  = genDefHandler(http.StatusUnprocessableEntity)
)

func genDefHandler(code int) http.HandlerFunc {
	return ResponseCodeWithMessage(code, fmt.Sprintf("%d %s", code, http.StatusText(code)))
}

func Redirect(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// rfc2616, only "HEAD" and "GET"
		switch r.Method {
		case http.MethodHead, http.MethodGet:
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func ResponseCode(code int) http.HandlerFunc {
	return ResponseCodeWithMessage(code, "")
}

func ResponseCodeWithMessage(code int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		fmt.Fprint(w, message)
	}
}
