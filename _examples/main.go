package main

import (
	"fmt"
	"net/http"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-router/defhandler"
	"github.com/payfazz/go-router/method"
	"github.com/payfazz/go-router/path"
	"github.com/payfazz/go-router/segment"
)

func main() {
	panic(http.ListenAndServe(":8080", path.H{
		"/test":  defhandler.ResponseCodeWithMessage(200, "test"),
		"/hello": hello(),
	}.C()))
}

func hello() http.HandlerFunc {
	return middleware.C(
		method.Must("GET"),
		path.H{
			"/":        helloWorld(),
			"/:key":    helloWithKey(),
			"/:key/ch": helloWithKeyCh(),
		}.C(),
	)
}

func helloWorld() http.HandlerFunc {
	return middleware.C(
		segment.MustEnd,
		defhandler.ResponseCodeWithMessage(200, "Hello World\n"),
	)
}

func helloWithKey() http.HandlerFunc {
	return middleware.C(
		segment.MustEnd,
		func(w http.ResponseWriter, r *http.Request) {
			name := segment.Param(r, "key")
			fmt.Fprintf(w, "Hello %s\n", name)
		},
	)
}

func helloWithKeyCh() http.HandlerFunc {
	return middleware.C(
		segment.MustEnd,
		func(w http.ResponseWriter, r *http.Request) {
			name := segment.Param(r, "key")
			fmt.Fprintf(w, "Hello, 世界, %s\n", name)
		},
	)
}
