package gee

import (
	"fmt"
)

type Router struct {
	router map[string]HandleFunc
}

func NewRouter() *Router {
	return &Router{
		router: make(map[string]HandleFunc),
	}
}

func (r *Router) AddRoute(method string, patten string, handle HandleFunc) {
	r.router[method+"-"+patten] = handle
}

func (r *Router) Handle(c *Context) {
	key := c.Req.Method + "-" + c.Req.URL.Path
	if handle, ok := r.router[key]; ok {
		handle(c)
	} else {
		fmt.Fprintf(c.Writer, "404 NOT FOUND: %s\n", c.Path)
	}
}
