package shifter

import (
	"context"
	"net/http"
	"strings"
)

type ctxType struct{}

var defCtxKey ctxType

type Shifter struct {
	tag  map[string]int
	list []string
	next int
}

func From(r *http.Request) (*Shifter, *http.Request) {
	return With(r, nil, nil)
}

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

func (s *Shifter) Shift() (string, bool) {
	if s.next == len(s.list) {
		return "", true
	}
	ret := s.list[s.next]
	s.next++
	return ret, s.next == len(s.list)
}

func (s *Shifter) Unshift() {
	if s.next == 0 {
		return
	}
	s.next--
}

func (s *Shifter) Get(i int) string {
	if i < 0 || i >= len(s.list) {
		return ""
	}
	return s.list[i]
}

func (s *Shifter) GetRelative(d int) string {
	return s.Get(s.CurrentIndex() + d)
}

func (s *Shifter) Size() int {
	return len(s.list)
}

func (s *Shifter) CurrentIndex() int {
	return s.next - 1
}

func (s *Shifter) End() bool {
	return s.next == len(s.list)
}

func (s *Shifter) Split() (done []string, rest []string) {
	done = make([]string, s.next)
	rest = make([]string, len(s.list)-s.next)
	copy(done, s.list[:s.next])
	copy(rest, s.list[s.next:])
	return done, rest
}

func (s *Shifter) Tag(tag string) {
	s.TagIndex(s.CurrentIndex(), tag)
}

func (s *Shifter) TagIndex(i int, tag string) {
	if i < 0 || i >= len(s.list) {
		return
	}
	s.tag[tag] = i
}

func (s *Shifter) TagRelative(d int, tag string) {
	s.TagIndex(s.CurrentIndex()+d, tag)
}

func (s *Shifter) GetByTag(tag string) (string, bool) {
	i, ok := s.tag[tag]
	if !ok {
		return "", false
	}
	return s.list[i], true
}
