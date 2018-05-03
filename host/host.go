// Package host provide host based routing
package host

import "net/http"

// H is type for mapping host and its handler
type H map[string]http.HandlerFunc

// Compile into single http.HandlerFunc
func Compile(h H, def http.HandlerFunc) http.HandlerFunc {
	panic("host: implementation is not complete yet")
}
