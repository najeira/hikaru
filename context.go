package hikaru

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/najeira/goutils/nlog"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type context struct {
	envContext
	*http.Request
	ResponseWriter http.ResponseWriter
	values         url.Values
	handlers       []HandlerFunc
	handlerIndex   int8
}

type Context struct {
	*context
	parent *Context
	key    interface{}
	value  interface{}
}

// Context should be http.ResponseWriter
var _ http.ResponseWriter = (*Context)(nil)

var (
	ErrKeyNotExist = errors.New("not exist")
	contextPool    sync.Pool
)

// Returns the Context.
func getContext() *Context {
	// try getting a Context from a pool.
	if v := contextPool.Get(); v != nil {
		return v.(*Context)
	}
	return &Context{context: &context{}}
}

// Release a Context.
func releaseContext(c *Context) {
	if c.parent != nil {
		releaseContext(c.parent)
	}
	c.parent = nil
	c.key = nil
	c.value = nil
	c.context.Request = nil
	c.context.ResponseWriter = nil
	c.context.values = nil
	c.context.handlers = nil
	c.context.handlerIndex = 0
	c.envContext.release()
	contextPool.Put(c)
}

func (c *Context) init(w http.ResponseWriter, r *http.Request, h []HandlerFunc) {
	if len(h) >= 127 {
		panic("handlers shold be less then 127")
	}
	c.parent = nil
	c.key = nil
	c.value = nil
	c.context.Request = r
	c.context.ResponseWriter = w
	c.context.values = nil
	c.context.handlers = h
	c.context.handlerIndex = 0
	c.envContext.init()
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
	// TODO: HTTP_X_FORWARDED_SSL, HTTP_X_FORWARDED_SCHEME, HTTP_X_FORWARDED_PROTO
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

func (c *Context) addParams(params httprouter.Params) {
	if params != nil {
		for _, v := range params {
			c.Add(v.Key, v.Value)
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
	c.SetHeader("Content-Type", "text/plain; charset=utf-8")
	return io.WriteString(c.ResponseWriter, body)
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
	c.WriteHeader(http.StatusNotModified)
}

// Sets response to HTTP 401 Unauthorized.
func (c *Context) Unauthorized() {
	c.WriteHeader(http.StatusUnauthorized)
}

// Sets response to HTTP 403 Forbidden.
func (c *Context) Forbidden() {
	c.WriteHeader(http.StatusForbidden)
}

// Sets response to HTTP 404 Not Found.
func (c *Context) NotFound() {
	c.WriteHeader(http.StatusNotFound)
}

// Sets response to HTTP 500 Internal Server Error.
func (c *Context) Fail(err interface{}) {
	c.WriteHeader(http.StatusInternalServerError)
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

func (c *Context) execute() {
	defer func() {
		if err := recover(); err != nil {
			c.WriteHeader(http.StatusInternalServerError)
			c.logPanic(err)
		}
	}()
	c.next()
}

func (c *Context) next() {
	if c.isGenLogEnabled(nlog.Debug) {
		// wrap ResponseWriter to handle status code.
		rw := responseWriter{c.ResponseWriter, http.StatusOK}
		c.ResponseWriter = &rw
		start := time.Now()
		c.Next()
		elapsed := time.Now().Sub(start)
		c.debugf("%3d | %12v | %4s %-7s", rw.statusCode, elapsed, c.Method, c.URL.Path)
	} else {
		c.Next()
	}
}

func (c *Context) WithValue(key interface{}, value interface{}) *Context {
	nc := getContext()
	nc.context = c.context
	nc.parent = c
	nc.key = c.key
	nc.value = c.value
	return nc
}

func (c *Context) Value(key interface{}) interface{} {
	if c.key == key {
		return c.value
	} else if c.parent != nil {
		return c.parent.Value(key)
	}
	return nil
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

var _ http.ResponseWriter = (*responseWriter)(nil)

func (r *responseWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
