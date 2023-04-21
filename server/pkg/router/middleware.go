package router

import "net/http"

type Middleware interface {
	Execute(http.ResponseWriter, *http.Request) error
	Next(Middleware)
}
