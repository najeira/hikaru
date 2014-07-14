// +build !appengine

package hikaru

import (
	"bytes"
	"net/http"
	"sync"
)

type Context struct {
	Application *Application
	Request     *http.Request
	Values      Values

	handlers     []HandlerFunc
	handlerIndex int
	statusCode   int
	res          http.ResponseWriter
	body         *bytes.Buffer
	mu           sync.RWMutex
	closed       bool
}

// Initializes the Context.
func (c *Context) init(a *Application, w http.ResponseWriter, r *http.Request, h []HandlerFunc) {
	c.Application = a
	c.Request = r
	c.Values = Values(r.URL.Query())
	c.handlers = h
	c.statusCode = http.StatusOK
	c.res = w
}
