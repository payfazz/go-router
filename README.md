# go-router

[![GoDoc](https://godoc.org/github.com/payfazz/go-router?status.svg)](https://godoc.org/github.com/payfazz/go-router)

Fast golang router.

This project is about simple http router that preserve `http.Handler` and `http.HanderFunc` signature.

Because every routing strategy (`path.Compile`, `method.Compile`, `segment.Compile`) will produce `http.HandlerFunc` (which is also `http.Handler`), it can be mix with another router (make sure to strip prefix when using `path`).

If you like to manually control your routing (using `if` of `switch`), `segment/shifter` package may help you.

It heavily use clousure and tail call, so it will be faster when tail-cail-optimization implemented on golang. The routing decission tree is precompute, so it should be faster.

for usage see examples directory

see also https://github.com/payfazz/go-middleware for middleware

see also https://gist.github.com/win-t/8a243301bd227cca6135374cf94d9e98 for example usage of go-middleware and go-router

## TODO

* More documentation and examples
* create testing. when i wrote this, this project was part of bigger project, all of the testing was done there.
