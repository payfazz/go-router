package segment_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	segmentInternal "github.com/payfazz/go-router/internal/segment"
	"github.com/payfazz/go-router/segment"
)

func respWriter(text string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a, b := segment.Split(r)
		sa := strings.Join(a, "/")
		sb := strings.Join(b, "/")
		fmt.Fprintf(w, "%s|%s|%s", text, sa, sb)
	}
}

func doTest(h http.HandlerFunc, url, expected string) func(t *testing.T) {
	return func(t *testing.T) {
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		h(w, r)

		b := w.Body.String()
		if b != expected {
			t.Fatalf("expected '%s', found '%s'", expected, b)
		}
	}
}

func Test1(t *testing.T) {
	respWriter2 := func(text string, count int) http.HandlerFunc {
		tmp := respWriter(text)
		return func(w http.ResponseWriter, r *http.Request) {
			segmentInternal.SetShifterIndex(r, count)
			tmp(w, r)
		}
	}
	h := segment.H{
		"c": segment.H{
			"x": respWriter2("/c/x", 0),
			"y": respWriter2("/c/y", 1),
			"z": respWriter2("/c/z", 3),
			"a": respWriter2("/c/a", 999),
		}.Compile(respWriter("nf")),
	}.Compile(respWriter("nf"))

	t.Run("1", doTest(h, "/c/x/a/b/c/d", "/c/x||c/x/a/b/c/d"))
	t.Run("2", doTest(h, "/c/y/a/b/c/d", "/c/y|c|y/a/b/c/d"))
	t.Run("3", doTest(h, "/c/z/a/b/c/d", "/c/z|c/z/a|b/c/d"))
	t.Run("4", doTest(h, "/c/a/a/b/c/d", "/c/a|c/a/a/b/c/d|"))
}
