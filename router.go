package ifviva

import (
	"net/http"
)

type Handler func(Context)

type route struct {
	method  string
	url     string
	handler Handler
}

type Router struct {
	routes []route
}

type Context struct {
	Req *http.Request
	Res http.ResponseWriter
}

func (r *Router) All(url string, h Handler) {
	r.routes = append(r.routes, route{"all", url, h})
}

func (r *Router) Match(url string) (Handler, error) {
	return r.routes[0].handler, nil
}
