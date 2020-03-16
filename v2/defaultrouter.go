package router

import (
	"context"
	"net/http"
)

type stateKeyT struct{}

var stateKey stateKeyT

// DefaultInjector return middleware to inject default router state into context
func DefaultInjector() func(Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			var state State
			state.Init(r.URL.EscapedPath())
			next(w, r.WithContext(
				context.WithValue(r.Context(), stateKey, &state),
			))
		}
	}
}

// Default is default Router for state that injected by DefaultInjector
var Default Router = defaultRouter

func defaultRouter(r *http.Request) *State {
	// if you got panic here, it mean that you are not using DefaultInjector
	return r.Context().Value(stateKey).(*State)
}
