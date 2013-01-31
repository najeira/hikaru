package hikaru

import (
	"appengine"
	"appengine_internal"
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

type Context struct {
	Method           string
	Application      *Application
	HttpRequest      *http.Request
	AppEngineContext appengine.Context
	ResponseWriter   http.ResponseWriter
	RouteData        *RouteData
	Result           Result
}

// Creates and returns a new Context.
func NewContext(app *Application, w http.ResponseWriter, r *http.Request) *Context {
	ac := appengine.NewContext(r)
	c := &Context{
		Method:           r.Method,
		Application:      app,
		HttpRequest:      r,
		AppEngineContext: ac,
		ResponseWriter:   w,
	}
	return c
}

func (c *Context) Call(service, method string, in, out proto.Message, opts *appengine_internal.CallOptions) error {
	return c.AppEngineContext.Call(service, method, in, out, opts)
}

func (c *Context) FullyQualifiedAppID() string {
	return c.AppEngineContext.FullyQualifiedAppID()
}

func (c *Context) Request() interface{} {
	return c.AppEngineContext.Request()
}

// Returns whether the request has the given key
// in route values and query.
func (c *Context) Has(key string) bool {
	_, ok := c.RouteData.Params[key]
	if ok {
		return true
	}
	_, ok = c.HttpRequest.URL.Query()[key]
	return ok
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
	vs, ok2 := c.HttpRequest.URL.Query()[key]
	if ok2 && len(vs) >= 1 {
		return vs[0]
	}
	return ""
}

// Returns the list of values associated with the given key
// from route values and query.
// If there are no values associated with the key, returns empty slice.
func (c *Context) Vals(key string) []string {
	v, ok := c.RouteData.Params[key]
	if ok {
		return []string{v}
	}
	vs, ok2 := c.HttpRequest.URL.Query()[key]
	if ok2 {
		return vs
	}
	return []string{}
}

// Returns the first value associated with the given key from form.
// If there are no values associated with the key, returns "".
// To access multiple values of a key, use Forms.
func (c *Context) Form(key string) string {
	if c.HttpRequest.Form == nil {
		c.HttpRequest.ParseForm()
	}
	vs, ok := c.HttpRequest.Form[key]
	if !ok || len(vs) <= 0 {
		return ""
	}
	return vs[0]
}

// Returns the list of values associated with the given key from form.
// If there are no values associated with the key, returns empty slice.
func (c *Context) Forms(key string) []string {
	if c.HttpRequest.Form == nil {
		c.HttpRequest.ParseForm()
	}
	vs, _ := c.HttpRequest.Form[key]
	return vs
}

// Creates and returns a new Result with raw string and content type.
func (c *Context) Raw(body string, content_type string) Result {
	result := NewResult()
	result.statusCode = http.StatusOK
	result.body.WriteString(body)
	if content_type != "" {
		result.header.Set("Content-Type", content_type)
	}
	return result
}

// Creates and returns a new Result with text string.
// The content type should be "text/plain; charset=utf-8".
func (c *Context) Text(body string) Result {
	return c.Raw(body, "text/plain; charset=utf-8")
}

// Creates and returns a new Result with HTTP 302 Found.
func (c *Context) Redirect(path string) Result {
	return c.redirectCode(path, http.StatusFound)
}

// Creates and returns a new Result with HTTP 302 Found.
func (c *Context) RedirectFound(path string) Result {
	return c.Redirect(path)
}

// Creates and returns a new Result with HTTP 301 Moved Permanently.
func (c *Context) Redirect301(path string) Result {
	return c.redirectCode(path, http.StatusMovedPermanently)
}

// Creates and returns a new Result with HTTP 301 Moved Permanently.
func (c *Context) RedirectPermanently(path string) Result {
	return c.Redirect301(path)
}

func (c *Context) redirectCode(path string, code int) Result {
	result := NewResult()
	result.statusCode = code
	result.header.Set("Location", path)
	return result
}

// Creates and returns a new Result with HTTP 404 Not Found.
func (c *Context) NotFound() Result {
	return c.AbortCode(http.StatusNotFound)
}

// Creates and returns a new Result with the given code.
func (c *Context) AbortCode(code int) Result {
	result := NewResult()
	result.statusCode = code
	return result
}

// Creates and returns a new Result with the given error
// and HTTP 500 Internal Server Error.
func (c *Context) Abort(err interface{}) Result {
	result := NewResult()
	result.statusCode = http.StatusInternalServerError
	result.err = err
	return result
}

func (c *Context) Html(name string, data interface{}) Result {
	app := c.Application

	// TODO: middlewares

	text := app.Renderer.Render(name, data)

	result := NewResult()
	result.statusCode = http.StatusOK
	result.body.WriteString(text)
	result.header.Set("Content-Type", "text/html; charset=utf-8")

	// TODO: middlewares

	return result
}

func (c *Context) executeRoute() bool {
	c.RouteData = c.matchRoute()
	return c.RouteData != nil
}

func (c *Context) matchRoute() *RouteData {
	var rd *RouteData
	for _, route := range c.Application.Routes {
		rd = route.Match(c.HttpRequest)
		if rd != nil {
			return rd
		}
	}
	return nil
}

func (c *Context) executeNotFound() {
	c.Result = c.NotFound()
}

func (c *Context) executeContext() {
	//TODO: before handler middlewares
	c.executeHandler()
	//TODO: after handler middlewares
}

func (c *Context) executeRecover() {
	if err := recover(); err != nil {
		c.Errorln(err)
		c.Result = c.resultPanic(err)
	}
}

func (c *Context) resultPanic(err interface{}) Result {
	var buf bytes.Buffer
	buf.Write(debug.Stack())
	stack := buf.String()
	err_msg := fmt.Sprintf("%v\n%s", err, stack)
	c.Errorf(err_msg)
	result := NewResult()
	result.statusCode = http.StatusInternalServerError
	result.err = err
	if c.Application.Debug {
		result.body.WriteString(err_msg)
	}
	return result
}

func (c *Context) executeHandler() {
	rd := c.RouteData
	r := rd.Route

	var to <-chan time.Time
	if r.Timeout <= 0 {
		to = make(<-chan time.Time) // no timeout
	} else {
		to = time.After(r.Timeout)
	}

	done := make(chan bool)
	go func() {
		c.executeHandlerWithRecover()
		done <- true
	}()

	select {
	case <-done:
		// succeeded
	case <-to:
		// timeouted
		c.Result = c.AbortCode(500)
	}
}

func (c *Context) executeHandlerWithRecover() {
	defer c.executeRecover()
	c.Result = c.RouteData.Route.Handler(c)
}

func (c *Context) executeResult() {
	c.Result.Execute(c)
}
