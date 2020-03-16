package router_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/payfazz/go-router/v2"
)

func Example() {
	genHandler := func(text string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, text)
		}
	}

	paramKey := "X-Some-Param"

	genPrintParam := func(prefix string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s%s", prefix, r.Header.Get(paramKey))
		}
	}

	handler := router.DefaultInjector()(
		router.Default.BySegment(router.Hmap{
			"":  genHandler("/"),
			"a": genHandler("/a"),
			"b": genHandler("/b"),
			"c": router.Default.ByParam(
				router.SetParamIntoHeader(paramKey),
				genHandler("/c"),
				genPrintParam("/c/"),
			),
		}),
	)

	assertHandler := func(url string, expected string) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		handler(rec, req)
		result := rec.Body.String()
		if result != expected {
			panic(fmt.Errorf("expecting '%s', but got '%s'", expected, result))
		}
	}

	assertHandler("/", "/")

	assertHandler("/a", "/a")
	assertHandler("/a/", "/a")
	assertHandler("/a/a", "/a")

	assertHandler("/b", "/b")
	assertHandler("/b/", "/b")
	assertHandler("/b/b", "/b")

	assertHandler("/c", "/c")
	assertHandler("/c/", "/c")

	assertHandler("/c/a", "/c/a")
	assertHandler("/c/b", "/c/b")
	assertHandler("/c/c", "/c/c")

	assertHandler("/c/c/c", "/c/c")

	assertHandler("/c//c", "/c/") // because empty segment is not last, it will treat it as parameter
}

func TestExample(t *testing.T) {
	Example()
}
