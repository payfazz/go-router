package router

import "net/http"

// Handler is type alias for http.HandlerFunc
type Handler = http.HandlerFunc

// Hmap is type alias for mapping string to handler
type Hmap = map[string]Handler

// Router is func to get routing state,
// this state will be used for routing decission based on next segment available
type Router func(*http.Request) *State
