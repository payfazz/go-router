package router

import (
	"net/http"

	"github.com/payfazz/go-router/defhandler"
)

// ShifterGetter is callback for getting shifter from the given request
type ShifterGetter func(r *http.Request) *Shifter

// BySegment is same with BySegmentWithDef(handler, nil)
func (sg ShifterGetter) BySegment(handler map[string]http.HandlerFunc) http.HandlerFunc {
	return sg.BySegmentWithDef(handler, nil)
}

// BySegmentWithDef will return handler for routing via segment,
//
// if then handler for spesific segment is not found it will return def,
//
// if def is nil defhandler.StatusNotFound is used
func (sg ShifterGetter) BySegmentWithDef(handler map[string]http.HandlerFunc, def http.HandlerFunc) http.HandlerFunc {

	if def == nil {
		def = defhandler.StatusNotFound
	}

	return func(w http.ResponseWriter, r *http.Request) {
		shifter := sg(r)
		alreadyEnd := shifter.end()

		var next http.HandlerFunc
		if tmp, ok := handler[shifter.next()]; ok {
			next = tmp
		} else {
			next = def
			if !alreadyEnd {
				shifter.prev()
			}
		}

		next(w, r)
	}
}

// ByParam return handler for handling parameter in segment,
//
// empty will be called if the param is empty, otherwise handler will be called.
//
// if empty is nil, handler will be used.
func (sg ShifterGetter) ByParam(setParam ParamSetter, empty http.HandlerFunc, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shifter := sg(r)
		_, rest := shifter.state()

		param := shifter.next()
		next := empty

		if param != "" || rest > 1 {
			next = handler
		}

		if next == nil {
			next = handler
		}
		setParam(r, param)
		next(w, r)
	}
}

// SegmentMustEnd return middleware to make sure that the handler is the last segment
//
// if the segment is not the end, defhandler.StatusNotFound is used
func (sg ShifterGetter) SegmentMustEnd() func(http.HandlerFunc) http.HandlerFunc {
	return sg.SegmentMustEndOr(nil)
}

// SegmentMustEndOr same with SegmentMustEnd but using def as handler when the segment is not end,
//
// if def is nil, defhandler.StatusNotFound is used
func (sg ShifterGetter) SegmentMustEndOr(def http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	if def == nil {
		def = defhandler.StatusNotFound
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			shifter := sg(r)
			_, restN := shifter.state()
			if restN == 0 || (restN == 1 && shifter.next() == "") {
				next(w, r)
				return
			}

			def(w, r)
		}
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
