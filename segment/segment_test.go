package segment_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	internalsegment "github.com/payfazz/go-router/internal/segment"
	"github.com/payfazz/go-router/segment"
)

type data struct {
	h        http.HandlerFunc
	path     string
	expected string
}

func respWriter(text string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a, b := segment.Split(r)
		sa := strings.Join(a, "/")
		sb := strings.Join(b, "/")
		fmt.Fprintf(w, "%s|%d:%s|%d:%s", text, len(a), sa, len(b), sb)
	}
}

func doTest(t *testing.T, data []data) {
	t.Parallel()

	for i := 0; i < len(data); i++ {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			path := data[i].path
			h := data[i].h
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
	h := segment.H{
		"":  respWriter(""),
		"a": respWriter("a"),
		"b": segment.H{
			"a": respWriter("b/a"),
			"b": respWriter("b/b"),
		}.Compile(respWriter("nfb")),
		"c": segment.H{
			"":  respWriter("c"),
			"a": respWriter("c/a"),
			"b": respWriter("c/b"),
		}.Compile(respWriter("nfc")),
	}.Compile(respWriter("nf"))

	doTest(t, []data{
		{h, "/", "|1:|0:"},
		{h, "/a", "a|1:a|0:"},
		{h, "/b", "nfb|1:b|0:"},
		{h, "/b/", "nfb|1:b|1:"},
		{h, "/b/c", "nfb|1:b|1:c"},
		{h, "/b/c/d", "nfb|1:b|2:c/d"},
		{h, "/b/c/d/", "nfb|1:b|3:c/d/"},
		{h, "/b/a", "b/a|2:b/a|0:"},
		{h, "/b/a/", "b/a|2:b/a|1:"},
		{h, "/b/b", "b/b|2:b/b|0:"},
		{h, "/c", "c|1:c|0:"},
		{h, "/c/", "c|2:c/|0:"},
		{h, "/c/a", "c/a|2:c/a|0:"},
		{h, "/c/a/", "c/a|2:c/a|1:"},
		{h, "/c/b", "c/b|2:c/b|0:"},
		{h, "/c/c", "nfc|1:c|1:c"},
		{h, "/c/c/a/b/c", "nfc|1:c|4:c/a/b/c"},
		{h, "/c/c/a/b/c/", "nfc|1:c|5:c/a/b/c/"},
		{h, "/d", "nf|0:|1:d"},
		{h, "/d/", "nf|0:|2:d/"},
		{h, "/d/c/a/b/c", "nf|0:|5:d/c/a/b/c"},
		{h, "/d/c/a/b/c", "nf|0:|5:d/c/a/b/c"},
	})
}

func Test2(t *testing.T) {
	respWriter2 := func(text string, count int) http.HandlerFunc {
		tmp := respWriter(text)
		return func(w http.ResponseWriter, r *http.Request) {
			s, r2 := internalsegment.TryShifterFrom(r)
			s.SetNext(count)
			tmp(w, r2)
		}
	}
	h := segment.H{
		"c": segment.H{
			"x": respWriter2("c/x", 0),
			"y": respWriter2("c/y", 1),
			"z": respWriter2("c/z", 3),
			"a": respWriter2("c/a", 99999),
		}.Compile(respWriter("nfc")),
	}.Compile(respWriter("nf"))

	doTest(t, []data{
		{h, "/c/x/a/b/c/d", "c/x|0:|6:c/x/a/b/c/d"},
		{h, "/c/y/a/b/c/d", "c/y|1:c|5:y/a/b/c/d"},
		{h, "/c/z/a/b/c/d", "c/z|3:c/z/a|3:b/c/d"},
		{h, "/c/p", "nfc|1:c|1:p"},
		{h, "/c/p/", "nfc|1:c|2:p/"},
		{h, "/d/p", "nf|0:|2:d/p"},
		{h, "/d/p/", "nf|0:|3:d/p/"},
		{h, "/c/a/a/b/c/d", "c/a|6:c/a/a/b/c/d|0:"},
		{h, "/c/a/a/b/c/d/", "c/a|7:c/a/a/b/c/d/|0:"},
	})
}
