package router

import (
	"strings"
)

// State for routing.
//
// Where you store this state is up to you, usually the state is saved in
// the request's context.
type State struct {
	segment []string
	cursor  int
}

// NewState from provided path, usually path is r.URL.EscapedPath()
func NewState(path string) *State {
	return &State{
		segment: strings.Split(
			strings.TrimPrefix(path, "/"), "/",
		),
		cursor: 0,
	}
}

func (s *State) next() string {
	if s.cursor == len(s.segment) {
		return ""
	}
	segment := s.segment[s.cursor]
	s.cursor++
	return segment
}

func (s *State) prev() {
	if s.cursor == 0 {
		return
	}
	s.cursor--
}

func (s *State) progressCursor() (int, int) {
	return s.cursor, len(s.segment) - s.cursor
}

// Progress return processed segment and rest of them
func (s *State) Progress() (done, rest string) {
	doneN, restN := s.progressCursor()
	doneList := make([]string, doneN)
	restList := make([]string, restN)
	copy(doneList, s.segment[:s.cursor])
	copy(restList, s.segment[s.cursor:])
	return "/" + strings.Join(doneList, "/"), strings.Join(restList, "/")
}
