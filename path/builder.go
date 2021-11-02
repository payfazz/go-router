package path

import (
	"net/http"
	"strings"

	internalsegment "github.com/payfazz/go-router/internal/segment"
	segmentpkg "github.com/payfazz/go-router/segment"
)

type pathPart interface{ pathPartTag() }
type segment string
type param string

func (segment) pathPartTag() {}
func (param) pathPartTag()   {}

type handlerPart interface{ handlerPartTag() }
type handler http.HandlerFunc
type tree map[pathPart]handlerPart

func (handler) handlerPartTag() {}
func (tree) handlerPartTag()    {}

func (root tree) add(path string, h http.HandlerFunc) {
	current := root
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	paths := strings.Split(path, "/")

	// for intermediate segment
	for i := 0; i < len(paths)-1; i++ {
		p := paths[i]
		var item pathPart
		if strings.HasPrefix(p, ":") {
			item = param(p[1:])
		} else {
			item = segment(p)
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
			panic("path: Non-exhaustive switch")
		}
	}

	// for last segment
	if len(paths) > 0 {
		p := paths[len(paths)-1]
		var item pathPart
		if strings.HasPrefix(p, ":") {
			item = param(p[1:])
		} else {
			item = segment(p)
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
				panic("path: Non-exhaustive switch")
			}
		}
	}
}

func (root tree) compile(def http.HandlerFunc) http.HandlerFunc {
	return compile(root, def, 0)
}

func compile(root handlerPart, def http.HandlerFunc, count int) http.HandlerFunc {
	switch root := root.(type) {
	case handler:
		return http.HandlerFunc(root)
	case tree:
		hMap := make(segmentpkg.H)
		var paramHandler http.HandlerFunc
		var paramTag string

		if item, ok := root[segment("")]; ok {
			tmp := compile(item, def, count+1)
			def = func(w http.ResponseWriter, rOld *http.Request) {
				s, r := internalsegment.TryShifterFrom(rOld)
				cur := s.CurrentIndex()
				for i := count; i <= cur; i++ {
					s.ClearTag(i)
				}
				s.SetNext(count)
				tmp(w, r)
			}
		}

		for key, item := range root {
			switch key := key.(type) {
			case segment:
				hMap[string(key)] = compile(item, def, count+1)
			case param:
				switch paramHandler {
				case nil:
					paramTag = string(key)
					paramHandler = segmentpkg.Tag(paramTag, compile(item, def, count+1))
				default:
					if paramTag != string(key) {
						panic("path: multiple param name, :" + paramTag + " != :" + string(key))
					}
				}
			default:
				panic("path: Non-exhaustive switch")
			}
		}

		if paramHandler == nil {
			paramHandler = def
		}

		return hMap.Compile(paramHandler)
	default:
		panic("path: Non-exhaustive switch")
	}
}
