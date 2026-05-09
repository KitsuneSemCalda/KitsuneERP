package router

func (r *Router) Use(mws ...Middleware) {
	r.middlewares = append(r.middlewares, mws...)
}
