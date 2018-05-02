package main

import (
	"fmt"
	"net/http"

	"github.com/payfazz/go-router/method"

	"github.com/payfazz/go-router/segment"
)

func main() {
	r := segment.Compile(segment.H{
		"": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "(1) ALL /")
		},
		"test": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "(2) ALL /test")
		},
		"lala": segment.Compile(nil,
			segment.Compile(segment.H{
				"hai": func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintf(w, "(3) ALL /lala/%s/hai\n", segment.GetRelative(r, -1))
				},
			}, method.Compile(method.H{
				http.MethodPost: func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, "(4) POST /lala")
				},
				http.MethodDelete: func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, "(5) DELETE /lala")
				},
			}, nil)),
		),
		"asdf": segment.Compile(segment.H{
			"haha": func(w http.ResponseWriter, r *http.Request) {
				done, rest := segment.Split(r)
				fmt.Fprintf(w, "(6) ALL /%s/%s\n", segment.GetRelative(r, -1), segment.Get(r))
				fmt.Fprintf(w, "%s\n", done)
				fmt.Fprintf(w, "%s\n", rest)
			},
		}, nil),
	}, func(w http.ResponseWriter, r *http.Request) {
		done, rest := segment.Split(r)
		fmt.Fprintf(w, "(7) ALL /%s\n", segment.Get(r))
		fmt.Fprintf(w, "%s\n", done)
		fmt.Fprintf(w, "%s\n", rest)
	})

	/*
		some feature that not implemented yet

		r := host.Compile(host.H{
			"example.com": exampleHandler(),
			"sub.example.com": subExampleHandler(),
		}, nil)

		r := path.Compile(path.H{
			"/": rootHandler(),
			"/test": testHandler(),
			"/lala/:name/hai": haiHandler(),
		}, nil)
	*/

	panic(http.ListenAndServe(":8080", r))
}
