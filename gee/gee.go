package gee

import (
	"net/http"
	"strings"
)

type HandleFunc func(*Context)

type Engine struct {
	*RouterGroup
	route  *Router
	groups []*RouterGroup
}

func (e *Engine) addRoute(method string, patten string, handle HandleFunc) {
	e.route.AddRoute(method, patten, handle)
}

func New() *Engine {
	engine := &Engine{
		route: NewRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) GET(patten string, handle HandleFunc) {
	e.addRoute("GET", patten, handle)
}

func (e *Engine) POST(patten string, handle HandleFunc) {
	e.addRoute("POST", patten, handle)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandleFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleware...)
		}
	}
	context := NewContext(w, req)
	context.handlers = middlewares
	e.route.handle(context)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

type RouterGroup struct {
	prefix     string
	middleware []HandleFunc
	parent     *RouterGroup
	engine     *Engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandleFunc) {
	patten := group.prefix + comp
	group.engine.route.AddRoute(method, patten, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middleware = append(group.middleware, middlewares...)
}
