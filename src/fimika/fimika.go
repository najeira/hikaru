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
	Debug          bool
	LogLevel       int
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
	Method           string
	Application      *Application
	Request          *Request
	AppEngineContext appengine.Context
	View             map[string]interface{}
	RouteData        *RouteData
	Result           *Result
	ResponseWriter   http.ResponseWriter
	Log              *Logger
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

type Handler func(*Context) *Result

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
	app.LogLevel = LogLevelInfo
	return app
}

func (app *Application) Start() {
	http.Handle("/", app)
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(app, w, r)

	c.RouteData = app.matchRoute(r)

	if c.RouteData == nil {
		app.ErrorNotFound(c)
	} else {
		app.executeContext(c)
	}

	app.executeResult(c)
}

func copyHeader(src, dst http.Header) {
	if src != nil {
		for k, vs := range src {
			if len(vs) >= 2 {
				for _, v := range vs {
					if v != "" {
						dst.Add(k, v)
					}
				}
			} else {
				v := vs[0]
				if v != "" {
					dst.Set(k, v)
				}
			}
		}
	}
}

func (app *Application) executeResult(c *Context) {
	result := c.Result

	copyHeader(result.Header, c.ResponseWriter.Header())

	if result.StatusCode > 0 {
		c.ResponseWriter.WriteHeader(result.StatusCode)
	}

	if result.Body.Len() > 0 {
		result.Body.WriteTo(c.ResponseWriter)
	}
}

func (app *Application) matchRoute(r *http.Request) *RouteData {
	var rd *RouteData
	for _, route := range app.Routes {
		rd = route.Match(r.URL.Path)
		if rd != nil {
			return rd
		}
	}
	return nil
}

func (app *Application) executeContext(c *Context) {
	defer func() {
		err_result := app.recoverPanic(c)
		if err_result != nil {
			c.Result = err_result
		}
	}()

	//TODO: before handler middlewares

	c.Result = c.RouteData.Route.Handler(c)

	//TODO: after handler middlewares
}

func (app *Application) recoverPanic(c *Context) *Result {
	err := recover()
	if err == nil {
		return nil
	}

	var buf bytes.Buffer
	buf.Write(debug.Stack())
	stack := buf.String()

	err_msg := fmt.Sprintf("%v\n%s", err, stack)

	c.Log.Errorf(err_msg)

	result := c.Error(err)
	if app.Debug {
		result.Body.WriteString(err_msg)
	}
	return result
}

func (app *Application) ErrorNotFound(c *Context) {
}

func (app *Application) Route(pattern string, handler Handler) {
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

func (r *Request) Val(key string) string {
	vs, ok := r.Params[key]
	if !ok || len(vs) <= 0 {
		return ""
	}
	return vs[0]
}

func (r *Request) Vals(key string) []string {
	vs, _ := r.Params[key]
	return vs
}

// Creates and returns a new Result.
func NewResult() *Result {
	result := new(Result)
	result.Header = make(http.Header)
	return result
}

// Creates and returns a new Context.
func NewContext(app *Application, w http.ResponseWriter, r *http.Request) *Context {
	req := NewRequest(r)
	ac := appengine.NewContext(r)
	lg := NewLogger(app.LogLevel, ac)
	c := &Context{
		Method:           r.Method,
		Application:      app,
		Request:          req,
		AppEngineContext: ac,
		ResponseWriter:   w,
		Log:              lg,
	}
	return c
}

// Returns whether the request has the given key
// in route values and query.
func (c *Context) Has(key string) bool {
	_, ok := c.RouteData.Params[key]
	if ok {
		return true
	}
	return c.Request.Has(key)
}

// Returns the first value associated with the given key
// from route values and query.
// If there are no values associated with the key, returns "".
// To access multiple values of a key, use Vals.
func (c *Context) Val(key string) string {
	v, ok := c.RouteData.Params[key]
	if ok {
		return v
	}
	return c.Request.Val(key)
}

// Returns the list of values associated with the given key
// from route values and query.
// If there are no values associated with the key, returns empty slice.
func (c *Context) Vals(key string) []string {
	v, ok := c.RouteData.Params[key]
	if ok {
		return []string{v}
	}
	return c.Request.Vals(key)
}

// Returns the first value associated with the given key from form.
// If there are no values associated with the key, returns "".
// To access multiple values of a key, use Forms.
func (c *Context) Form(key string) string {
	vs, ok := c.Request.Form[key]
	if !ok || len(vs) <= 0 {
		return ""
	}
	return vs[0]
}

// Returns the list of values associated with the given key from form.
// If there are no values associated with the key, returns empty slice.
func (c *Context) Forms(key string) []string {
	vs, _ := c.Request.Form[key]
	return vs
}

// Creates and returns a new Result with raw string and content type.
func (c *Context) Raw(body string, content_type string) *Result {
	result := NewResult()
	result.StatusCode = http.StatusOK
	result.Body.WriteString(body)
	if content_type != "" {
		result.Header.Set("Content-Type", content_type)
	}
	return result
}

// Creates and returns a new Result with text string.
// The content type should be "text/plain; charset=utf-8".
func (c *Context) Text(body string) *Result {
	return c.Raw(body, "text/plain; charset=utf-8")
}

// Creates and returns a new Result with HTTP 302 Found.
func (c *Context) Redirect(path string) *Result {
	return c.redirectCode(path, http.StatusFound)
}

// Creates and returns a new Result with HTTP 302 Found.
func (c *Context) RedirectFound(path string) *Result {
	return c.Redirect(path)
}

// Creates and returns a new Result with HTTP 301 Moved Permanently.
func (c *Context) Redirect301(path string) *Result {
	return c.redirectCode(path, http.StatusMovedPermanently)
}

// Creates and returns a new Result with HTTP 301 Moved Permanently.
func (c *Context) RedirectPermanently(path string) *Result {
	return c.Redirect301(path)
}

func (c *Context) redirectCode(path string, code int) *Result {
	result := NewResult()
	result.StatusCode = code
	result.Header.Set("Location", path)
	return result
}

func (c *Context) doRedirect(result *Result) {
	http.Redirect(c.ResponseWriter, c.Request.Request, result.Header.Get("Location"), result.StatusCode)
}

// Creates and returns a new Result with HTTP 404 Not Found.
func (c *Context) NotFound() *Result {
	return c.Abort(http.StatusNotFound)
}

// Creates and returns a new Result with the given code.
func (c *Context) Abort(code int) *Result {
	result := NewResult()
	result.StatusCode = code
	return result
}

// Creates and returns a new Result with the given error
// and HTTP 500 Internal Server Error.
func (c *Context) Error(err interface{}) *Result {
	result := NewResult()
	result.StatusCode = http.StatusInternalServerError
	result.Error = err
	return result
}

func (r *Result) SetHeader(key string, val string) {
}

func (r *Result) SetCookie(cookie *http.Cookie) {
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
