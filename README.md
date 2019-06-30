# go-router

[![GoDoc](https://godoc.org/github.com/payfazz/go-router?status.svg)](https://godoc.org/github.com/payfazz/go-router)

Simple golang router.

This project is about simple http router that preserve `http.Handler` and `http.HanderFunc` signature.

Because every routing strategy (`path.H.Compile`, `method.H.Compile`, `segment.H.Compile`) will produce `http.HandlerFunc` (which is also `http.Handler`), it can be mix with another library.

It heavily use clousure and tail call, so it will be faster when tail-cail-optimization implemented on golang. The routing decission tree is precompute, so it should be faster.

see also https://github.com/payfazz/go-middleware for middleware


### Host based routing
`host` package provide host based routing

example:

```go
http.ListenAndServe(":8080", host.H{
  "payfazz.com":     handler1,
  "api.payfazz.com": handler2,
}.C())
```

### Method based routing
`method` pakcage provide method based routing, by itself method based routing is not useful, but because it generate `http.HandlerFunc` you can mix it with others

example:

```go
http.ListenAndServe(":8080", path.H{
  "/about": handler1,
  "/data":  method.H{
    "GET":  handler2,
    "POST": handler3,
  }.C(),
}.C())
```

### Segment based routing
`segment` provide segment based routing

example:

```go
http.ListenAndServe(":8080", segment.H{
  "a": handler1,
  "b": segment.H{
    "c": handler2,
    "d": handler3,
  }.C(),
}.C())
```

### Path based routing
`path` package is just tool to compose `segment` based routing.
```go
h := path.H{
  "/a":   handler1,
  "/b/c": handler2,
  "/b/d": handler3,
}.C()
```
will be same with
```go
h := segment.H{
  "a": handler1,
  "b": segment.H{
    "c": handler2,
    "d": handler3,
  }.C(),
}.C()
```

`path` also provide parameter in path, because internally `path` use `segment`, this parameter can be access via `segment`
```go
http.ListenAndServe(":8080", path.H{
  "a/:b/c/:d": func(w http.ResponseWriter, r *http.Request) {
    b, _ := segment.Get(r, "b")
    d, _ := segment.Get(r, "d")
    // ...
  },
}.C())
```

## Quick Note
Routing is done by prefix segment matching, so
```go
h := path.H{
  "/a/b": handler1,
}.C()
```
will be still handling request to `/a/b/c/d/e`, if you need to only handle `/a/b` you need to use `segment.MustEnd` middleware
```go
h := path.H{
  "/a/b": segment.MustEnd(handler1),
}.C()
```

This is intentionally, useful for path grouping
```go
func main() {
  http.ListenAndServe(":8080", root())
}

func root() http.HandlerFunc {
  return path.H{
    "/api":  api(),
    "/blog": blog(),
  }.C()
}

func api() http.HandlerFunc {
  return path.H{
    "/order": orderHandler,
    "/user":  userHandler,
  }.C()
}
```
request to `/api/order` will be handled by `orderHandler`

If you use https://github.com/payfazz/go-middleware, you can easily compose with another middleware, for example `method.Must`
```go
h := path.H{
  "/a/b": middleware.C(
    segment.MustEnd,
    method.Must("GET"),
    handler1,
  ),
}.C()
```

## TODO

* More documentation and examples
* create more testing
* make it faster by zero-allocation, current bottleneck:
  * net/http.(*Request).WithContext
  * strings.Split
