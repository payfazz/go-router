package router

import (
	"net/http"

	"github.com/payfazz/go-router/defhandler"
)

// ShifterGetter is callback for getting shifter from context
type ShifterGetter func(r *http.Request) *SegmentShifter

// ParameterSetter is callback for setting parameter
type ParameterSetter func(r *http.Request, param string)

// BySegment will return handler for routing via segment
//
// if getShifter is nil, it will default shifter from ShifterInjector
//
// will panic if getShifter is nil but you not use ShifterInjector before it
func BySegment(getShifter ShifterGetter, def http.HandlerFunc, next map[string]http.HandlerFunc) http.HandlerFunc {
	if getShifter == nil {
		getShifter = defaultShifter
	}
	if def == nil {
		def = defhandler.StatusNotFound
	}

	return func(w http.ResponseWriter, r *http.Request) {
		shifter := getShifter(r)
		alreadyEnd := shifter.End()

		var handler http.HandlerFunc = def

		if nextHandler, ok := next[shifter.Next()]; ok {
			handler = nextHandler
		} else {
			handler = def
			if !alreadyEnd {
				shifter.Prev()
			}
		}

		handler(w, r)
	}
}

// ByParam return composite handler
//
// it will shift current segment and and call setParam with it
//
// empty will be called if current segment is empty, otherwise handler will be called
//
// if getShifter is nil, it will default shifter from ShifterInjector
//
// will panic if getShifter is nil but you not use ShifterInjector before it
func ByParam(getShifter ShifterGetter, setParam ParameterSetter, empty http.HandlerFunc, handler http.HandlerFunc) http.HandlerFunc {
	if getShifter == nil {
		getShifter = defaultShifter
	}

	return func(w http.ResponseWriter, r *http.Request) {
		param := getShifter(r).Next()
		setParam(r, param)
		if param == "" && empty != nil {
			empty(w, r)
		} else {
			handler(w, r)
		}
	}
}

// SegmentMustEnd return middleware to make sure that the handler is the last segment
//
// if getShifter is nil, it will default shifter from ShifterInjector
//
// will panic if getShifter is nil but you not use ShifterInjector before it
func SegmentMustEnd(getShifter ShifterGetter) func(http.HandlerFunc) http.HandlerFunc {
	if getShifter == nil {
		getShifter = defaultShifter
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if getShifter(r).Next() != "" {
				defhandler.StatusNotFound(w, r)
				return
			}
			next(w, r)
		}
	}
}
