package main

import (
	"fmt"
	"net/http"

	"github.com/payfazz/go-router/path"
	"github.com/payfazz/go-router/segment"
)

func main() {
	gen := func(str string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, str)
		}
	}

	panic(http.ListenAndServe(":8080", path.H{
		"/":          gen("1"),
		"/:key":      gen("2"),
		"/asdf/":     gen("3"),
		"/lala/:key": gen("4"),
		"/lala/:key/aa": func(w http.ResponseWriter, r *http.Request) {
			key, _ := segment.Get(r, "key")
			fmt.Fprintln(w, "5")
			fmt.Fprintln(w, key)
		},
		"/test-slash":   path.WithTrailingSlash(gen("aa")),
		"/test-slash-2": path.WithoutTrailingSlash(gen("bb")),
	}.Compile(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "MY NOT FOUND HANDLER")
		done, rest := segment.Split(r)
		fmt.Fprintln(w, done, rest)
	})))

}
