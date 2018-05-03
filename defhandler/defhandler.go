// Package defhandler provide default handler construction
package defhandler

import "net/http"

var (
	StatusBadRequest           = ResponseCode(http.StatusBadRequest)
	StatusUnauthorized         = ResponseCode(http.StatusUnauthorized)
	StatusForbidden            = ResponseCode(http.StatusForbidden)
	StatusNotFound             = ResponseCode(http.StatusNotFound)
	StatusMethodNotAllowed     = ResponseCode(http.StatusMethodNotAllowed)
	StatusUnsupportedMediaType = ResponseCode(http.StatusUnsupportedMediaType)
	StatusUnprocessableEntity  = ResponseCode(http.StatusUnprocessableEntity)
)

func Error(err string, code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err, code)
	}
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
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	}
}
