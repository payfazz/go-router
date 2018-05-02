// Package defhandler provide default handler construction
package defhandler

import "net/http"

var (
	StatusBadRequest                 = ResponseCode(http.StatusBadRequest)
	StatusUnauthorized               = ResponseCode(http.StatusUnauthorized)
	StatusForbidden                  = ResponseCode(http.StatusForbidden)
	StatusNotFound                   = ResponseCode(http.StatusNotFound)
	StatusMethodNotAllowed           = ResponseCode(http.StatusMethodNotAllowed)
	StatusNotAcceptable              = ResponseCode(http.StatusNotAcceptable)
	StatusConflict                   = ResponseCode(http.StatusConflict)
	StatusGone                       = ResponseCode(http.StatusGone)
	StatusUnsupportedMediaType       = ResponseCode(http.StatusUnsupportedMediaType)
	StatusUnprocessableEntity        = ResponseCode(http.StatusUnprocessableEntity)
	StatusLocked                     = ResponseCode(http.StatusLocked)
	StatusUnavailableForLegalReasons = ResponseCode(http.StatusUnavailableForLegalReasons)
)

func Error(error string, code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, error, code)
	}
}

func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}
}

func Redirect(url string, code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, code)
	}
}

func ResponseCode(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	}
}

func FileServer(fs http.FileSystem) http.HandlerFunc {
	return http.FileServer(fs).ServeHTTP
}

func FileServerRelative(fs http.FileSystem) http.HandlerFunc {
	panic("defhandler: implementation is not complete yet")
}

func FileServerDir(dir string) http.HandlerFunc {
	return FileServer(http.Dir(dir))
}

func FileServerDirRelative(dir string) http.HandlerFunc {
	return FileServerRelative(http.Dir(dir))
}
