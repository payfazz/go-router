// Package path provide path based routing
package path

import "net/http"

type H map[string]http.HandlerFunc

func Compile(h H, def http.HandlerFunc) http.HandlerFunc {
	panic("path: implementation is not complete yet")
}
