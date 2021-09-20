// Package defhandler provide default handler construction
package defhandler

import (
	"fmt"
	"net/http"
)

var (
	// StatusBadRequest is http.HandlerFunc that always send HTTP status BadRequest
	StatusBadRequest = ResponseCode(http.StatusBadRequest)

	// StatusUnauthorized is http.HandlerFunc that always send HTTP status Unauthorized
	StatusUnauthorized = ResponseCode(http.StatusUnauthorized)

	// StatusForbidden is http.HandlerFunc that always send HTTP status Forbidden
	StatusForbidden = ResponseCode(http.StatusForbidden)

	// StatusNotFound is http.HandlerFunc that always send HTTP status NotFound
	StatusNotFound = ResponseCode(http.StatusNotFound)

	// StatusMethodNotAllowed is http.HandlerFunc that always send HTTP status MethodNotAllowed
	StatusMethodNotAllowed = ResponseCode(http.StatusMethodNotAllowed)

	// StatusUnsupportedMediaType is http.HandlerFunc that always send HTTP status UnsupportedMediaType
	StatusUnsupportedMediaType = ResponseCode(http.StatusUnsupportedMediaType)

	// StatusUnprocessableEntity is http.HandlerFunc that always send HTTP status UnprocessableEntity
	StatusUnprocessableEntity = ResponseCode(http.StatusUnprocessableEntity)

	// StatusNotImplemented is http.HandlerFunc that always send HTTP status NotImplemented
	StatusNotImplemented = ResponseCode(http.StatusNotImplemented)

	// StatusInternalServerError is http.HandlerFunc that always send HTTP status InternalServerError
	StatusInternalServerError = ResponseCode(http.StatusInternalServerError)
)

// Redirect return http.HandlerFunc that always redirect to url
func Redirect(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// https://httpwg.org/specs/rfc7231.html#status.301, only "HEAD" and "GET"
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
	return ResponseCodeWithMessage(code, fmt.Sprintf("%d %s", code, http.StatusText(code)))
}

// ResponseCodeWithMessage return http.HandlerFunc that always send HTTP response
// with defined status code and text/plain message
func ResponseCodeWithMessage(code int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if message != "" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}
		w.WriteHeader(code)
		if message != "" {
			fmt.Fprint(w, message)
		}
	}
}
