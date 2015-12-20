package ifviva

import (
	"fmt"
	"net/http"
	"regexp"
)

type Handler func(Context)

type route struct {
	method  string
	url     string
	handler Handler
	regex   *regexp.Regexp
}

func (r route) MatchMethod(method string) bool {
	return r.method == "*" || method == r.method || (method == "HEAD" && r.method == "GET")
}

func (r route) Match(method string, url string) (bool, map[string]string) {
	if !r.MatchMethod(method) {
		return false, nil
	}

	matches := r.regex.FindStringSubmatch(url)
	if len(matches) > 0 && matches[0] == url {
		params := make(map[string]string)
		for i, name := range r.regex.SubexpNames() {
			if len(name) > 0 {
				params[name] = matches[i]
			}
		}
		return true, params
	}
	return false, nil
}

type Router struct {
	routes []route
}

type Context struct {
	Req    *http.Request
	Res    http.ResponseWriter
	Params map[string]string
}

func (r *Router) addRoute(method string, url string, handler Handler) {
	reg := regexp.MustCompile(`:[^/#?()\.\\]+`)
	pattern := reg.ReplaceAllStringFunc(url, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})
	regDir := regexp.MustCompile(`\*\*`)
	var index int
	pattern = regDir.ReplaceAllStringFunc(pattern, func(m string) string {
		index++
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, index)
	})
	pattern += `\/?`
	route := route{
		method:  method,
		url:     url,
		handler: handler,
		regex:   regexp.MustCompile(pattern),
	}
	r.routes = append(r.routes, route)
}

func (r *Router) All(url string, h Handler) {
	r.addRoute("*", url, h)
}

func (r *Router) Post(url string, h Handler) {
	r.addRoute("POST", url, h)
}

func (r *Router) Get(url string, h Handler) {
	r.addRoute("GET", url, h)
}

func (r *Router) Put(url string, h Handler) {
	r.addRoute("PUT", url, h)
}

func (r *Router) Patch(url string, h Handler) {
	r.addRoute("PATCH", url, h)
}

func (r *Router) Delete(url string, h Handler) {
	r.addRoute("DELETE", url, h)
}

func (r *Router) Head(url string, h Handler) {
	r.addRoute("HEAD", url, h)
}

func (r *Router) Options(url string, h Handler) {
	r.addRoute("OPTIONS", url, h)
}

func (r *Router) Match(method string, url string) (bool, Handler, map[string]string) {
	for _, route := range r.routes {
		if match, params := route.Match(method, url); match == true {
			return true, route.handler, params
		}
	}
	return false, nil, nil
}
