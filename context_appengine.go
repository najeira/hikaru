// +build appengine

package hikaru

import (
	"appengine"
)

type Context struct {
	context
	appengine.Context
}

func (c *Context) initEnv() {
	// set the appengine.Context
	c.Context = appengine.NewContext(c.Request)
}
