package hikaru

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"net/http"
)

var (
	_ http.ResponseWriter = (*Context)(nil)
	_ http.Hijacker       = (*Context)(nil)
	_ http.Flusher        = (*Context)(nil)
	_ http.CloseNotifier  = (*Context)(nil)
)

// Returns response headers.
func (c *Context) Header() http.Header {
	return c.ResponseWriter.Header()
}

// GetHeader gets a response header.
func (c *Context) GetHeader(key string) string {
	return c.Header().Get(key)
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.Header().Set(key, value)
}

// Adds a response header.
func (c *Context) AddHeader(key, value string) {
	c.Header().Add(key, value)
}

// Adds a cookie header.
func (c *Context) SetCookie(cookie *http.Cookie) {
	c.Header().Set("Set-Cookie", cookie.String())
}

func (c *Context) SetContentType(value string) {
	c.Header().Set("Content-Type", value)
}

func (c *Context) Status() int {
	return c.status
}

func (c *Context) Size() int {
	return c.size
}

func (c *Context) Written() bool {
	return c.size >= 0
}

func (c *Context) WriteHeader(code int) {
	if code > 0 && c.status != code {
		c.status = code
	}
}

func (c *Context) WriteHeaderAndSend(code int) {
	c.WriteHeader(code)
	c.writeHeaderIfNotSent()
}

func (c *Context) writeHeaderIfNotSent() {
	if !c.Written() {
		c.size = 0
		c.ResponseWriter.WriteHeader(c.status)
	}
}

func (c *Context) Write(msg []byte) (int, error) {
	c.writeHeaderIfNotSent()
	n, err := c.ResponseWriter.Write(msg)
	c.size += n
	return n, err
}

// Writes bytes and content type.
func (c *Context) WriteBody(body []byte, contentType string) (int, error) {
	if contentType != "" {
		c.SetContentType(contentType)
	}
	return c.Write(body)
}

// Writes a text string.
// The content type should be "text/plain; charset=utf-8".
func (c *Context) Text(body string) (int, error) {
	c.SetContentType("text/plain; charset=utf-8")
	return io.WriteString(c, body)
}

func (c *Context) Json(value interface{}) error {
	c.SetContentType("application/json; charset=utf-8")
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.Write(body)
	return nil
}

// Sets response to HTTP 3xx.
func (c *Context) Redirect(path string, code int) {
	c.SetHeader("Location", path)
	http.Redirect(c, c.Request, path, code)
	c.writeHeaderIfNotSent()
}

// Sets response to HTTP 304 Not Modified.
func (c *Context) NotModified() {
	c.WriteHeaderAndSend(http.StatusNotModified)
}

// Sets response to HTTP 401 Unauthorized.
func (c *Context) Unauthorized() {
	c.WriteHeaderAndSend(http.StatusUnauthorized)
}

// Sets response to HTTP 403 Forbidden.
func (c *Context) Forbidden() {
	c.WriteHeaderAndSend(http.StatusForbidden)
}

// Sets response to HTTP 404 Not Found.
func (c *Context) NotFound() {
	c.WriteHeaderAndSend(http.StatusNotFound)
}

// Sets response to HTTP 500 Internal Server Erroc.
func (c *Context) Fail(err interface{}) {
	c.WriteHeaderAndSend(http.StatusInternalServerError)
}

// Hijack lets the caller take over the connection.
func (c *Context) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if c.size < 0 {
		c.size = 0
	}
	return c.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify returns a channel that receives a single value
// when the client connection has gone away.
func (c *Context) CloseNotify() <-chan bool {
	return c.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush sends any buffered data to the client.
func (c *Context) Flush() {
	c.ResponseWriter.(http.Flusher).Flush()
}
