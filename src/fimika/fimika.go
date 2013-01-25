package fimika

import (
	"fmt"
	"regexp"
	"strings"
	"net/http"
	"appengine"
)

type Application struct {
	Routes		[]*Route
}

type Request struct {
	*http.Request
	Params		map[string][]string
	Host		string
	Path		string
	Fragment	string
}

type Response struct {
	http.ResponseWriter
}

type Context struct {
	Method		string
	Params		map[string][]string
	Form		map[string][]string
	Request 	*Request
	Response	*Response
	Context		*appengine.Context
	View		map[string]interface{}
}

type Route struct {
	Pattern		string
	Regexp		*regexp.Regexp
	Handler		Handler
}

type RouteData struct {
	Path		string
	Route		*Route
	Params		map[string]string
}

type Handler func(*Context) error

var (
	ApplicationInstance = new(Application)
	routeParam *regexp.Regexp = regexp.MustCompile("<[^>]+>")
)

func init() {
	http.HandleFunc("/", defaultHandler)
}

// match regexp with string, and return a named group map
// Example:
//   regexp: "(?P<name>[A-Za-z]+)-(?P<age>\\d+)"
//   string: "CGC-30"
//   return: map[string]string{ "name":"CGC", "age":"30" }
func NamedRegexpGroup(str string, reg *regexp.Regexp) map[string]string {
	rst := reg.FindStringSubmatch(str)
	len_rst := len(rst)
	if len_rst <= 0 {
		return nil
	}
	ng := make(map[string]string)
	sn := reg.SubexpNames()
	for k, v := range sn {
		if k == 0 || v == "" {
			continue
		}
		if k + 1 > len_rst {
			break
		}
		ng[v] = rst[k]
	}
	return ng
}

func NewRequest(r *http.Request) *Request {
	req := &Request{r, r.URL.Query(), r.URL.Host, r.URL.Path, r.URL.Fragment}
	return req
}

func (r *Request) Get(key string) string {
	vs, ok := r.Params[key]
	if !ok || len(vs) <= 0 {
		return ""
	}
	return vs[0]
}

func (r *Request) Gets(key string) []string {
	vs, ok := r.Params[key]
	if !ok {
		return nil
	}
	return vs
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{w}
}

func NewContext(w http.ResponseWriter, r *http.Request, rd *RouteData) *Context {
	req := NewRequest(r)
	res := NewResponse(w)
	c := &Context{
		Method: r.Method,
		Request: req,
		Response: res,
		Params: req.Params,
	}
	for k, vs := range rd.Params {
		c.Params[k] = []string{vs}
	}
	return c
}

func NewRoute(pattern string) *Route {
	r := &Route{Pattern: pattern}
	r.Regexp = compilePattern(r.Pattern)
	return r
}

func compilePattern(pattern string) *regexp.Regexp {
	// /blog/<id>           =>   /blog/(?P<id>[^\?#/]+)
	// /blog/<id:[0-9]+>    =>   /blog/(?P<id>[0-9]+)
	ret := routeParam.ReplaceAllStringFunc(pattern, func(s string) string {
		var name, reg string
		body := s[1:len(s)-1]
		if (strings.Contains(body, ":")) {
			parts := strings.SplitN(body, ":", 2)
			name, reg = parts[0], parts[1]
		} else {
			name = body
			reg = "[^\\?#/]+"
		}
		return fmt.Sprintf("(?P<%s>%s)", regexp.QuoteMeta(name), reg)
	})
	return regexp.MustCompile("^" + ret + "$")
}

func AddRoute(pattern string, handler Handler) {
	route := NewRoute(pattern)
	route.Handler = handler
	ApplicationInstance.Routes = append(ApplicationInstance.Routes, route)
}

func (r *Route) Match(path string) *RouteData {
	params := NamedRegexpGroup(path, r.Regexp)
	if params == nil {
		return nil
	}
	rd := &RouteData{Path: path, Route: r, Params: params}
	return rd
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	var handler Handler
	var rd *RouteData
	for _, route := range ApplicationInstance.Routes {
		rd = route.Match(r.URL.Path)
		if rd != nil {
			handler = rd.Route.Handler
			break
		}
	}
	if handler != nil {
		c := NewContext(w, r, rd)
		handler(c)
	}
}
