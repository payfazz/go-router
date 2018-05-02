// Package host provide host based routing
package host

import "net/http"

type H map[string]http.HandlerFunc

func Compile(h H, def http.HandlerFunc) http.HandlerFunc {
	panic("host: implementation is not complete yet")
}
