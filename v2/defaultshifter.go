package router

import (
	"context"
	"net/http"
)

type shifterKeyT struct{}

var shifterKey shifterKeyT

// DefaultShifterInjector return middleware to inject default shifter into context
func DefaultShifterInjector() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r.WithContext(
				context.WithValue(r.Context(), shifterKey,
					NewShifter(r.URL.EscapedPath()),
				),
			))
		}
	}
}

// DefaultShifter is default ShifterGetter for shifter that injected by DefaultShifterInjector
var DefaultShifter ShifterGetter = defaultShifter

func defaultShifter(r *http.Request) *Shifter {
	// if you got panic here, it mean that you are not using DefaultShifterInjector
	return r.Context().Value(shifterKey).(*Shifter)
}
