package hikaru

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
)

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
	ret := ""
	ips := strings.Split(c.Request.RemoteAddr, ":")
	if ips != nil && len(ips) > 0 {
		if ips[0] != "[" {
			ret = ips[0]
		}
	}
	if strings.Contains(c.Application.Config.ProxyAddr, ret) {
		addrs := c.ProxyAddrs()
		if addrs != nil && len(addrs) > 0 && addrs[0] != "" {
			ret = strings.Split(addrs[0], ":")[0]
		}
	}
	return ret
}

func (c *Context) ProxyAddrs() []string {
	names := []string{"X-Forwarded-For", "X-Real-IP"}
	for _, name := range names {
		if ips := c.Request.Header.Get(name); ips != "" {
			return strings.Split(ips, ",")
		}
	}
	return nil
}

// Writes raw bytes and content type.
func (c *Context) Write(body []byte, content_type string) error {
	_, err := c.WriteBody(body)
	if err != nil {
		return err
	}
	if content_type != "" {
		c.SetHeader("Content-Type", content_type)
	}
	return nil
}

// Writes a text string.
// The content type should be "text/plain; charset=utf-8".
func (c *Context) Text(body string) error {
	return c.Write([]byte(body), "text/plain; charset=utf-8")
}

func (c *Context) Json(value interface{}) error {
	e := json.NewEncoder(c.body)
	if err := e.Encode(value); err != nil {
		//c.Fail(err)
		return err
	}
	c.SetHeader("Content-Type", "application/json; charset=utf-8")
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
	c.statusCode = code
	c.SetHeader("Location", path)
}

// Sets response to HTTP 304 Not Modified.
func (c *Context) NotModified() {
	c.Abort(http.StatusNotModified)
}

// Sets response to HTTP 401 Unauthorized.
func (c *Context) Unauthorized() {
	c.Abort(http.StatusUnauthorized)
}

// Sets response to HTTP 403 Forbidden.
func (c *Context) Forbidden() {
	c.Abort(http.StatusForbidden)
}

// Sets response to HTTP 404 Not Found.
func (c *Context) NotFound() {
	c.Abort(http.StatusNotFound)
}

// Sets response by the given code.
func (c *Context) Abort(code int) {
	c.statusCode = code
}

// Sets response to HTTP 500 Internal Server Error.
func (c *Context) Fail(err interface{}) {
	c.statusCode = http.StatusInternalServerError
}

func (c *Context) SetStatusCode(code int) {
	c.statusCode = code
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
	c.write(http.StatusInternalServerError, []byte(errMsg))
}

func (c *Context) String() string {
	return fmt.Sprintf("&{Context(Request=%s)}", c.Request)
}

func (c *Context) execute() {
	defer c.recover()
	c.logDebugf("[hikaru] execute: url is %v", c.Request.URL)
	c.Next()
	c.flush()
}
