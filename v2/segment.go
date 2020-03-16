package router

import (
	"net/http"

	"github.com/payfazz/go-router/defhandler"
)

// BySegment generate a handler that take routing decission based on provided hmap
//
// if segment handler is not found in hmap, will generate http status 404
func (router Router) BySegment(hmap Hmap) Handler {
	return router.BySegmentWithDef(hmap, defhandler.StatusNotFound)
}

// BySegmentWithDef is same with BySegment, but you can provide custom default handler
// instead of generating http status 404
func (router Router) BySegmentWithDef(hmap Hmap, def Handler) Handler {
	return (func(w http.ResponseWriter, r *http.Request) {
		state := router(r)
		_, rest := state.Progress()

		var next Handler

		if rest == "" {
			next = hmap[""]
		} else {
			var ok bool
			next, ok = hmap[state.next()]
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

// ByParam generate handler that take segment as parameter
//
// root will be called if the param is empty string and that param is the last segment,
// otherwise param will be called.
//
// if root is nil, then param will be used.
func (router Router) ByParam(setParam ParamSetter, root Handler, param Handler) Handler {
	return (func(w http.ResponseWriter, r *http.Request) {
		state := router(r)
		_, rest := state.Progress()

		var next Handler

		if rest == "" {
			next = root
		} else {
			segment := state.next()
			if segment == "" && rest == "/" {
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
func (router Router) SegmentMustEnd() func(Handler) Handler {
	return router.SegmentMustEndOr(defhandler.StatusNotFound)
}

// SegmentMustEndOr same with SegmentMustEnd, but you can provide custom default handler
// instead of generating http status 404
func (router Router) SegmentMustEndOr(def Handler) func(Handler) Handler {
	return func(handler Handler) Handler {
		return (func(w http.ResponseWriter, r *http.Request) {
			state := router(r)
			_, rest := state.Progress()

			var next Handler

			if rest == "" || (state.next() == "" && rest == "/") {
				next = handler
				if rest != "" {
					state.prev()
				}
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
