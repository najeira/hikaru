package hikaru

import (
	"appengine"
	"bytes"
	"fmt"
	"net/http"
	"net/http/pprof"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type Context struct {
	Method           string
	Application      *Application
	Request          *http.Request
	AppEngineContext appengine.Context
	View             map[string]interface{}
	RouteData        *RouteData
	Result           Result
	ResponseWriter   http.ResponseWriter
}

// Creates and returns a new Context.
func NewContext(app *Application, w http.ResponseWriter, r *http.Request) *Context {
	ac := appengine.NewContext(r)
	c := &Context{
		Method:           r.Method,
		Application:      app,
		Request:          req,
		AppEngineContext: ac,
		ResponseWriter:   w,
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
	return c.Abort(http.StatusNotFound)
}

// Creates and returns a new Result with the given code.
func (c *Context) Abort(code int) Result {
	result := NewResult()
	result.statusCode = code
	return result
}

// Creates and returns a new Result with the given error
// and HTTP 500 Internal Server Error.
func (c *Context) Error(err interface{}) Result {
	result := NewResult()
	result.statusCode = http.StatusInternalServerError
	result.err = err
	return result
}

func (c *Context) Render(template string) {
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
		c.Result = r.Handler(c)
		done <- true
	}()

	select {
	case <-done:
		break
	case <-to:
		c.Result = c.ErrorTimeout()
		break
	}
}

func (c *Context) executeResult() {
	c.Result.Execute(c)
}
