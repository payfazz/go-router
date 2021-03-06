package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/payfazz/go-router/segment/shifter"
)

func main() {
	panic(http.ListenAndServe(":8080", http.HandlerFunc(handler)))
}

func handler(w http.ResponseWriter, r *http.Request) {
	shifter, r := shifterFor(r)
	segment, end := shifter.Shift()
	fmt.Fprintf(w, "segment=%v, end=%v\n", segment, end)
	if !end {
		handler(w, r)
	} else {
		shifter.TagRelative(0, "latest")
		shifter.Unshift()
		shifter.TagRelative(0, "latest-1")
		shifter.Unshift()
		shifter.TagRelative(0, "latest-2")
		shifter.Unshift()
		fmt.Fprintln(w)
		endHandler(w, r)
	}
}

func endHandler(w http.ResponseWriter, r *http.Request) {
	var shifter *shifter.Shifter
	shifter, r = shifterFor(r)
	_, end := shifter.Shift()
	fmt.Fprintf(w, "segment=%v, end=%v\n", shifter.Get(shifter.CurrentIndex()), end)
	for i := -2; i <= 2; i++ {
		fmt.Fprintf(w, "GetRelative(%d)=%s\n", i, shifter.GetRelative(i))
	}
	fmt.Fprintln(w)
	type stat struct {
		value string
		ok    bool
	}
	get := func(s string) stat {
		val, ok := shifter.GetByTag(s)
		return stat{val, ok}
	}
	ll := []stat{
		get("latest"),
		get("latest-1"),
		get("latest-2"),
	}
	fmt.Fprintf(w, "latest=%v, latest-1=%v, latest-2=%v\n", ll[0], ll[1], ll[2])
}

type ctxKeyType struct{}

var ctxKey ctxKeyType

// shifterFor .
func shifterFor(r *http.Request) (*shifter.Shifter, *http.Request) {
	tmp := r.Context().Value(ctxKey)
	if tmp != nil {
		return tmp.(*shifter.Shifter), r
	}
	list := strings.Split(
		strings.TrimPrefix(r.URL.EscapedPath(), "/"), "/",
	)

	s := shifter.New(list)
	r = r.WithContext(context.WithValue(
		r.Context(), ctxKey, s,
	))
	return s, r
}
