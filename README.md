# go-router

Fast golang router.

Internally it compiled into tree structure (nested `map[string]interface{}`), see `examples/2_segment_and_method/main.go`.

## TODO

* More documentation and example
* Implement path based routing with segment based routing as building block. Segment based block is pretty hard to use, see `examples/2_segment_and_method`. Path based routing will look like:
    ```go
    r := path.Compile(path.H{
        "/a/b": handler1,
        "/asdf/:paramhere/asge": handler2,
        "/asdf/:paramhere/:param2here/gege": handler3,
    }, nil)
    ```
