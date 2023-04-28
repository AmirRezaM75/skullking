package router

import (
	"context"
	"net/http"
	"regexp"
)

type Router struct {
	routes      []*Route
	middlewares []Middleware
}

func NewRouter() *Router {
	return &Router{}
}

func (router *Router) Middleware(m Middleware) *Router {
	router.middlewares = append(router.middlewares, m)
	return router
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range router.routes {
		params := route.match(r)

		if params == nil {
			continue
		}

		// Execute global middlewares
		for _, m := range router.middlewares {
			err := m.Execute(w, r)
			if err != nil {
				return
			}
		}

		// Handles preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ctx := context.WithValue(r.Context(), "params", params)
		r = r.WithContext(ctx)

		// Execute route middlewares
		for _, m := range route.middlewares {
			err := m.Execute(w, r)
			if err != nil {
				return
			}
		}

		route.handler(w, r)

		return
	}

	http.NotFound(w, r)
}

func (router *Router) handle(path, method string, handler http.HandlerFunc) *Route {
	r := regexp.MustCompile(`:(\w+)`)
	path = r.ReplaceAllString(path, "(?P<$1>(?:=|\\w)+)")

	route := &Route{
		path:    path,
		method:  method,
		handler: handler,
	}

	router.routes = append(router.routes, route)

	return route
}

func (router *Router) Get(path string, handler http.HandlerFunc) *Route {
	return router.handle(path, "GET", handler)
}

func (router *Router) Post(path string, handler http.HandlerFunc) *Route {
	return router.handle(path, "POST", handler)
}

func (router *Router) Patch(path string, handler http.HandlerFunc) *Route {
	return router.handle(path, "PATCH", handler)
}

func (router *Router) Put(path string, handler http.HandlerFunc) *Route {
	return router.handle(path, "PUT", handler)
}
