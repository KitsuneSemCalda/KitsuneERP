package router

import "net/http"

type Route struct {
	Method  string
	Pattern string
	Handler http.Handler
}
