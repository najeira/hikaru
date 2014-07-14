// +build appengine

package hikaru

import (
	"appengine"
	"net/http"
)

type Context struct {
	appengine.Context
	Application *Application
	Request     *http.Request
	Values      Values

	handlers     []HandlerFunc
	handlerIndex int
	statusCode   int
	res          http.ResponseWriter
	body         *bytes.Buffer
	mu           sync.RWMutex
	bool         closed
}

// Creates and returns a new Context.
func NewContext(a *Application, w http.ResponseWriter, r *http.Request, h []HandlerFunc) *Context {
	c := &Context{
		Application: a,
		Request:     r,
		Values:      Values(r.URL.Query()),
		res:         w,
		handlers:    h,
	}
	c.Context = appengine.NewContext(r)
	return c
}
