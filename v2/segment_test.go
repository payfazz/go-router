package router_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/payfazz/go-router/v2"
)

func assert(t *testing.T, cond bool, msgf string, msgs ...interface{}) {
	if !cond {
		t.Fatalf(msgf, msgs...)
	}
}

func genHandler(text string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, text)
	}
}
func TestRouterBySegment(t *testing.T) {
	handler := router.DefaultInjector()(router.Default.BySegmentWithDef(
		router.Hmap{
			"a": genHandler("/a"),
			"b": router.Default.SegmentMustEndOr(genHandler("404(2)"))(genHandler("/b")),
		},
		genHandler("404(1)"),
	))

	doTest := func(target string, data string, msg string) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", target, nil)
		handler.ServeHTTP(rec, req)
		assert(t, rec.Body.String() == data, msg)
	}

	doTest("/a", "/a", "normal a")
	doTest("/a/", "/a", "a with trailing slash")
	doTest("/a/a", "/a", "a with extra path")

	doTest("/c/a", "404(1)", "c with extra path")

	doTest("/b", "/b", "normal b")
	doTest("/b/", "/b", "b with trailing slash")
	doTest("/b/b", "404(2)", "b with extra path")
	doTest("/b//", "404(2)", "b with extra path but the first path is empty string (1)")
	doTest("/b//b", "404(2)", "b with extra path but the first path is empty string (2)")
}

func TestRouterByParam(t *testing.T) {
	handler := router.DefaultInjector()(
		router.Default.BySegment(router.Hmap{
			"lala": router.Default.ByParam(
				router.SetParamIntoHeader("X-Param"),
				genHandler("root"),
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintf(w, "param:%s", r.Header.Get("X-Param"))
				}),
			),
		}),
	)

	doTest := func(target string, data string, msg string) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", target, nil)
		handler.ServeHTTP(rec, req)
		tmp := rec.Body.String()
		assert(t, tmp == data, fmt.Sprintf("%s|%s|%s", msg, tmp, data))
	}

	doTest("/lala", "root", "root without trailing slash")
	doTest("/lala/", "root", "root with trailing slash")
	doTest("/lala/1", "param:1", "1 without trailing slash")
	doTest("/lala/2/", "param:2", "2 with trailing slash")
	doTest("/lala//aa", "param:", "empty param(1)")
	doTest("/lala//", "param:", "empty param(2)")
}

func TestRouterBySegment2(t *testing.T) {
	handler := router.DefaultInjector()(router.Default.BySegmentWithDef(
		router.Hmap{
			"":  genHandler("/"),
			"a": genHandler("/a"),
			"b": router.Default.BySegmentWithDef(
				router.Hmap{
					"":  genHandler("/b"),
					"a": genHandler("/b/a"),
					"b": genHandler("/b/b"),
				},
				genHandler("404(2)"),
			),
		},
		genHandler("404(1)"),
	))

	doTest := func(target string, data string, msg string) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", target, nil)
		handler.ServeHTTP(rec, req)
		assert(t, rec.Body.String() == data, msg)
	}

	doTest("/", "/", "test root")
	doTest("/a", "/a", "test /a")
	doTest("/a/", "/a", "test /a/")
	doTest("/a/a", "/a", "test /a/a")
	doTest("/b", "/b", "test /b")
	doTest("/b/a", "/b/a", "test /b/a")
	doTest("/b/a/a", "/b/a", "test /b/a/a")
	doTest("/b/a/b", "/b/a", "test /b/a/b")
	doTest("/b/b", "/b/b", "test /b/b")
	doTest("/b/d", "404(2)", "test /b/d")
	doTest("/c", "404(1)", "test /c")
}
