package hikaru

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"sync"
)

var (
	ErrResponseClosed = errors.New("hikaru: response closed")
	bytesBufferPool   sync.Pool
)

func getBytesBuffer() *bytes.Buffer {
	if v := bytesBufferPool.Get(); v != nil {
		b := v.(*bytes.Buffer)
		b.Reset()
		return b
	}
	return &bytes.Buffer{}
}

func putBytesBuffer(b *bytes.Buffer) {
	b.Reset()
	bytesBufferPool.Put(b)
}

func (c *Context) Header() http.Header {
	return c.res.Header()
}

// Sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.res.Header().Set(key, value)
}

// Adds a response header.
func (c *Context) AddHeader(key, value string) {
	c.res.Header().Add(key, value)
}

// Adds a cookie header.
func (c *Context) SetCookie(cookie *http.Cookie) {
	c.res.Header().Set("Set-Cookie", cookie.String())
}

func (c *Context) WriteBody(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return 0, ErrResponseClosed
	}
	if c.body == nil {
		c.body = getBytesBuffer()
	}
	return c.body.Write(b)
}

func (c *Context) setClosed() bool {
	c.mu.Lock()
	closed := c.closed
	c.closed = true
	c.mu.Unlock()
	return !closed
}

func (c *Context) close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	if c.body != nil {
		putBytesBuffer(c.body)
		c.body = nil
	}
}

func (c *Context) write(code int, msg []byte) {
	if !c.setClosed() {
		c.logInfoln("[hikaru] write: already closed")
		return // already closed
	}
	c.res.WriteHeader(code)
	if msg != nil {
		c.res.Write(msg)
	}
	c.close()
}

func (c *Context) writeToWriter(out io.Writer) {
	if !c.setClosed() {
		c.logInfoln("[hikaru] writeToWriter: already closed")
		return // already closed
	}
	c.res.WriteHeader(c.statusCode)
	if c.body != nil && c.body.Len() > 0 {
		c.body.WriteTo(out)
	}
	c.close()
}

func (c *Context) writeRedirect(location string) {
	if !c.setClosed() {
		c.logInfoln("[hikaru] writeRedirect: already closed")
		return // already closed
	}
	http.Redirect(c.res, c.Request, location, c.statusCode)
	c.close()
}

func (c *Context) flush() {
	location := c.res.Header().Get("Location")
	if location != "" {
		c.writeRedirect(location)
	} else {
		c.writeToWriter(c.res)
	}
}
