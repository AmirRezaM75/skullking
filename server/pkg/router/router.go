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

		ctx := context.WithValue(r.Context(), "params", params)
		r = r.WithContext(ctx)

		// It merges global middlewares and route middlewares
		// To make it easier to loop through in reverse order.
		var middlewares []Middleware
		middlewares = append(middlewares, router.middlewares...)
		middlewares = append(middlewares, route.middlewares...)

		h := http.Handler(route.handler)

		for i := range middlewares {
			h = middlewares[len(middlewares)-1-i].Handle(h)
		}

		h.ServeHTTP(w, r)

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
