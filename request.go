package hikaru

import (
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
)

// ReqHeader gets a request header.
func (c *Context) ReqHeader(key string) string {
	return c.Request.Header.Get(key)
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

func getAddressWithoutPort(addr string) string {
	if index := strings.LastIndexByte(addr, ':'); index >= 0 {
		addr = addr[0:index]
	}
	return strings.TrimSpace(addr)
}

func (c *Context) RemoteAddr() string {
	return strings.TrimSpace(c.Request.RemoteAddr)
}

func (c *Context) ClientAddr() string {
	ip := c.ForwardedAddr()
	if len(ip) > 0 {
		return ip
	}
	return c.RemoteAddr()
}

func (c *Context) ForwardedAddr() string {
	if addrs := c.ForwardedAddrs(); len(addrs) > 0 {
		return addrs[0]
	}
	return ""
}

func (c *Context) ForwardedAddrs() []string {
	rets := make([]string, 0)
	names := []string{"X-Forwarded-For", "X-Real-IP"}
	for _, name := range names {
		if ips := c.Request.Header.Get(name); len(ips) > 0 {
			if arr := strings.Split(ips, ","); len(arr) > 0 {
				for _, ip := range arr {
					ip = strings.TrimSpace(ip)
					if len(ip) > 0 {
						rets = append(rets, ip)
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
	v, _ := c.getQuery(key)
	return v
}

func (c *Context) getQuery(key string) (string, bool) {
	if values := c.QueryValues()[key]; len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func (c *Context) Form(key string) string {
	v, _ := c.getForm(key)
	return v
}

func (c *Context) getForm(key string) (string, bool) {
	req := c.Request
	req.ParseMultipartForm(32 << 20) // 32 MB
	if values := req.PostForm[key]; len(values) > 0 {
		return values[0], true
	}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			return values[0], true
		}
	}
	return "", false
}

func (c *Context) File(key string) (multipart.File, *multipart.FileHeader, error) {
	req := c.Request
	req.ParseMultipartForm(32 << 20) // 32 MB
	return req.FormFile(key)
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
		if s := c.params.ByName(key); s != "" {
			return s, true
		}
	}
	if v, ok := c.getQuery(key); ok {
		return v, true
	}
	if v, ok := c.getForm(key); ok {
		return v, true
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
