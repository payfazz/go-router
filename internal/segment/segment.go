package segment

import (
	"net/http"

	"github.com/payfazz/go-router/defhandler"
	"github.com/payfazz/go-router/segment/shifter"
)

// H .
type H map[string]HandlerFunc

func compile(h H, def HandlerFunc) HandlerFunc {
	if def == nil {
		def = FromStd(defhandler.StatusNotFound)
	}
	for k, v := range h {
		if v == nil {
			h[k] = FromStd(defhandler.StatusNotImplemented)
		}
	}
	return func(s *shifter.Shifter, w http.ResponseWriter, r *http.Request) {
		var next HandlerFunc
		end := s.End()
		seg, _ := s.Shift()
		next, ok := h[seg]
		if !ok {
			next = def
			if !end {
				s.Unshift()
			}
		}
		next(s, w, r)
	}
}

// Compile .
func (h H) Compile(def HandlerFunc) HandlerFunc {
	return compile(h, def)
}

// Tag .
func Tag(tag string, next HandlerFunc) HandlerFunc {
	return func(s *shifter.Shifter, w http.ResponseWriter, r *http.Request) {
		if !s.End() {
			s.Shift()
			s.Tag(tag)
		}
		next(s, w, r)
	}
}
