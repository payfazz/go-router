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
	handler := router.DefaultShifterInjector()(router.DefaultShifter.BySegmentWithDef(
		map[string]http.HandlerFunc{
			"a": genHandler("/a"),
			"b": router.DefaultShifter.SegmentMustEndOr(genHandler("404(2)"))(genHandler("/b")),
		},
		genHandler("404(1)"),
	))

	doTest := func(target string, data string, msg string) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", target, nil)
		handler(rec, req)
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
	handler := router.DefaultShifterInjector()(
		router.DefaultShifter.BySegment(map[string]http.HandlerFunc{
			"lala": router.DefaultShifter.ByParam(
				router.SetParamIntoHeader("X-Param"),
				genHandler("root"),
				func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintf(w, "param:%s", r.Header.Get("X-Param"))
				},
			),
		}),
	)

	doTest := func(target string, data string, msg string) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", target, nil)
		handler(rec, req)
		assert(t, rec.Body.String() == data, msg)
	}

	doTest("/lala", "root", "root without trailing slash")
	doTest("/lala/", "root", "root with trailing slash")
	doTest("/lala/1", "param:1", "1 without trailing slash")
	doTest("/lala/2/", "param:2", "2 with trailing slash")
	doTest("/lala//aa", "param:", "empty param(1)")
	doTest("/lala//", "param:", "empty param(2)")
}
