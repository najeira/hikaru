package hikaru

import (
	"errors"
	"net/http"
	"net/url"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type Context struct {
	// request
	*http.Request
	params httprouter.Params
	query  url.Values

	// response
	http.ResponseWriter
	status int
	size   int
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
		c = &Context{}
	}
	c.init(w, r, params)
	return c
}

// Release a Context.
func releaseContext(c *Context) {
	c.init(nil, nil, nil)
	contextPool.Put(c)
}

func (c *Context) init(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	c.Request = r
	c.params = params
	c.query = nil
	c.ResponseWriter = w
	c.status = http.StatusOK
	c.size = -1
}
