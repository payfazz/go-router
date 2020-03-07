package router

import (
	"context"
	"net/http"
	"strings"
)

// SegmentShifter state for routing
type SegmentShifter struct {
	segment []string
	cursor  int
}

// NewSegmentShifter from provided path, usually path is r.URL.EscapedPath()
func NewSegmentShifter(path string) *SegmentShifter {
	return &SegmentShifter{
		segment: strings.Split(
			strings.TrimPrefix(path, "/"), "/",
		),
		cursor: 0,
	}
}

// Next will return current segment, then increase the cursor
//
// will return empty string if already reached the end
func (s *SegmentShifter) Next() string {
	if s.End() {
		return ""
	}
	segment := s.segment[s.cursor]
	s.cursor++
	return segment
}

// Prev decrease the cursor
func (s *SegmentShifter) Prev() {
	if s.cursor == 0 {
		return
	}
	s.cursor--
}

// End indicate that already reached the end
func (s *SegmentShifter) End() bool {
	return s.cursor == len(s.segment)
}

// Split return processed segment and rest of them
func (s *SegmentShifter) Split() (done, rest string) {
	doneList := make([]string, s.cursor)
	restList := make([]string, len(s.segment)-s.cursor)
	copy(doneList, s.segment[:s.cursor])
	copy(restList, s.segment[s.cursor:])
	return "/" + strings.Join(doneList, "/"), strings.Join(restList, "/")
}

type shifterKeyT struct{}

var shifterKey shifterKeyT

// ShifterInjector return middleware to inject default shifter into context
func ShifterInjector() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r.WithContext(
				context.WithValue(r.Context(), shifterKey,
					NewSegmentShifter(r.URL.EscapedPath()),
				),
			))
		}
	}
}

func defaultShifter(r *http.Request) *SegmentShifter {
	return r.Context().Value(shifterKey).(*SegmentShifter)
}
