// Package path provide path based routing
package path

import (
	"net/http"
)

type H map[string]http.HandlerFunc

func Compile(h H, notfoundHandler http.HandlerFunc) http.HandlerFunc {
	b := &builderT{make(tree)}
	for k, v := range h {
		b.add(k, v)
	}
	return b.compile(notfoundHandler)
}
