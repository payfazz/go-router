package main

import (
	"fmt"
	"net/http"

	"github.com/payfazz/go-router/method"
	"github.com/payfazz/go-router/segment"
)

func main() {
	panic(http.ListenAndServe(":8080", segment.Compile(segment.H{
		"": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "(1) ALL /")
		},
		"test": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "(2) ALL /test | /test/**")
		},
		"test2": segment.C(segment.H{
			"": func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "(2) ALL /test2 | /test2/")
			},
			"aa": func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "(3) ALL /test2/aa | /test2/aa/**")
			},
			"bb": func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "(4) ALL /test2/bb | /test2/bb/**")
			},
		}),
		"lala": segment.Tag("mytag", func(w http.ResponseWriter, r *http.Request) {
			s, ok := segment.Get(r, "mytag")
			if ok {
				fmt.Fprintf(w, "(5) ALL /lala/%v | /lala/%v/**\n", s, s)
			} else {
				fmt.Fprintln(w, "(5) ALL /lala")
			}
		}),
		"lala2": segment.Tag("mytag", segment.C(segment.H{
			"": func(w http.ResponseWriter, r *http.Request) {
				s, ok := segment.Get(r, "mytag")
				if ok {
					fmt.Fprintf(w, "(6) ALL /lala2/%v | /lala/%v/\n", s, s)
				} else {
					fmt.Fprintln(w, "(6) ALL /lala2")
				}
			},
			"aa": func(w http.ResponseWriter, r *http.Request) {
				s, _ := segment.Get(r, "mytag")
				fmt.Fprintf(w, "(7) ALL /test2/%v/aa | /test2/%v/aa/**\n", s, s)
			},
			"bb": func(w http.ResponseWriter, r *http.Request) {
				s, _ := segment.Get(r, "mytag")
				fmt.Fprintf(w, "(8) ALL /test2/%v/bb | /test2/%v/bb/**\n", s, s)
			},
		})),
		"lala3": segment.Stripper(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "(9) ALL /lala3 | /lala3/**")
			fmt.Fprintf(w, "New URL = %v\n", r.URL.String())
		}),
	}, method.C(method.H{
		http.MethodGet: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "(10) GET *default")
			fmt.Fprintf(w, "%#v\n", segment.Rest(r))
		},
		http.MethodPost: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "(11) POST *default")
			fmt.Fprintf(w, "%#v\n", segment.Rest(r))
		},
	}))))
}
