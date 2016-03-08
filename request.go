package hikaru

import (
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// RequestHeader gets a request header.
func (c *Context) RequestHeader(key string) string {
	return c.Request.Header.Get(key)
}

// IsPost returns true if the request method is POST.
func (c *Context) IsPost() bool {
	m := strings.ToUpper(c.Request.Method)
	return m == "POST"
}

// IsGet returns true if the request method is GET.
func (c *Context) IsGet() bool {
	m := strings.ToUpper(c.Request.Method)
	return m == "GET"
}

// IsAjax returns true if the request is XMLHttpRequest.
func (c *Context) IsAjax() bool {
	h := c.Request.Header.Get("X-Requested-With")
	return h == "XMLHttpRequest"
}

// IsSecure returns true if the request is scure.
func (c *Context) IsSecure() bool {
	// TODO: HTTP_X_FORWARDED_SSL, HTTP_X_FORWARDED_SCHEME, HTTP_X_FORWARDED_PROTO
	return c.Request.URL.Scheme == "https"
}

// IsUpload returns true if the request has files.
func (c *Context) IsUpload() bool {
	ct := c.Request.Header.Get("Content-Type")
	return strings.Contains(ct, "multipart/form-data")
}

// RemoteAddr returns the address of the request.
func (c *Context) RemoteAddr() string {
	return strings.TrimSpace(c.Request.RemoteAddr)
}

// ClientAddr returns the address of the client.
func (c *Context) ClientAddr() string {
	if ip := c.ForwardedAddr(); len(ip) > 0 {
		return ip
	}
	return c.RemoteAddr()
}

// ForwardedAddr returns the address that is in
// X-Real-IP and X-Forwarded-For headers of the request.
func (c *Context) ForwardedAddr() string {
	if addrs := c.ForwardedAddrs(); len(addrs) > 0 {
		return addrs[0]
	}
	return ""
}

// ForwardedAddrs returns the addresses that are in
// X-Real-IP and X-Forwarded-For headers of the request.
func (c *Context) ForwardedAddrs() []string {
	rets := make([]string, 0)
	names := []string{"X-Real-IP", "X-Forwarded-For"}
	for _, name := range names {
		if ips := c.Request.Header.Get(name); len(ips) > 0 {
			if arr := strings.Split(ips, ","); len(arr) > 0 {
				for _, ip := range arr {
					if ip = strings.TrimSpace(ip); len(ip) > 0 {
						rets = append(rets, ip)
					}
				}
			}
		}
	}
	return rets
}

// Query returns the URL-encoded query values.
func (c *Context) Query() url.Values {
	if c.query == nil {
		c.query = c.Request.URL.Query()
	}
	return c.query
}

// queryValue returns the first value for the named component
// of the query and the value was found or not.
func (c *Context) queryValue(key string) (string, bool) {
	if values := c.Query()[key]; len(values) > 0 {
		return values[0], true
	}
	return "", false
}

// postFormValue returns the first value for the named component
// of the POST or PUT request body and the value was found or not.
// URL query parameters are ignored.
func (c *Context) postFormValue(key string) (string, bool) {
	if c.Request.PostForm == nil {
		c.ParseMultipartForm(32 << 20) // 32 MB
	}
	if values := c.PostForm[key]; len(values) > 0 {
		return values[0], true
	}
	return "", false
}

// FormFile returns the first file for the provided form key.
func (c *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	f, h, err := c.Request.FormFile(key)
	if err == http.ErrMissingFile {
		err = nil
	}
	return f, h, err
}

// Has returns true if the key is in the request.
func (c *Context) Has(key string) bool {
	_, ok := c.getString(key)
	return ok
}

// String gets the first value associated with the given key.
// If there are no values associated with the key,
// String returns the failover string.
// To access multiple values, use the map directly.
func (c *Context) String(key string, args ...string) string {
	if ret, err := c.TryString(key); err == nil {
		return ret
	} else if len(args) > 0 {
		return args[0]
	}
	return ""
}

func (c *Context) getString(key string) (string, bool) {
	if c.params != nil {
		if s := c.params.ByName(key); s != "" {
			return s, true
		}
	}
	if v, ok := c.queryValue(key); ok {
		return v, true
	} else if v, ok := c.postFormValue(key); ok {
		return v, true
	}
	return "", false
}

// TryString gets the first value associated with the given key.
// If there are no values associated with the key, returns the ErrKeyNotExist.
func (c *Context) TryString(key string) (string, error) {
	if s, ok := c.getString(key); ok {
		return s, nil
	}
	return "", ErrKeyNotExist
}

func (c *Context) Int(key string, args ...int64) int64 {
	if ret, err := c.TryInt(key); err == nil {
		return ret
	} else if len(args) > 0 {
		return args[0]
	}
	return 0
}

func (c *Context) TryInt(key string) (int64, error) {
	s, err := c.TryString(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(s, 10, 64)
}

func (c *Context) Float(key string, args ...float64) float64 {
	if ret, err := c.TryFloat(key); err == nil {
		return ret
	} else if len(args) > 0 {
		return args[0]
	}
	return 0
}

func (c *Context) TryFloat(key string) (float64, error) {
	s, err := c.TryString(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

func (c *Context) Bool(key string, args ...bool) bool {
	if ret, err := c.TryBool(key); err == nil {
		return ret
	} else if len(args) > 0 {
		return args[0]
	}
	return false
}

func (c *Context) TryBool(key string) (bool, error) {
	s, err := c.TryString(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(s)
}
