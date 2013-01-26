package fimika

import (
	"appengine"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Application struct {
	Routes []*Route
}

type Request struct {
	*http.Request
	Params   map[string][]string
	Host     string
	Path     string
	Fragment string
}

type Response struct {
	http.ResponseWriter
}

type Context struct {
	Method           string
	Params           map[string][]string
	Form             map[string][]string
	Request          *Request
	Response         *Response
	Application      *Application
	AppEngineContext *appengine.Context
	View             map[string]interface{}
}

type Route struct {
	Pattern string
	Regexp  *regexp.Regexp
	Handler Handler
}

type RouteData struct {
	Path   string
	Route  *Route
	Params map[string]string
}

type Handler func(*Context)

var (
	routeParam *regexp.Regexp = regexp.MustCompile("<[^>]+>")
)

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
		if k+1 > len_rst {
			break
		}
		ng[v] = rst[k]
	}
	return ng
}

func NewApplication() *Application {
	app := new(Application)
	return app
}

func (app *Application) Start() {
	http.Handle("/", app)
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rd := app.FindRoute(r)
	c := NewContext(w, r, rd)
	if rd == nil {
		app.ErrorNotFound(c)
	} else {
		rd.Route.Handler(c)
	}
}

func (app *Application) FindRoute(r *http.Request) *RouteData {
	var rd *RouteData
	for _, route := range app.Routes {
		rd = route.Match(r.URL.Path)
		if rd != nil {
			return rd
		}
	}
	return nil
}

func (app *Application) ErrorNotFound(c *Context) {
}

func (app *Application) AddRoute(pattern string, handler Handler) {
	route := NewRoute(pattern)
	route.Handler = handler
	app.Routes = append(app.Routes, route)
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
		Method:   r.Method,
		Request:  req,
		Response: res,
		Params:   req.Params,
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
		body := s[1 : len(s)-1]
		if strings.Contains(body, ":") {
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

func (r *Route) Match(path string) *RouteData {
	params := NamedRegexpGroup(path, r.Regexp)
	if params == nil {
		return nil
	}
	rd := &RouteData{Path: path, Route: r, Params: params}
	return rd
}
