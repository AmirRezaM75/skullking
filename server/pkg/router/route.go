package router

import (
	"net/http"
	"regexp"
)

type Route struct {
	path        string
	method      string
	handler     http.HandlerFunc
	middlewares []Middleware
}

func (route *Route) Middleware(m Middleware) *Route {
	route.middlewares = append(route.middlewares, m)
	return route
}

func (route Route) match(request *http.Request) map[string]string {
	if request.Method != http.MethodOptions && route.method != request.Method {
		// TODO: 405
		return nil
	}

	r := regexp.MustCompile("^" + route.path + "$")
	matches := r.FindStringSubmatch(request.URL.Path)

	if matches == nil {
		return nil
	}

	params := make(map[string]string)

	names := r.SubexpNames()

	for i, match := range matches {
		// SubexpNames returns the names of the parenthesized subexpressions in this Regexp.
		// Since the Regexp as a whole cannot be named, names[0] is always the empty string.
		if i == 0 {
			continue
		}
		params[names[i]] = match
	}

	return params
}

func (route Route) chainMiddlewares() {
	for k, m := range route.middlewares {
		if len(route.middlewares)-2 == k {
			break
		}

		m.Next(route.middlewares[k+1])
	}
}

func (route Route) runMiddlewares(w http.ResponseWriter, r *http.Request) error {
	if len(route.middlewares) == 0 {
		return nil
	}

	if len(route.middlewares) > 1 {
		route.chainMiddlewares()
	}

	m := route.middlewares[0]
	return m.Execute(w, r)
}
