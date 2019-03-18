package path_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/payfazz/go-router/path"
	"github.com/payfazz/go-router/segment"
)

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

func respWriter(text string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a, b := segment.Split(r)
		sa := strings.Join(a, "/")
		sb := strings.Join(b, "/")
		p, pok := segment.Get(r, "p")
		fmt.Fprintf(w, "%s|%s|%s|%s|%t", text, sa, sb, p, pok)
	}
}

func Test1(t *testing.T) {
	h := path.H{
		"/a":            respWriter("/a"),
		"/a/:p":         respWriter("/a/:p"),
		"/a/:p/a":       respWriter("/a/:p/a"),
		"/b":            respWriter("/b"),
		"/b/a":          respWriter("/b/a"),
		"/d/a/b/c/d/e":  respWriter("/d/a/b/c/d/e"),
		"/d/a/b/:p/d/e": respWriter("/d/a/b/:p/d/e"),
		"/e/:p":         respWriter("/e/:p"),
	}.Compile(respWriter("nf"))

	t.Run("1", doTest(h, "/a", "/a|a|||false"))
	t.Run("2", doTest(h, "/a/", "/a|a/|||false"))
	t.Run("3", doTest(h, "/a/a", "/a/:p|a/a||a|true"))
	t.Run("4", doTest(h, "/a/a/b", "/a/:p|a/a|b|a|true"))
	t.Run("5", doTest(h, "/a/a/a", "/a/:p/a|a/a/a||a|true"))
	t.Run("6", doTest(h, "/x/y/z", "nf||x/y/z||false"))
	t.Run("7", doTest(h, "/b", "/b|b|||false"))
	t.Run("8", doTest(h, "/c", "nf||c||false"))
	t.Run("9", doTest(h, "/b/a", "/b/a|b/a|||false"))
	t.Run("10", doTest(h, "/d/a/b/c/d/e", "/d/a/b/c/d/e|d/a/b/c/d/e|||false"))
	t.Run("11", doTest(h, "/d/a/b/p/d/e", "/d/a/b/:p/d/e|d/a/b/p/d/e||p|true"))
	t.Run("12", doTest(h, "/d/a/b//d/e", "/d/a/b/:p/d/e|d/a/b//d/e|||true"))
	t.Run("13", doTest(h, "/e", "/e/:p|e|||false"))
	t.Run("14", doTest(h, "/e/", "/e/:p|e/|||true"))
}

func Test2(t *testing.T) {
	h := path.H{
		"/a/b/c/d": respWriter("/a/b/c/d"),
		"/":        respWriter("/"),
	}.Compile(respWriter("nf"))

	t.Run("1", doTest(h, "/a/b/c/d/e/f/g", "/a/b/c/d|a/b/c/d|e/f/g||false"))
	t.Run("2", doTest(h, "/b/b/c/d/e/f/g", "/||b/b/c/d/e/f/g||false"))
}

func Test3(t *testing.T) {
	h := path.H{
		"/a/b/c/d": respWriter("/a/b/c/d"),
		"/d/e/f/g": respWriter("/d/e/f/g"),
	}.Compile(respWriter("nf"))

	t.Run("1", doTest(h, "/a/b/c/d/e/f/g", "/a/b/c/d|a/b/c/d|e/f/g||false"))
	t.Run("2", doTest(h, "/b/b/c/d/e/f/g", "nf||b/b/c/d/e/f/g||false"))
	t.Run("3", doTest(h, "/d/e/h/i", "nf|d/e|h/i||false"))
}

func Test4(t *testing.T) {
	h := path.H{
		"/a/b/c/d": segment.Stripper(respWriter("/a/b/c/d")),
		"/a":       respWriter("/a"),
		"/":        respWriter("/"),
	}.Compile(respWriter("nf"))

	t.Run("1", doTest(h, "/a/b/c/d/g/h/i/j", "/a/b/c/d||g/h/i/j||false"))
	t.Run("2", doTest(h, "/a/b/c", "/a|a|b/c||false"))
	t.Run("3", doTest(h, "/c/d/e", "/||c/d/e||false"))
}
