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

	panic(http.ListenAndServe(":8080", path.Compile(path.H{
		"/":          gen("1"),
		"/:key":      gen("2"),
		"/asdf/":     gen("3"),
		"/lala/:key": gen("4"),
		"/lala/:key/aa": func(w http.ResponseWriter, r *http.Request) {
			key, _ := segment.Get(r, "key")
			fmt.Fprintln(w, key)
		},
	}, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "MY NOT FOUND HANDLER")
	})))

}