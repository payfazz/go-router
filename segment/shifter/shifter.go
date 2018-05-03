// Package shifter provide simple routing by dividing path into its segment
package shifter

import (
	"context"
	"net/http"
	"strings"
)

type ctxType struct{}

var defCtxKey ctxType

// Shifter hold state of shifting segment in the path
type Shifter struct {
	tag  map[string]int
	list []string
	next int
}

// From create Shifter from http.Request.
// It modify the context value, so old http.Request should not be used
func From(r *http.Request) (*Shifter, *http.Request) {
	return With(r, nil, nil)
}

// With create Shifter.
// It modify the context value, so old http.Request should not be used
//
// key is needed when there are multiple instance of shifter attached to current context.
//
// list is the segment of path. If nil, it will derived from r.URL.EscapedPath()
func With(r *http.Request, key interface{}, list []string) (*Shifter, *http.Request) {
	if key == nil {
		key = defCtxKey
	}

	if tmp := r.Context().Value(key); tmp != nil {
		return tmp.(*Shifter), r
	}

	if list == nil {
		list = strings.Split(
			strings.TrimPrefix(r.URL.EscapedPath(), "/"), "/",
		)
	}

	s := &Shifter{make(map[string]int), list, 0}

	r = r.WithContext(context.WithValue(
		r.Context(), key, s,
	))

	return s, r
}

// Shift to next segment, also telling if already in last segment
func (s *Shifter) Shift() (string, bool) {
	if s.next == len(s.list) {
		return "", true
	}
	ret := s.list[s.next]
	s.next++
	return ret, s.next == len(s.list)
}

// Unshift do reverse of Shift
func (s *Shifter) Unshift() {
	if s.next == 0 {
		return
	}
	s.next--
}

// Get i-th segment
func (s *Shifter) Get(i int) string {
	if i < 0 || i >= len(s.list) {
		return ""
	}
	return s.list[i]
}

// GetRelative is same with Get, but relative to current segment
func (s *Shifter) GetRelative(d int) string {
	return s.Get(s.CurrentIndex() + d)
}

// Size return the size of segment in path
func (s *Shifter) Size() int {
	return len(s.list)
}

// CurrentIndex of shifter state
func (s *Shifter) CurrentIndex() int {
	return s.next - 1
}

// End indicated end segment in the path
func (s *Shifter) End() bool {
	return s.next == len(s.list)
}

// Split return processed segment and rest of them
func (s *Shifter) Split() (done []string, rest []string) {
	done = make([]string, s.next)
	rest = make([]string, len(s.list)-s.next)
	copy(done, s.list[:s.next])
	copy(rest, s.list[s.next:])
	return done, rest
}

// Tag current segment
func (s *Shifter) Tag(tag string) {
	s.TagIndex(s.CurrentIndex(), tag)
}

// TagIndex will tag i-th segment
func (s *Shifter) TagIndex(i int, tag string) {
	if i < 0 || i >= len(s.list) {
		return
	}
	s.tag[tag] = i
}

// TagRelative is same with TagIndex, but relative to current segment
func (s *Shifter) TagRelative(d int, tag string) {
	s.TagIndex(s.CurrentIndex()+d, tag)
}

// GetByTag return tagged segment
func (s *Shifter) GetByTag(tag string) (string, bool) {
	i, ok := s.tag[tag]
	if !ok {
		return "", false
	}
	return s.list[i], true
}
