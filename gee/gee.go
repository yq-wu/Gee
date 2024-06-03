package gee

import (
	"net/http"
)

type HandleFunc func(*Context)

type Engine struct {
	route *Router
}

func (e *Engine) addRoute(method string, patten string, handle HandleFunc) {
	e.route.AddRoute(method, patten, handle)
}

func New() *Engine {
	return &Engine{
		route: NewRouter(),
	}
}

func (e *Engine) GET(patten string, handle HandleFunc) {
	e.addRoute("GET", patten, handle)
}

func (e *Engine) POST(patten string, handle HandleFunc) {
	e.addRoute("POST", patten, handle)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := NewContext(w, req)
	e.route.Handle(context)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
