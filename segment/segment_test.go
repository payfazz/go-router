package segment_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	h := segment.H{
		"":  respWriter("/"),
		"a": respWriter("/a"),
		"b": segment.H{
			"a": respWriter("/b/a"),
			"b": respWriter("/b/b"),
		}.Compile(respWriter("nf")),
		"c": segment.H{
			"":  respWriter("/c"),
			"a": respWriter("/c/a"),
			"b": respWriter("/c/b"),
			"x": func() http.HandlerFunc {
				tmp := respWriter("/c/x")
				return func(w http.ResponseWriter, r *http.Request) {
					segment.UnshiftInternalShifter(r, 2)
					tmp(w, r)
				}
			}(),
		}.Compile(respWriter("nf")),
	}.Compile(respWriter("nf"))

	t.Run("1", doTest(h, "/", "/||"))
	t.Run("2", doTest(h, "/a", "/a|a|"))
	t.Run("3", doTest(h, "/b", "nf|b|"))
	t.Run("4", doTest(h, "/b/a", "/b/a|b/a|"))
	t.Run("5", doTest(h, "/b/b", "/b/b|b/b|"))
	t.Run("6", doTest(h, "/c", "/c|c|"))
	t.Run("7", doTest(h, "/c/a", "/c/a|c/a|"))
	t.Run("8", doTest(h, "/c/b", "/c/b|c/b|"))
	t.Run("9", doTest(h, "/c/c/a/b/c", "nf|c|c/a/b/c"))
	t.Run("10", doTest(h, "/d/c/a/b/c", "nf||d/c/a/b/c"))
	t.Run("11", doTest(h, "/c/x/d/e/f", "/c/x||c/x/d/e/f"))
}
