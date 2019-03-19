package path

import (
	"net/http"
	"strings"

	segmentpkgInternal "github.com/payfazz/go-router/internal/segment"
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

	for i := 0; i < len(paths)-1; i++ {
		segStr := paths[i]
		var item interface{}
		if strings.HasPrefix(segStr, ":") {
			item = param(segStr[1:])
		} else {
			item = segment(segStr)
		}
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
			panic("path: DEADCODE: (BUG) invalid tree")
		}
	}

	segStr := paths[len(paths)-1]
	var item interface{}
	if strings.HasPrefix(segStr, ":") {
		item = param(segStr[1:])
	} else {
		item = segment(segStr)
	}
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
			panic("path: DEADCODE: (BUG) invalid tree")
		}
	}
}

func (builder *builderT) compile(def http.HandlerFunc) http.HandlerFunc {
	return builderCompile(builder.root, def, 0)
}

func builderCompile(root interface{}, def http.HandlerFunc, count int) http.HandlerFunc {
	switch root := root.(type) {
	case handler:
		return http.HandlerFunc(root)
	case tree:
		hMap := make(segmentpkg.H)
		var paramHandler http.HandlerFunc
		var paramTag string

		if item, ok := root[segment("")]; ok {
			tmp := builderCompile(item, def, count+1)
			def = func(w http.ResponseWriter, r *http.Request) {
				segmentpkgInternal.SetShifterIndex(r, count)
				tmp(w, r)
			}
		}

		for key, item := range root {
			switch key := key.(type) {
			case segment:
				hMap[string(key)] = builderCompile(item, def, count+1)
			case param:
				switch paramHandler {
				case nil:
					paramTag = string(key)
					paramHandler = segmentpkg.Tag(paramTag, builderCompile(item, def, count+1))
				default:
					if paramTag != string(key) {
						panic("path: multiple param name, :" + paramTag + " != :" + string(key))
					}
				}
			default:
				panic("path: DEADCODE: (BUG) invalid tree")
			}
		}

		if paramHandler == nil {
			paramHandler = def
		}

		return hMap.Compile(paramHandler)
	default:
		panic("path: DEADCODE: (BUG) invalid tree")
	}
}
