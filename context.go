package hikaru

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"sync"
	"errors"
)

type context struct {
	Application *Application
	Request     *http.Request
	Values      Values
	ResponseWriter          http.ResponseWriter
	handlers     []HandlerFunc
	handlerIndex int
}

// Context should be http.ResponseWriter
var _ http.ResponseWriter = (*Context)(nil)
var _ http.ResponseWriter = (*context)(nil)

var (
	contextPool        sync.Pool
)

// Returns the Context.
func getContext(a *Application, w http.ResponseWriter, r *http.Request, h []HandlerFunc) *context {
	var c *context
	// try getting a Context from a pool.
	if v := contextPool.Get(); v != nil {
		c = v.(*Context)
	} else {
		c = &Context{}
		c.context = context{}
	}
	c.init(a, w, r, h)
	c.initEnv()
	return c
}

// Release a Context.
func releaseContext(c *Context) {
	c.Application = nil
	c.Request = nil
	c.handlers = nil
	c.handlerIndex = 0
	c.ResponseWriter = nil
	c.statusCode = http.StatusOK
	c.responseWrote = false
	contextPool.Put(c)
}

func (c *Context) init(a *Application, w http.ResponseWriter, r *http.Request, h []HandlerFunc) {
	c.Application = a
	c.Request = r
	if c.Values != nil {
		// sets the URL and clear old Values if the Values already allocated.
		// reuse the Values make less allocations.
		c.Values.u = r.URL
		c.Values.v = nil
	} else {
		// allocate a new Values
		c.Values = NewValues(r.URL)
	}
	c.ResponseWriter = w
	c.handlers = h
	c.handlerIndex = 0
}

// Returns the request URL.
func (c *Context) URL() *url.URL {
	return c.Request.URL
}

// Returns the request path.
func (c *Context) Path() string {
	return c.Request.URL.Path
}

// Returns the request method.
func (c *Context) Method() string {
	return c.Request.Method
}

// Returns whether the request method is POST or not.
func (c *Context) IsPost() bool {
	return strings.ToUpper(c.Request.Method) == "POST"
}

// Returns whether the request method is GET or not.
func (c *Context) IsGet() bool {
	return strings.ToUpper(c.Request.Method) == "GET"
}

func (c *Context) IsAjax() bool {
	return c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (c *Context) IsSecure() bool {
	//HTTP_X_FORWARDED_SSL, HTTP_X_FORWARDED_SCHEME, HTTP_X_FORWARDED_PROTO
	return c.Request.URL.Scheme == "XMLHttpRequest"
}

func (c *Context) IsUpload() bool {
	return strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data")
}

func (c *Context) GetForm() (Values, error) {
	err := c.Request.ParseForm()
	if err != nil {
		return nil, err
	}
	return Values(c.Request.PostForm), nil
}

func (c *Context) GetMultipartForm() (*multipart.Form, error) {
	err := c.Request.ParseMultipartForm(1024 * 1024 * 1024)
	if err != nil {
		return nil, err
	}
	return c.Request.MultipartForm, nil
}

func (c *Context) RemoteAddr() string {
	ips := strings.Split(c.Request.RemoteAddr, ":")
	if ips != nil && len(ips) > 0 && ips[0] != "" && ips[0] != "[" {
		return ips[0]
	}
	return ""
}

func (c *Context) ForwardedAddr() string {
	addrs := c.ForwardedAddrs()
	if addrs != nil && len(addrs) > 0 {
		return addrs[0]
	}
	return ""
}

func (c *Context) ForwardedAddrs() []string {
	rets := make([]string, 0)
	names := []string{"X-Forwarded-For", "X-Real-IP"}
	for _, name := range names {
		if ips := c.Request.Header.Get(name); ips != "" {
			arr := strings.Split(ips, ",")
			if arr != nil {
				for _, ip := range arr {
					parts := strings.Split(ip, ":")
					if parts != nil && len(parts) > 0 && parts[0] != "" {
						rets = append(rets, parts[0])
					}
				}
			}
		}
	}
	return rets
}

// Returns response headers.
func (c *Context) Header() http.Header {
	return c.ResponseWriter.Header()
}

// Sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

// Adds a response header.
func (c *Context) AddHeader(key, value string) {
	c.ResponseWriter.Header().Add(key, value)
}

// Adds a cookie header.
func (c *Context) SetCookie(cookie *http.Cookie) {
	c.ResponseWriter.Header().Set("Set-Cookie", cookie.String())
}

func (c *Context) WriteHeader(code int) {
	return c.ResponseWriter.WriteHeader(code)
}

func (c *Context) SetStatusCode(code int) {
	return c.WriteHeader(code)
}

func (c *Context) Write(msg []byte) (int64, error) {
	return c.ResponseWriter.Write(msg)
}

// Writes raw bytes and content type.
func (c *Context) Raw(body []byte, contentType string) (int64, error) {
	if contentType != "" {
		c.SetHeader("Content-Type", contentType)
	}
	return c.ResponseWriter.Write(body)
}

// Writes a text string.
// The content type should be "text/plain; charset=utf-8".
func (c *Context) Text(body string) (int64, error) {
	return c.Raw([]byte(body), "text/plain; charset=utf-8")
}

func (c *Context) Json(value interface{}) error {
	c.SetHeader("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(c.ResponseWriter)
	if err := e.Encode(value); err != nil {
		return err
	}
	return nil
}

// Sets response to HTTP 302 Found.
func (c *Context) RedirectFound(path string) {
	c.Redirect(path, http.StatusFound)
}

// Sets response to HTTP 301 Moved Permanently.
func (c *Context) RedirectMoved(path string) {
	c.Redirect(path, http.StatusMovedPermanently)
}

// Sets response to HTTP 3xx.
func (c *Context) Redirect(path string, code int) {
	c.SetHeader("Location", path)
	http.Redirect(c.ResponseWriter, c.Request, path, code)
}

// Sets response to HTTP 304 Not Modified.
func (c *Context) NotModified() {
	c.ResponseWriter.WriteHeader(http.StatusNotModified)
}

// Sets response to HTTP 401 Unauthorized.
func (c *Context) Unauthorized() {
	c.ResponseWriter.WriteHeader(http.StatusUnauthorized)
}

// Sets response to HTTP 403 Forbidden.
func (c *Context) Forbidden() {
	c.ResponseWriter.WriteHeader(http.StatusForbidden)
}

// Sets response to HTTP 404 Not Found.
func (c *Context) NotFound() {
	c.ResponseWriter.WriteHeader(http.StatusNotFound)
}

// Sets response to HTTP 500 Internal Server Error.
func (c *Context) Fail(err interface{}) {
	c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
}

func (c *Context) Next() {
	if c.handlerIndex < 0 {
		return
	}
	s := len(c.handlers)
	for c.handlerIndex < s {
		i := c.handlerIndex
		c.handlerIndex++
		c.handlers[i](c)
	}
}

func (c *Context) recover() {
	if err := recover(); err != nil {
		c.handlePanic(err)
	}
}

func (c *Context) handlePanic(err interface{}) {
	var buf bytes.Buffer
	buf.Write(debug.Stack())
	stack := buf.String()
	errMsg := fmt.Sprintf("%v\n%s", err, stack)
	c.Errorln(errMsg)
	c.SetHeader("Content-Type", "text/plain; charset=utf-8")
	c.writeToResponse(http.StatusInternalServerError, []byte(errMsg))
}

func (c *Context) String() string {
	return fmt.Sprintf("&{Context(Request=%s)}", c.Request)
}

func (c *Context) execute() {
	defer c.recover()
	c.logDebugf("execute: url is %v", c.Request.URL)
	c.Next()
}
