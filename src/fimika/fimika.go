package fimika

import (
	"appengine"
	"bytes"
	"fmt"
	"net/http"
	"net/http/pprof"
	"regexp"
	"runtime/debug"
	"strings"
	"time"
)

type Application struct {
	Routes []*Route

	RootDir        string // project root dir
	StaticPath     string // static file dir, "static" if empty
	ViewPath       string // view file dir, "views" if empty
	HandlerTimeout time.Duration
	Logger         *Logger
	LogLevel       int
	Debug          bool
}

type Request struct {
	*http.Request
	Params   map[string][]string
	Path     string
	Fragment string
}

type Result struct {
	StatusCode int
	Header     http.Header
	Body       bytes.Buffer
	Error      interface{}
}

type Context struct {
	Method             string
	Application        *Application
	Request            *Request
	AppEngineContext   appengine.Context
	View               map[string]interface{}
	RouteData          *RouteData
	Result             *Result
	ResponseWriter     http.ResponseWriter
	ResponseStatusCode int // e.g. 200
}

type Route struct {
	Pattern string
	Regexp  *regexp.Regexp
	Handler Handler
}

type RouteData struct {
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
	app.Logger = NewLogger()
	return app
}

func (app *Application) Start() {
	http.Handle("/", app)
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rd := app.FindRoute(r)
	c := NewContext(app, w, r, rd)
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

func (app *Application) handlePanic(c *Context) *Result {
	// handle all the error
	err := recover()
	if err == nil {
		return nil
	}

	var buf bytes.Buffer
	buf.Write(debug.Stack())
	stack := buf.String()

	err_msg := fmt.Sprintf("%v\n%s", err, stack)

	app.Logger.Errorln(err_msg)

	result := c.Error(err)
	if app.Debug {
		result.Body.WriteString(err_msg)
	}
	return result
}

func (app *Application) ErrorNotFound(c *Context) {
}

func (app *Application) AddRoute(pattern string, handler Handler) {
	route := NewRoute(pattern)
	route.Handler = handler
	app.Routes = append(app.Routes, route)
}

func (app *Application) handlePprof(c *Context) bool {
	w := c.ResponseWriter
	r := c.Request.Request
	if !app.Debug {
		return false
	}
	switch r.RequestURI {
	case "/debug/pprof/cmdline":
		pprof.Cmdline(w, r)
		return true
	case "/debug/pprof/profile":
		pprof.Profile(w, r)
		return true
	case "/debug/pprof/heap":
		h := pprof.Handler("heap")
		h.ServeHTTP(w, r)
		return true
	case "/debug/pprof/symbol":
		pprof.Symbol(w, r)
		return true
	default:
		return false
	}
	return false
}

func NewRequest(r *http.Request) *Request {
	req := &Request{r, r.URL.Query(), r.URL.Path, r.URL.Fragment}
	return req
}

func (r *Request) Has(key string) bool {
	_, ok := r.Params[key]
	return ok
}

func (r *Request) Get(key string) string {
	vs, ok := r.Params[key]
	if !ok || len(vs) <= 0 {
		return ""
	}
	return vs[0]
}

func (r *Request) Gets(key string) []string {
	vs, _ := r.Params[key]
	return vs
}

func NewResult() *Result {
	result := new(Result)
	result.Header = make(http.Header)
	return result
}

func NewContext(app *Application, w http.ResponseWriter, r *http.Request, rd *RouteData) *Context {
	req := NewRequest(r)
	ac := appengine.NewContext(r)
	c := &Context{
		Method:             r.Method,
		Application:        app,
		Request:            req,
		AppEngineContext:   ac,
		RouteData:          rd,
		ResponseWriter:     w,
		ResponseStatusCode: http.StatusOK,
	}
	return c
}

func (c *Context) Has(key string) bool {
	_, ok := c.RouteData.Params[key]
	if ok {
		return true
	}
	return c.Request.Has(key)
}

func (c *Context) Get(key string) string {
	v, ok := c.RouteData.Params[key]
	if ok {
		return v
	}
	return c.Request.Get(key)
}

func (c *Context) Gets(key string) []string {
	v, ok := c.RouteData.Params[key]
	if ok {
		return []string{v}
	}
	return c.Request.Gets(key)
}

func (c *Context) Form(key string) string {
	vs, ok := c.Request.Form[key]
	if !ok || len(vs) <= 0 {
		return ""
	}
	return vs[0]
}

func (c *Context) Forms(key string) []string {
	vs, _ := c.Request.Form[key]
	return vs
}

func (c *Context) Raw(body string, content_type string) *Result {
	result := NewResult()
	result.StatusCode = http.StatusOK
	result.Body.WriteString(body)
	if content_type != "" {
		result.Header.Set("Content-Type", content_type)
	}
	return result
}

func (c *Context) Redirect(path string) *Result {
	return c.RedirectCode(path, http.StatusFound)
}

func (c *Context) RedirectFound(path string) *Result {
	return c.Redirect(path)
}

func (c *Context) RedirectPermanent(path string) *Result {
	return c.RedirectCode(path, http.StatusMovedPermanently)
}

func (c *Context) Redirect301(path string) *Result {
	return c.RedirectPermanent(path)
}

func (c *Context) RedirectCode(path string, code int) *Result {
	result := NewResult()
	result.StatusCode = code
	result.Header.Set("Location", path)
	return result
}

func (c *Context) doRedirect(result *Result) {
	http.Redirect(c.ResponseWriter, c.Request.Request, result.Header.Get("Location"), result.StatusCode)
}

func (c *Context) NotFound() *Result {
	return c.Abort(http.StatusNotFound)
}

func (c *Context) Abort(code int) *Result {
	result := NewResult()
	result.StatusCode = code
	return result
}

func (c *Context) Error(err interface{}) *Result {
	result := NewResult()
	result.StatusCode = http.StatusInternalServerError
	result.Error = err
	return result
}

func (c *Context) SetHeader(key string, val string) {
	c.ResponseWriter.Header().Set(key, val)
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.ResponseWriter, cookie)
}

func (c *Context) Render(template string) {
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
	rd := &RouteData{Route: r, Params: params}
	return rd
}
