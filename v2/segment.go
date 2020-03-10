package router

import (
	"net/http"

	"github.com/payfazz/go-router/defhandler"
)

// HandlerMapping is type alias for mapping string to handler
type HandlerMapping = map[string]http.HandlerFunc

// Router is func to get routing state,
// this state will be used for routing decission based on next segment available
type Router func(*http.Request) *State

// BySegment generate a handler that take routing decission based on provided segment handler
//
// if segment handler is not found in handler, will generate http status 404
func (router Router) BySegment(handler HandlerMapping) http.HandlerFunc {
	return router.BySegmentWithDef(handler, defhandler.StatusNotFound)
}

// BySegmentWithDef is same with BySegment, but you can provide custom default handler
// instead of generating http status 404
func (router Router) BySegmentWithDef(handler HandlerMapping, def http.HandlerFunc) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		state := router(r)
		_, rest := state.progressCursor()

		var next http.HandlerFunc

		if rest == 0 {
			next = handler[""]
		} else {
			var ok bool
			next, ok = handler[state.next()]
			if !ok {
				state.prev()
			}
		}

		if next == nil {
			next = def
		}

		next(w, r)
	})
}

// ByParam generate handler that take next segment as parameter
//
// root will be called if the param is empty string and that param is the last segment,
// otherwise param will be called.
//
// if root is nil, then param will be used.
func (router Router) ByParam(setParam ParamSetter, root http.HandlerFunc, param http.HandlerFunc) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		state := router(r)
		_, rest := state.progressCursor()

		var next http.HandlerFunc

		if rest == 0 {
			next = root
		} else {
			segment := state.next()
			if segment == "" && rest == 1 { // treat trailing slash as end
				next = root
			} else {
				setParam(r, segment)
				next = param
			}
		}

		if next == nil {
			next = param
		}

		next(w, r)
	})
}

// SegmentMustEnd return middleware to make sure that the handler is the last segment
//
// if the segment is not last segment, will generate http status 404
func (router Router) SegmentMustEnd() func(http.HandlerFunc) http.HandlerFunc {
	return router.SegmentMustEndOr(defhandler.StatusNotFound)
}

// SegmentMustEndOr same with SegmentMustEnd, but you can provide custom default handler
// instead of generating http status 404
func (router Router) SegmentMustEndOr(def http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return (func(w http.ResponseWriter, r *http.Request) {
			state := router(r)
			_, rest := state.progressCursor()

			var next http.HandlerFunc

			if rest == 0 || (state.next() == "" && rest == 1) {
				next = handler
			} else {
				next = def
			}

			next(w, r)
		})
	}
}

// ParamSetter is callback for setting parameter
type ParamSetter func(r *http.Request, param string)

// SetParamIntoHeader return ParamSetter that can be used to set header
// based on parameter on segment
func SetParamIntoHeader(key string) ParamSetter {
	return func(r *http.Request, param string) {
		r.Header.Set(key, param)
	}
}
