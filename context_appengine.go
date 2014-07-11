// +build appengine

package hikaru

import (
	"appengine"
	"net/http"
)

type Context struct {
	appengine.Context
	Application    *Application
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Values         Values
	Result         *Result
	handlers       []HandlerFunc
	handlerIndex   int
}

// Creates and returns a new Context.
func NewContext(a *Application, w http.ResponseWriter, r *http.Request, h []HandlerFunc) *Context {
	c := &Context{
		Application:    a,
		Request:        r,
		ResponseWriter: w,
		Values:         Values(r.URL.Query()),
		Result:         NewResult(),
		handlers:       h,
	}
	c.Context = appengine.NewContext(r)
	return c
}
