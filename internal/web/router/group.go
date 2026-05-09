package router

func (r *Router) Group(prefix string, fn func(*Router)) {
	sub := &Router{
		mux:         r.mux,
		prefix:      join(r.prefix, prefix),
		middlewares: append([]Middleware(nil), r.middlewares...),
		routes:      r.routes,
	}

	fn(sub)
}
