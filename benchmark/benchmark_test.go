package benchmark_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/labstack/echo"
	"github.com/payfazz/go-router/path"
	"github.com/payfazz/go-router/segment"

	"github.com/gorilla/mux"

	"github.com/gin-gonic/gin"
)

var (
	hGoRouter   http.HandlerFunc
	hGorillaMux http.HandlerFunc
	hGin        http.HandlerFunc
	hEcho       http.HandlerFunc

	w *httptest.ResponseRecorder
	r *http.Request
)

func init() {
	func() {
		hGoRouter = path.H{
			"/a":       respWriter([]byte("1")),
			"/b/c":     respWriter([]byte("2")),
			"/c/d/e":   respWriter([]byte("3")),
			"/d/e/f":   respWriter([]byte("4")),
			"/f/g/h/i": respWriter([]byte("5")),
			"/g/:h/i/:j/k": func(w http.ResponseWriter, r *http.Request) {
				h, _ := segment.Get(r, "h")
				hi, _ := strconv.Atoi(h)
				j, _ := segment.Get(r, "j")
				ji, _ := strconv.Atoi(j)
				fmt.Fprint(w, hi+ji)
			},
		}.C()
	}()

	func() {
		hh := mux.NewRouter()
		hh.HandleFunc("/a", respWriter([]byte("1")))
		hh.HandleFunc("/b/c", respWriter([]byte("2")))
		hh.HandleFunc("/c/d/e", respWriter([]byte("3")))
		hh.HandleFunc("/d/e/f", respWriter([]byte("4")))
		hh.HandleFunc("/f/g/h/i", respWriter([]byte("5")))
		hh.HandleFunc("/g/{h}/i/{j}/k", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			h := vars["h"]
			hi, _ := strconv.Atoi(h)
			j := vars["j"]
			ji, _ := strconv.Atoi(j)
			fmt.Fprint(w, hi+ji)
		})
		hGorillaMux = hh.ServeHTTP
	}()

	func() {
		gin.SetMode(gin.ReleaseMode)

		hh := gin.New()
		hh.GET("/a", respWriterGin([]byte("1")))
		hh.GET("/b/c", respWriterGin([]byte("2")))
		hh.GET("/c/d/e", respWriterGin([]byte("3")))
		hh.GET("/d/e/f", respWriterGin([]byte("4")))
		hh.GET("/f/g/h/i", respWriterGin([]byte("5")))
		hh.GET("/g/:h/i/:j/k", func(c *gin.Context) {
			h := c.Param("h")
			hi, _ := strconv.Atoi(h)
			j := c.Param("j")
			ji, _ := strconv.Atoi(j)
			fmt.Fprint(c.Writer, hi+ji)
		})

		hGin = hh.ServeHTTP
	}()

	func() {
		hh := echo.New()
		hh.GET("/a", respWriterEcho([]byte("1")))
		hh.GET("/b/c", respWriterEcho([]byte("2")))
		hh.GET("/c/d/e", respWriterEcho([]byte("3")))
		hh.GET("/d/e/f", respWriterEcho([]byte("4")))
		hh.GET("/f/g/h/i", respWriterEcho([]byte("5")))
		hh.GET("/g/:h/i/:j/k", func(c echo.Context) error {
			h := c.Param("h")
			hi, _ := strconv.Atoi(h)
			j := c.Param("j")
			ji, _ := strconv.Atoi(j)
			fmt.Fprint(c.Response().Writer, hi+ji)
			return nil
		})

		hEcho = hh.ServeHTTP
	}()

	w = httptest.NewRecorder()
	w.Body.Grow(1024)

	r = &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Host:       "localhost",
	}
}

func respWriter(data []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}
}

func respWriterGin(data []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Write(data)
	}
}

func respWriterEcho(data []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Writer.Write(data)
		return nil
	}
}

func check(b *testing.B, h http.HandlerFunc, path string, data []byte) {
	w.Body.Reset()
	r.URL.Path = path
	h(w, r)
	if !bytes.Equal(w.Body.Bytes(), data) {
		b.Fail()
	}
}

func BenchmarkGoRouter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		check(b, hGoRouter, "/a", []byte("1"))
		check(b, hGoRouter, "/b/c", []byte("2"))
		check(b, hGoRouter, "/c/d/e", []byte("3"))
		check(b, hGoRouter, "/d/e/f", []byte("4"))
		check(b, hGoRouter, "/f/g/h/i", []byte("5"))
		check(b, hGoRouter, "/g/3/i/3/k", []byte("6"))
	}
}

func BenchmarkGorillaMux(b *testing.B) {
	for i := 0; i < b.N; i++ {
		check(b, hGorillaMux, "/a", []byte("1"))
		check(b, hGorillaMux, "/b/c", []byte("2"))
		check(b, hGorillaMux, "/c/d/e", []byte("3"))
		check(b, hGorillaMux, "/d/e/f", []byte("4"))
		check(b, hGorillaMux, "/f/g/h/i", []byte("5"))
		check(b, hGorillaMux, "/g/3/i/3/k", []byte("6"))
	}

}

func BenchmarkGin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		check(b, hGin, "/a", []byte("1"))
		check(b, hGin, "/b/c", []byte("2"))
		check(b, hGin, "/c/d/e", []byte("3"))
		check(b, hGin, "/d/e/f", []byte("4"))
		check(b, hGin, "/f/g/h/i", []byte("5"))
		check(b, hGin, "/g/3/i/3/k", []byte("6"))
	}
}

func BenchmarkEcho(b *testing.B) {
	for i := 0; i < b.N; i++ {
		check(b, hEcho, "/a", []byte("1"))
		check(b, hEcho, "/b/c", []byte("2"))
		check(b, hEcho, "/c/d/e", []byte("3"))
		check(b, hEcho, "/d/e/f", []byte("4"))
		check(b, hEcho, "/f/g/h/i", []byte("5"))
		check(b, hEcho, "/g/3/i/3/k", []byte("6"))
	}
}
