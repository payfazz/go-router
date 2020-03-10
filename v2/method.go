package router

import (
	"net/http"

	v1method "github.com/payfazz/go-router/method"
)

// ByMethod will return handler for routing via method
func (Router) ByMethod(handler HandlerMapping) http.HandlerFunc {
	return ((v1method.H)(handler)).C()
}
