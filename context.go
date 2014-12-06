package hikaru

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

type context struct {
	*http.Request
	http.ResponseWriter
	Application  *Application
	values       url.Values
	handlers     []HandlerFunc
	handlerIndex int8
}

// Context should be http.ResponseWriter
var _ http.ResponseWriter = (*Context)(nil)

var (
	ErrKeyNotExist = errors.New("not exist")
	contextPool    sync.Pool
)

// Returns the Context.
func getContext(a *Application, w http.ResponseWriter, r *http.Request, h []HandlerFunc) *Context {
	var c *Context
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
	c.ResponseWriter = nil
	c.values = nil
	c.handlers = nil
	c.handlerIndex = 0
	contextPool.Put(c)
}

func (c *Context) init(a *Application, w http.ResponseWriter, r *http.Request, h []HandlerFunc) {
	if len(h) >= 128 {
		panic("handlers shold be less then 128")
	}
	c.Application = a
	c.Request = r
	c.ResponseWriter = w
	c.values = nil
	c.handlers = h
	c.handlerIndex = 0
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
	return c.Request.URL.Scheme == "https"
}

func (c *Context) IsUpload() bool {
	return strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data")
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

func (c *Context) Values() url.Values {
	if c.values == nil {
		c.Request.ParseForm()
		c.values = c.Request.Form
	}
	return c.values
}

func (c *Context) File(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
}

// Has returns whether the request has the given key in the route values and
// the query.
func (c *Context) Has(key string) bool {
	_, ok := c.Values()[key]
	return ok
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key,
// Get returns the failover string.
// To access multiple values, use the map directly.
func (c *Context) String(key string, failover string) string {
	ret, err := c.TryString(key)
	if err != nil {
		return failover
	}
	return ret
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key,
// Get returns the ErrKeyNotExist.
func (c *Context) TryString(key string) (string, error) {
	ss, ok := c.Values()[key]
	if ok && ss != nil && len(ss) > 0 {
		return ss[0], nil
	}
	return "", ErrKeyNotExist
}

func (c *Context) Int(key string, failover int64) int64 {
	ret, err := c.TryInt(key)
	if err != nil {
		return failover
	}
	return ret
}

func (c *Context) TryInt(key string) (int64, error) {
	s, err := c.TryString(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(s, 10, 64)
}

func (c *Context) Float(key string, failover float64) float64 {
	ret, err := c.TryFloat(key)
	if err != nil {
		return failover
	}
	return ret
}

func (c *Context) TryFloat(key string) (float64, error) {
	s, err := c.TryString(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

func (c *Context) Bool(key string, failover bool) bool {
	ret, err := c.TryBool(key)
	if err != nil {
		return failover
	}
	return ret
}

func (c *Context) TryBool(key string) (bool, error) {
	s, err := c.TryString(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(s)
}

// Set sets the key to value. It replaces any existing values.
func (c *Context) Set(key, value string) {
	c.Values()[key] = []string{value}
}

// Add adds the key to value. It appends to any existing values associated
// with key.
func (c *Context) Add(key, value string) {
	c.Values()[key] = append(c.Values()[key], value)
}

// Del deletes the values associated with key.
func (c *Context) Del(key string) {
	delete(c.Values(), key)
}

func (c *Context) Update(v url.Values) {
	for key, ss := range v {
		if ss != nil && len(ss) > 0 {
			for _, s := range ss {
				c.Values().Add(key, s)
			}
		}
	}
}

// Returns response headers.
func (c *Context) Header() http.Header {
	return c.ResponseWriter.Header()
}

// GetHeader gets a response header.
func (c *Context) GetHeader(key string) string {
	return c.ResponseWriter.Header().Get(key)
}

// SetHeader sets a response header.
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
	c.ResponseWriter.WriteHeader(code)
}

func (c *Context) SetStatusCode(code int) {
	c.WriteHeader(code)
}

func (c *Context) Write(msg []byte) (int, error) {
	return c.ResponseWriter.Write(msg)
}

// Writes raw bytes and content type.
func (c *Context) Raw(body []byte, contentType string) (int, error) {
	if contentType != "" {
		c.SetHeader("Content-Type", contentType)
	}
	return c.ResponseWriter.Write(body)
}

// Writes a text string.
// The content type should be "text/plain; charset=utf-8".
func (c *Context) Text(body string) (int, error) {
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
	s := int8(len(c.handlers))
	for c.handlerIndex < s {
		i := c.handlerIndex
		c.handlerIndex++
		c.handlers[i](c)
	}
}

func (c *Context) logPanic(err interface{}) {
	var buf bytes.Buffer
	buf.Write(debug.Stack())
	stack := buf.String()
	errMsg := fmt.Sprintf("%v\n%s", err, stack)
	c.Errorln(errMsg)
}

func (c *Context) execute() {
	defer func() {
		if err := recover(); err != nil {
			c.WriteHeader(http.StatusInternalServerError)
			c.logPanic(err)
		}
	}()
	c.Next()
}
