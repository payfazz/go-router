package router

import (
	"testing"
)

func assert(t *testing.T, cond bool, msgf string, msgs ...interface{}) {
	if !cond {
		t.Fatalf(msgf, msgs...)
	}
}

func stateIsEnd(s *State) bool {
	_, rest := s.Progress()
	return rest == ""
}

func TestShifterNextPrev(t *testing.T) {
	var s State
	s.Init("/a/b/c")
	assert(t, s.next() == "a", "must return a")
	assert(t, s.next() == "b", "must return b")
	s.prev()
	assert(t, s.next() == "b", "must return b")
	assert(t, s.next() == "c", "must return c")
	assert(t, s.next() == "", "must return emptystring")
	s.prev()
	s.prev()
	s.prev()
	s.prev()
	assert(t, s.next() == "a", "must return a")
	assert(t, s.next() == "b", "must return b")
	assert(t, s.next() == "c", "must return c")
	assert(t, s.next() == "", "must return emptystring")
}

func TestShifterNextEnd(t *testing.T) {
	var s State
	s.Init("/a/b/c")
	assert(t, s.next() == "a", "must return a")
	assert(t, s.next() == "b", "must return b")
	assert(t, s.next() == "c", "must return c")
	assert(t, stateIsEnd(&s), "must end")
}

func TestShifterNextEndTrailingSlash(t *testing.T) {
	var s State
	s.Init("/a/b/c/")
	assert(t, s.next() == "a", "must return a")
	assert(t, s.next() == "b", "must return b")
	assert(t, s.next() == "c", "must return c")
	assert(t, !stateIsEnd(&s), "must not end")
	assert(t, s.next() == "", "emptystring")
	assert(t, stateIsEnd(&s), "must end")
}

func TestShifterStateAndSplit(t *testing.T) {
	var done, rest string
	var s State
	s.Init("/a/b/c")

	done, rest = s.Progress()
	assert(t, done == "" && rest == "/a/b/c", "invalid progress 1")

	s.next()

	done, rest = s.Progress()
	assert(t, done == "/a" && rest == "/b/c", "invalid progress 2")

	s.next()

	done, rest = s.Progress()
	assert(t, done == "/a/b" && rest == "/c", "invalid progress 3")

	s.next()

	done, rest = s.Progress()
	assert(t, done == "/a/b/c" && rest == "", "invalid progress 4")

	s.next()

	done, rest = s.Progress()
	assert(t, done == "/a/b/c" && rest == "", "invalid progress 5")
}

func TestShifterStateAndSplitWithTrailingSlash(t *testing.T) {
	var done, rest string
	var s State
	s.Init("/a/b/c/")

	done, rest = s.Progress()
	assert(t, done == "" && rest == "/a/b/c/", "invalid progress 1")

	s.next()

	done, rest = s.Progress()
	assert(t, done == "/a" && rest == "/b/c/", "invalid progress 2")

	s.next()

	done, rest = s.Progress()
	assert(t, done == "/a/b" && rest == "/c/", "invalid progress 3")

	s.next()

	done, rest = s.Progress()
	assert(t, done == "/a/b/c" && rest == "/", "invalid progress 4")

	s.next()

	done, rest = s.Progress()
	assert(t, done == "/a/b/c/" && rest == "", "invalid progress 5")

	s.next()

	done, rest = s.Progress()
	assert(t, done == "/a/b/c/" && rest == "", "invalid progress 6")
}
