package router

import "net/http"

type Router struct {
	mux         *http.ServeMux
	prefix      string
	middlewares []Middleware
	routes      []Route
}

func New() *Router {
	return &Router{
		mux:    http.NewServeMux(),
		routes: make([]Route, 0),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) Get(pattern string, h http.Handler) {
	r.handle(http.MethodGet, pattern, h)
}

func (r *Router) Post(pattern string, h http.Handler) {
	r.handle(http.MethodPost, pattern, h)
}

func (r *Router) Put(pattern string, h http.Handler) {
	r.handle(http.MethodPut, pattern, h)
}

func (r *Router) Delete(pattern string, h http.Handler) {
	r.handle(http.MethodDelete, pattern, h)
}

func (r *Router) Patch(pattern string, h http.Handler) {
	r.handle(http.MethodPatch, pattern, h)
}

func (r *Router) Static(pattern, dir string) {
	fullPath := join(r.prefix, pattern)

	fs := http.StripPrefix(
		fullPath,
		http.FileServer(http.Dir(dir)),
	)

	r.handle(http.MethodGet, fullPath+"/", fs)
}

func (r *Router) handle(method, pattern string, h http.Handler) {
	fullPath := join(r.prefix, pattern)

	fullPattern := fullPath

	if method != "" {
		fullPattern = method + " " + fullPath
	}

	handler := chain(h, r.middlewares)

	r.routes = append(r.routes, Route{
		Method:  method,
		Pattern: fullPath,
		Handler: handler,
	})

	r.mux.Handle(fullPattern, handler)
}
