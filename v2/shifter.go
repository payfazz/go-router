package router

import (
	"strings"
)

// Shifter is state for routing.
//
// Where you store this state is up to you, usually the state is saved in
// the request's context.
type Shifter struct {
	segment []string
	cursor  int
}

// NewShifter from provided path, usually path is r.URL.EscapedPath()
func NewShifter(path string) *Shifter {
	return &Shifter{
		segment: strings.Split(
			strings.TrimPrefix(path, "/"), "/",
		),
		cursor: 0,
	}
}

func (s *Shifter) next() string {
	if s.cursor == len(s.segment) {
		return ""
	}
	segment := s.segment[s.cursor]
	s.cursor++
	return segment
}

func (s *Shifter) prev() {
	if s.cursor == 0 {
		return
	}
	s.cursor--
}

func (s *Shifter) state() (done, rest int) {
	return s.cursor, len(s.segment) - s.cursor
}

func (s *Shifter) end() bool {
	return s.cursor == len(s.segment)
}

// Split return processed segment and rest of them
func (s *Shifter) Split() (done, rest string) {
	doneN, restN := s.state()
	doneList := make([]string, doneN)
	restList := make([]string, restN)
	copy(doneList, s.segment[:s.cursor])
	copy(restList, s.segment[s.cursor:])
	return "/" + strings.Join(doneList, "/"), strings.Join(restList, "/")
}
