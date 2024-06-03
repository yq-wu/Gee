package gee

import (
	"fmt"
	"net/http"
)

type HandleFunc func(w http.ResponseWriter, req *http.Request)

type engine struct {
	route map[string]HandleFunc
}

func (e *engine) addRoute(method string, patten string, handle HandleFunc) {
	e.route[method+"-"+patten] = handle
}

func New() *engine {
	return &engine{
		route: make(map[string]HandleFunc),
	}
}

func (e *engine) GET(patten string, handle HandleFunc) {
	e.addRoute("GET", patten, handle)
}

func (e *engine) POST(patten string, handle HandleFunc) {
	e.addRoute("POST", patten, handle)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handle, ok := e.route[key]; ok {
		handle(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func (e *engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
