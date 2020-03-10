package router

import (
	"strings"
)

// State for routing.
//
// Where you store this state is up to you, usually the state is saved in
// the request's context.
//
// you should embed this struct into another struct and call Init to initialize the
// routing state
type State struct {
	path string
	idx  int
}

// Init from provided path, usually path is r.URL.EscapedPath()
func (s *State) Init(path string) {
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	s.path = path
	s.idx = 0
}

func (s *State) next() string {
	if s.idx == len(s.path) {
		return ""
	}
	s.idx++
	end := strings.IndexByte(s.path[s.idx:], '/')
	if end == -1 {
		end = len(s.path)
	} else {
		end += s.idx
	}
	segment := s.path[s.idx:end]
	s.idx = end
	return segment
}

func (s *State) prev() {
	if s.idx == 0 {
		return
	}
	s.idx = strings.LastIndexByte(s.path[:s.idx], '/')
}

// Progress return processed segment and rest of them
func (s *State) Progress() (done, rest string) {
	return s.path[:s.idx], s.path[s.idx:]
}
