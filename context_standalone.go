// +build !appengine

package hikaru

type Context struct {
	context
}

func (c *Context) initEnv() {
	// nothing for standalone environment
}
