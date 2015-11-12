package hikaru

import (
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type Context struct {
	*envContext

	// request
	*http.Request
	params httprouter.Params
	query  url.Values
	form   url.Values

	// response
	Response *Response
}

var (
	ErrKeyNotExist = errors.New("not exist")
	contextPool    sync.Pool
)

// Returns the Context.
func getContext(w http.ResponseWriter, r *http.Request, params httprouter.Params) *Context {
	var c *Context = nil
	if v := contextPool.Get(); v != nil {
		c = v.(*Context)
	} else {
		c = &Context{
			envContext: &envContext{},
			Response:   &Response{},
		}
	}
	c.envContext.init(r)
	c.init(w, r, params)
	return c
}

// Release a Context.
func releaseContext(c *Context) {
	c.init(nil, nil, nil)
	c.envContext.release()
	contextPool.Put(c)
}

func (c *Context) init(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	c.Request = r
	c.params = params
	c.query = nil
	c.form = nil
	c.Response.init(w, r)
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

func (c *Context) QueryValues() url.Values {
	if c.query == nil {
		c.query = c.Request.URL.Query()
	}
	return c.query
}

func (c *Context) Query(key string) string {
	return c.QueryValues().Get(key)
}

func (c *Context) FormValues() url.Values {
	if c.form == nil {
		c.Request.ParseForm()
		c.form = c.Request.Form
	}
	return c.form
}

func (c *Context) Form(key string) string {
	return c.FormValues().Get(key)
}

func (c *Context) File(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
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

func (c *Context) getString(key string) (string, bool) {
	if c.params != nil {
		s := c.params.ByName(key)
		if s != "" {
			return s, true
		}
	}
	if ss, ok := c.QueryValues()[key]; ok && len(ss) > 0 {
		return ss[0], true
	}
	if ss, ok := c.FormValues()[key]; ok && len(ss) > 0 {
		return ss[0], true
	}
	return "", false
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key,
// Get returns the ErrKeyNotExist.
func (c *Context) TryString(key string) (string, error) {
	s, ok := c.getString(key)
	if ok {
		return s, nil
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
