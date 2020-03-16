package router

import (
	v1method "github.com/payfazz/go-router/method"
)

// ByMethod will return hmap for routing via method
func (router Router) ByMethod(hmap Hmap) Handler {
	return ((v1method.H)(hmap)).C()
}
