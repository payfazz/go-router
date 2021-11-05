package path_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/payfazz/go-router/path"
	"github.com/payfazz/go-router/segment"
)

type data struct {
	path     string
	expected string
}

func respWriter(text string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a, b := segment.Split(r)
		sa := strings.Join(a, "/")
		sb := strings.Join(b, "/")
		p, pok := segment.Get(r, "p")
		q, qok := segment.Get(r, "q")
		fmt.Fprintf(w, "%s|%d:%s|%d:%s|%t:%s|%t:%s", text, len(a), sa, len(b), sb, pok, p, qok, q)
	}
}

func doTest(t *testing.T, h http.HandlerFunc, data []data) {
	t.Parallel()

	for i := 0; i < len(data); i++ {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			path := data[i].path
			expected := data[i].expected

			t.Parallel()

			r := httptest.NewRequest("GET", path, nil)
			w := httptest.NewRecorder()

			h(w, r)

			b := w.Body.String()
			if b != expected {
				t.Fatalf("expected '%s', found '%s'", expected, b)
			}
		})
	}
}

func Test1(t *testing.T) {
	h := path.H{
		"/a":            respWriter("a"),
		"/a/:p":         respWriter("a/:p"),
		"/a/:p/a":       respWriter("a/:p/a"),
		"/a/:p/a/:q/a":  respWriter("a/:p/a/:q/a"),
		"/b":            respWriter("b"),
		"/b/a":          respWriter("b/a"),
		"/d/a/b/c/d/e":  respWriter("d/a/b/c/d/e"),
		"/d/a/b/:p/d/e": respWriter("d/a/b/:p/d/e"),
		"/e/:p":         respWriter("e/:p"),
	}.Compile(respWriter("nf"))

	doTest(t, h, []data{
		{"/a", "a|1:a|0:|false:|false:"},
		{"/a/", "a|2:a/|0:|false:|false:"},
		{"/a/a", "a/:p|2:a/a|0:|true:a|false:"},
		{"/a/a/", "a/:p|3:a/a/|0:|true:a|false:"},
		{"/a/a/b", "a/:p|2:a/a|1:b|true:a|false:"},
		{"/a/a/b/", "a/:p|2:a/a|2:b/|true:a|false:"},
		{"/a/a/a", "a/:p/a|3:a/a/a|0:|true:a|false:"},
		{"/a/a/a/", "a/:p/a|4:a/a/a/|0:|true:a|false:"},
		{"/a/pp/a/qq/a", "a/:p/a/:q/a|5:a/pp/a/qq/a|0:|true:pp|true:qq"},
		{"/a/pp/a/qq/b", "a/:p/a|3:a/pp/a|2:qq/b|true:pp|false:"},
		{"/a/pp/a/qq/b/", "a/:p/a|3:a/pp/a|3:qq/b/|true:pp|false:"},
		{"/x/y/z", "nf|0:|3:x/y/z|false:|false:"},
		{"/b", "b|1:b|0:|false:|false:"},
		{"/b/", "b|2:b/|0:|false:|false:"},
		{"/c", "nf|0:|1:c|false:|false:"},
		{"/c/", "nf|0:|2:c/|false:|false:"},
		{"/b/a", "b/a|2:b/a|0:|false:|false:"},
		{"/b/a/", "b/a|2:b/a|1:|false:|false:"},
		{"/b/a/c/d", "b/a|2:b/a|2:c/d|false:|false:"},
		{"/b/a/c/d/", "b/a|2:b/a|3:c/d/|false:|false:"},
		{"/d/a/b/c/d/e", "d/a/b/c/d/e|6:d/a/b/c/d/e|0:|false:|false:"},
		{"/d/a/b/c/d/e/", "d/a/b/c/d/e|6:d/a/b/c/d/e|1:|false:|false:"},
		{"/d/a/b/p/d/e", "d/a/b/:p/d/e|6:d/a/b/p/d/e|0:|true:p|false:"},
		{"/d/a/b/p/d/e/", "d/a/b/:p/d/e|6:d/a/b/p/d/e|1:|true:p|false:"},
		{"/d/a/b//d/e", "d/a/b/:p/d/e|6:d/a/b//d/e|0:|true:|false:"},
		{"/d/a/b//d/e/", "d/a/b/:p/d/e|6:d/a/b//d/e|1:|true:|false:"},
		{"/e", "e/:p|1:e|0:|false:|false:"},
		{"/e/", "e/:p|2:e/|0:|true:|false:"},
		{"/e/x", "e/:p|2:e/x|0:|true:x|false:"},
		{"/e/x/", "e/:p|2:e/x|1:|true:x|false:"},
	})
}

func Test2(t *testing.T) {
	h := path.H{
		"/a/b/c/d": respWriter("a/b/c/d"),
		"/":        respWriter(""),
	}.Compile(respWriter("nf"))

	doTest(t, h, []data{
		{"/a/b/c/d/e/f/g", "a/b/c/d|4:a/b/c/d|3:e/f/g|false:|false:"},
		{"/b/b/c/d/e/f/g", "|0:|7:b/b/c/d/e/f/g|false:|false:"},
	})
}

func Test3a(t *testing.T) {
	h := path.H{
		"/a/b/c/d": respWriter("a/b/c/d"),
		"/d/e/f/g": respWriter("d/e/f/g"),
	}.Compile(respWriter("nf"))

	doTest(t, h, []data{
		{"/a/b/c/d/e/f/g", "a/b/c/d|4:a/b/c/d|3:e/f/g|false:|false:"},
		{"/b/b/c/d/e/f/g", "nf|0:|7:b/b/c/d/e/f/g|false:|false:"},
		{"/d/e/h/i", "nf|2:d/e|2:h/i|false:|false:"},
		{"/d/e/h/i/", "nf|2:d/e|3:h/i/|false:|false:"},
	})
}

func Test3b(t *testing.T) {
	h := path.H{
		"/a/b/c/d": respWriter("a/b/c/d"),
		"/d/e/f/g": respWriter("d/e/f/g"),
		"/":        respWriter(""),
	}.Compile(respWriter("nf"))

	doTest(t, h, []data{
		{"/a/b/c/d/e/f/g", "a/b/c/d|4:a/b/c/d|3:e/f/g|false:|false:"},
		{"/b/b/c/d/e/f/g", "|0:|7:b/b/c/d/e/f/g|false:|false:"},
		{"/d/e/h/i", "|0:|4:d/e/h/i|false:|false:"},
		{"/d/e/h/i/", "|0:|5:d/e/h/i/|false:|false:"},
	})
}

func Test4a(t *testing.T) {
	h := path.H{
		"/a/b/c/d": segment.Stripper(respWriter("a/b/c/d")),
		"/a":       respWriter("a"),
		"/":        respWriter(""),
	}.Compile(respWriter("nf"))

	doTest(t, h, []data{
		{"/a/b/c/d/g/h/i/j", "a/b/c/d|0:|4:g/h/i/j|false:|false:"},
		{"/a/b/c", "a|1:a|2:b/c|false:|false:"},
		{"/a/b/c/", "a|1:a|3:b/c/|false:|false:"},
		{"/c/d/e", "|0:|3:c/d/e|false:|false:"},
	})
}

func Test4b(t *testing.T) {
	h := path.H{
		"/a/b/c/d": segment.Stripper(respWriter("a/b/c/d")),
		"/a":       respWriter("a"),
	}.Compile(respWriter("nf"))

	doTest(t, h, []data{
		{"/a/b/c/d/g/h/i/j", "a/b/c/d|0:|4:g/h/i/j|false:|false:"},
		{"/a/b/c", "a|1:a|2:b/c|false:|false:"},
		{"/a/b/c/", "a|1:a|3:b/c/|false:|false:"},
		{"/c/d/e", "nf|0:|3:c/d/e|false:|false:"},
	})
}

func Test5(t *testing.T) {
	h := path.H{
		"/a/exact": respWriter("a/exact"),
		"/a/:p":    respWriter("a/:p"),
	}.Compile(respWriter("nf"))

	doTest(t, h, []data{
		{"/a/exact", "a/exact|2:a/exact|0:|false:|false:"},
		{"/a/exact/", "a/exact|2:a/exact|1:|false:|false:"},
		{"/a/exact/b/c/d", "a/exact|2:a/exact|3:b/c/d|false:|false:"},
		{"/a/exact/b/c/d/", "a/exact|2:a/exact|4:b/c/d/|false:|false:"},
		{"/a/something", "a/:p|2:a/something|0:|true:something|false:"},
		{"/a/something/", "a/:p|2:a/something|1:|true:something|false:"},
		{"/a/something/b/c/d", "a/:p|2:a/something|3:b/c/d|true:something|false:"},
		{"/a/something/b/c/d/", "a/:p|2:a/something|4:b/c/d/|true:something|false:"},
	})
}

func TestStripper(t *testing.T) {
	h := path.H{
		"/a/b": segment.Stripper(respWriter("")),
	}.Compile(respWriter("nf"))
	doTest(t, h, []data{
		{"/a/b/c/d", "|0:|2:c/d|false:|false:"},
	})
}

func TestDuplicateParamName(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("duplicate parameter name should panic")
		}
	}()

	path.H{
		"/asdf/:gege/haha": nil,
		"/asdf/:lele/haha": nil,
	}.C()
}

func TestDuplicateHandler(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("duplicate handler should panic")
		}
	}()

	path.H{
		"/asdf/": nil,
		"/asdf":  nil,
	}.C()
}
