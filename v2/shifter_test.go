package router

import (
	"testing"
)

func assert(t *testing.T, cond bool, msgf string, msgs ...interface{}) {
	if !cond {
		t.Fatalf(msgf, msgs...)
	}
}

func TestShifterNextPrev(t *testing.T) {
	s := NewShifter("/a/b/c")
	assert(t, s.next() == "a", "must return a")
	assert(t, s.next() == "b", "must return b")
	s.prev()
	assert(t, s.next() == "b", "must return b")
	assert(t, s.next() == "c", "must return c")
	assert(t, s.next() == "", "must return emptystring")
	s.prev()
	s.prev()
	assert(t, s.next() == "b", "must return b")
	assert(t, s.next() == "c", "must return c")
	assert(t, s.next() == "", "must return emptystring")
}

func TestShifterNextEnd(t *testing.T) {
	s := NewShifter("/a/b/c")
	assert(t, s.next() == "a", "must return a")
	assert(t, s.next() == "b", "must return b")
	assert(t, s.next() == "c", "must return c")
	assert(t, s.end(), "must end")
}

func TestShifterNextEndTrailingSlash(t *testing.T) {
	s := NewShifter("/a/b/c/")
	assert(t, s.next() == "a", "must return a")
	assert(t, s.next() == "b", "must return b")
	assert(t, s.next() == "c", "must return c")
	assert(t, !s.end(), "must not end")
	assert(t, s.next() == "", "emptystring")
	assert(t, s.end(), "must end")
}

func TestShifterStateAndSplit(t *testing.T) {
	var doneN, restN int
	var done, rest string
	s := NewShifter("/a/b/c")

	doneN, restN = s.state()
	done, rest = s.Split()
	assert(t, doneN == 0 && restN == 3 && done == "/" && rest == "a/b/c", "invalid state and split 1")

	s.next()

	doneN, restN = s.state()
	done, rest = s.Split()
	assert(t, doneN == 1 && restN == 2 && done == "/a" && rest == "b/c", "invalid state and split 2")

	s.next()

	doneN, restN = s.state()
	done, rest = s.Split()
	assert(t, doneN == 2 && restN == 1 && done == "/a/b" && rest == "c", "invalid state and split 3")
}

func TestShifterStateAndSplitWithTrailingSlash(t *testing.T) {
	var doneN, restN int
	var done, rest string
	s := NewShifter("/a/b/c/")

	doneN, restN = s.state()
	done, rest = s.Split()
	assert(t, doneN == 0 && restN == 4 && done == "/" && rest == "a/b/c/", "invalid state and split 1")

	s.next()

	doneN, restN = s.state()
	done, rest = s.Split()
	assert(t, doneN == 1 && restN == 3 && done == "/a" && rest == "b/c/", "invalid state and split 2")

	s.next()

	doneN, restN = s.state()
	done, rest = s.Split()
	assert(t, doneN == 2 && restN == 2 && done == "/a/b" && rest == "c/", "invalid state and split 3")
}
