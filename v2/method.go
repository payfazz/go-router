package router

import (
	"net/http"

	v1method "github.com/payfazz/go-router/method"
)

// ByMethod will return handler for routing via method
func ByMethod(next map[string]http.HandlerFunc) http.HandlerFunc {
	return ((v1method.H)(next)).C()
}
