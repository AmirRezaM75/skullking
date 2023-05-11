package router

import "net/http"

type Middleware interface {
	Handle(http.Handler) http.Handler
}
