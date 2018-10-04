// Package defhandler provide default handler construction
package defhandler

import (
	"fmt"
	"net/http"
)

var (
	// StatusBadRequest is http.HandlerFunc that always send HTTP status BadRequest
	StatusBadRequest = genDefHandler(http.StatusBadRequest)

	// StatusUnauthorized is http.HandlerFunc that always send HTTP status Unauthorized
	StatusUnauthorized = genDefHandler(http.StatusUnauthorized)

	// StatusForbidden is http.HandlerFunc that always send HTTP status Forbidden
	StatusForbidden = genDefHandler(http.StatusForbidden)

	// StatusNotFound is http.HandlerFunc that always send HTTP status NotFound
	StatusNotFound = genDefHandler(http.StatusNotFound)

	// StatusMethodNotAllowed is http.HandlerFunc that always send HTTP status MethodNotAllowed
	StatusMethodNotAllowed = genDefHandler(http.StatusMethodNotAllowed)

	// StatusUnsupportedMediaType is http.HandlerFunc that always send HTTP status UnsupportedMediaType
	StatusUnsupportedMediaType = genDefHandler(http.StatusUnsupportedMediaType)

	// StatusUnprocessableEntity is http.HandlerFunc that always send HTTP status UnprocessableEntity
	StatusUnprocessableEntity = genDefHandler(http.StatusUnprocessableEntity)
)

func genDefHandler(code int) http.HandlerFunc {
	return ResponseCodeWithMessage(code, fmt.Sprintf("%d %s", code, http.StatusText(code)))
}

// Redirect return http.HandlerFunc that always redirect to url
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

// ResponseCode return http.HandlerFunc that always send empty HTTP response with
// defined status code
func ResponseCode(code int) http.HandlerFunc {
	return ResponseCodeWithMessage(code, "")
}

// ResponseCodeWithMessage return http.HandlerFunc that always send HTTP response
// with defined status code and text/plain message
func ResponseCodeWithMessage(code int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if message != "" {
			w.Header().Set("Content-Type", "text/plain")
		}
		w.WriteHeader(code)
		if message != "" {
			fmt.Fprint(w, message)
		}
	}
}
