package benchmark_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
)

func init() {
	func() {
		hGoRouter = path.H{
			"/a":       respWriter("1"),
			"/b/c":     respWriter("2"),
			"/c/d/e":   respWriter("3"),
			"/d/e/f":   respWriter("4"),
			"/f/g/h/i": respWriter("5"),
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
		hh.HandleFunc("/a", respWriter("1"))
		hh.HandleFunc("/b/c", respWriter("2"))
		hh.HandleFunc("/c/d/e", respWriter("3"))
		hh.HandleFunc("/d/e/f", respWriter("4"))
		hh.HandleFunc("/f/g/h/i", respWriter("5"))
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
		hh.GET("/a", respWriterGin("1"))
		hh.GET("/b/c", respWriterGin("2"))
		hh.GET("/c/d/e", respWriterGin("3"))
		hh.GET("/d/e/f", respWriterGin("4"))
		hh.GET("/f/g/h/i", respWriterGin("5"))
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
		hh.GET("/a", respWriterEcho("1"))
		hh.GET("/b/c", respWriterEcho("2"))
		hh.GET("/c/d/e", respWriterEcho("3"))
		hh.GET("/d/e/f", respWriterEcho("4"))
		hh.GET("/f/g/h/i", respWriterEcho("5"))
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
}

func respWriter(text string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(text))
	}
}

func respWriterGin(text string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Write([]byte(text))
	}
}

func respWriterEcho(text string) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Writer.Write([]byte(text))
		return nil
	}
}

func check(b *testing.B, h http.HandlerFunc, path string, data string) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", path, nil)
	h(w, r)
	if w.Body.String() != data {
		b.Fail()
	}
}

func BenchmarkGoRouter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		check(b, hGoRouter, "/a", "1")
		check(b, hGoRouter, "/b/c", "2")
		check(b, hGoRouter, "/c/d/e", "3")
		check(b, hGoRouter, "/d/e/f", "4")
		check(b, hGoRouter, "/f/g/h/i", "5")
		check(b, hGoRouter, "/g/3/i/3/k", "6")
	}
}

func BenchmarkGorillaMux(b *testing.B) {
	for i := 0; i < b.N; i++ {
		check(b, hGorillaMux, "/a", "1")
		check(b, hGorillaMux, "/b/c", "2")
		check(b, hGorillaMux, "/c/d/e", "3")
		check(b, hGorillaMux, "/d/e/f", "4")
		check(b, hGorillaMux, "/f/g/h/i", "5")
		check(b, hGorillaMux, "/g/3/i/3/k", "6")
	}

}

func BenchmarkGin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		check(b, hGin, "/a", "1")
		check(b, hGin, "/b/c", "2")
		check(b, hGin, "/c/d/e", "3")
		check(b, hGin, "/d/e/f", "4")
		check(b, hGin, "/f/g/h/i", "5")
		check(b, hGin, "/g/3/i/3/k", "6")
	}
}

func BenchmarkEcho(b *testing.B) {
	for i := 0; i < b.N; i++ {
		check(b, hEcho, "/a", "1")
		check(b, hEcho, "/b/c", "2")
		check(b, hEcho, "/c/d/e", "3")
		check(b, hEcho, "/d/e/f", "4")
		check(b, hEcho, "/f/g/h/i", "5")
		check(b, hEcho, "/g/3/i/3/k", "6")
	}
}
