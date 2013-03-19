package hikaru

import (
	"appengine"
	"appengine_internal"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type Context interface {
	appengine.Context
	AppEngineContext() appengine.Context
	Method() string
	Application() *Application
	HttpRequest() *http.Request
	ResponseWriter() http.ResponseWriter
	RouteData() *RouteData
	Result() Result
	SetResult(Result)
	Execute()
	Has(key string) bool
	Val(key string) string
	Vals(key string) []string
	Form(key string) string
	Forms(key string) []string
	IsPost() bool
	IsGet() bool
	Raw(body []byte, content_type string) Result
	Text(body string) Result
	Redirect(path string) Result
	RedirectFound(path string) Result
	Redirect301(path string) Result
	RedirectPermanently(path string) Result
	NotFound() Result
	AbortCode(code int) Result
	Abort(err interface{}) Result
	Html(args ...interface{}) Result
	//Json(args ...interface{}) Result
}

type HikaruContext struct {
	application      *Application
	httpRequest      *http.Request
	appEngineContext appengine.Context
	responseWriter   http.ResponseWriter
	routeData        *RouteData
	result           Result
}

// Creates and returns a new Context.
func NewContext(app *Application, w http.ResponseWriter, r *http.Request) *HikaruContext {
	ac := appengine.NewContext(r)
	c := &HikaruContext{
		application:      app,
		httpRequest:      r,
		appEngineContext: ac,
		responseWriter:   w,
	}
	return c
}

func (c *HikaruContext) Call(service, method string, in, out appengine_internal.ProtoMessage, opts *appengine_internal.CallOptions) error {
	return c.AppEngineContext().Call(service, method, in, out, opts)
}

func (c *HikaruContext) FullyQualifiedAppID() string {
	return c.AppEngineContext().FullyQualifiedAppID()
}

func (c *HikaruContext) Request() interface{} {
	return c.AppEngineContext().Request()
}

func (c *HikaruContext) AppEngineContext() appengine.Context {
	return c.appEngineContext
}

// Returns the request method.
func (c *HikaruContext) Method() string {
	return c.httpRequest.Method
}

// Returns the application.
func (c *HikaruContext) Application() *Application {
	return c.application
}

// Returns the *http.Request.
func (c *HikaruContext) HttpRequest() *http.Request {
	return c.httpRequest
}

// Returns the ResponseWriter.
func (c *HikaruContext) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (c *HikaruContext) RouteData() *RouteData {
	return c.routeData
}

func (c *HikaruContext) Result() Result {
	return c.result
}

func (c *HikaruContext) SetResult(result Result) {
	c.result = result
}

func (c *HikaruContext) Execute() {
	c.executeContext()
	c.executeResult()
}

func (c *HikaruContext) executeContext() {
	defer c.executeRecover()
	c.executeRequestFilters()
	if c.result == nil {
		c.executeRoute()
		if c.routeData != nil {
			c.executeHandler()
		} else {
			c.executeNotFound()
		}
	}
}

// Returns whether the request has the given key
// in route values and query.
func (c *HikaruContext) Has(key string) bool {
	if c.routeData != nil {
		_, ok := c.routeData.Params[key]
		if ok {
			return true
		}
	}
	_, ok := c.httpRequest.URL.Query()[key]
	return ok
}

// Returns the first value associated with the given key
// from route values and query.
// If there are no values associated with the key, returns "".
// To access multiple values of a key, use Vals.
func (c *HikaruContext) Val(key string) string {
	if c.routeData != nil {
		v, ok := c.routeData.Params[key]
		if ok {
			return v
		}
	}
	vs, ok2 := c.httpRequest.URL.Query()[key]
	if ok2 && len(vs) >= 1 {
		return vs[0]
	}
	return ""
}

// Returns the list of values associated with the given key
// from route values and query.
// If there are no values associated with the key, returns empty slice.
func (c *HikaruContext) Vals(key string) []string {
	if c.routeData != nil {
		v, ok := c.routeData.Params[key]
		if ok {
			return []string{v}
		}
	}
	vs, ok2 := c.httpRequest.URL.Query()[key]
	if ok2 {
		return vs
	}
	return []string{}
}

// Returns the first value associated with the given key from form.
// If there are no values associated with the key, returns "".
// To access multiple values of a key, use Forms.
func (c *HikaruContext) Form(key string) string {
	if c.httpRequest.Form == nil {
		c.httpRequest.ParseForm()
	}
	vs, ok := c.httpRequest.Form[key]
	if !ok || len(vs) <= 0 {
		return ""
	}
	return vs[0]
}

// Returns the list of values associated with the given key from form.
// If there are no values associated with the key, returns empty slice.
func (c *HikaruContext) Forms(key string) []string {
	if c.httpRequest.Form == nil {
		c.httpRequest.ParseForm()
	}
	vs, _ := c.httpRequest.Form[key]
	return vs
}

// Returns whether the request method is POST or not.
func (c *HikaruContext) IsPost() bool {
	return strings.ToUpper(c.Method()) == "POST"
}

// Returns whether the request method is GET or not.
func (c *HikaruContext) IsGet() bool {
	return strings.ToUpper(c.Method()) == "GET"
}

// Creates and returns a new Result with raw string and content type.
func (c *HikaruContext) Raw(body []byte, content_type string) Result {
	result := NewResult()
	result.statusCode = http.StatusOK
	result.body.Write(body)
	if content_type != "" {
		result.header.Set("Content-Type", content_type)
	}
	return result
}

// Creates and returns a new Result with text string.
// The content type should be "text/plain; charset=utf-8".
func (c *HikaruContext) Text(body string) Result {
	return c.Raw([]byte(body), "text/plain; charset=utf-8")
}

// Creates and returns a new Result with HTTP 302 Found.
func (c *HikaruContext) Redirect(path string) Result {
	return c.redirectCode(path, http.StatusFound)
}

// Creates and returns a new Result with HTTP 302 Found.
func (c *HikaruContext) RedirectFound(path string) Result {
	return c.Redirect(path)
}

// Creates and returns a new Result with HTTP 301 Moved Permanently.
func (c *HikaruContext) Redirect301(path string) Result {
	return c.redirectCode(path, http.StatusMovedPermanently)
}

// Creates and returns a new Result with HTTP 301 Moved Permanently.
func (c *HikaruContext) RedirectPermanently(path string) Result {
	return c.Redirect301(path)
}

func (c *HikaruContext) redirectCode(path string, code int) Result {
	result := NewResult()
	result.statusCode = code
	result.header.Set("Location", path)
	return result
}

// Creates and returns a new Result with HTTP 404 Not Found.
func (c *HikaruContext) NotFound() Result {
	return c.AbortCode(http.StatusNotFound)
}

// Creates and returns a new Result with the given code.
func (c *HikaruContext) AbortCode(code int) Result {
	result := NewResult()
	result.statusCode = code
	return result
}

// Creates and returns a new Result with the given error
// and HTTP 500 Internal Server Error.
func (c *HikaruContext) Abort(err interface{}) Result {
	result := NewResult()
	result.statusCode = http.StatusInternalServerError
	result.err = err
	return result
}

// Renders html template and returns a new Result with the
// rendered html content.
func (c *HikaruContext) Html(args ...interface{}) Result {
	// TODO: middlewares
	r := c.application.GetRenderer("html")
	if r == nil {
		panic(errors.New("no html renderer"))
	}
	result := r.Render(args...)
	// TODO: middlewares
	return result
}

func (c *HikaruContext) executeRoute() {
	c.routeData = c.application.Match(c.httpRequest)
}

func (c *HikaruContext) executeNotFound() {
	c.result = c.NotFound()
}

func (c *HikaruContext) executeRequestFilters() {
	for _, f := range c.application.RequestFilters {
		if res := f(c.httpRequest); res != nil {
			c.result = res
		}
	}
}

func (c *HikaruContext) executeHandlerFilters() {
	for _, f := range c.application.HandlerFilters {
		if res := f(c, c.routeData.Route.Handler()); res != nil {
			c.result = res
		}
	}
}

func (c *HikaruContext) executeErrorFilters() {
	for _, f := range c.application.ErrorFilters {
		if res := f(c, c.result); res != nil {
			c.result = res
		}
	}
}

func (c *HikaruContext) executeResponseFilters() {
	for _, f := range c.application.ResponseFilters {
		if res := f(c, c.result); res != nil {
			c.result = res
		}
	}
}

func (c *HikaruContext) executeRecover() {
	if err := recover(); err != nil {
		c.Errorln(err)
		c.result = c.resultPanic(err)
	}
}

func (c *HikaruContext) resultPanic(err interface{}) Result {
	var buf bytes.Buffer
	buf.Write(debug.Stack())
	stack := buf.String()
	err_msg := fmt.Sprintf("%v\n%s", err, stack)
	c.Errorf(err_msg)
	result := NewResult()
	result.statusCode = http.StatusInternalServerError
	result.err = err
	if c.application.Debug {
		result.body.WriteString(err_msg)
	}
	return result
}

func (c *HikaruContext) executeHandler() {
	rd := c.routeData
	r := rd.Route

	var to <-chan time.Time
	if r.Timeout() <= 0 {
		to = make(<-chan time.Time) // no timeout
	} else {
		to = time.After(r.Timeout())
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
		c.result = c.AbortCode(500)
	}
}

func (c *HikaruContext) executeHandlerWithRecover() {
	defer c.executeRecover()
	c.result = c.routeData.Route.Handler()(c)
}

func (c *HikaruContext) executeResult() {
	c.result.Execute(c)
}
