// Package router provides a basic HTTP-router which wraps
// [http.ServerMux].
package router

import (
	"net/http"
	"path"
	"slices"
	"strings"
)

// Router is similar to [http.ServeMux], except is allows
// registering middleware and creating "sub" routers.
//
// Handler-patterns support the same features and syntax as
// [http.ServeMux].
type Router struct {
	routes     []route
	middleware []func(http.Handler) http.Handler
	subRouters []subRouter
}

type subRouter struct {
	prefix    string
	subrouter *Router
}

type route struct {
	pattern string
	handler http.Handler
}

// New allocates a new [Router].
func New() *Router {
	return &Router{}
}

// New creates a sub-router. All handlers registered to the sub-router
// will use the parent router's middleware. All handlers registered to the
// sub-router will be scoped to the passed-in path-prefix (which may be
// empty).
//
// All parent-middleware will be applied before all sub-router middleware.
func (r *Router) New(prefix string) *Router {
	s := &Router{}
	r.subRouters = append(r.subRouters, subRouter{
		prefix:    prefix,
		subrouter: s,
	})
	return s
}

// Handle registers a handler for a pattern. Patterns support the same
// syntax as [http.ServerMux].
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.routes = append(r.routes, route{pattern: pattern, handler: handler})
}

// HandleFunc registers a handler for a pattern. Patterns support the same
// syntax as [http.ServerMux].
func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.Handle(pattern, http.HandlerFunc(handler))
}

// Use applies a middle-ware function to all handlers (including handlers
// in sub-routers).
func (r *Router) Use(middelware func(http.Handler) http.Handler) {
	r.middleware = append(r.middleware, middelware)
}

// Handler produces an [http.Handler] to server requests using the handlers
// and middleware registered to the router.
func (r *Router) Handler() http.Handler {
	m := http.NewServeMux()
	r.register(m)
	return m
}

func (r *Router) register(m *http.ServeMux) {
	r.registerChildRouter(m, "", [](func(http.Handler) http.Handler){})
}

func (r *Router) registerChildRouter(m *http.ServeMux, prefix string, parentMiddleware []func(http.Handler) http.Handler) {
	middleware := slices.Concat(parentMiddleware, r.middleware)

	for _, route := range r.routes {
		h := applyMiddleware(middleware, route.handler)
		pattern := applyPrefixToPattern(prefix, route.pattern)
		m.Handle(pattern, h)
	}

	for _, s := range r.subRouters {
		subPrefix := path.Join(prefix, s.prefix)
		s.subrouter.registerChildRouter(m, subPrefix, middleware)
	}
}

func applyMiddleware(m []func(http.Handler) http.Handler, h http.Handler) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

func applyPrefixToPattern(prefix string, pattern string) string {
	before, after, found := strings.Cut(pattern, " ")
	if !found {
		return path.Join(prefix, pattern)
	}
	return before + " " + path.Join(prefix, after)
}
