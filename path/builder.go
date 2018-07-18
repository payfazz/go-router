package path

import (
	"fmt"
	"net/http"
	"strings"

	segmentpkg "github.com/payfazz/go-router/segment"
)

type segment string
type param string

type handler http.HandlerFunc

type tree map[interface{}]interface{} // map[(segment or param)](tree or handler)

type builderT struct {
	root tree
}

func (builder *builderT) add(path string, h http.HandlerFunc) {
	current := builder.root
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	paths := strings.Split(path, "/")
	for i, segStr := range paths {
		var item interface{}
		if strings.HasPrefix(segStr, ":") {
			item = param(segStr[1:])
		} else {
			item = segment(segStr)
		}
		if i != len(paths)-1 {
			next, ok := current[item]
			if !ok {
				newTree := make(tree)
				current[item] = newTree
				next = newTree
			}
			switch next := next.(type) {
			case tree:
				current = next
			case handler:
				newTree := tree{segment(""): next}
				current[item] = newTree
				current = newTree
			default:
				panic("path: (BUG) invalid tree")
			}
		} else {
			next, ok := current[item]
			if !ok {
				current[item] = handler(h)
			} else {
				switch next := next.(type) {
				case tree:
					next[segment("")] = handler(h)
				case handler:
					if path == "" {
						path = "/"
					} else {
						path = "/" + path + "/"
					}
					panic("path: duplicate handler: " + path)
				default:
					panic("path: (BUG) invalid tree")
				}
			}
		}
	}
}

func (builder *builderT) compile(def http.HandlerFunc) http.HandlerFunc {
	return compile(builder.root, def)
}

func compile(root interface{}, def http.HandlerFunc) http.HandlerFunc {
	switch root := root.(type) {
	case handler:
		return http.HandlerFunc(root)
	case tree:
		hMap := make(segmentpkg.H)
		var paramHandler http.HandlerFunc
		var paramTag string
		for key, item := range root {
			switch key := key.(type) {
			case segment:
				hMap[string(key)] = compile(item, def)
			case param:
				switch paramHandler {
				case nil:
					paramTag = string(key)
					paramHandler = segmentpkg.Tag(paramTag, compile(item, def))
				default:
					if paramTag != string(key) {
						panic(fmt.Sprintf("path: multiple param name, :%s != :%s", paramTag, string(key)))
					}
				}
			default:
				panic("path: (BUG) invalid tree")
			}

		}
		if paramHandler == nil {
			paramHandler = def
		}
		return segmentpkg.Compile(hMap, paramHandler)
	default:
		panic("path: (BUG) invalid tree")
	}
}
